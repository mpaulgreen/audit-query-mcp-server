package parsing

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"audit-query-mcp-server/utils"
)

// EnhancedParserConfig holds configuration for the enhanced parser
type EnhancedParserConfig struct {
	UseJSONParsing   bool          `json:"use_json_parsing"`
	EnableFallback   bool          `json:"enable_fallback"`
	MaxLineLength    int           `json:"max_line_length"`
	MaxParseErrors   int           `json:"max_parse_errors"`
	Timeout          time.Duration `json:"timeout"`
	EnableValidation bool          `json:"enable_validation"`
	EnableMetrics    bool          `json:"enable_metrics"`
	JQAvailable      bool          `json:"jq_available"`
	FallbackToGrep   bool          `json:"fallback_to_grep"`
}

// DefaultEnhancedParserConfig returns the default enhanced parser configuration
func DefaultEnhancedParserConfig() EnhancedParserConfig {
	return EnhancedParserConfig{
		UseJSONParsing:   true,
		EnableFallback:   true,
		MaxLineLength:    100000, // 100KB
		MaxParseErrors:   1000,
		Timeout:          30 * time.Second,
		EnableValidation: true,
		EnableMetrics:    true,
		JQAvailable:      false, // Will be detected
		FallbackToGrep:   true,
	}
}

// EnhancedParseResult represents the result of enhanced parsing
type EnhancedParseResult struct {
	Entries          []AuditLogEntry  `json:"entries"`
	TotalLines       int              `json:"total_lines"`
	ParsedLines      int              `json:"parsed_lines"`
	ErrorLines       int              `json:"error_lines"`
	ParseErrors      []string         `json:"parse_errors"`
	ParseTime        time.Duration    `json:"parse_time"`
	Performance      ParsePerformance `json:"performance"`
	JSONParsedLines  int              `json:"json_parsed_lines"`
	GrepParsedLines  int              `json:"grep_parsed_lines"`
	FallbackUsed     bool             `json:"fallback_used"`
	AccuracyEstimate float64          `json:"accuracy_estimate"`
}

// EnhancedParser represents the enhanced parser with JSON-aware capabilities
type EnhancedParser struct {
	config EnhancedParserConfig
}

// NewEnhancedParser creates a new enhanced parser
func NewEnhancedParser(config EnhancedParserConfig) *EnhancedParser {
	// Auto-detect jq availability if not set
	if !config.JQAvailable {
		config.JQAvailable = checkJQAvailability()
	}

	return &EnhancedParser{
		config: config,
	}
}

// ParseAuditLogsEnhanced parses audit logs with enhanced JSON-aware capabilities
func (ep *EnhancedParser) ParseAuditLogsEnhanced(lines []string) EnhancedParseResult {
	startTime := time.Now()
	result := EnhancedParseResult{
		Entries:     make([]AuditLogEntry, 0, len(lines)),
		TotalLines:  len(lines),
		ParseErrors: make([]string, 0),
	}

	// Performance tracking
	var totalLineSize int
	errorCount := 0

	for i, line := range lines {
		// Check timeout
		if time.Since(startTime) > ep.config.Timeout {
			result.ParseErrors = append(result.ParseErrors, fmt.Sprintf("Parsing timeout after %v", ep.config.Timeout))
			break
		}

		// Check line length
		if len(line) > ep.config.MaxLineLength {
			result.ErrorLines++
			errorCount++
			result.ParseErrors = append(result.ParseErrors, fmt.Sprintf("Line %d: exceeds max length (%d > %d)", i+1, len(line), ep.config.MaxLineLength))
			continue
		}

		// Parse individual line with enhanced capabilities
		entry, parseMethod, err := ep.parseAuditLogLineEnhanced(line)
		if err != nil {
			result.ErrorLines++
			errorCount++
			if errorCount <= ep.config.MaxParseErrors {
				result.ParseErrors = append(result.ParseErrors, fmt.Sprintf("Line %d: %v", i+1, err))
			}
			continue
		}

		result.Entries = append(result.Entries, entry)
		result.ParsedLines++
		totalLineSize += len(line)

		// Track parsing method
		if parseMethod == "json" {
			result.JSONParsedLines++
		} else if parseMethod == "grep" {
			result.GrepParsedLines++
		}
	}

	// Calculate performance metrics
	result.ParseTime = time.Since(startTime)
	if result.ParseTime > 0 {
		result.Performance.LinesPerSecond = float64(result.ParsedLines) / result.ParseTime.Seconds()
	}
	if result.ParsedLines > 0 {
		result.Performance.AverageLineSize = totalLineSize / result.ParsedLines
	}

	// Calculate accuracy estimate
	result.AccuracyEstimate = ep.calculateAccuracyEstimate(result)

	return result
}

// parseAuditLogLineEnhanced parses a single audit log line with enhanced capabilities
func (ep *EnhancedParser) parseAuditLogLineEnhanced(line string) (AuditLogEntry, string, error) {
	entry := AuditLogEntry{
		RawLine:   line,
		ParseTime: time.Now(),
	}

	// Try JSON parsing first if enabled
	if ep.config.UseJSONParsing {
		if err := ep.parseJSONLineEnhanced(line, &entry); err == nil {
			// Validate entry if enabled
			if ep.config.EnableValidation {
				if err := validateEntry(&entry); err != nil {
					entry.ParseErrors = append(entry.ParseErrors, err.Error())
				}
			}
			return entry, "json", nil
		}
	}

	// Try enhanced structured parsing
	if err := ep.parseStructuredLineEnhanced(line, &entry); err == nil {
		// Validate entry if enabled
		if ep.config.EnableValidation {
			if err := validateEntry(&entry); err != nil {
				entry.ParseErrors = append(entry.ParseErrors, err.Error())
			}
		}
		return entry, "structured", nil
	}

	// Fallback to grep-based parsing if enabled
	if ep.config.FallbackToGrep {
		if err := ep.parseGrepLine(line, &entry); err == nil {
			// Validate entry if enabled
			if ep.config.EnableValidation {
				if err := validateEntry(&entry); err != nil {
					entry.ParseErrors = append(entry.ParseErrors, err.Error())
				}
			}
			return entry, "grep", nil
		}
	}

	return entry, "failed", fmt.Errorf("all parsing methods failed")
}

// parseJSONLineEnhanced attempts to parse the line as JSON with enhanced error handling
func (ep *EnhancedParser) parseJSONLineEnhanced(line string, entry *AuditLogEntry) error {
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(line), &rawData); err != nil {
		return fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	// Enhanced field extraction with better error handling
	ep.extractTimestampEnhanced(rawData, entry)
	ep.extractUserInfoEnhanced(rawData, entry)
	ep.extractVerbEnhanced(rawData, entry)
	ep.extractObjectRefEnhanced(rawData, entry)
	ep.extractResponseStatusEnhanced(rawData, entry)
	ep.extractRequestInfoEnhanced(rawData, entry)
	ep.extractAuthenticationInfoEnhanced(rawData, entry)
	ep.extractAdditionalInfoEnhanced(rawData, entry)

	return nil
}

// extractTimestampEnhanced extracts timestamp with enhanced error handling
func (ep *EnhancedParser) extractTimestampEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	// Try multiple timestamp field names
	timestampFields := []string{
		utils.AuditLogFields["RequestReceivedTimestamp"],
		"requestReceivedTimestamp",
		"timestamp",
		"time",
		"created",
	}

	for _, field := range timestampFields {
		if timestamp, ok := rawData[field].(string); ok && timestamp != "" {
			entry.Timestamp = timestamp
			return
		}
	}
}

// extractUserInfoEnhanced extracts user information with enhanced error handling
func (ep *EnhancedParser) extractUserInfoEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	// Try multiple user field structures
	userFields := []string{
		utils.AuditLogFields["User"],
		"user",
		"userInfo",
		"requestUser",
	}

	for _, field := range userFields {
		if userData, ok := rawData[field].(map[string]interface{}); ok {
			// Extract username
			usernameFields := []string{
				utils.AuditLogFields["Username"],
				"username",
				"name",
				"user",
			}
			for _, usernameField := range usernameFields {
				if username, ok := userData[usernameField].(string); ok && username != "" {
					entry.Username = username
					break
				}
			}

			// Extract UID
			uidFields := []string{
				utils.AuditLogFields["UID"],
				"uid",
				"id",
			}
			for _, uidField := range uidFields {
				if uid, ok := userData[uidField].(string); ok && uid != "" {
					entry.UID = uid
					break
				}
			}

			// Extract groups
			groupsFields := []string{
				utils.AuditLogFields["Groups"],
				"groups",
				"group",
			}
			for _, groupsField := range groupsFields {
				if groups, ok := userData[groupsField].([]interface{}); ok {
					entry.Groups = make([]string, len(groups))
					for i, group := range groups {
						if groupStr, ok := group.(string); ok {
							entry.Groups[i] = groupStr
						}
					}
					break
				}
			}

			// Extract extra information
			if extra, ok := userData[utils.AuditLogFields["Extra"]].(map[string]interface{}); ok {
				entry.Extra = extra
			}

			break
		}
	}

	// Try direct username fields
	directUsernameFields := []string{
		"impersonatedUser",
		"requestUser",
		"authenticatedUser",
	}

	for _, field := range directUsernameFields {
		if username, ok := rawData[field].(string); ok && username != "" {
			entry.Username = username
			break
		}
	}
}

// extractVerbEnhanced extracts verb with enhanced error handling
func (ep *EnhancedParser) extractVerbEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	verbFields := []string{
		utils.AuditLogFields["Verb"],
		"verb",
		"method",
		"action",
		"operation",
		"requestMethod",
		"httpMethod",
	}

	for _, field := range verbFields {
		if verb, ok := rawData[field].(string); ok && verb != "" {
			entry.Verb = verb
			return
		}
	}
}

// extractObjectRefEnhanced extracts object reference with enhanced error handling
func (ep *EnhancedParser) extractObjectRefEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	objectRefFields := []string{
		utils.AuditLogFields["ObjectRef"],
		"objectRef",
		"object",
		"resource",
	}

	for _, field := range objectRefFields {
		if objRef, ok := rawData[field].(map[string]interface{}); ok {
			// Extract resource
			resourceFields := []string{
				utils.AuditLogFields["Resource"],
				"resource",
				"kind",
				"type",
			}
			for _, resourceField := range resourceFields {
				if resource, ok := objRef[resourceField].(string); ok && resource != "" {
					entry.Resource = resource
					break
				}
			}

			// Extract namespace
			namespaceFields := []string{
				utils.AuditLogFields["Namespace"],
				"namespace",
				"ns",
			}
			for _, namespaceField := range namespaceFields {
				if namespace, ok := objRef[namespaceField].(string); ok && namespace != "" {
					entry.Namespace = namespace
					break
				}
			}

			// Extract name
			nameFields := []string{
				utils.AuditLogFields["Name"],
				"name",
				"id",
			}
			for _, nameField := range nameFields {
				if name, ok := objRef[nameField].(string); ok && name != "" {
					entry.Name = name
					break
				}
			}

			// Extract API group
			apiGroupFields := []string{
				utils.AuditLogFields["APIGroup"],
				"apiGroup",
				"group",
			}
			for _, apiGroupField := range apiGroupFields {
				if apiGroup, ok := objRef[apiGroupField].(string); ok && apiGroup != "" {
					entry.APIGroup = apiGroup
					break
				}
			}

			// Extract API version
			apiVersionFields := []string{
				utils.AuditLogFields["APIVersion"],
				"apiVersion",
				"version",
			}
			for _, apiVersionField := range apiVersionFields {
				if apiVersion, ok := objRef[apiVersionField].(string); ok && apiVersion != "" {
					entry.APIVersion = apiVersion
					break
				}
			}

			break
		}
	}
}

// extractResponseStatusEnhanced extracts response status with enhanced error handling
func (ep *EnhancedParser) extractResponseStatusEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	responseStatusFields := []string{
		utils.AuditLogFields["ResponseStatus"],
		"responseStatus",
		"status",
		"response",
	}

	for _, field := range responseStatusFields {
		if responseStatus, ok := rawData[field].(map[string]interface{}); ok {
			// Extract status code
			codeFields := []string{
				utils.AuditLogFields["Code"],
				"code",
				"statusCode",
				"httpCode",
			}
			for _, codeField := range codeFields {
				if code, ok := responseStatus[codeField].(float64); ok {
					entry.StatusCode = int(code)
					break
				} else if codeStr, ok := responseStatus[codeField].(string); ok {
					if code, err := strconv.Atoi(codeStr); err == nil {
						entry.StatusCode = code
						break
					}
				}
			}

			// Extract status message
			messageFields := []string{
				utils.AuditLogFields["Message"],
				"message",
				"statusMessage",
				"reason",
			}
			for _, messageField := range messageFields {
				if message, ok := responseStatus[messageField].(string); ok && message != "" {
					entry.StatusMessage = message
					break
				}
			}

			// Extract status reason
			reasonFields := []string{
				utils.AuditLogFields["Reason"],
				"reason",
				"statusReason",
			}
			for _, reasonField := range reasonFields {
				if reason, ok := responseStatus[reasonField].(string); ok && reason != "" {
					entry.StatusReason = reason
					break
				}
			}

			break
		}
	}
}

// extractRequestInfoEnhanced extracts request information with enhanced error handling
func (ep *EnhancedParser) extractRequestInfoEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	// Extract request URI
	requestURIFields := []string{
		utils.AuditLogFields["RequestURI"],
		"requestURI",
		"uri",
		"path",
		"url",
	}

	for _, field := range requestURIFields {
		if requestURI, ok := rawData[field].(string); ok && requestURI != "" {
			entry.RequestURI = requestURI
			break
		}
	}

	// Extract user agent
	userAgentFields := []string{
		utils.AuditLogFields["UserAgent"],
		"userAgent",
		"user-agent",
		"agent",
	}

	for _, field := range userAgentFields {
		if userAgent, ok := rawData[field].(string); ok && userAgent != "" {
			entry.UserAgent = userAgent
			break
		}
	}

	// Extract source IPs
	sourceIPsFields := []string{
		utils.AuditLogFields["SourceIPs"],
		"sourceIPs",
		"sourceIP",
		"clientIP",
		"remoteAddr",
	}

	for _, field := range sourceIPsFields {
		if sourceIPs, ok := rawData[field].([]interface{}); ok {
			entry.SourceIPs = make([]string, len(sourceIPs))
			for i, ip := range sourceIPs {
				if ipStr, ok := ip.(string); ok {
					entry.SourceIPs[i] = ipStr
				}
			}
			break
		} else if sourceIP, ok := rawData[field].(string); ok && sourceIP != "" {
			entry.SourceIPs = []string{sourceIP}
			break
		}
	}
}

// extractAuthenticationInfoEnhanced extracts authentication information with enhanced error handling
func (ep *EnhancedParser) extractAuthenticationInfoEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	// Extract authentication decision
	authDecisionFields := []string{
		utils.AuditLogFields["AuthenticationDecision"],
		"authenticationDecision",
		"authDecision",
		"auth_decision",
	}

	for _, field := range authDecisionFields {
		if authDecision, ok := rawData[field].(string); ok && authDecision != "" {
			entry.AuthDecision = authDecision
			break
		}
	}

	// Extract authorization decision
	authzDecisionFields := []string{
		utils.AuditLogFields["AuthorizationDecision"],
		"authorizationDecision",
		"authzDecision",
		"authz_decision",
	}

	for _, field := range authzDecisionFields {
		if authzDecision, ok := rawData[field].(string); ok && authzDecision != "" {
			entry.AuthzDecision = authzDecision
			break
		}
	}

	// Extract impersonated user
	impersonatedUserFields := []string{
		utils.AuditLogFields["ImpersonatedUser"],
		"impersonatedUser",
		"impersonated_user",
		"impersonate",
	}

	for _, field := range impersonatedUserFields {
		if impersonatedUser, ok := rawData[field].(string); ok && impersonatedUser != "" {
			entry.ImpersonatedUser = impersonatedUser
			break
		}
	}
}

// extractAdditionalInfoEnhanced extracts additional information with enhanced error handling
func (ep *EnhancedParser) extractAdditionalInfoEnhanced(rawData map[string]interface{}, entry *AuditLogEntry) {
	// Extract annotations
	annotationsFields := []string{
		utils.AuditLogFields["Annotations"],
		"annotations",
		"metadata",
	}

	for _, field := range annotationsFields {
		if annotations, ok := rawData[field].(map[string]interface{}); ok {
			entry.Annotations = annotations
			break
		}
	}

	// Extract headers
	headersFields := []string{
		"headers",
		"requestHeaders",
		"responseHeaders",
	}

	for _, field := range headersFields {
		if headers, ok := rawData[field].(map[string]interface{}); ok {
			entry.Headers = headers
			break
		}
	}
}

// parseStructuredLineEnhanced parses non-JSON lines using enhanced regex patterns
func (ep *EnhancedParser) parseStructuredLineEnhanced(line string, entry *AuditLogEntry) error {
	// Enhanced regex patterns with multiple field name variations
	patterns := ep.getEnhancedRegexPatterns()

	for fieldName, pattern := range patterns {
		if match := pattern.FindStringSubmatch(line); len(match) > 1 {
			ep.setFieldValue(entry, fieldName, match[1])
		}
	}

	return nil
}

// parseGrepLine parses lines using basic grep patterns as fallback
func (ep *EnhancedParser) parseGrepLine(line string, entry *AuditLogEntry) error {
	// Basic grep-based extraction for non-JSON lines
	// This is a simplified fallback method

	// Extract timestamp
	if timestamp := ep.extractWithGrep(line, `(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)`); timestamp != "" {
		entry.Timestamp = timestamp
	}

	// Extract username
	if username := ep.extractWithGrep(line, `"username":"([^"]+)"`); username != "" {
		entry.Username = username
	}

	// Extract verb
	if verb := ep.extractWithGrep(line, `"verb":"([^"]+)"`); verb != "" {
		entry.Verb = verb
	}

	// Extract resource
	if resource := ep.extractWithGrep(line, `"resource":"([^"]+)"`); resource != "" {
		entry.Resource = resource
	}

	// Extract namespace
	if namespace := ep.extractWithGrep(line, `"namespace":"([^"]+)"`); namespace != "" {
		entry.Namespace = namespace
	}

	// Extract status code
	if statusCode := ep.extractWithGrep(line, `"code":(\d+)`); statusCode != "" {
		if code, err := strconv.Atoi(statusCode); err == nil {
			entry.StatusCode = code
		}
	}

	return nil
}

// extractWithGrep extracts a value using a regex pattern
func (ep *EnhancedParser) extractWithGrep(line, pattern string) string {
	re := regexp.MustCompile(pattern)
	if match := re.FindStringSubmatch(line); len(match) > 1 {
		return match[1]
	}
	return ""
}

// getEnhancedRegexPatterns returns enhanced regex patterns for structured parsing
func (ep *EnhancedParser) getEnhancedRegexPatterns() map[string]*regexp.Regexp {
	return map[string]*regexp.Regexp{
		"timestamp":     regexp.MustCompile(`"requestReceivedTimestamp":"([^"]+)"`),
		"username":      regexp.MustCompile(`"username":"([^"]+)"`),
		"verb":          regexp.MustCompile(`"verb":"([^"]+)"`),
		"resource":      regexp.MustCompile(`"resource":"([^"]+)"`),
		"namespace":     regexp.MustCompile(`"namespace":"([^"]+)"`),
		"name":          regexp.MustCompile(`"name":"([^"]+)"`),
		"statusCode":    regexp.MustCompile(`"code":(\d+)`),
		"statusMessage": regexp.MustCompile(`"message":"([^"]+)"`),
		"requestURI":    regexp.MustCompile(`"requestURI":"([^"]+)"`),
		"userAgent":     regexp.MustCompile(`"userAgent":"([^"]+)"`),
	}
}

// setFieldValue sets a field value on the audit log entry
func (ep *EnhancedParser) setFieldValue(entry *AuditLogEntry, fieldName, value string) {
	switch fieldName {
	case "timestamp":
		entry.Timestamp = value
	case "username":
		entry.Username = value
	case "verb":
		entry.Verb = value
	case "resource":
		entry.Resource = value
	case "namespace":
		entry.Namespace = value
	case "name":
		entry.Name = value
	case "statusMessage":
		entry.StatusMessage = value
	case "requestURI":
		entry.RequestURI = value
	case "userAgent":
		entry.UserAgent = value
	case "statusCode":
		if code, err := strconv.Atoi(value); err == nil {
			entry.StatusCode = code
		}
	}
}

// calculateAccuracyEstimate calculates an accuracy estimate based on parsing method
func (ep *EnhancedParser) calculateAccuracyEstimate(result EnhancedParseResult) float64 {
	if result.ParsedLines == 0 {
		return 0.0
	}

	// JSON parsing is considered 95% accurate
	jsonAccuracy := float64(result.JSONParsedLines) * 0.95

	// Structured parsing is considered 85% accurate
	structuredAccuracy := float64(result.ParsedLines-result.JSONParsedLines-result.GrepParsedLines) * 0.85

	// Grep parsing is considered 70% accurate
	grepAccuracy := float64(result.GrepParsedLines) * 0.70

	totalAccuracy := jsonAccuracy + structuredAccuracy + grepAccuracy
	return totalAccuracy / float64(result.ParsedLines)
}

// checkJQAvailability checks if jq is available in the system
func checkJQAvailability() bool {
	// This would be implemented to check jq availability
	// For now, return false to use fallback
	return false
}
