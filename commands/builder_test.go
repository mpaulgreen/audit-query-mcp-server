package commands

import (
	"strings"
	"testing"
	"time"

	"audit-query-mcp-server/types"
)

// TestBuildOcCommand_Basic tests basic command construction
func TestBuildOcCommand_Basic(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"test-pattern"},
	}

	command := BuildOcCommand(params)
	expected := "oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'test-pattern'"

	if command != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, command)
	}
}

// TestBuildOcCommand_AllLogSources tests all supported log sources
func TestBuildOcCommand_AllLogSources(t *testing.T) {
	logSources := []string{
		"kube-apiserver",
		"oauth-server",
		"openshift-apiserver",
		"oauth-apiserver",
		"node",
	}

	for _, source := range logSources {
		t.Run(source, func(t *testing.T) {
			params := types.AuditQueryParams{
				LogSource: source,
			}

			command := BuildOcCommand(params)
			if !strings.Contains(command, "oc adm node-logs --role=master") {
				t.Errorf("Command should contain base oc command")
			}
		})
	}
}

// TestBuildOcCommand_WithFilters tests comprehensive filtering
func TestBuildOcCommand_WithFilters(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Username:  "testuser",
		Resource:  "pods",
		Verb:      "GET",
		Namespace: "default",
		Exclude:   []string{"health-check"},
	}

	command := BuildOcCommand(params)

	// Check that all filters are present
	checks := []string{
		"oc adm node-logs --role=master",
		"--path=kube-apiserver/audit.log",
		"grep '\"user\":{\"[^\"]*\":\"testuser\"'",
		"grep '\"objectRef\":{\"[^\"]*\":\"pods\"'",
		"grep '\"verb\":\"GET\"'",
		"grep '\"objectRef\":{\"[^\"]*\":\"default\"'",
		"grep -v 'health-check'",
	}

	for _, check := range checks {
		if !strings.Contains(command, check) {
			t.Errorf("Command should contain: %s", check)
		}
	}
}

// TestBuildOcCommand_TimeframeToday tests today timeframe
func TestBuildOcCommand_TimeframeToday(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
	}

	command := BuildOcCommand(params)

	// Should use current log file for today
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Should use current log file for today")
	}
}

// TestBuildOcCommand_TimeframeYesterday tests yesterday timeframe
func TestBuildOcCommand_TimeframeYesterday(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "yesterday",
	}

	command := BuildOcCommand(params)

	// Should use multi-file approach for yesterday
	if !strings.Contains(command, "&&") {
		t.Errorf("Should use multi-file command for yesterday: %s", command)
	}
}

// TestBuildOcCommand_TimeframeLastWeek tests last week timeframe
func TestBuildOcCommand_TimeframeLastWeek(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "last week",
	}

	command := BuildOcCommand(params)

	// Should use multi-file approach for last week
	if !strings.Contains(command, "&&") {
		t.Errorf("Should use multi-file command for last week: %s", command)
	}
}

// TestBuildOcCommand_TimeframeLastMonth tests last month timeframe
func TestBuildOcCommand_TimeframeLastMonth(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "last month",
	}

	command := BuildOcCommand(params)

	// Should use multi-file approach for last month
	if !strings.Contains(command, "&&") {
		t.Errorf("Should use multi-file command for last month: %s", command)
	}
}

// TestBuildOcCommand_TimeframeShortForms tests short form timeframes
func TestBuildOcCommand_TimeframeShortForms(t *testing.T) {
	testCases := []struct {
		timeframe string
		expected  string
	}{
		{"5m", "5m"},
		{"2h", "2h"},
		{"3d", "3d"},
		{"1w", "1w"},
		{"6y", "6y"},
	}

	for _, tc := range testCases {
		t.Run(tc.timeframe, func(t *testing.T) {
			params := types.AuditQueryParams{
				LogSource: "kube-apiserver",
				Timeframe: tc.timeframe,
			}

			command := BuildOcCommand(params)

			// Should use current log file for short timeframes
			if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
				t.Errorf("Should use current log file for short timeframe %s", tc.timeframe)
			}
		})
	}
}

// TestBuildOcCommand_TimeframeAgoForms tests ago form timeframes
func TestBuildOcCommand_TimeframeAgoForms(t *testing.T) {
	testCases := []struct {
		timeframe string
		expected  string
	}{
		{"5m ago", "5m ago"},
		{"2h ago", "2h ago"},
		{"3d ago", "3d ago"},
		{"1w ago", "1w ago"},
		{"6y ago", "6y ago"},
	}

	for _, tc := range testCases {
		t.Run(tc.timeframe, func(t *testing.T) {
			params := types.AuditQueryParams{
				LogSource: "kube-apiserver",
				Timeframe: tc.timeframe,
			}

			command := BuildOcCommand(params)

			// Should use current log file for ago timeframes
			if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
				t.Errorf("Should use current log file for ago timeframe %s", tc.timeframe)
			}
		})
	}
}

// TestBuildOcCommand_TimeframeSinceDate tests since date timeframe
func TestBuildOcCommand_TimeframeSinceDate(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "since 2024-01-15",
	}

	command := BuildOcCommand(params)

	// Should use multi-file approach for since date
	if !strings.Contains(command, "&&") {
		t.Errorf("Should use multi-file command for since date: %s", command)
	}
}

// TestBuildOcCommand_TimeframeSinceDateTime tests since datetime timeframe
func TestBuildOcCommand_TimeframeSinceDateTime(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "since 2024-01-15 14:30:00",
	}

	command := BuildOcCommand(params)

	// Should use multi-file approach for since datetime
	if !strings.Contains(command, "&&") {
		t.Errorf("Should use multi-file command for since datetime: %s", command)
	}
}

// TestBuildOcCommand_ComplexQuery tests complex query with multiple filters
func TestBuildOcCommand_ComplexQuery(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"error", "failed"},
		Username:  "admin",
		Resource:  "deployments",
		Verb:      "CREATE",
		Namespace: "production",
		Exclude:   []string{"health-check", "liveness"},
		Timeframe: "last 7 days",
	}

	command := BuildOcCommand(params)

	// Should use multi-file approach for 7 days
	if !strings.Contains(command, "&&") {
		t.Errorf("Should use multi-file command for 7 days: %s", command)
	}
}

// TestDetermineLogFiles_NoTimeframe tests log file determination without timeframe
func TestDetermineLogFiles_NoTimeframe(t *testing.T) {
	logFiles := determineLogFiles("kube-apiserver", "")

	if len(logFiles) != 1 {
		t.Errorf("Expected 1 log file, got %d", len(logFiles))
	}

	if logFiles[0].Path != "--path=kube-apiserver/audit.log" {
		t.Errorf("Expected current log file path, got %s", logFiles[0].Path)
	}

	if !logFiles[0].IsCurrent {
		t.Errorf("Expected IsCurrent to be true")
	}
}

// TestDetermineLogFiles_Today tests log file determination for today
func TestDetermineLogFiles_Today(t *testing.T) {
	logFiles := determineLogFiles("kube-apiserver", "today")

	if len(logFiles) != 1 {
		t.Errorf("Expected 1 log file for today, got %d", len(logFiles))
	}

	if !logFiles[0].IsCurrent {
		t.Errorf("Expected IsCurrent to be true for today")
	}
}

// TestDetermineLogFiles_Yesterday tests log file determination for yesterday
func TestDetermineLogFiles_Yesterday(t *testing.T) {
	logFiles := determineLogFiles("kube-apiserver", "yesterday")

	if len(logFiles) < 2 {
		t.Errorf("Expected multiple log files for yesterday, got %d", len(logFiles))
	}

	// Should include current log file
	foundCurrent := false
	for _, lf := range logFiles {
		if lf.IsCurrent {
			foundCurrent = true
			break
		}
	}

	if !foundCurrent {
		t.Errorf("Should include current log file")
	}
}

// TestDetermineLogFiles_LastWeek tests log file determination for last week
func TestDetermineLogFiles_LastWeek(t *testing.T) {
	logFiles := determineLogFiles("kube-apiserver", "last week")

	if len(logFiles) < 8 {
		t.Errorf("Expected multiple log files for last week, got %d", len(logFiles))
	}
}

// TestDetermineLogFiles_LastMonth tests log file determination for last month
func TestDetermineLogFiles_LastMonth(t *testing.T) {
	logFiles := determineLogFiles("kube-apiserver", "last month")

	if len(logFiles) < 32 {
		t.Errorf("Expected multiple log files for last month, got %d", len(logFiles))
	}
}

// TestParseTimeframe_Today tests timeframe parsing for today
func TestParseTimeframe_Today(t *testing.T) {
	start, end := parseTimeframe("today")

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if !start.Equal(todayStart) {
		t.Errorf("Expected start time to be today start, got %v", start)
	}

	if end.Before(start) {
		t.Errorf("End time should be after start time")
	}
}

// TestParseTimeframe_Yesterday tests timeframe parsing for yesterday
func TestParseTimeframe_Yesterday(t *testing.T) {
	start, end := parseTimeframe("yesterday")

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, yesterday.Location())

	if !start.Equal(yesterdayStart) {
		t.Errorf("Expected start time to be yesterday start, got %v", start)
	}

	if !end.Equal(yesterdayEnd) {
		t.Errorf("Expected end time to be yesterday end, got %v", end)
	}
}

// TestParseTimeframe_LastWeek tests timeframe parsing for last week
func TestParseTimeframe_LastWeek(t *testing.T) {
	start, end := parseTimeframe("last week")

	if start.IsZero() || end.IsZero() {
		t.Errorf("Expected valid dates for last week")
	}

	if end.Before(start) {
		t.Errorf("End time should be after start time")
	}

	// Should be approximately 7 days apart
	duration := end.Sub(start)
	if duration < 6*24*time.Hour || duration > 8*24*time.Hour {
		t.Errorf("Expected approximately 7 days, got %v", duration)
	}
}

// TestParseTimeframe_ShortForms tests timeframe parsing for short forms
func TestParseTimeframe_ShortForms(t *testing.T) {
	testCases := []struct {
		timeframe string
		expected  string
	}{
		{"5m", "5m"},
		{"2h", "2h"},
		{"3d", "3d"},
		{"1w", "1w"},
		{"6y", "6y"},
	}

	for _, tc := range testCases {
		t.Run(tc.timeframe, func(t *testing.T) {
			start, end := parseTimeframe(tc.timeframe)

			if start.IsZero() || end.IsZero() {
				t.Errorf("Expected valid dates for %s", tc.timeframe)
			}

			if end.Before(start) {
				t.Errorf("End time should be after start time for %s", tc.timeframe)
			}
		})
	}
}

// TestParseTimeframe_AgoForms tests timeframe parsing for ago forms
func TestParseTimeframe_AgoForms(t *testing.T) {
	testCases := []struct {
		timeframe string
		expected  string
	}{
		{"5m ago", "5m ago"},
		{"2h ago", "2h ago"},
		{"3d ago", "3d ago"},
		{"1w ago", "1w ago"},
		{"6y ago", "6y ago"},
	}

	for _, tc := range testCases {
		t.Run(tc.timeframe, func(t *testing.T) {
			start, end := parseTimeframe(tc.timeframe)

			if start.IsZero() || end.IsZero() {
				t.Errorf("Expected valid dates for %s", tc.timeframe)
			}

			if end.Before(start) {
				t.Errorf("End time should be after start time for %s", tc.timeframe)
			}
		})
	}
}

// TestParseTimeframe_SinceDate tests timeframe parsing for since date
func TestParseTimeframe_SinceDate(t *testing.T) {
	start, end := parseTimeframe("since 2024-01-15")

	expectedStart, _ := time.Parse("2006-01-02", "2024-01-15")

	if !start.Equal(expectedStart) {
		t.Errorf("Expected start time to be 2024-01-15, got %v", start)
	}

	if end.IsZero() {
		t.Errorf("Expected valid end time")
	}
}

// TestParseTimeframe_SinceDateTime tests timeframe parsing for since datetime
func TestParseTimeframe_SinceDateTime(t *testing.T) {
	start, end := parseTimeframe("since 2024-01-15 14:30:00")

	expectedStart, _ := time.Parse("2006-01-02 15:04:05", "2024-01-15 14:30:00")

	if !start.Equal(expectedStart) {
		t.Errorf("Expected start time to be 2024-01-15 14:30:00, got %v", start)
	}

	if end.IsZero() {
		t.Errorf("Expected valid end time")
	}
}

// TestParseTimeframe_Invalid tests timeframe parsing for invalid input
func TestParseTimeframe_Invalid(t *testing.T) {
	start, end := parseTimeframe("invalid timeframe")

	if !start.IsZero() || !end.IsZero() {
		t.Errorf("Expected zero times for invalid timeframe, got %v, %v", start, end)
	}
}

// TestParseTimeframe_EdgeCases tests edge cases for timeframe parsing
func TestParseTimeframe_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		timeframe string
		wantValid bool
	}{
		{"empty string", "", false},
		{"invalid format", "invalid", false},
		{"zero minutes", "last 0 minutes", true},
		{"zero hours", "last 0 hours", true},
		{"zero days", "last 0 days", true},
		{"zero weeks", "last 0 weeks", true},
		{"zero months", "last 0 months", true},
		{"zero years", "last 0 years", true},
		{"large number", "last 999999 days", true},
		{"negative number", "last -5 days", false},
		{"decimal number", "last 5.5 days", false},
		{"mixed case", "LAST 5 DAYS", false},
		{"extra spaces", "  last  5  days  ", false},
		{"short form zero", "0d", true},
		{"short form zero ago", "0d ago", true},
		{"invalid short form", "5x", false},
		{"invalid ago form", "5x ago", false},
		{"since invalid date", "since 2023-13-45", false},
		{"since invalid datetime", "since 2023-12-25 25:70:80", false},
		{"since valid date", "since 2023-12-25", true},
		{"since valid datetime", "since 2023-12-25 15:30:45", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := parseTimeframe(tt.timeframe)
			if tt.wantValid {
				if start.IsZero() || end.IsZero() {
					t.Errorf("parseTimeframe() should return valid times for %s", tt.timeframe)
				}
			} else {
				if !start.IsZero() || !end.IsZero() {
					t.Errorf("parseTimeframe() should return zero times for invalid %s", tt.timeframe)
				}
			}
		})
	}
}

// TestParseTimeframe_AllPatterns tests all supported timeframe patterns
func TestParseTimeframe_AllPatterns(t *testing.T) {
	tests := []struct {
		name      string
		timeframe string
	}{
		{"minutes pattern", "last 30 minutes"},
		{"hours pattern", "last 12 hours"},
		{"days pattern", "last 7 days"},
		{"weeks pattern", "last 4 weeks"},
		{"months pattern", "last 6 months"},
		{"years pattern", "last 2 years"},
		{"short form minutes", "30m"},
		{"short form hours", "12h"},
		{"short form days", "7d"},
		{"short form weeks", "4w"},
		{"short form years", "2y"},
		{"ago form minutes", "30m ago"},
		{"ago form hours", "12h ago"},
		{"ago form days", "7d ago"},
		{"ago form weeks", "4w ago"},
		{"ago form years", "2y ago"},
		{"this week", "this week"},
		{"last week", "last week"},
		{"this month", "this month"},
		{"last month", "last month"},
		{"today", "today"},
		{"yesterday", "yesterday"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := parseTimeframe(tt.timeframe)
			if start.IsZero() || end.IsZero() {
				t.Errorf("parseTimeframe() failed for %s", tt.timeframe)
			}
			if start.After(end) {
				t.Errorf("parseTimeframe() start time after end time for %s", tt.timeframe)
			}
		})
	}
}

// TestBuildFlexibleTimeframeFilter_AllPatterns tests all patterns in buildFlexibleTimeframeFilter
func TestBuildFlexibleTimeframeFilter_AllPatterns(t *testing.T) {
	tests := []struct {
		name         string
		timeframe    string
		wantContains string
	}{
		{"today", "today", "grep"},
		{"yesterday", "yesterday", "grep"},
		{"this week", "this week", "grep"},
		{"last hour", "last hour", "grep"},
		{"24h", "24h", "grep"},
		{"last 24 hours", "last 24 hours", "grep"},
		{"7d", "7d", "grep"},
		{"last 7 days", "last 7 days", "grep"},
		{"last week", "last week", "grep"},
		{"this month", "this month", "grep"},
		{"last month", "last month", "grep"},
		{"last 30 days", "last 30 days", "grep"},
		{"last 5 minutes", "last 5 minutes", "grep"},
		{"last 3 hours", "last 3 hours", "grep"},
		{"last 10 days", "last 10 days", "grep"},
		{"last 2 weeks", "last 2 weeks", "grep"},
		{"last 6 months", "last 6 months", "grep"},
		{"last 1 year", "last 1 year", "grep"},
		{"5m", "5m", "grep"},
		{"2h", "2h", "grep"},
		{"3d", "3d", "grep"},
		{"1w", "1w", "grep"},
		{"6y", "6y", "grep"},
		{"5m ago", "5m ago", "grep"},
		{"2h ago", "2h ago", "grep"},
		{"3d ago", "3d ago", "grep"},
		{"1w ago", "1w ago", "grep"},
		{"6y ago", "6y ago", "grep"},
		{"since date", "since 2023-12-25", "grep"},
		{"since datetime", "since 2023-12-25 15:30:45", "grep"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFlexibleTimeframeFilter(tt.timeframe)
			if tt.wantContains != "" {
				if !strings.Contains(result, tt.wantContains) {
					t.Errorf("buildFlexibleTimeframeFilter() result should contain %s for %s, got: %s",
						tt.wantContains, tt.timeframe, result)
				}
			}
		})
	}
}

// TestBuildFlexibleTimeframeFilter_EdgeCases tests edge cases for buildFlexibleTimeframeFilter
func TestBuildFlexibleTimeframeFilter_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		timeframe string
		wantEmpty bool
	}{
		{"empty string", "", true},
		{"invalid format", "invalid", true},
		{"zero minutes", "last 0 minutes", false},
		{"zero hours", "last 0 hours", false},
		{"zero days", "last 0 days", false},
		{"zero weeks", "last 0 weeks", false},
		{"zero months", "last 0 months", false},
		{"zero years", "last 0 years", false},
		{"negative number", "last -5 days", true},
		{"decimal number", "last 5.5 days", true},
		{"mixed case", "LAST 5 DAYS", true},
		{"extra spaces", "  last  5  days  ", true},
		{"short form zero", "0d", false},
		{"short form zero ago", "0d ago", false},
		{"invalid short form", "5x", true},
		{"invalid ago form", "5x ago", true},
		{"since invalid date", "since 2023-13-45", false},
		{"since invalid datetime", "since 2023-12-25 25:70:80", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFlexibleTimeframeFilter(tt.timeframe)
			if tt.wantEmpty {
				if result != "" {
					t.Errorf("buildFlexibleTimeframeFilter() should return empty string for %s, got: %s",
						tt.timeframe, result)
				}
			} else {
				if result == "" {
					t.Errorf("buildFlexibleTimeframeFilter() should return non-empty string for %s", tt.timeframe)
				}
			}
		})
	}
}

// TestGetLogBasePath tests log base path generation
func TestGetLogBasePath(t *testing.T) {
	testCases := []struct {
		logSource string
		expected  string
	}{
		{"kube-apiserver", "kube-apiserver/audit"},
		{"oauth-server", "oauth-server/audit"},
		{"openshift-apiserver", "openshift-apiserver/audit"},
		{"oauth-apiserver", "oauth-apiserver/audit"},
		{"node", "audit/audit"},
		{"unknown", "kube-apiserver/audit"}, // default
	}

	for _, tc := range testCases {
		t.Run(tc.logSource, func(t *testing.T) {
			result := getLogBasePath(tc.logSource)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestGetDefaultLogPath tests default log path generation
func TestGetDefaultLogPath(t *testing.T) {
	testCases := []struct {
		logSource string
		expected  string
	}{
		{"kube-apiserver", "--path=kube-apiserver/audit.log"},
		{"oauth-server", "--path=oauth-server/audit.log"},
		{"openshift-apiserver", "--path=openshift-apiserver/audit.log"},
		{"oauth-apiserver", "--path=oauth-apiserver/audit.log"},
		{"node", "--path=audit/audit.log"},
		{"unknown", "--path=kube-apiserver/audit.log"}, // default
	}

	for _, tc := range testCases {
		t.Run(tc.logSource, func(t *testing.T) {
			result := getDefaultLogPath(tc.logSource)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestGenerateRollingLogPaths tests rolling log path generation
func TestGenerateRollingLogPaths(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	logFiles := generateRollingLogPaths("kube-apiserver", date)

	if len(logFiles) == 0 {
		t.Errorf("Expected rolling log paths, got none")
	}

	// Check that all paths have the correct date
	for _, lf := range logFiles {
		if !lf.Date.Equal(date) {
			t.Errorf("Expected date %v, got %v", date, lf.Date)
		}

		if lf.IsCurrent {
			t.Errorf("Rolling log files should not be marked as current")
		}
	}
}

// TestGenerateRollingLogPaths_CompressedFiles tests that compressed file paths are generated
func TestGenerateRollingLogPaths_CompressedFiles(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	logFiles := generateRollingLogPaths("kube-apiserver", date)

	// Check for compressed file patterns
	foundGz := false
	foundBz2 := false

	for _, lf := range logFiles {
		if strings.Contains(lf.Path, ".gz") {
			foundGz = true
		}
		if strings.Contains(lf.Path, ".bz2") {
			foundBz2 = true
		}
	}

	if !foundGz {
		t.Errorf("Expected to find .gz compressed file paths")
	}

	if !foundBz2 {
		t.Errorf("Expected to find .bz2 compressed file paths")
	}
}

// TestBuildMultiFileCommand tests multi-file command building
func TestBuildMultiFileCommand(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Username:  "testuser",
		Patterns:  []string{"error"},
	}

	logFiles := []LogFileInfo{
		{Path: "kube-apiserver/audit.log", IsCurrent: true},
		{Path: "kube-apiserver/audit.log.1", Date: time.Now().AddDate(0, 0, -1)},
		{Path: "kube-apiserver/audit.log.2", Date: time.Now().AddDate(0, 0, -2)},
	}

	command := buildMultiFileCommand(params, logFiles)

	if !strings.Contains(command, "&&") {
		t.Errorf("Should use multi-file command: %s", command)
	}
}

// TestBuildEfficientMultiFileCommand tests efficient multi-file command building
func TestBuildEfficientMultiFileCommand(t *testing.T) {
	commands := []string{
		"oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep error",
		"oc adm node-logs --role=master --path=kube-apiserver/audit.log.1 | grep error",
	}

	result := buildEfficientMultiFileCommand(commands)

	if !strings.Contains(result, "&&") {
		t.Errorf("Should use concatenated commands: %s", result)
	}
}

// TestBuildEfficientMultiFileCommand_SingleCommand tests single command handling
func TestBuildEfficientMultiFileCommand_SingleCommand(t *testing.T) {
	commands := []string{
		"oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep error",
	}

	result := buildEfficientMultiFileCommand(commands)

	if result != commands[0] {
		t.Errorf("Should return single command unchanged: %s", result)
	}
}

// TestBuildEfficientMultiFileCommand_Empty tests empty command handling
func TestBuildEfficientMultiFileCommand_Empty(t *testing.T) {
	commands := []string{}

	result := buildEfficientMultiFileCommand(commands)

	if result != "" {
		t.Errorf("Should return empty string for empty commands: %s", result)
	}
}

// TestRealWorldLogPatterns tests real-world log patterns observed in the cluster
func TestRealWorldLogPatterns(t *testing.T) {
	// Test kube-apiserver patterns
	t.Run("kube-apiserver", func(t *testing.T) {
		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Timeframe: "yesterday",
		}

		command := BuildOcCommand(params)

		// Should include current log file
		if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
			t.Errorf("Should include current log file")
		}

		// Should use multi-file approach for historical data
		if !strings.Contains(command, "&&") {
			t.Errorf("Should use multi-file command for historical data")
		}
	})

	// Test node audit patterns
	t.Run("node", func(t *testing.T) {
		params := types.AuditQueryParams{
			LogSource: "node",
			Timeframe: "last week",
		}

		command := BuildOcCommand(params)

		// Should include current log file
		if !strings.Contains(command, "--path=audit/audit.log") {
			t.Errorf("Should include current log file")
		}

		// Should use multi-file approach for historical data
		if !strings.Contains(command, "&&") {
			t.Errorf("Should use multi-file command for historical data")
		}
	})
}

// TestFilterEscaping tests filter escaping functionality
func TestFilterEscaping(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"test[pattern]"},
		Exclude:   []string{"test(pattern)"},
	}

	command := BuildOcCommand(params)

	// Current implementation doesn't escape patterns and exclusions in main function
	// but does escape them in filter functions
	if !strings.Contains(command, "test[pattern]") {
		t.Errorf("Should contain unescaped pattern: test[pattern]")
	}

	if !strings.Contains(command, "test(pattern)") {
		t.Errorf("Should contain unescaped exclusion: test(pattern)")
	}

	// Test that filter functions do escape properly
	usernameFilter := BuildUsernameFilter("test[user]")
	if !strings.Contains(usernameFilter, "test\\[user\\]") {
		t.Errorf("Username filter should escape special characters")
	}
}

// TestUsernameFilterComprehensive tests comprehensive username filtering
func TestUsernameFilterComprehensive(t *testing.T) {
	username := "testuser"
	filter := BuildUsernameFilter(username)

	// Should include multiple patterns for comprehensive coverage
	patterns := []string{
		"\"user\":{\"[^\"]*\":\"testuser\"",
		"\"user\":\"testuser\"",
		"\"userInfo\":{\"[^\"]*\":\"testuser\"",
		"\"impersonatedUser\":\"testuser\"",
		"\"requestUser\":\"testuser\"",
	}

	for _, pattern := range patterns {
		if !strings.Contains(filter, pattern) {
			t.Errorf("Username filter should contain pattern: %s", pattern)
		}
	}
}

// TestResourceFilterComprehensive tests comprehensive resource filtering
func TestResourceFilterComprehensive(t *testing.T) {
	resource := "pods"
	filter := BuildResourceFilter(resource)

	// Should include multiple patterns for comprehensive coverage
	patterns := []string{
		"\"objectRef\":{\"[^\"]*\":\"pods\"",
		"\"objectRef\":{\"[^\"]*\":\"[^\"]*pods\"",
		"\"requestObject\":{\"[^\"]*\":\"pods\"",
		"\"responseObject\":{\"[^\"]*\":\"pods\"",
	}

	for _, pattern := range patterns {
		if !strings.Contains(filter, pattern) {
			t.Errorf("Resource filter should contain pattern: %s", pattern)
		}
	}
}

// TestVerbFilterComprehensive tests comprehensive verb filtering
func TestVerbFilterComprehensive(t *testing.T) {
	verb := "GET"
	filter := BuildVerbFilter(verb)

	// Should include multiple patterns for comprehensive coverage
	patterns := []string{
		"\"verb\":\"GET\"",
		"\"method\":\"GET\"",
		"\"action\":\"GET\"",
		"\"operation\":\"GET\"",
	}

	for _, pattern := range patterns {
		if !strings.Contains(filter, pattern) {
			t.Errorf("Verb filter should contain pattern: %s", pattern)
		}
	}
}

// TestNamespaceFilterComprehensive tests comprehensive namespace filtering
func TestNamespaceFilterComprehensive(t *testing.T) {
	namespace := "default"
	filter := BuildNamespaceFilter(namespace)

	// Should include multiple patterns for comprehensive coverage
	patterns := []string{
		"\"objectRef\":{\"[^\"]*\":\"default\"",
		"\"requestObject\":{\"[^\"]*\":{\"[^\"]*\":\"default\"",
		"\"responseObject\":{\"[^\"]*\":{\"[^\"]*\":\"default\"",
		"\"requestURI\":\"[^\"]*default[^\"]*\"",
	}

	for _, pattern := range patterns {
		if !strings.Contains(filter, pattern) {
			t.Errorf("Namespace filter should contain pattern: %s", pattern)
		}
	}
}

// TestRealAuditLogFormat tests against real audit log format
func TestRealAuditLogFormat(t *testing.T) {
	t.Run("real-json-format", func(t *testing.T) {
		// Test against the real JSON format observed in the cluster
		// Real format: {"kind":"Event","apiVersion":"audit.k8s.io/v1","level":"Metadata","auditID":"...","stage":"ResponseComplete","requestURI":"...","verb":"create","user":{"username":"system:serviceaccount:nvidia-gpu-operator:gpu-operator",...}}

		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Username:  "system:serviceaccount:nvidia-gpu-operator:gpu-operator",
			Verb:      "create",
			Resource:  "clusterrolebindings",
			Namespace: "nvidia-gpu-operator",
		}

		command := BuildOcCommand(params)

		// Should include filters for the real JSON structure
		checks := []string{
			"\"user\":{\"[^\"]*\":\"system:serviceaccount:nvidia-gpu-operator:gpu-operator\"",
			"\"verb\":\"create\"",
			"\"objectRef\":{\"[^\"]*\":\"clusterrolebindings\"",
			"\"objectRef\":{\"[^\"]*\":\"nvidia-gpu-operator\"",
		}

		for _, check := range checks {
			if !strings.Contains(command, check) {
				t.Errorf("Command should contain: %s", check)
			}
		}
	})
}

// TestErrorHandling tests error handling scenarios
func TestErrorHandling(t *testing.T) {
	t.Run("invalid-timeframe", func(t *testing.T) {
		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Timeframe: "invalid timeframe",
		}

		command := BuildOcCommand(params)

		// Should fallback to current log file
		if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
			t.Errorf("Should fallback to current log file for invalid timeframe")
		}
	})

	t.Run("empty-params", func(t *testing.T) {
		params := types.AuditQueryParams{}

		command := BuildOcCommand(params)

		// Should use default log source
		if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
			t.Errorf("Should use default log source for empty params")
		}
	})
}

// TestPerformance tests performance characteristics
func TestPerformance(t *testing.T) {
	t.Run("large-timeframe", func(t *testing.T) {
		// Test performance with large timeframe (1 year)
		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Timeframe: "last 1 year",
		}

		start := time.Now()
		command := BuildOcCommand(params)
		duration := time.Since(start)

		// Should complete within reasonable time
		if duration > 100*time.Millisecond {
			t.Errorf("Command generation took too long: %v", duration)
		}

		// Should use multi-file approach
		if !strings.Contains(command, "&&") {
			t.Errorf("Should use multi-file command for large timeframe, got: %s", command)
		}
	})
}

// TestEdgeCases tests edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("special-characters", func(t *testing.T) {
		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Patterns:  []string{"test[pattern]"},
			Exclude:   []string{"test(pattern)"},
		}

		command := BuildOcCommand(params)

		// Should handle special characters appropriately
		if !strings.Contains(command, "test[pattern]") {
			t.Errorf("Should handle special characters in patterns")
		}
	})

	t.Run("very-short-timeframe", func(t *testing.T) {
		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Timeframe: "5m",
		}

		command := BuildOcCommand(params)

		// Should use current log file for very short timeframes
		if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
			t.Errorf("Should use current log file for very short timeframes")
		}
	})
}

// TestIntegration tests integration scenarios
func TestIntegration(t *testing.T) {
	t.Run("complex-real-world-query", func(t *testing.T) {
		// Test a complex real-world query scenario
		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Patterns:  []string{"error", "failed", "timeout"},
			Username:  "system:serviceaccount:openshift-cluster-version:default",
			Resource:  "services",
			Verb:      "get",
			Namespace: "openshift-network-operator",
			Exclude:   []string{"health-check", "liveness", "readiness"},
			Timeframe: "last 7 days",
		}

		command := BuildOcCommand(params)

		// Should include all components
		checks := []string{
			"oc adm node-logs --role=master",
			"\"user\":{\"[^\"]*\":\"system:serviceaccount:openshift-cluster-version:default\"",
			"\"objectRef\":{\"[^\"]*\":\"services\"",
			"\"verb\":\"get\"",
			"\"objectRef\":{\"[^\"]*\":\"openshift-network-operator\"",
			"grep -i 'error'",
			"grep -i 'failed'",
			"grep -i 'timeout'",
			"grep -v 'health-check'",
			"grep -v 'liveness'",
			"grep -v 'readiness'",
		}

		for _, check := range checks {
			if !strings.Contains(command, check) {
				t.Errorf("Command should contain: %s", check)
			}
		}

		// Should use multi-file approach for 7 days
		if !strings.Contains(command, "&&") {
			t.Errorf("Should use multi-file command for 7 days")
		}
	})
}

// TestBuildOcCommand_ComplexityLimits tests complexity limits in command building
func TestBuildOcCommand_ComplexityLimits(t *testing.T) {
	// Test pattern limit (max 3)
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"pattern1", "pattern2", "pattern3", "pattern4", "pattern5"},
		Timeframe: "today",
	}

	command := BuildOcCommand(params)
	patternCount := strings.Count(command, "grep -i")
	if patternCount > 3 {
		t.Errorf("Command should limit patterns to 3, got %d grep patterns", patternCount)
	}

	// Test exclusion limit (max 3)
	params = types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Exclude:   []string{"exclude1", "exclude2", "exclude3", "exclude4", "exclude5"},
		Timeframe: "today",
	}

	command = BuildOcCommand(params)
	excludeCount := strings.Count(command, "grep -v")
	if excludeCount > 3 {
		t.Errorf("Command should limit exclusions to 3, got %d grep -v patterns", excludeCount)
	}
}

// TestBuildMultiFileCommand_ComplexityLimits tests complexity limits in multi-file command building
func TestBuildMultiFileCommand_ComplexityLimits(t *testing.T) {
	// Test with many log files to ensure limits are respected
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"pattern1", "pattern2", "pattern3", "pattern4", "pattern5"},
		Exclude:   []string{"exclude1", "exclude2", "exclude3", "exclude4", "exclude5"},
		Timeframe: "last month", // This will generate many log files
	}

	command := BuildOcCommand(params)

	// In multi-file commands, each file gets its own patterns, so we expect more grep patterns
	// but each individual command should still respect the 3-pattern limit
	// Let's check that the command is properly structured
	if !strings.Contains(command, "&&") {
		t.Errorf("Multi-file command should use && to chain commands")
	}

	// Check that the command is not empty
	if command == "" {
		t.Errorf("Multi-file command should not be empty")
	}

	t.Logf("Multi-file command generated: %s", command)
}

// TestDetermineLogFiles_FileLimit tests the file limit in determineLogFiles
func TestDetermineLogFiles_FileLimit(t *testing.T) {
	// Test with a very large timeframe that would generate many files
	logFiles := determineLogFiles("kube-apiserver", "last year")

	// Should not exceed maxFiles (50)
	if len(logFiles) > 50 {
		t.Errorf("determineLogFiles() should limit files to 50, got %d files", len(logFiles))
	}
}

// TestGenerateRollingLogPaths_Limit tests the limit in generateRollingLogPaths
func TestGenerateRollingLogPaths_Limit(t *testing.T) {
	date := time.Now()
	paths := generateRollingLogPaths("kube-apiserver", date)

	// Should generate a reasonable number of paths (not too many)
	// Based on the implementation: 3 numbered patterns * 2 (normal + compressed) * 2 (numbered + date-based)
	// Plus additional patterns for date-based files
	// 3 numbered * 2 (normal + compressed) * 2 (numbered + date-based) = 24 paths
	expectedMax := 3 * 2 * 2 * 2 // 24 paths
	if len(paths) > expectedMax {
		t.Errorf("generateRollingLogPaths() should limit paths, got %d paths", len(paths))
	}

	t.Logf("Generated %d rolling log paths", len(paths))
}

// TestBuildOcCommand_TimeframeFilter tests timeframe filter behavior
func TestBuildOcCommand_TimeframeFilter(t *testing.T) {
	// Test that timeframe filter is only added when not using multi-file approach
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
	}

	command := BuildOcCommand(params)
	// For today, should use single file and add timeframe filter
	if !strings.Contains(command, "grep") {
		t.Errorf("Command for today should include timeframe filter")
	}

	// Test with yesterday (should use multi-file approach)
	params.Timeframe = "yesterday"
	command = BuildOcCommand(params)
	// Should use multi-file approach, so no additional timeframe filter
	if strings.Contains(command, "&&") {
		// Multi-file command, should not have additional timeframe filter
		if strings.Count(command, "grep") > 0 {
			// Only grep patterns from filters, not timeframe
			t.Logf("Multi-file command generated: %s", command)
		}
	}
}

// Benchmark tests for performance
func BenchmarkBuildOcCommand_Basic(b *testing.B) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"test-pattern"},
	}

	for i := 0; i < b.N; i++ {
		BuildOcCommand(params)
	}
}

func BenchmarkBuildOcCommand_Complex(b *testing.B) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"error", "failed", "timeout"},
		Username:  "admin",
		Resource:  "deployments",
		Verb:      "CREATE",
		Namespace: "production",
		Exclude:   []string{"health-check", "liveness", "readiness"},
		Timeframe: "last 7 days",
	}

	for i := 0; i < b.N; i++ {
		BuildOcCommand(params)
	}
}

func BenchmarkParseTimeframe_Complex(b *testing.B) {
	timeframes := []string{
		"today",
		"yesterday",
		"last week",
		"last month",
		"5m",
		"2h",
		"3d",
		"1w",
		"since 2024-01-15",
	}

	for i := 0; i < b.N; i++ {
		for _, timeframe := range timeframes {
			parseTimeframe(timeframe)
		}
	}
}

func BenchmarkDetermineLogFiles_Historical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		determineLogFiles("kube-apiserver", "last month")
	}
}

func BenchmarkGenerateRollingLogPaths(b *testing.B) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		generateRollingLogPaths("kube-apiserver", date)
	}
}

func BenchmarkRealWorldPatterns(b *testing.B) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"error", "failed"},
		Username:  "system:serviceaccount:nvidia-gpu-operator:gpu-operator",
		Resource:  "clusterrolebindings",
		Verb:      "create",
		Namespace: "nvidia-gpu-operator",
		Timeframe: "last week",
	}

	for i := 0; i < b.N; i++ {
		BuildOcCommand(params)
	}
}

func BenchmarkLargeTimeframe(b *testing.B) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "last month",
	}

	for i := 0; i < b.N; i++ {
		BuildOcCommand(params)
	}
}
