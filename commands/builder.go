package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"audit-query-mcp-server/types"
	"os/exec"
)

// LogFileInfo represents information about a log file
type LogFileInfo struct {
	Path      string
	Date      time.Time
	IsCurrent bool
}

// CommandBuilder represents the new command builder for Phase 1 migration
type CommandBuilder struct {
	Config    types.AuditQueryConfig
	Migration types.MigrationConfig
	Discovery types.DiscoveryConfig
	Cache     *types.FileDiscoveryCache
	Circuit   *types.CircuitBreaker
}

// NewCommandBuilder creates a new command builder with default configuration
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		Config:    types.DefaultAuditQueryConfig(),
		Migration: types.DefaultMigrationConfig(),
		Discovery: types.DefaultDiscoveryConfig(),
		Cache: &types.FileDiscoveryCache{
			Cache: make(map[string][]string),
			TTL:   5 * time.Minute,
		},
		Circuit: &types.CircuitBreaker{
			FailureThreshold: 3,
			ResetTimeout:     30 * time.Second,
			State:            types.CircuitStateClosed,
		},
	}
}

// BuildOcCommand constructs the oc command based on parameters with support for rolling logs
// This is the legacy function that maintains backward compatibility
func BuildOcCommand(params types.AuditQueryParams) string {
	// Use the new command builder with default configuration
	builder := NewCommandBuilder()
	return builder.BuildOptimalCommand(params)
}

// BuildOcCommandWithConfig constructs the oc command with custom configuration
func BuildOcCommandWithConfig(params types.AuditQueryParams, config types.AuditQueryConfig) string {
	builder := NewCommandBuilder()
	builder.Config = config
	return builder.BuildOptimalCommand(params)
}

// BuildOptimalCommand builds the optimal command based on parameters and configuration
func (cb *CommandBuilder) BuildOptimalCommand(params types.AuditQueryParams) string {
	// Check circuit breaker state
	if cb.Circuit.State == types.CircuitStateOpen {
		if time.Since(cb.Circuit.LastFailureTime) < cb.Circuit.ResetTimeout {
			// Use fallback command when circuit breaker is open
			return cb.buildFallbackCommand(params)
		}
		// Reset to half-open
		cb.Circuit.State = types.CircuitStateHalfOpen
	}

	// Always start with simple approach for reliability (Phase 1 fix)
	if cb.shouldUseSimpleCommand(params) {
		return cb.buildSimpleCommand(params)
	}

	// Only use multi-file if specifically requested and safe
	if cb.shouldUseMultiFile(params) {
		return cb.buildErrorTolerantMultiFileCommand(params)
	}

	return cb.buildFallbackCommand(params)
}

// shouldUseSimpleCommand determines if we should use simple command approach
func (cb *CommandBuilder) shouldUseSimpleCommand(params types.AuditQueryParams) bool {
	// Use simple for recent timeframes or when specifically requested
	return params.Timeframe == "today" ||
		params.Timeframe == "1h" ||
		params.Timeframe == "" ||
		cb.Config.ForceSimple ||
		cb.Migration.PreserveOldBehavior
}

// shouldUseMultiFile determines if we should use multi-file approach
func (cb *CommandBuilder) shouldUseMultiFile(params types.AuditQueryParams) bool {
	// Only use multi-file if explicitly enabled and safe
	return cb.Migration.EnableNewBuilder &&
		!cb.Config.ForceSimple &&
		cb.Migration.MaxFiles > 1
}

// buildSimpleCommand builds a simple, reliable command
func (cb *CommandBuilder) buildSimpleCommand(params types.AuditQueryParams) string {
	// Check if JSON parsing is enabled and available
	if cb.Config.UseJSONParsing && cb.checkJQAvailability() {
		return cb.buildJSONAwareCommand(params)
	}

	var parts []string

	// Base command
	parts = append(parts, "oc adm node-logs --role=master")
	parts = append(parts, getDefaultLogPath(params.LogSource))

	// Add filters with complexity control
	if len(params.Patterns) > 0 {
		// Limit to first 3 patterns to avoid complexity
		maxPatterns := 3
		if len(params.Patterns) > maxPatterns {
			params.Patterns = params.Patterns[:maxPatterns]
		}
		for _, pattern := range params.Patterns {
			parts = append(parts, fmt.Sprintf("| grep -i '%s'", pattern))
		}
	}

	// Add username filter
	if params.Username != "" {
		usernameFilter := BuildUsernameFilter(params.Username)
		if usernameFilter != "" {
			parts = append(parts, usernameFilter)
		}
	}

	// Add resource filter
	if params.Resource != "" {
		resourceFilter := BuildResourceFilter(params.Resource)
		if resourceFilter != "" {
			parts = append(parts, resourceFilter)
		}
	}

	// Add verb filter
	if params.Verb != "" {
		verbFilter := BuildVerbFilter(params.Verb)
		if verbFilter != "" {
			parts = append(parts, verbFilter)
		}
	}

	// Add namespace filter
	if params.Namespace != "" {
		namespaceFilter := BuildNamespaceFilter(params.Namespace)
		if namespaceFilter != "" {
			parts = append(parts, namespaceFilter)
		}
	}

	// Add exclusions with complexity control
	if len(params.Exclude) > 0 {
		// Limit to first 3 exclusions to avoid complexity
		maxExclusions := 3
		if len(params.Exclude) > maxExclusions {
			params.Exclude = params.Exclude[:maxExclusions]
		}
		for _, exclude := range params.Exclude {
			parts = append(parts, fmt.Sprintf("| grep -v '%s'", exclude))
		}
	}

	// Add timeframe filter for simple commands
	if params.Timeframe != "" {
		timeframeFilter := buildTimeframeFilter(params.Timeframe)
		if timeframeFilter != "" {
			parts = append(parts, timeframeFilter)
		}
	}

	return strings.Join(parts, " ")
}

// buildJSONAwareCommand builds a JSON-aware command using jq for better accuracy
func (cb *CommandBuilder) buildJSONAwareCommand(params types.AuditQueryParams) string {
	baseCommand := "oc adm node-logs --role=master " + getDefaultLogPath(params.LogSource)

	// Build jq filters for JSON-aware filtering
	var jqFilters []string

	// Add username filter
	if params.Username != "" {
		escapedUsername := escapeForJQ(params.Username)
		jqFilters = append(jqFilters, fmt.Sprintf(`(.user.username // .userInfo.username // .impersonatedUser // .requestUser) | test("%s"; "i")`, escapedUsername))
	}

	// Add verb filter
	if params.Verb != "" {
		escapedVerb := escapeForJQ(params.Verb)
		jqFilters = append(jqFilters, fmt.Sprintf(`.verb | test("%s"; "i")`, escapedVerb))
	}

	// Add resource filter
	if params.Resource != "" {
		escapedResource := escapeForJQ(params.Resource)
		jqFilters = append(jqFilters, fmt.Sprintf(`(.objectRef.resource // .objectRef.apiVersion // .requestObject.kind // .responseObject.kind) | test("%s"; "i")`, escapedResource))
	}

	// Add namespace filter
	if params.Namespace != "" {
		escapedNamespace := escapeForJQ(params.Namespace)
		jqFilters = append(jqFilters, fmt.Sprintf(`(.objectRef.namespace // .requestObject.metadata.namespace // .responseObject.metadata.namespace) | test("%s"; "i")`, escapedNamespace))
	}

	// Add pattern filters
	if len(params.Patterns) > 0 {
		maxPatterns := 3
		if len(params.Patterns) > maxPatterns {
			params.Patterns = params.Patterns[:maxPatterns]
		}
		for _, pattern := range params.Patterns {
			escapedPattern := escapeForJQ(pattern)
			jqFilters = append(jqFilters, fmt.Sprintf(`tostring | test("%s"; "i")`, escapedPattern))
		}
	}

	// Add exclusion filters
	if len(params.Exclude) > 0 {
		maxExclusions := 3
		if len(params.Exclude) > maxExclusions {
			params.Exclude = params.Exclude[:maxExclusions]
		}
		for _, exclude := range params.Exclude {
			escapedExclude := escapeForJQ(exclude)
			jqFilters = append(jqFilters, fmt.Sprintf(`(tostring | test("%s"; "i") | not)`, escapedExclude))
		}
	}

	// Add timeframe filter
	if params.Timeframe != "" {
		timeframeFilter := buildJSONTimeframeFilter(params.Timeframe)
		if timeframeFilter != "" {
			jqFilters = append(jqFilters, timeframeFilter)
		}
	}

	// Build the complete jq expression
	var jqExpression string
	if len(jqFilters) > 0 {
		jqExpression = fmt.Sprintf(`select(%s)`, strings.Join(jqFilters, " and "))
	} else {
		jqExpression = "."
	}

	// Add output formatting for better readability
	jqExpression += ` | {
		timestamp: .requestReceivedTimestamp,
		username: (.user.username // .userInfo.username // "unknown"),
		verb: .verb,
		resource: (.objectRef.resource // "unknown"),
		namespace: (.objectRef.namespace // "unknown"),
		name: (.objectRef.name // "unknown"),
		statusCode: (.responseStatus.code // 0),
		statusMessage: (.responseStatus.message // ""),
		requestURI: .requestURI,
		userAgent: .userAgent,
		sourceIPs: (.sourceIPs // [])
	}`

	return fmt.Sprintf("%s | jq -r '%s'", baseCommand, jqExpression)
}

// buildJSONTimeframeFilter creates a JSON-aware timeframe filter
func buildJSONTimeframeFilter(timeframe string) string {
	now := time.Now()

	switch timeframe {
	case "today":
		today := now.Format("2006-01-02")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, today)
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, yesterday)
	case "this week":
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, now.Format("2006-01-02"))
	case "last hour":
		lastHour := now.Add(-1 * time.Hour).Format("2006-01-02")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, lastHour)
	case "24h", "last 24 hours":
		last24h := now.AddDate(0, 0, -1).Format("2006-01-02")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, last24h)
	case "7d", "last 7 days":
		last7d := now.AddDate(0, 0, -7).Format("2006-01-02")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, last7d)
	case "last week":
		lastWeek := now.AddDate(0, 0, -7).Format("2006-01-02")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, lastWeek)
	case "this month":
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, now.Format("2006-01"))
	case "last month":
		lastMonth := now.AddDate(0, -1, 0).Format("2006-01")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, lastMonth)
	case "last 30 days":
		last30d := now.AddDate(0, 0, -30).Format("2006-01-02")
		return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, last30d)
	}

	// Handle "last X minutes/hours/days" patterns
	if matched, _ := regexp.MatchString(`^last (\d+) minute(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) minute(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			minutes, _ := strconv.Atoi(matches[1])
			lastMinutes := now.Add(-time.Duration(minutes) * time.Minute).Format("2006-01-02")
			return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, lastMinutes)
		}
	}

	if matched, _ := regexp.MatchString(`^last (\d+) hour(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) hour(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			hours, _ := strconv.Atoi(matches[1])
			lastHours := now.Add(-time.Duration(hours) * time.Hour).Format("2006-01-02")
			return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, lastHours)
		}
	}

	if matched, _ := regexp.MatchString(`^last (\d+) day(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) day(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			days, _ := strconv.Atoi(matches[1])
			lastDays := now.AddDate(0, 0, -days).Format("2006-01-02")
			return fmt.Sprintf(`.requestReceivedTimestamp | test("%s")`, lastDays)
		}
	}

	return ""
}

// checkJQAvailability checks if jq is available in the system
func (cb *CommandBuilder) checkJQAvailability() bool {
	cmd := exec.Command("jq", "--version")
	err := cmd.Run()
	return err == nil
}

// escapeForJQ escapes special characters for safe jq usage
func escapeForJQ(input string) string {
	// Escape special jq characters: " \ / [ ] { } ( ) * + ? | ^ $ . ~
	escaped := input
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	escaped = strings.ReplaceAll(escaped, "/", "\\/")
	escaped = strings.ReplaceAll(escaped, "[", "\\[")
	escaped = strings.ReplaceAll(escaped, "]", "\\]")
	escaped = strings.ReplaceAll(escaped, "{", "\\{")
	escaped = strings.ReplaceAll(escaped, "}", "\\}")
	escaped = strings.ReplaceAll(escaped, "(", "\\(")
	escaped = strings.ReplaceAll(escaped, ")", "\\)")
	escaped = strings.ReplaceAll(escaped, "*", "\\*")
	escaped = strings.ReplaceAll(escaped, "+", "\\+")
	escaped = strings.ReplaceAll(escaped, "?", "\\?")
	escaped = strings.ReplaceAll(escaped, "|", "\\|")
	escaped = strings.ReplaceAll(escaped, "^", "\\^")
	escaped = strings.ReplaceAll(escaped, "$", "\\$")
	escaped = strings.ReplaceAll(escaped, ".", "\\.")
	escaped = strings.ReplaceAll(escaped, "~", "\\~")

	return escaped
}

// buildErrorTolerantMultiFileCommand builds an error-tolerant multi-file command
func (cb *CommandBuilder) buildErrorTolerantMultiFileCommand(params types.AuditQueryParams) string {
	// Get available log files
	logFiles := cb.getAvailableLogFiles(params.LogSource, params.Timeframe)

	if len(logFiles) <= 1 {
		return cb.buildSimpleCommand(params)
	}

	var commands []string

	for _, logFile := range logFiles {
		// Build individual command for this file
		fileCommand := cb.buildSingleFileCommand(params, logFile)

		// Add error tolerance: continue on failure
		errorTolerantCommand := fmt.Sprintf("(%s) || true", fileCommand)
		commands = append(commands, errorTolerantCommand)
	}

	// Use semicolon instead of && for error tolerance
	return strings.Join(commands, " ; ")
}

// buildSingleFileCommand builds a command for a single log file
func (cb *CommandBuilder) buildSingleFileCommand(params types.AuditQueryParams, logFile types.LogFileInfo) string {
	var parts []string

	// Base command
	parts = append(parts, "oc adm node-logs --role=master")

	// Handle path
	if strings.HasPrefix(logFile.Path, "--path=") {
		parts = append(parts, logFile.Path)
	} else {
		parts = append(parts, fmt.Sprintf("--path=%s", logFile.Path))
	}

	// Add filters with complexity control
	if len(params.Patterns) > 0 {
		maxPatterns := 3
		if len(params.Patterns) > maxPatterns {
			params.Patterns = params.Patterns[:maxPatterns]
		}
		for _, pattern := range params.Patterns {
			parts = append(parts, fmt.Sprintf("| grep -i '%s'", pattern))
		}
	}

	if params.Username != "" {
		usernameFilter := BuildUsernameFilter(params.Username)
		if usernameFilter != "" {
			parts = append(parts, usernameFilter)
		}
	}

	if params.Resource != "" {
		resourceFilter := BuildResourceFilter(params.Resource)
		if resourceFilter != "" {
			parts = append(parts, resourceFilter)
		}
	}

	if params.Verb != "" {
		verbFilter := BuildVerbFilter(params.Verb)
		if verbFilter != "" {
			parts = append(parts, verbFilter)
		}
	}

	if params.Namespace != "" {
		namespaceFilter := BuildNamespaceFilter(params.Namespace)
		if namespaceFilter != "" {
			parts = append(parts, namespaceFilter)
		}
	}

	// Add exclusions with complexity control
	if len(params.Exclude) > 0 {
		maxExclusions := 3
		if len(params.Exclude) > maxExclusions {
			params.Exclude = params.Exclude[:maxExclusions]
		}
		for _, exclude := range params.Exclude {
			parts = append(parts, fmt.Sprintf("| grep -v '%s'", exclude))
		}
	}

	// Add date-specific timeframe filter for rolling logs
	if !logFile.IsCurrent && !logFile.Date.IsZero() {
		dateFilter := fmt.Sprintf("| grep '%s'", logFile.Date.Format("2006-01-02"))
		parts = append(parts, dateFilter)
	}

	return strings.Join(parts, " ")
}

// buildFallbackCommand builds a fallback command when other approaches fail
func (cb *CommandBuilder) buildFallbackCommand(params types.AuditQueryParams) string {
	// Use the most reliable approach - single file with current log
	fallbackParams := params
	fallbackParams.Timeframe = "" // Force single file approach
	return cb.buildSimpleCommand(fallbackParams)
}

// getAvailableLogFiles gets available log files with caching
func (cb *CommandBuilder) getAvailableLogFiles(logSource, timeframe string) []types.LogFileInfo {
	// Check cache first
	if time.Since(cb.Cache.LastCheck) < cb.Cache.TTL {
		if cached, exists := cb.Cache.Cache[logSource]; exists {
			return cb.convertToLogFileInfo(cached, timeframe)
		}
	}

	// Discover files if enabled
	var availableFiles []string
	if cb.Discovery.EnableDiscovery {
		availableFiles = cb.discoverAvailableLogFiles(logSource)
	} else {
		// Use fallback files
		availableFiles = cb.Discovery.FallbackFiles
	}

	// Update cache
	cb.Cache.Cache[logSource] = availableFiles
	cb.Cache.LastCheck = time.Now()

	return cb.convertToLogFileInfo(availableFiles, timeframe)
}

// discoverAvailableLogFiles discovers available log files from the cluster
func (cb *CommandBuilder) discoverAvailableLogFiles(logSource string) []string {
	// Use oc adm node-logs --list-files to discover available files
	cmd := exec.Command("oc", "adm", "node-logs", "--role=master", "--list-files")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to known patterns
		return cb.getDefaultLogPatterns(logSource)
	}

	// Parse output and filter for audit files
	var availableFiles []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, logSource) && strings.Contains(line, "audit") {
			availableFiles = append(availableFiles, strings.TrimSpace(line))
		}
	}

	// Limit the number of files to check
	if len(availableFiles) > cb.Discovery.MaxFilesToCheck {
		availableFiles = availableFiles[:cb.Discovery.MaxFilesToCheck]
	}

	return availableFiles
}

// getDefaultLogPatterns returns default log file patterns
func (cb *CommandBuilder) getDefaultLogPatterns(logSource string) []string {
	basePath := getLogBasePath(logSource)
	return []string{
		fmt.Sprintf("%s.log", basePath),
		fmt.Sprintf("%s.log.1", basePath),
		fmt.Sprintf("%s.log.2", basePath),
	}
}

// convertToLogFileInfo converts file paths to LogFileInfo structs
func (cb *CommandBuilder) convertToLogFileInfo(files []string, timeframe string) []types.LogFileInfo {
	var logFiles []types.LogFileInfo

	for i, file := range files {
		logFile := types.LogFileInfo{
			Path:      file,
			IsCurrent: i == 0, // First file is current
			Exists:    true,   // Assume exists for now
		}

		// Parse date from filename if possible
		if !logFile.IsCurrent {
			logFile.Date = cb.parseDateFromFilename(file)
		}

		logFiles = append(logFiles, logFile)
	}

	return logFiles
}

// parseDateFromFilename attempts to parse a date from a log filename
func (cb *CommandBuilder) parseDateFromFilename(filename string) time.Time {
	// Try to extract date from filename patterns
	patterns := []string{
		`(\d{4}-\d{2}-\d{2})`, // YYYY-MM-DD
		`(\d{8})`,             // YYYYMMDD
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(filename)
		if len(matches) > 1 {
			dateStr := matches[1]
			if len(dateStr) == 8 {
				// YYYYMMDD format
				if date, err := time.Parse("20060102", dateStr); err == nil {
					return date
				}
			} else if len(dateStr) == 10 {
				// YYYY-MM-DD format
				if date, err := time.Parse("2006-01-02", dateStr); err == nil {
					return date
				}
			}
		}
	}

	return time.Time{} // Return zero time if no date found
}

// ExecuteCommand executes a command with circuit breaker protection
func (cb *CommandBuilder) ExecuteCommand(command string) (string, error) {
	if cb.Circuit.State == types.CircuitStateOpen {
		return "", fmt.Errorf("circuit breaker is open")
	}

	// Execute command
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		cb.recordFailure()
		return "", err
	}

	// Reset circuit breaker on success
	cb.recordSuccess()

	return string(output), nil
}

// recordFailure records a failure in the circuit breaker
func (cb *CommandBuilder) recordFailure() {
	cb.Circuit.FailureCount++
	cb.Circuit.LastFailureTime = time.Now()

	if cb.Circuit.FailureCount >= cb.Circuit.FailureThreshold {
		cb.Circuit.State = types.CircuitStateOpen
	}
}

// recordSuccess records a success in the circuit breaker
func (cb *CommandBuilder) recordSuccess() {
	cb.Circuit.FailureCount = 0
	if cb.Circuit.State == types.CircuitStateHalfOpen {
		cb.Circuit.State = types.CircuitStateClosed
	}
}

// determineLogFiles determines which log files to query based on timeframe
// Updated for Phase 1: Always use simple approach for reliability
func determineLogFiles(logSource, timeframe string) []LogFileInfo {
	// Phase 1 fix: Always use simple approach for reliability
	// This prevents the complex multi-file commands that were causing failures
	return []LogFileInfo{
		{Path: getDefaultLogPath(logSource), IsCurrent: true},
	}
}

// buildMultiFileCommand builds a command that queries multiple log files
func buildMultiFileCommand(params types.AuditQueryParams, logFiles []LogFileInfo) string {
	var commands []string

	for _, logFile := range logFiles {
		// Create a copy of params for this log file
		fileParams := params

		// Build command for this specific log file
		parts := []string{"oc adm node-logs --role=master"}

		// Handle path - if it already has --path= prefix, use it directly
		if strings.HasPrefix(logFile.Path, "--path=") {
			parts = append(parts, logFile.Path)
		} else {
			parts = append(parts, fmt.Sprintf("--path=%s", logFile.Path))
		}

		// Add filters with complexity control
		if len(fileParams.Patterns) > 0 {
			// Limit to first 3 patterns to avoid complexity
			maxPatterns := 3
			if len(fileParams.Patterns) > maxPatterns {
				fileParams.Patterns = fileParams.Patterns[:maxPatterns]
			}
			for _, pattern := range fileParams.Patterns {
				parts = append(parts, fmt.Sprintf("| grep -i '%s'", pattern))
			}
		}

		if fileParams.Username != "" {
			usernameFilter := BuildUsernameFilter(fileParams.Username)
			if usernameFilter != "" {
				parts = append(parts, usernameFilter)
			}
		}

		if fileParams.Resource != "" {
			resourceFilter := BuildResourceFilter(fileParams.Resource)
			if resourceFilter != "" {
				parts = append(parts, resourceFilter)
			}
		}

		if fileParams.Verb != "" {
			verbFilter := BuildVerbFilter(fileParams.Verb)
			if verbFilter != "" {
				parts = append(parts, verbFilter)
			}
		}

		if fileParams.Namespace != "" {
			namespaceFilter := BuildNamespaceFilter(fileParams.Namespace)
			if namespaceFilter != "" {
				parts = append(parts, namespaceFilter)
			}
		}

		// Add exclusions with complexity control
		if len(fileParams.Exclude) > 0 {
			// Limit to first 3 exclusions to avoid complexity
			maxExclusions := 3
			if len(fileParams.Exclude) > maxExclusions {
				fileParams.Exclude = fileParams.Exclude[:maxExclusions]
			}
			for _, exclude := range fileParams.Exclude {
				parts = append(parts, fmt.Sprintf("| grep -v '%s'", exclude))
			}
		}

		// Add date-specific timeframe filter for rolling logs
		if !logFile.IsCurrent && !logFile.Date.IsZero() {
			dateFilter := fmt.Sprintf("| grep '%s'", logFile.Date.Format("2006-01-02"))
			parts = append(parts, dateFilter)
		}

		commands = append(commands, strings.Join(parts, " "))
	}

	// Combine all commands
	if len(commands) == 1 {
		return commands[0]
	}

	// Use the efficient multi-file command builder
	return buildEfficientMultiFileCommand(commands)
}

// buildEfficientMultiFileCommand builds an efficient command for multiple files
func buildEfficientMultiFileCommand(commands []string) string {
	if len(commands) == 0 {
		return ""
	}

	if len(commands) == 1 {
		return commands[0]
	}

	// Use && to chain commands - this ensures all commands must succeed
	// and is safer than ; which allows commands to continue even if previous ones fail
	return fmt.Sprintf("(%s)", strings.Join(commands, " && "))
}

// generateRollingLogPaths generates possible rolling log file paths for a given date
func generateRollingLogPaths(logSource string, date time.Time) []LogFileInfo {
	var paths []LogFileInfo

	// Simplified rolling log patterns - only the most common ones
	patterns := []string{
		"%s.log.%s",
		"%s-%s.log",
	}

	// Add compressed file patterns
	compressedPatterns := []string{
		"%s.log.%s.gz",
		"%s-%s.log.gz",
		"%s.log.%s.bz2",
		"%s-%s.log.bz2",
	}

	basePath := getLogBasePath(logSource)
	dateStr := date.Format("2006-01-02")

	// Generate numbered rolling files (1-3 only to prevent complexity)
	for i := 1; i <= 3; i++ {
		for _, pattern := range patterns {
			path := fmt.Sprintf(pattern, basePath, fmt.Sprintf("%d", i))
			paths = append(paths, LogFileInfo{
				Path: path,
				Date: date,
			})
		}
		// Add compressed numbered files
		for _, pattern := range compressedPatterns {
			path := fmt.Sprintf(pattern, basePath, fmt.Sprintf("%d", i))
			paths = append(paths, LogFileInfo{
				Path: path,
				Date: date,
			})
		}
	}

	// Generate date-based rolling files (only 2 patterns)
	for _, pattern := range patterns {
		path := fmt.Sprintf(pattern, basePath, dateStr)
		paths = append(paths, LogFileInfo{
			Path: path,
			Date: date,
		})
	}
	// Add compressed date-based files
	for _, pattern := range compressedPatterns {
		path := fmt.Sprintf(pattern, basePath, dateStr)
		paths = append(paths, LogFileInfo{
			Path: path,
			Date: date,
		})
	}

	return paths
}

// getLogBasePath returns the base path for a log source
func getLogBasePath(logSource string) string {
	switch logSource {
	case "kube-apiserver":
		return "kube-apiserver/audit"
	case "oauth-server":
		return "oauth-server/audit"
	case "openshift-apiserver":
		return "openshift-apiserver/audit"
	case "oauth-apiserver":
		return "oauth-apiserver/audit"
	case "node":
		return "audit/audit"
	default:
		return "kube-apiserver/audit"
	}
}

// getDefaultLogPath returns the default log path for a log source
func getDefaultLogPath(logSource string) string {
	switch logSource {
	case "kube-apiserver":
		return "--path=kube-apiserver/audit.log"
	case "oauth-server":
		return "--path=oauth-server/audit.log"
	case "openshift-apiserver":
		return "--path=openshift-apiserver/audit.log"
	case "oauth-apiserver":
		return "--path=oauth-apiserver/audit.log"
	case "node":
		return "--path=audit/audit.log"
	default:
		return "--path=kube-apiserver/audit.log"
	}
}

// parseTimeframe parses a timeframe string and returns start and end dates
func parseTimeframe(timeframe string) (time.Time, time.Time) {
	now := time.Now()

	// Handle special cases first
	switch timeframe {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return start, now
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		start := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
		end := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, yesterday.Location())
		return start, end
	case "this week":
		// Start of current week (Monday)
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		start := now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		return start, now
	case "last week":
		// Previous week
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		lastWeekStart := now.AddDate(0, 0, -weekday-6)
		lastWeekStart = time.Date(lastWeekStart.Year(), lastWeekStart.Month(), lastWeekStart.Day(), 0, 0, 0, 0, lastWeekStart.Location())
		lastWeekEnd := lastWeekStart.AddDate(0, 0, 6)
		lastWeekEnd = time.Date(lastWeekEnd.Year(), lastWeekEnd.Month(), lastWeekEnd.Day(), 23, 59, 59, 999999999, lastWeekEnd.Location())
		return lastWeekStart, lastWeekEnd
	case "this month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, now
	case "last month":
		lastMonth := now.AddDate(0, -1, 0)
		start := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
		end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-time.Nanosecond)
		if end.Before(start) {
			// Fix for edge case where end is before start
			end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		}
		return start, end
	}

	// Parse "last X minutes"
	if matched, _ := regexp.MatchString(`^last (\d+) minute(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) minute(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			minutes, _ := strconv.Atoi(matches[1])
			start := now.Add(-time.Duration(minutes) * time.Minute)
			return start, now
		}
	}

	// Parse "last X hours"
	if matched, _ := regexp.MatchString(`^last (\d+) hour(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) hour(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			hours, _ := strconv.Atoi(matches[1])
			start := now.Add(-time.Duration(hours) * time.Hour)
			return start, now
		}
	}

	// Parse "last X days"
	if matched, _ := regexp.MatchString(`^last (\d+) day(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) day(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			days, _ := strconv.Atoi(matches[1])
			start := now.AddDate(0, 0, -days)
			return start, now
		}
	}

	// Parse "last X weeks"
	if matched, _ := regexp.MatchString(`^last (\d+) week(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) week(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			weeks, _ := strconv.Atoi(matches[1])
			start := now.AddDate(0, 0, -7*weeks)
			return start, now
		}
	}

	// Parse "last X months"
	if matched, _ := regexp.MatchString(`^last (\d+) month(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) month(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			months, _ := strconv.Atoi(matches[1])
			start := now.AddDate(0, -months, 0)
			return start, now
		}
	}

	// Parse "last X years"
	if matched, _ := regexp.MatchString(`^last (\d+) year(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) year(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			years, _ := strconv.Atoi(matches[1])
			start := now.AddDate(-years, 0, 0)
			return start, now
		}
	}

	// Parse short forms like "5m", "2h", "3d", "1w", "6y"
	if matched, _ := regexp.MatchString(`^(\d+)([mhdwy])$`, timeframe); matched {
		re := regexp.MustCompile(`^(\d+)([mhdwy])$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 2 {
			value, _ := strconv.Atoi(matches[1])
			unit := matches[2]
			var start time.Time
			switch unit {
			case "m":
				start = now.Add(-time.Duration(value) * time.Minute)
			case "h":
				start = now.Add(-time.Duration(value) * time.Hour)
			case "d":
				start = now.AddDate(0, 0, -value)
			case "w":
				start = now.AddDate(0, 0, -7*value)
			case "y":
				start = now.AddDate(-value, 0, 0)
			}
			return start, now
		}
	}

	// Parse "X ago" forms like "5m ago", "2h ago", "3d ago"
	if matched, _ := regexp.MatchString(`^(\d+)([mhdwy]) ago$`, timeframe); matched {
		re := regexp.MustCompile(`^(\d+)([mhdwy]) ago$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 2 {
			value, _ := strconv.Atoi(matches[1])
			unit := matches[2]
			var start time.Time
			switch unit {
			case "m":
				start = now.Add(-time.Duration(value) * time.Minute)
			case "h":
				start = now.Add(-time.Duration(value) * time.Hour)
			case "d":
				start = now.AddDate(0, 0, -value)
			case "w":
				start = now.AddDate(0, 0, -7*value)
			case "y":
				start = now.AddDate(-value, 0, 0)
			}
			return start, now
		}
	}

	// Parse "since" dates
	if matched, _ := regexp.MatchString(`^since (\d{4}-\d{2}-\d{2})$`, timeframe); matched {
		re := regexp.MustCompile(`^since (\d{4}-\d{2}-\d{2})$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			date, err := time.Parse("2006-01-02", matches[1])
			if err == nil {
				return date, now
			}
		}
	}

	// Parse "since" dates with time
	if matched, _ := regexp.MatchString(`^since (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})$`, timeframe); matched {
		re := regexp.MustCompile(`^since (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			datetime, err := time.Parse("2006-01-02 15:04:05", matches[1])
			if err == nil {
				return datetime, now
			}
		}
	}

	// Return zero times for invalid timeframe
	return time.Time{}, time.Time{}
}

// buildTimeframeFilter creates a timeframe filter for the command (legacy support)
func buildTimeframeFilter(timeframe string) string {
	// Handle all timeframes using flexible parsing
	return buildFlexibleTimeframeFilter(timeframe)
}

// buildFlexibleTimeframeFilter handles dynamic timeframe patterns (legacy support)
func buildFlexibleTimeframeFilter(timeframe string) string {
	now := time.Now()

	// Handle special cases first
	switch timeframe {
	case "today":
		return fmt.Sprintf("| grep '%s'", now.Format("2006-01-02"))
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		return fmt.Sprintf("| grep '%s'", yesterday.Format("2006-01-02"))
	case "this week":
		return fmt.Sprintf("| grep '%s'", now.Format("2006-01-02"))
	case "last hour":
		lastHour := now.Add(-1 * time.Hour)
		return fmt.Sprintf("| grep '%s'", lastHour.Format("2006-01-02"))
	case "24h", "last 24 hours":
		last24h := now.AddDate(0, 0, -1)
		return fmt.Sprintf("| grep '%s'", last24h.Format("2006-01-02"))
	case "7d", "last 7 days":
		last7d := now.AddDate(0, 0, -7)
		return fmt.Sprintf("| grep '%s'", last7d.Format("2006-01-02"))
	case "last week":
		lastWeek := now.AddDate(0, 0, -7)
		return fmt.Sprintf("| grep '%s'", lastWeek.Format("2006-01-02"))
	case "this month":
		return fmt.Sprintf("| grep '%s'", now.Format("2006-01"))
	case "last month":
		lastMonth := now.AddDate(0, -1, 0)
		return fmt.Sprintf("| grep '%s'", lastMonth.Format("2006-01"))
	case "last 30 days":
		last30d := now.AddDate(0, 0, -30)
		return fmt.Sprintf("| grep '%s'", last30d.Format("2006-01-02"))
	}

	// Parse "last X minutes"
	if matched, _ := regexp.MatchString(`^last (\d+) minute(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) minute(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			minutes, _ := strconv.Atoi(matches[1])
			lastMinutes := now.Add(-time.Duration(minutes) * time.Minute)
			return fmt.Sprintf("| grep '%s'", lastMinutes.Format("2006-01-02"))
		}
	}

	// Parse "last X hours"
	if matched, _ := regexp.MatchString(`^last (\d+) hour(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) hour(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			hours, _ := strconv.Atoi(matches[1])
			lastHours := now.Add(-time.Duration(hours) * time.Hour)
			return fmt.Sprintf("| grep '%s'", lastHours.Format("2006-01-02"))
		}
	}

	// Parse "last X days"
	if matched, _ := regexp.MatchString(`^last (\d+) day(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) day(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			days, _ := strconv.Atoi(matches[1])
			lastDays := now.AddDate(0, 0, -days)
			return fmt.Sprintf("| grep '%s'", lastDays.Format("2006-01-02"))
		}
	}

	// Parse "last X weeks"
	if matched, _ := regexp.MatchString(`^last (\d+) week(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) week(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			weeks, _ := strconv.Atoi(matches[1])
			lastWeeks := now.AddDate(0, 0, -7*weeks)
			return fmt.Sprintf("| grep '%s'", lastWeeks.Format("2006-01-02"))
		}
	}

	// Parse "last X months"
	if matched, _ := regexp.MatchString(`^last (\d+) month(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) month(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			months, _ := strconv.Atoi(matches[1])
			lastMonths := now.AddDate(0, -months, 0)
			return fmt.Sprintf("| grep '%s'", lastMonths.Format("2006-01"))
		}
	}

	// Parse "last X years"
	if matched, _ := regexp.MatchString(`^last (\d+) year(s)?$`, timeframe); matched {
		re := regexp.MustCompile(`^last (\d+) year(s)?$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			years, _ := strconv.Atoi(matches[1])
			lastYears := now.AddDate(-years, 0, 0)
			return fmt.Sprintf("| grep '%s'", lastYears.Format("2006-01"))
		}
	}

	// Parse short forms like "5m", "2h", "3d", "1w", "6y"
	if matched, _ := regexp.MatchString(`^(\d+)([mhdwy])$`, timeframe); matched {
		re := regexp.MustCompile(`^(\d+)([mhdwy])$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 2 {
			value, _ := strconv.Atoi(matches[1])
			unit := matches[2]
			switch unit {
			case "m":
				lastMinutes := now.Add(-time.Duration(value) * time.Minute)
				return fmt.Sprintf("| grep '%s'", lastMinutes.Format("2006-01-02"))
			case "h":
				lastHours := now.Add(-time.Duration(value) * time.Hour)
				return fmt.Sprintf("| grep '%s'", lastHours.Format("2006-01-02"))
			case "d":
				lastDays := now.AddDate(0, 0, -value)
				return fmt.Sprintf("| grep '%s'", lastDays.Format("2006-01-02"))
			case "w":
				lastWeeks := now.AddDate(0, 0, -7*value)
				return fmt.Sprintf("| grep '%s'", lastWeeks.Format("2006-01-02"))
			case "y":
				lastYears := now.AddDate(-value, 0, 0)
				return fmt.Sprintf("| grep '%s'", lastYears.Format("2006-01"))
			}
		}
	}

	// Parse "X ago" forms like "5m ago", "2h ago", "3d ago"
	if matched, _ := regexp.MatchString(`^(\d+)([mhdwy]) ago$`, timeframe); matched {
		re := regexp.MustCompile(`^(\d+)([mhdwy]) ago$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 2 {
			value, _ := strconv.Atoi(matches[1])
			unit := matches[2]
			switch unit {
			case "m":
				lastMinutes := now.Add(-time.Duration(value) * time.Minute)
				return fmt.Sprintf("| grep '%s'", lastMinutes.Format("2006-01-02"))
			case "h":
				lastHours := now.Add(-time.Duration(value) * time.Hour)
				return fmt.Sprintf("| grep '%s'", lastHours.Format("2006-01-02"))
			case "d":
				lastDays := now.AddDate(0, 0, -value)
				return fmt.Sprintf("| grep '%s'", lastDays.Format("2006-01-02"))
			case "w":
				lastWeeks := now.AddDate(0, 0, -7*value)
				return fmt.Sprintf("| grep '%s'", lastWeeks.Format("2006-01-02"))
			case "y":
				lastYears := now.AddDate(-value, 0, 0)
				return fmt.Sprintf("| grep '%s'", lastYears.Format("2006-01"))
			}
		}
	}

	// Parse "since" dates
	if matched, _ := regexp.MatchString(`^since (\d{4}-\d{2}-\d{2})$`, timeframe); matched {
		re := regexp.MustCompile(`^since (\d{4}-\d{2}-\d{2})$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			date := matches[1]
			return fmt.Sprintf("| grep '%s'", date)
		}
	}

	// Parse "since" dates with time
	if matched, _ := regexp.MatchString(`^since (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})$`, timeframe); matched {
		re := regexp.MustCompile(`^since (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})$`)
		matches := re.FindStringSubmatch(timeframe)
		if len(matches) > 1 {
			datetime := matches[1]
			return fmt.Sprintf("| grep '%s'", datetime)
		}
	}

	return ""
}
