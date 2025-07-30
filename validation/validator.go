package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"audit-query-mcp-server/types"
	"audit-query-mcp-server/utils"
)

// Validator provides validation functionality for audit query parameters
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateQueryParams validates all query parameters
func ValidateQueryParams(params types.AuditQueryParams) error {
	// Validate log source
	if !utils.Contains(utils.ValidLogSources, params.LogSource) {
		return fmt.Errorf("invalid log source: %s", params.LogSource)
	}

	// Validate timeframe
	if params.Timeframe != "" {
		if !isValidTimeframe(params.Timeframe) {
			return fmt.Errorf("invalid timeframe: %s", params.Timeframe)
		}
	}

	// Validate resource types
	if params.Resource != "" {
		if !utils.Contains(utils.ValidResources, params.Resource) {
			return fmt.Errorf("invalid resource: %s", params.Resource)
		}
	}

	// Validate verbs
	if params.Verb != "" {
		if !isValidVerbPattern(params.Verb) {
			return fmt.Errorf("invalid verb: %s", params.Verb)
		}
	}

	// Validate namespace patterns
	if params.Namespace != "" {
		if !isValidNamespace(params.Namespace) {
			return fmt.Errorf("invalid namespace pattern: %s", params.Namespace)
		}
	}

	// Validate username patterns
	if params.Username != "" {
		if !isValidUsername(params.Username) {
			return fmt.Errorf("invalid username pattern: %s", params.Username)
		}
	}

	return nil
}

// ValidateGeneratedCommand performs final safety validation
func ValidateGeneratedCommand(command string) error {
	// Check for dangerous commands with whitelist exception for multi-file commands
	for _, pattern := range utils.DangerousPatterns {
		if strings.Contains(command, pattern) {
			// Check if this is a safe multi-file command pattern
			if isSafeMultiFileCommand(command, pattern) {
				continue // Allow this specific pattern
			}
			return fmt.Errorf("command contains dangerous pattern: %s", pattern)
		}
	}

	// Check for dangerous command substitution (but allow safe date commands)
	if strings.Contains(command, "$(") {
		// Allow only safe date commands (both Linux and macOS syntax)
		isSafe := false
		for _, safePattern := range utils.SafeDatePatterns {
			if strings.Contains(command, safePattern) {
				isSafe = true
				break
			}
		}
		if !isSafe {
			return fmt.Errorf("command contains dangerous command substitution")
		}
	}

	// Ensure it starts with oc adm node-logs (handle both single and multi-file commands)
	trimmedCommand := strings.TrimSpace(command)
	if !strings.HasPrefix(trimmedCommand, "oc adm node-logs") && !strings.HasPrefix(trimmedCommand, "(oc adm node-logs") {
		return fmt.Errorf("command must start with 'oc adm node-logs'")
	}

	return nil
}

// isSafeMultiFileCommand validates if a command with dangerous patterns is a safe multi-file command
func isSafeMultiFileCommand(command, dangerousPattern string) bool {
	// Only allow && and ; patterns for multi-file commands
	if dangerousPattern != "&&" && dangerousPattern != ";" {
		return false
	}

	// Must be wrapped in parentheses
	if !strings.HasPrefix(strings.TrimSpace(command), "(") || !strings.HasSuffix(strings.TrimSpace(command), ")") {
		return false
	}

	// Split by the dangerous pattern to get individual commands
	commands := strings.Split(command, dangerousPattern)
	if len(commands) < 2 {
		return false
	}

	// Validate each sub-command
	for _, subCommand := range commands {
		// Clean up the sub-command
		subCommand = strings.TrimSpace(subCommand)
		subCommand = strings.TrimPrefix(subCommand, "(")
		subCommand = strings.TrimSuffix(subCommand, ")")

		if subCommand == "" {
			continue // Skip empty commands
		}

		// Each sub-command must be a safe oc adm node-logs command
		if !isSafeOcAdmNodeLogsCommand(subCommand) {
			return false
		}
	}

	return true
}

// isSafeOcAdmNodeLogsCommand validates if a command is a safe oc adm node-logs command
func isSafeOcAdmNodeLogsCommand(command string) bool {
	command = strings.TrimSpace(command)

	// Must start with oc adm node-logs
	if !strings.HasPrefix(command, "oc adm node-logs") {
		return false
	}

	// Must contain --role=master
	if !strings.Contains(command, "--role=master") {
		return false
	}

	// Must contain --path= with a valid log path
	if !strings.Contains(command, "--path=") {
		return false
	}

	// Check for valid log paths only
	validLogPaths := []string{
		"--path=kube-apiserver/",
		"--path=oauth-server/",
		"--path=openshift-apiserver/",
		"--path=oauth-apiserver/",
		"--path=audit/",
	}

	hasValidPath := false
	for _, validPath := range validLogPaths {
		if strings.Contains(command, validPath) {
			hasValidPath = true
			break
		}
	}

	if !hasValidPath {
		return false
	}

	// Must not contain any dangerous patterns (recursive check, but exclude the current pattern)
	dangerousPatterns := []string{"oc delete", "oc apply", "oc create", "oc patch", "oc replace", "kubectl delete", "kubectl apply", "kubectl create", "kubectl patch", "kubectl replace", "`"}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(command, pattern) {
			return false
		}
	}

	// Must not contain command substitution except for safe date patterns
	if strings.Contains(command, "$(") {
		isSafe := false
		for _, safePattern := range utils.SafeDatePatterns {
			if strings.Contains(command, safePattern) {
				isSafe = true
				break
			}
		}
		if !isSafe {
			return false
		}
	}

	return true
}

// ValidateAuditResult validates an AuditResult instance
func ValidateAuditResult(result types.AuditResult) error {
	var errors []string

	// Validate QueryID
	if result.QueryID == "" {
		errors = append(errors, "QueryID is required")
	} else if len(result.QueryID) > 255 {
		errors = append(errors, "QueryID is too long (max 255 chars)")
	}

	// Validate Timestamp
	if result.Timestamp == "" {
		errors = append(errors, "Timestamp is required")
	} else {
		if _, err := time.Parse(time.RFC3339, result.Timestamp); err != nil {
			errors = append(errors, fmt.Sprintf("invalid timestamp format: %v", err))
		}
	}

	// Validate Command
	if result.Command == "" && result.Error == "" {
		errors = append(errors, "Command is required when no error is present")
	}
	if len(result.Command) > 10000 {
		errors = append(errors, "Command is too long (max 10000 chars)")
	}

	// Validate RawOutput
	if len(result.RawOutput) > 1000000 {
		errors = append(errors, "RawOutput is too large (max 1MB)")
	}

	// Validate ParsedData
	if result.ParsedData == nil {
		errors = append(errors, "ParsedData cannot be nil (use empty slice instead)")
	} else if len(result.ParsedData) > 100000 {
		errors = append(errors, "ParsedData has too many entries (max 100000)")
	}

	// Validate Summary
	if len(result.Summary) > 10000 {
		errors = append(errors, "Summary is too long (max 10000 chars)")
	}

	// Validate Error
	if len(result.Error) > 5000 {
		errors = append(errors, "Error message is too long (max 5000 chars)")
	}

	// Validate ExecutionTime
	if result.ExecutionTime < 0 {
		errors = append(errors, "ExecutionTime cannot be negative")
	}
	if result.ExecutionTime > 3600000 {
		errors = append(errors, "ExecutionTime is unreasonably high (max 1 hour in ms)")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ValidateAuditResultStrict performs strict validation with additional logical checks
func ValidateAuditResultStrict(result types.AuditResult) error {
	// First perform basic validation
	if err := ValidateAuditResult(result); err != nil {
		return err
	}

	var errors []string

	// Strict logical consistency checks
	if result.Error != "" && result.RawOutput != "" && !strings.Contains(result.Error, "timed out") {
		errors = append(errors, "RawOutput should be empty when there is an error (except for timeouts)")
	}

	if result.Error != "" && len(result.ParsedData) > 0 {
		errors = append(errors, "ParsedData should be empty when there is an error")
	}

	if result.Error == "" && result.Command == "" {
		errors = append(errors, "Command should be present when there is no error")
	}

	if result.Error == "" && result.Summary == "" && len(result.ParsedData) > 0 {
		errors = append(errors, "Summary should be present when there is parsed data and no error")
	}

	if len(errors) > 0 {
		return fmt.Errorf("strict validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ValidateStatusCode validates HTTP status codes
func ValidateStatusCode(statusCode string) bool {
	// Convert to integer for validation
	if code, err := strconv.Atoi(statusCode); err == nil {
		// Check if it's a valid HTTP status code range
		return code >= 100 && code <= 599
	}
	return false
}

// ValidateStatusCodeRange validates if a status code falls within a specific range
func ValidateStatusCodeRange(statusCode string, rangeName string) bool {
	if code, err := strconv.Atoi(statusCode); err == nil {
		if ranges, exists := utils.StatusCodeRanges[rangeName]; exists {
			for _, validCode := range ranges {
				if code == validCode {
					return true
				}
			}
		}
	}
	return false
}

// ValidateIPAddress validates IP address patterns
func ValidateIPAddress(ipAddress string) bool {
	for _, pattern := range IPAddressPatterns {
		matched, _ := regexp.MatchString(pattern, ipAddress)
		if matched {
			return true
		}
	}
	return false
}

// ValidateResourceName validates Kubernetes resource name patterns
func ValidateResourceName(name string) bool {
	// Check length constraints
	if len(name) < 1 || len(name) > 253 {
		return false
	}

	for _, pattern := range ResourceNamePatterns {
		matched, _ := regexp.MatchString(pattern, name)
		if matched {
			return true
		}
	}
	return false
}

// ValidateAPIGroup validates Kubernetes API group patterns
func ValidateAPIGroup(apiGroup string) bool {
	for _, pattern := range APIGroupPatterns {
		matched, _ := regexp.MatchString(pattern, apiGroup)
		if matched {
			return true
		}
	}
	return false
}

// ValidateAPIVersion validates Kubernetes API version patterns
func ValidateAPIVersion(apiVersion string) bool {
	for _, pattern := range APIVersionPatterns {
		matched, _ := regexp.MatchString(pattern, apiVersion)
		if matched {
			return true
		}
	}
	return false
}

// ValidateAuditLogField validates audit log field names
func ValidateAuditLogField(fieldName string) bool {
	for _, field := range utils.AuditLogFields {
		if field == fieldName {
			return true
		}
	}
	return false
}

// ValidateTimeFrameConstant validates timeframe constants
func ValidateTimeFrameConstant(timeframe string) bool {
	for _, constant := range utils.TimeFrameConstants {
		if constant == timeframe {
			return true
		}
	}
	return false
}

// isValidUsername validates username patterns for OpenShift authentication
func isValidUsername(username string) bool {
	for _, pattern := range UsernamePatterns {
		matched, _ := regexp.MatchString(pattern, username)
		if matched {
			return true
		}
	}
	return false
}

// isValidTimeframe checks if a timeframe string is valid using flexible parsing
func isValidTimeframe(timeframe string) bool {
	for _, pattern := range TimeframePatterns {
		matched, _ := regexp.MatchString(pattern, timeframe)
		if matched {
			return true
		}
	}
	return false
}

// isValidVerbPattern validates verb patterns, including pipe-separated patterns
func isValidVerbPattern(verb string) bool {
	// Handle pipe-separated verb patterns like "create|update|patch|delete"
	if strings.Contains(verb, "|") {
		verbs := strings.Split(verb, "|")
		for _, v := range verbs {
			if !utils.Contains(utils.ValidVerbs, strings.TrimSpace(v)) {
				return false
			}
		}
		return true
	}

	// Single verb validation
	return utils.Contains(utils.ValidVerbs, verb)
}

// isValidNamespace validates namespace patterns for Kubernetes/OpenShift
func isValidNamespace(namespace string) bool {
	// Check length constraints first (1-63 characters for namespaces)
	if len(namespace) < 1 || len(namespace) > 63 {
		return false
	}

	// Check for valid characters and format using consolidated patterns
	for _, pattern := range NamespacePatterns {
		matched, _ := regexp.MatchString(pattern, namespace)
		if matched {
			return true
		}
	}

	return false
}
