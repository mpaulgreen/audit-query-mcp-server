package parsing

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"audit-query-mcp-server/utils"
)

// AuditLogEntry represents a structured audit log entry with proper typing
type AuditLogEntry struct {
	// Core fields
	Timestamp  string   `json:"timestamp,omitempty"`
	Username   string   `json:"username,omitempty"`
	UID        string   `json:"uid,omitempty"`
	Groups     []string `json:"groups,omitempty"`
	Verb       string   `json:"verb,omitempty"`
	Resource   string   `json:"resource,omitempty"`
	Namespace  string   `json:"namespace,omitempty"`
	Name       string   `json:"name,omitempty"`
	APIGroup   string   `json:"api_group,omitempty"`
	APIVersion string   `json:"api_version,omitempty"`
	RequestURI string   `json:"request_uri,omitempty"`
	UserAgent  string   `json:"user_agent,omitempty"`
	SourceIPs  []string `json:"source_ips,omitempty"`

	// Response fields
	StatusCode    int    `json:"status_code,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
	StatusReason  string `json:"status_reason,omitempty"`

	// Authentication fields
	AuthDecision     string `json:"auth_decision,omitempty"`
	AuthzDecision    string `json:"authz_decision,omitempty"`
	ImpersonatedUser string `json:"impersonated_user,omitempty"`

	// Additional fields
	Annotations map[string]interface{} `json:"annotations,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
	Headers     map[string]interface{} `json:"headers,omitempty"`

	// Metadata
	RawLine     string    `json:"raw_line,omitempty"`
	ParseErrors []string  `json:"parse_errors,omitempty"`
	ParseTime   time.Time `json:"parse_time,omitempty"`
}

// ParseResult represents the result of parsing audit logs
type ParseResult struct {
	Entries     []AuditLogEntry  `json:"entries"`
	TotalLines  int              `json:"total_lines"`
	ParsedLines int              `json:"parsed_lines"`
	ErrorLines  int              `json:"error_lines"`
	ParseErrors []string         `json:"parse_errors"`
	ParseTime   time.Duration    `json:"parse_time"`
	Performance ParsePerformance `json:"performance"`
}

// ParsePerformance tracks parsing performance metrics
type ParsePerformance struct {
	LinesPerSecond  float64 `json:"lines_per_second"`
	AverageLineSize int     `json:"average_line_size"`
	MemoryUsage     int64   `json:"memory_usage_bytes"`
}

// ParserConfig holds configuration for the parser
type ParserConfig struct {
	MaxLineLength    int           `json:"max_line_length"`
	MaxParseErrors   int           `json:"max_parse_errors"`
	Timeout          time.Duration `json:"timeout"`
	EnableValidation bool          `json:"enable_validation"`
	EnableMetrics    bool          `json:"enable_metrics"`
}

// DefaultParserConfig returns the default parser configuration
func DefaultParserConfig() ParserConfig {
	return ParserConfig{
		MaxLineLength:    100000, // 100KB
		MaxParseErrors:   1000,
		Timeout:          30 * time.Second,
		EnableValidation: true,
		EnableMetrics:    true,
	}
}

// ParseAuditLogs parses multiple audit log lines with enhanced error handling and performance tracking
func ParseAuditLogs(lines []string, config ParserConfig) ParseResult {
	startTime := time.Now()
	result := ParseResult{
		Entries:     make([]AuditLogEntry, 0, len(lines)),
		TotalLines:  len(lines),
		ParseErrors: make([]string, 0),
	}

	// Performance tracking
	var totalLineSize int
	errorCount := 0

	for i, line := range lines {
		// Check timeout
		if time.Since(startTime) > config.Timeout {
			result.ParseErrors = append(result.ParseErrors, fmt.Sprintf("Parsing timeout after %v", config.Timeout))
			break
		}

		// Check line length
		if len(line) > config.MaxLineLength {
			result.ErrorLines++
			errorCount++
			result.ParseErrors = append(result.ParseErrors, fmt.Sprintf("Line %d: exceeds max length (%d > %d)", i+1, len(line), config.MaxLineLength))
			continue
		}

		// Parse individual line
		entry, err := ParseAuditLogLine(line, config)
		if err != nil {
			result.ErrorLines++
			errorCount++
			if errorCount <= config.MaxParseErrors {
				result.ParseErrors = append(result.ParseErrors, fmt.Sprintf("Line %d: %v", i+1, err))
			}
			continue
		}

		result.Entries = append(result.Entries, entry)
		result.ParsedLines++
		totalLineSize += len(line)
	}

	// Calculate performance metrics
	result.ParseTime = time.Since(startTime)
	if result.ParseTime > 0 {
		result.Performance.LinesPerSecond = float64(result.ParsedLines) / result.ParseTime.Seconds()
	}
	if result.ParsedLines > 0 {
		result.Performance.AverageLineSize = totalLineSize / result.ParsedLines
	}

	return result
}

// ParseAuditLogLine parses a single audit log line with JSON parsing and error handling
func ParseAuditLogLine(line string, config ParserConfig) (AuditLogEntry, error) {
	entry := AuditLogEntry{
		RawLine:   line,
		ParseTime: time.Now(),
	}

	// Try JSON parsing first
	jsonErr := parseJSONLine(line, &entry)
	if jsonErr == nil {
		// Validate entry if enabled
		if config.EnableValidation {
			if err := validateEntry(&entry); err != nil {
				entry.ParseErrors = append(entry.ParseErrors, err.Error())
			}
		}
		return entry, nil
	}

	// For truly malformed JSON, return error instead of falling back to regex
	if strings.Contains(line, "{") && strings.Contains(line, "}") {
		return entry, fmt.Errorf("malformed JSON: %v", jsonErr)
	}

	// Fallback to structured parsing for non-JSON lines
	if err := parseStructuredLine(line, &entry); err != nil {
		return entry, fmt.Errorf("failed to parse line: %v", err)
	}

	// Validate entry if enabled
	if config.EnableValidation {
		if err := validateEntry(&entry); err != nil {
			entry.ParseErrors = append(entry.ParseErrors, err.Error())
		}
	}

	return entry, nil
}

// parseJSONLine attempts to parse the line as JSON
func parseJSONLine(line string, entry *AuditLogEntry) error {
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(line), &rawData); err != nil {
		return fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	// Extract timestamp
	if timestamp, ok := rawData[utils.AuditLogFields["RequestReceivedTimestamp"]].(string); ok {
		entry.Timestamp = timestamp
	}

	// Extract user information
	if userData, ok := rawData[utils.AuditLogFields["User"]].(map[string]interface{}); ok {
		if username, ok := userData[utils.AuditLogFields["Username"]].(string); ok {
			entry.Username = username
		}
		if uid, ok := userData[utils.AuditLogFields["UID"]].(string); ok {
			entry.UID = uid
		}
		if groups, ok := userData[utils.AuditLogFields["Groups"]].([]interface{}); ok {
			entry.Groups = make([]string, len(groups))
			for i, group := range groups {
				if groupStr, ok := group.(string); ok {
					entry.Groups[i] = groupStr
				}
			}
		}
		if extra, ok := userData[utils.AuditLogFields["Extra"]].(map[string]interface{}); ok {
			entry.Extra = extra
		}
	}

	// Extract verb
	if verb, ok := rawData[utils.AuditLogFields["Verb"]].(string); ok {
		entry.Verb = verb
	}

	// Extract object reference
	if objRef, ok := rawData[utils.AuditLogFields["ObjectRef"]].(map[string]interface{}); ok {
		if resource, ok := objRef[utils.AuditLogFields["Resource"]].(string); ok {
			entry.Resource = resource
		}
		if namespace, ok := objRef[utils.AuditLogFields["Namespace"]].(string); ok {
			entry.Namespace = namespace
		}
		if name, ok := objRef[utils.AuditLogFields["Name"]].(string); ok {
			entry.Name = name
		}
		if apiGroup, ok := objRef[utils.AuditLogFields["APIGroup"]].(string); ok {
			entry.APIGroup = apiGroup
		}
		if apiVersion, ok := objRef[utils.AuditLogFields["APIVersion"]].(string); ok {
			entry.APIVersion = apiVersion
		}
	}

	// Extract response status
	if responseStatus, ok := rawData[utils.AuditLogFields["ResponseStatus"]].(map[string]interface{}); ok {
		if code, ok := responseStatus[utils.AuditLogFields["Code"]].(float64); ok {
			entry.StatusCode = int(code)
		}
		if message, ok := responseStatus[utils.AuditLogFields["Message"]].(string); ok {
			entry.StatusMessage = message
		}
		if reason, ok := responseStatus[utils.AuditLogFields["Reason"]].(string); ok {
			entry.StatusReason = reason
		}
	}

	// Extract request URI
	if requestURI, ok := rawData[utils.AuditLogFields["RequestURI"]].(string); ok {
		entry.RequestURI = requestURI
	}

	// Extract user agent
	if userAgent, ok := rawData[utils.AuditLogFields["UserAgent"]].(string); ok {
		entry.UserAgent = userAgent
	}

	// Extract source IPs
	if sourceIPs, ok := rawData[utils.AuditLogFields["SourceIPs"]].([]interface{}); ok {
		entry.SourceIPs = make([]string, len(sourceIPs))
		for i, ip := range sourceIPs {
			if ipStr, ok := ip.(string); ok {
				entry.SourceIPs[i] = ipStr
			}
		}
	}

	// Extract annotations
	if annotations, ok := rawData[utils.AuditLogFields["Annotations"]].(map[string]interface{}); ok {
		entry.Annotations = annotations
	}

	// Extract authentication decisions
	if authDecision, ok := rawData[utils.AuditLogFields["AuthenticationDecision"]].(string); ok {
		entry.AuthDecision = authDecision
	}
	if authzDecision, ok := rawData[utils.AuditLogFields["AuthorizationDecision"]].(string); ok {
		entry.AuthzDecision = authzDecision
	}
	if impersonatedUser, ok := rawData[utils.AuditLogFields["ImpersonatedUser"]].(string); ok {
		entry.ImpersonatedUser = impersonatedUser
	}

	return nil
}

// parseStructuredLine parses non-JSON lines using regex patterns (fallback method)
func parseStructuredLine(line string, entry *AuditLogEntry) error {
	// Extract timestamp using the field constant
	timestampRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["RequestReceivedTimestamp"]))
	if match := timestampRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.Timestamp = match[1]
	}

	// Extract username using the field constant
	usernameRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["Username"]))
	if match := usernameRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.Username = match[1]
	}

	// Extract verb using the field constant
	verbRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["Verb"]))
	if match := verbRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.Verb = match[1]
	}

	// Extract resource using the field constant
	resourceRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["Resource"]))
	if match := resourceRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.Resource = match[1]
	}

	// Extract namespace using the field constant
	namespaceRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["Namespace"]))
	if match := namespaceRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.Namespace = match[1]
	}

	// Extract response status code using the field constant
	statusCodeRegex := regexp.MustCompile(fmt.Sprintf(`"%s":\s*{\s*"%s":\s*(\d+)`, utils.AuditLogFields["ResponseStatus"], utils.AuditLogFields["Code"]))
	if match := statusCodeRegex.FindStringSubmatch(line); len(match) > 1 {
		if code, err := strconv.Atoi(match[1]); err == nil {
			entry.StatusCode = code
		}
	}

	// Extract response status message using the field constant
	statusMessageRegex := regexp.MustCompile(fmt.Sprintf(`"%s":\s*{\s*[^}]*"%s":\s*"([^"]+)"`, utils.AuditLogFields["ResponseStatus"], utils.AuditLogFields["Message"]))
	if match := statusMessageRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.StatusMessage = match[1]
	}

	// Extract request URI using the field constant
	requestURIRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["RequestURI"]))
	if match := requestURIRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.RequestURI = match[1]
	}

	// Extract user agent using the field constant
	userAgentRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["UserAgent"]))
	if match := userAgentRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.UserAgent = match[1]
	}

	// Extract source IPs using the field constant
	sourceIPsRegex := regexp.MustCompile(fmt.Sprintf(`"%s":\s*\[([^\]]+)\]`, utils.AuditLogFields["SourceIPs"]))
	if match := sourceIPsRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.SourceIPs = strings.Split(match[1], ",")
	}

	// Extract authentication decision using the field constant
	authDecisionRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["AuthenticationDecision"]))
	if match := authDecisionRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.AuthDecision = match[1]
	}

	// Extract authorization decision using the field constant
	authzDecisionRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["AuthorizationDecision"]))
	if match := authzDecisionRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.AuthzDecision = match[1]
	}

	// Extract impersonated user using the field constant
	impersonatedUserRegex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, utils.AuditLogFields["ImpersonatedUser"]))
	if match := impersonatedUserRegex.FindStringSubmatch(line); len(match) > 1 {
		entry.ImpersonatedUser = match[1]
	}

	return nil
}

// validateEntry validates the parsed entry
func validateEntry(entry *AuditLogEntry) error {
	var errors []string

	// Validate timestamp format
	if entry.Timestamp != "" {
		if _, err := time.Parse(time.RFC3339, entry.Timestamp); err != nil {
			errors = append(errors, fmt.Sprintf("invalid timestamp format: %s", entry.Timestamp))
		}
	}

	// Validate status code range
	if entry.StatusCode != 0 && (entry.StatusCode < 100 || entry.StatusCode > 599) {
		errors = append(errors, fmt.Sprintf("invalid status code: %d", entry.StatusCode))
	}

	// Validate required fields (if any)
	if entry.Username == "" && entry.UID == "" {
		errors = append(errors, "missing user identification")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GenerateSummary creates a human-readable summary of the results
func GenerateSummary(entries []AuditLogEntry, context map[string]interface{}) string {
	if len(entries) == 0 {
		return "No audit entries found matching the criteria."
	}

	summary := fmt.Sprintf("Found %d audit entries", len(entries))

	// Count by username
	userCounts := make(map[string]int)
	for _, entry := range entries {
		if entry.Username != "" {
			userCounts[entry.Username]++
		}
	}

	if len(userCounts) > 0 {
		summary += ". Users involved: "
		userList := make([]string, 0, len(userCounts))
		for user, count := range userCounts {
			userList = append(userList, fmt.Sprintf("%s (%d)", user, count))
		}
		summary += strings.Join(userList, ", ")
	}

	// Count by status code
	statusCounts := make(map[int]int)
	for _, entry := range entries {
		if entry.StatusCode != 0 {
			statusCounts[entry.StatusCode]++
		}
	}

	if len(statusCounts) > 0 {
		summary += ". Status codes: "
		statusList := make([]string, 0, len(statusCounts))
		for status, count := range statusCounts {
			statusList = append(statusList, fmt.Sprintf("%d (%d)", status, count))
		}
		summary += strings.Join(statusList, ", ")
	}

	// Count by verb
	verbCounts := make(map[string]int)
	for _, entry := range entries {
		if entry.Verb != "" {
			verbCounts[entry.Verb]++
		}
	}

	if len(verbCounts) > 0 {
		summary += ". Actions: "
		verbList := make([]string, 0, len(verbCounts))
		for verb, count := range verbCounts {
			verbList = append(verbList, fmt.Sprintf("%s (%d)", verb, count))
		}
		summary += strings.Join(verbList, ", ")
	}

	// Count by resource
	resourceCounts := make(map[string]int)
	for _, entry := range entries {
		if entry.Resource != "" {
			resourceCounts[entry.Resource]++
		}
	}

	if len(resourceCounts) > 0 {
		summary += ". Resources: "
		resourceList := make([]string, 0, len(resourceCounts))
		for resource, count := range resourceCounts {
			resourceList = append(resourceList, fmt.Sprintf("%s (%d)", resource, count))
		}
		summary += strings.Join(resourceList, ", ")
	}

	return summary
}

// ParseAuditLogField extracts a specific field from audit log JSON (legacy function)
func ParseAuditLogField(line string, fieldName string) (string, bool) {
	// Get the actual field name from constants
	actualFieldName, exists := utils.AuditLogFields[fieldName]
	if !exists {
		return "", false
	}

	// Try JSON parsing first
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(line), &rawData); err == nil {
		// Handle nested fields
		if fieldName == "Code" || fieldName == "Message" || fieldName == "Reason" {
			if responseStatus, ok := rawData[utils.AuditLogFields["ResponseStatus"]].(map[string]interface{}); ok {
				if value, ok := responseStatus[actualFieldName].(string); ok {
					return value, true
				}
				if value, ok := responseStatus[actualFieldName].(float64); ok {
					return fmt.Sprintf("%.0f", value), true
				}
			}
		} else if fieldName == "Username" || fieldName == "UID" || fieldName == "Groups" {
			if userData, ok := rawData[utils.AuditLogFields["User"]].(map[string]interface{}); ok {
				if value, ok := userData[actualFieldName].(string); ok {
					return value, true
				}
			}
		} else if fieldName == "Resource" || fieldName == "Namespace" || fieldName == "Name" || fieldName == "APIGroup" || fieldName == "APIVersion" {
			if objRef, ok := rawData[utils.AuditLogFields["ObjectRef"]].(map[string]interface{}); ok {
				if value, ok := objRef[actualFieldName].(string); ok {
					return value, true
				}
			}
		} else {
			// Direct field access
			if value, ok := rawData[actualFieldName].(string); ok {
				return value, true
			}
		}
	}

	// Fallback to regex parsing
	pattern := fmt.Sprintf(`"%s":"([^"]+)"`, actualFieldName)
	regex := regexp.MustCompile(pattern)

	if match := regex.FindStringSubmatch(line); len(match) > 1 {
		return match[1], true
	}

	return "", false
}

// ParseStatusCodes extracts and categorizes status codes from audit logs
func ParseStatusCodes(entries []AuditLogEntry) map[string]int {
	statusCounts := make(map[string]int)

	for _, entry := range entries {
		if entry.StatusCode != 0 {
			category := categorizeStatusCode(entry.StatusCode)
			statusCounts[category]++
		}
	}

	return statusCounts
}

// categorizeStatusCode categorizes HTTP status codes
func categorizeStatusCode(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "success"
	case code >= 400 && code < 500:
		return "client_error"
	case code >= 500 && code < 600:
		return "server_error"
	default:
		return "other"
	}
}

// ConvertLegacyEntries converts legacy map[string]interface{} entries to AuditLogEntry
func ConvertLegacyEntries(legacyEntries []map[string]interface{}) []AuditLogEntry {
	entries := make([]AuditLogEntry, len(legacyEntries))

	for i, legacy := range legacyEntries {
		entry := AuditLogEntry{
			RawLine:   fmt.Sprintf("%v", legacy),
			ParseTime: time.Now(),
		}

		// Convert fields
		if timestamp, ok := legacy["timestamp"].(string); ok {
			entry.Timestamp = timestamp
		}
		if username, ok := legacy["username"].(string); ok {
			entry.Username = username
		}
		if verb, ok := legacy["verb"].(string); ok {
			entry.Verb = verb
		}
		if resource, ok := legacy["resource"].(string); ok {
			entry.Resource = resource
		}
		if namespace, ok := legacy["namespace"].(string); ok {
			entry.Namespace = namespace
		}
		if statusCode, ok := legacy["status_code"].(string); ok {
			if code, err := strconv.Atoi(statusCode); err == nil {
				entry.StatusCode = code
			}
		}
		if statusMessage, ok := legacy["status_message"].(string); ok {
			entry.StatusMessage = statusMessage
		}
		if requestURI, ok := legacy["request_uri"].(string); ok {
			entry.RequestURI = requestURI
		}
		if userAgent, ok := legacy["user_agent"].(string); ok {
			entry.UserAgent = userAgent
		}
		if sourceIPs, ok := legacy["source_ips"].([]string); ok {
			entry.SourceIPs = sourceIPs
		}
		if authDecision, ok := legacy["auth_decision"].(string); ok {
			entry.AuthDecision = authDecision
		}
		if authzDecision, ok := legacy["authz_decision"].(string); ok {
			entry.AuthzDecision = authzDecision
		}
		if impersonatedUser, ok := legacy["impersonated_user"].(string); ok {
			entry.ImpersonatedUser = impersonatedUser
		}

		entries[i] = entry
	}

	return entries
}
