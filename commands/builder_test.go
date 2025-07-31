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

	// Phase 2: JSON-aware parsing is enabled by default
	// Should use jq for pattern matching instead of grep
	if !strings.Contains(command, "jq -r") {
		t.Errorf("Expected JSON-aware command with jq, got: %s", command)
	}

	if !strings.Contains(command, "test-pattern") {
		t.Errorf("Expected command to contain pattern 'test-pattern', got: %s", command)
	}

	if !strings.Contains(command, "oc adm node-logs --role=master --path=kube-apiserver/audit.log") {
		t.Errorf("Expected base oc command, got: %s", command)
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

	// Phase 2: JSON-aware parsing is enabled by default
	// Should use jq for filtering instead of grep
	checks := []string{
		"oc adm node-logs --role=master",
		"--path=kube-apiserver/audit.log",
		"jq -r",
		"testuser",
		"pods",
		"GET",
		"default",
		"health-check",
	}

	for _, check := range checks {
		if !strings.Contains(command, check) {
			t.Errorf("Command should contain: %s", check)
		}
	}

	// Should not contain old grep patterns
	grepChecks := []string{
		"grep '\"user\":{\"[^\"]*\":\"testuser\"'",
		"grep '\"objectRef\":{\"[^\"]*\":\"pods\"'",
		"grep '\"verb\":\"GET\"'",
		"grep '\"objectRef\":{\"[^\"]*\":\"default\"'",
		"grep -v 'health-check'",
	}

	for _, check := range grepChecks {
		if strings.Contains(command, check) {
			t.Errorf("Command should not contain old grep pattern: %s", check)
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

	// Should use simple approach by default for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability: %s", command)
	}

	// Should use current log file with date filtering
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Should use current log file: %s", command)
	}
}

// TestBuildOcCommand_TimeframeLastWeek tests last week timeframe
func TestBuildOcCommand_TimeframeLastWeek(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "last week",
	}

	command := BuildOcCommand(params)

	// Should use simple approach by default for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability: %s", command)
	}

	// Should use current log file with date filtering
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Should use current log file: %s", command)
	}
}

// TestBuildOcCommand_TimeframeLastMonth tests last month timeframe
func TestBuildOcCommand_TimeframeLastMonth(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "last month",
	}

	command := BuildOcCommand(params)

	// Should use simple approach by default for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability: %s", command)
	}

	// Should use current log file with date filtering
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Should use current log file: %s", command)
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

	// Should use simple approach by default for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability: %s", command)
	}

	// Should use current log file with date filtering
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Should use current log file: %s", command)
	}
}

// TestBuildOcCommand_TimeframeSinceDateTime tests since datetime timeframe
func TestBuildOcCommand_TimeframeSinceDateTime(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "since 2024-01-15 14:30:00",
	}

	command := BuildOcCommand(params)

	// Should use simple approach by default for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability: %s", command)
	}

	// Should use current log file with date filtering
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Should use current log file: %s", command)
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

	// Should use simple approach by default for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability: %s", command)
	}

	// Should use current log file with date filtering
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Should use current log file: %s", command)
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

	// Should use simple approach for reliability (Phase 1 fix)
	if len(logFiles) != 1 {
		t.Errorf("Expected 1 log file for yesterday (simple approach), got %d", len(logFiles))
	}

	// Should include current log file
	if !logFiles[0].IsCurrent {
		t.Errorf("Should include current log file")
	}
}

// TestDetermineLogFiles_LastWeek tests log file determination for last week
func TestDetermineLogFiles_LastWeek(t *testing.T) {
	logFiles := determineLogFiles("kube-apiserver", "last week")

	// Should use simple approach for reliability (Phase 1 fix)
	if len(logFiles) != 1 {
		t.Errorf("Expected 1 log file for last week (simple approach), got %d", len(logFiles))
	}
}

// TestDetermineLogFiles_LastMonth tests log file determination for last month
func TestDetermineLogFiles_LastMonth(t *testing.T) {
	logFiles := determineLogFiles("kube-apiserver", "last month")

	// Should use simple approach for reliability (Phase 1 fix)
	if len(logFiles) != 1 {
		t.Errorf("Expected 1 log file for last month (simple approach), got %d", len(logFiles))
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

		// Should use simple approach for reliability (Phase 1 fix)
		if strings.Contains(command, "&&") {
			t.Errorf("Should use simple command for reliability")
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

		// Should use simple approach for reliability (Phase 1 fix)
		if strings.Contains(command, "&&") {
			t.Errorf("Should use simple command for reliability")
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

	// Phase 2: JSON-aware parsing is enabled by default
	// Should use jq for pattern matching with proper escaping
	if !strings.Contains(command, "jq -r") {
		t.Errorf("Expected JSON-aware command with jq, got: %s", command)
	}

	// Should contain the patterns in the jq expression (may be escaped for jq)
	if !strings.Contains(command, "test") {
		t.Errorf("Should contain pattern base: test")
	}

	if !strings.Contains(command, "pattern") {
		t.Errorf("Should contain pattern part: pattern")
	}

	// Test that filter functions do escape properly for legacy compatibility
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

		// Phase 2: JSON-aware parsing is enabled by default
		// Should use jq for filtering instead of grep patterns
		checks := []string{
			"jq -r",
			"system:serviceaccount:nvidia-gpu-operator:gpu-operator",
			"create",
			"clusterrolebindings",
			"nvidia-gpu-operator",
		}

		for _, check := range checks {
			if !strings.Contains(command, check) {
				t.Errorf("Command should contain: %s", check)
			}
		}

		// Should not contain old grep patterns
		grepChecks := []string{
			"\"user\":{\"[^\"]*\":\"system:serviceaccount:nvidia-gpu-operator:gpu-operator\"",
			"\"verb\":\"create\"",
			"\"objectRef\":{\"[^\"]*\":\"clusterrolebindings\"",
			"\"objectRef\":{\"[^\"]*\":\"nvidia-gpu-operator\"",
		}

		for _, check := range grepChecks {
			if strings.Contains(command, check) {
				t.Errorf("Command should not contain old grep pattern: %s", check)
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

		// Should use simple approach for reliability (Phase 1 fix)
		if strings.Contains(command, "&&") {
			t.Errorf("Should use simple command for reliability, got: %s", command)
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

		// Phase 2: JSON-aware parsing is enabled by default
		// Should use jq for filtering instead of grep patterns
		checks := []string{
			"oc adm node-logs --role=master",
			"jq -r",
			"system:serviceaccount:openshift-cluster-version:default",
			"services",
			"get",
			"openshift-network-operator",
			"error",
			"failed",
			"timeout",
			"health-check",
			"liveness",
			"readiness",
		}

		for _, check := range checks {
			if !strings.Contains(command, check) {
				t.Errorf("Command should contain: %s", check)
			}
		}

		// Should not contain old grep patterns
		grepChecks := []string{
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

		for _, check := range grepChecks {
			if strings.Contains(command, check) {
				t.Errorf("Command should not contain old grep pattern: %s", check)
			}
		}

		// Should use simple approach for reliability (Phase 1 fix)
		if strings.Contains(command, "&&") {
			t.Errorf("Should use simple command for reliability")
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
	// Phase 2: JSON-aware parsing is enabled by default
	// Should use jq instead of grep patterns
	if !strings.Contains(command, "jq -r") {
		t.Errorf("Expected JSON-aware command with jq, got: %s", command)
	}

	// Should contain the patterns in the jq expression
	patternCount := strings.Count(command, "pattern")
	if patternCount < 3 {
		t.Errorf("Command should contain at least 3 patterns, got %d pattern references", patternCount)
	}

	// Test exclusion limit (max 3)
	params = types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Exclude:   []string{"exclude1", "exclude2", "exclude3", "exclude4", "exclude5"},
		Timeframe: "today",
	}

	command = BuildOcCommand(params)
	// Should use jq instead of grep -v
	if !strings.Contains(command, "jq -r") {
		t.Errorf("Expected JSON-aware command with jq, got: %s", command)
	}

	// Should contain the exclusions in the jq expression
	excludeCount := strings.Count(command, "exclude")
	if excludeCount < 3 {
		t.Errorf("Command should contain at least 3 exclusions, got %d exclusion references", excludeCount)
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

	// Should use simple approach for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability")
	}

	// Check that the command is properly structured
	if !strings.Contains(command, "oc adm node-logs --role=master") {
		t.Errorf("Command should contain base oc command")
	}

	// Check that the command is not empty
	if command == "" {
		t.Errorf("Command should not be empty")
	}

	t.Logf("Simple command generated: %s", command)
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
	// Phase 2: JSON-aware parsing is enabled by default
	// For today, should use single file and add timeframe filter via jq
	if !strings.Contains(command, "jq -r") {
		t.Errorf("Command for today should include JSON-aware timeframe filter")
	}

	// Test with yesterday (should use simple approach for reliability)
	params.Timeframe = "yesterday"
	command = BuildOcCommand(params)
	// Should use simple approach for reliability (Phase 1 fix)
	if strings.Contains(command, "&&") {
		t.Errorf("Should use simple command for reliability")
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

func TestBuildJSONAwareCommand_Basic(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Username:  "admin",
		Verb:      "create",
		Resource:  "pods",
		Namespace: "default",
	}

	command := builder.buildJSONAwareCommand(params)

	// Should contain jq command
	if !strings.Contains(command, "jq -r") {
		t.Errorf("Expected jq command, got: %s", command)
	}

	// Should contain base oc command
	if !strings.Contains(command, "oc adm node-logs --role=master") {
		t.Errorf("Expected oc command, got: %s", command)
	}

	// Should contain log path
	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Expected log path, got: %s", command)
	}

	// Should contain username filter
	if !strings.Contains(command, "username") {
		t.Errorf("Expected username filter, got: %s", command)
	}

	// Should contain verb filter
	if !strings.Contains(command, "verb") {
		t.Errorf("Expected verb filter, got: %s", command)
	}

	// Should contain resource filter
	if !strings.Contains(command, "resource") {
		t.Errorf("Expected resource filter, got: %s", command)
	}

	// Should contain namespace filter
	if !strings.Contains(command, "namespace") {
		t.Errorf("Expected namespace filter, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_WithPatterns(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"customresourcedefinition", "delete"},
		Username:  "admin",
	}

	command := builder.buildJSONAwareCommand(params)

	// Should contain pattern filters
	if !strings.Contains(command, "customresourcedefinition") {
		t.Errorf("Expected pattern filter, got: %s", command)
	}

	if !strings.Contains(command, "delete") {
		t.Errorf("Expected pattern filter, got: %s", command)
	}

	// Should contain tostring filter for patterns
	if !strings.Contains(command, "tostring") {
		t.Errorf("Expected tostring filter for patterns, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_WithExclusions(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Exclude:   []string{"system:", "kube-system"},
		Username:  "admin",
	}

	command := builder.buildJSONAwareCommand(params)

	// Should contain exclusion filters
	if !strings.Contains(command, "system:") {
		t.Errorf("Expected exclusion filter, got: %s", command)
	}

	if !strings.Contains(command, "kube-system") {
		t.Errorf("Expected exclusion filter, got: %s", command)
	}

	// Should contain not operator for exclusions
	if !strings.Contains(command, "not") {
		t.Errorf("Expected not operator for exclusions, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_WithTimeframe(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
		Username:  "admin",
	}

	command := builder.buildJSONAwareCommand(params)

	// Should contain timeframe filter
	if !strings.Contains(command, "requestReceivedTimestamp") {
		t.Errorf("Expected timeframe filter, got: %s", command)
	}

	// Should contain test function
	if !strings.Contains(command, "test") {
		t.Errorf("Expected test function, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_ComplexFilters(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Username:  "admin",
		Verb:      "delete",
		Resource:  "customresourcedefinitions",
		Namespace: "default",
		Patterns:  []string{"customer"},
		Exclude:   []string{"system:"},
		Timeframe: "yesterday",
	}

	command := builder.buildJSONAwareCommand(params)

	// Should contain all filters
	if !strings.Contains(command, "admin") {
		t.Errorf("Expected username filter, got: %s", command)
	}

	if !strings.Contains(command, "delete") {
		t.Errorf("Expected verb filter, got: %s", command)
	}

	if !strings.Contains(command, "customresourcedefinitions") {
		t.Errorf("Expected resource filter, got: %s", command)
	}

	if !strings.Contains(command, "default") {
		t.Errorf("Expected namespace filter, got: %s", command)
	}

	if !strings.Contains(command, "customer") {
		t.Errorf("Expected pattern filter, got: %s", command)
	}

	if !strings.Contains(command, "system:") {
		t.Errorf("Expected exclusion filter, got: %s", command)
	}

	// Should contain select with multiple conditions
	if !strings.Contains(command, "select(") {
		t.Errorf("Expected select function, got: %s", command)
	}

	// Should contain and operator for multiple conditions
	if !strings.Contains(command, " and ") {
		t.Errorf("Expected and operator, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_OutputFormatting(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Username:  "admin",
	}

	command := builder.buildJSONAwareCommand(params)

	// Should contain output formatting
	if !strings.Contains(command, "timestamp:") {
		t.Errorf("Expected timestamp field, got: %s", command)
	}

	if !strings.Contains(command, "username:") {
		t.Errorf("Expected username field, got: %s", command)
	}

	if !strings.Contains(command, "verb:") {
		t.Errorf("Expected verb field, got: %s", command)
	}

	if !strings.Contains(command, "resource:") {
		t.Errorf("Expected resource field, got: %s", command)
	}

	if !strings.Contains(command, "namespace:") {
		t.Errorf("Expected namespace field, got: %s", command)
	}

	if !strings.Contains(command, "name:") {
		t.Errorf("Expected name field, got: %s", command)
	}

	if !strings.Contains(command, "statusCode:") {
		t.Errorf("Expected statusCode field, got: %s", command)
	}

	if !strings.Contains(command, "statusMessage:") {
		t.Errorf("Expected statusMessage field, got: %s", command)
	}

	if !strings.Contains(command, "requestURI:") {
		t.Errorf("Expected requestURI field, got: %s", command)
	}

	if !strings.Contains(command, "userAgent:") {
		t.Errorf("Expected userAgent field, got: %s", command)
	}

	if !strings.Contains(command, "sourceIPs:") {
		t.Errorf("Expected sourceIPs field, got: %s", command)
	}
}

func TestBuildJSONTimeframeFilter_Today(t *testing.T) {
	filter := buildJSONTimeframeFilter("today")

	if filter == "" {
		t.Error("Expected non-empty filter for today")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_Yesterday(t *testing.T) {
	filter := buildJSONTimeframeFilter("yesterday")

	if filter == "" {
		t.Error("Expected non-empty filter for yesterday")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_LastHour(t *testing.T) {
	filter := buildJSONTimeframeFilter("last hour")

	if filter == "" {
		t.Error("Expected non-empty filter for last hour")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_Last24Hours(t *testing.T) {
	filter := buildJSONTimeframeFilter("24h")

	if filter == "" {
		t.Error("Expected non-empty filter for 24h")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_LastWeek(t *testing.T) {
	filter := buildJSONTimeframeFilter("last week")

	if filter == "" {
		t.Error("Expected non-empty filter for last week")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_LastMonth(t *testing.T) {
	filter := buildJSONTimeframeFilter("last month")

	if filter == "" {
		t.Error("Expected non-empty filter for last month")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_Last30Days(t *testing.T) {
	filter := buildJSONTimeframeFilter("last 30 days")

	if filter == "" {
		t.Error("Expected non-empty filter for last 30 days")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_LastMinutes(t *testing.T) {
	filter := buildJSONTimeframeFilter("last 5 minutes")

	if filter == "" {
		t.Error("Expected non-empty filter for last 5 minutes")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_LastHours(t *testing.T) {
	filter := buildJSONTimeframeFilter("last 2 hours")

	if filter == "" {
		t.Error("Expected non-empty filter for last 2 hours")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_LastDays(t *testing.T) {
	filter := buildJSONTimeframeFilter("last 7 days")

	if filter == "" {
		t.Error("Expected non-empty filter for last 7 days")
	}

	if !strings.Contains(filter, "requestReceivedTimestamp") {
		t.Errorf("Expected requestReceivedTimestamp field, got: %s", filter)
	}

	if !strings.Contains(filter, "test") {
		t.Errorf("Expected test function, got: %s", filter)
	}
}

func TestBuildJSONTimeframeFilter_InvalidTimeframe(t *testing.T) {
	filter := buildJSONTimeframeFilter("invalid timeframe")

	if filter != "" {
		t.Errorf("Expected empty filter for invalid timeframe, got: %s", filter)
	}
}

func TestEscapeForJQ_Basic(t *testing.T) {
	input := "simple text"
	escaped := escapeForJQ(input)

	if escaped != input {
		t.Errorf("Expected unchanged text, got: %s", escaped)
	}
}

func TestEscapeForJQ_SpecialCharacters(t *testing.T) {
	input := `test"with'quotes\and\slashes[and]braces{and}parentheses(and)operators*+?|^$.~`
	escaped := escapeForJQ(input)

	// Should escape special characters
	if !strings.Contains(escaped, `\"`) {
		t.Errorf("Expected escaped quotes, got: %s", escaped)
	}

	// The function escapes backslashes: \ becomes \\
	if !strings.Contains(escaped, `\\`) {
		t.Errorf("Expected escaped backslashes, got: %s", escaped)
	}

	// Test with input that actually contains forward slashes
	inputWithSlashes := `test/path/with/slashes`
	escapedWithSlashes := escapeForJQ(inputWithSlashes)
	if !strings.Contains(escapedWithSlashes, `\/`) {
		t.Errorf("Expected escaped forward slashes, got: %s", escapedWithSlashes)
	}

	if !strings.Contains(escaped, `\[`) {
		t.Errorf("Expected escaped brackets, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\]`) {
		t.Errorf("Expected escaped brackets, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\{`) {
		t.Errorf("Expected escaped braces, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\}`) {
		t.Errorf("Expected escaped braces, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\(`) {
		t.Errorf("Expected escaped parentheses, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\)`) {
		t.Errorf("Expected escaped parentheses, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\*`) {
		t.Errorf("Expected escaped asterisk, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\+`) {
		t.Errorf("Expected escaped plus, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\?`) {
		t.Errorf("Expected escaped question mark, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\|`) {
		t.Errorf("Expected escaped pipe, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\^`) {
		t.Errorf("Expected escaped caret, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\$`) {
		t.Errorf("Expected escaped dollar, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\.`) {
		t.Errorf("Expected escaped dot, got: %s", escaped)
	}

	if !strings.Contains(escaped, `\~`) {
		t.Errorf("Expected escaped tilde, got: %s", escaped)
	}
}

func TestEscapeForJQ_EmptyString(t *testing.T) {
	input := ""
	escaped := escapeForJQ(input)

	if escaped != "" {
		t.Errorf("Expected empty string, got: %s", escaped)
	}
}

func TestEscapeForJQ_Unicode(t *testing.T) {
	input := "test with unicode: 测试 🚀"
	escaped := escapeForJQ(input)

	// Unicode characters should not be escaped
	if !strings.Contains(escaped, "测试") {
		t.Errorf("Expected unicode characters to be preserved, got: %s", escaped)
	}

	if !strings.Contains(escaped, "🚀") {
		t.Errorf("Expected emoji to be preserved, got: %s", escaped)
	}
}

func TestCheckJQAvailability(t *testing.T) {
	builder := NewCommandBuilder()
	available := builder.checkJQAvailability()

	// This test depends on the system, so we just check it doesn't panic
	// In a real environment, this would check if jq is installed
	if available {
		t.Log("jq is available on this system")
	} else {
		t.Log("jq is not available on this system")
	}
}

func TestBuildSimpleCommand_JSONParsingEnabled(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Username:  "admin",
		Verb:      "create",
	}

	command := builder.buildSimpleCommand(params)

	// Should use JSON-aware command when jq is available
	// Note: This test assumes jq is not available in test environment
	// In a real environment with jq, it would use JSON parsing
	if strings.Contains(command, "jq -r") {
		t.Log("JSON parsing was used (jq available)")
	} else {
		t.Log("Fallback to grep parsing was used (jq not available)")
	}

	// Should still contain basic command structure
	if !strings.Contains(command, "oc adm node-logs --role=master") {
		t.Errorf("Expected oc command, got: %s", command)
	}

	if !strings.Contains(command, "--path=kube-apiserver/audit.log") {
		t.Errorf("Expected log path, got: %s", command)
	}
}

func TestBuildSimpleCommand_JSONParsingDisabled(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = false

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Username:  "admin",
		Verb:      "create",
	}

	command := builder.buildSimpleCommand(params)

	// Should not use JSON parsing
	if strings.Contains(command, "jq -r") {
		t.Errorf("Expected no jq command when JSON parsing is disabled, got: %s", command)
	}

	// Should use grep-based parsing
	if !strings.Contains(command, "grep") {
		t.Errorf("Expected grep-based parsing, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_NoFilters(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
	}

	command := builder.buildJSONAwareCommand(params)

	// Should contain base command
	if !strings.Contains(command, "oc adm node-logs --role=master") {
		t.Errorf("Expected oc command, got: %s", command)
	}

	// Should contain jq with just output formatting (no filters)
	if !strings.Contains(command, "jq -r") {
		t.Errorf("Expected jq command, got: %s", command)
	}

	// Should not contain select function when no filters
	if strings.Contains(command, "select(") {
		t.Errorf("Expected no select function when no filters, got: %s", command)
	}

	// Should contain output formatting
	if !strings.Contains(command, "timestamp:") {
		t.Errorf("Expected output formatting, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_MaxPatterns(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"pattern1", "pattern2", "pattern3", "pattern4", "pattern5"}, // More than max (3)
	}

	command := builder.buildJSONAwareCommand(params)

	// Should only include first 3 patterns
	if strings.Contains(command, "pattern4") {
		t.Errorf("Expected pattern4 to be excluded, got: %s", command)
	}

	if strings.Contains(command, "pattern5") {
		t.Errorf("Expected pattern5 to be excluded, got: %s", command)
	}

	// Should include first 3 patterns
	if !strings.Contains(command, "pattern1") {
		t.Errorf("Expected pattern1 to be included, got: %s", command)
	}

	if !strings.Contains(command, "pattern2") {
		t.Errorf("Expected pattern2 to be included, got: %s", command)
	}

	if !strings.Contains(command, "pattern3") {
		t.Errorf("Expected pattern3 to be included, got: %s", command)
	}
}

func TestBuildJSONAwareCommand_MaxExclusions(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Config.UseJSONParsing = true

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Exclude:   []string{"exclude1", "exclude2", "exclude3", "exclude4", "exclude5"}, // More than max (3)
	}

	command := builder.buildJSONAwareCommand(params)

	// Should only include first 3 exclusions
	if strings.Contains(command, "exclude4") {
		t.Errorf("Expected exclude4 to be excluded, got: %s", command)
	}

	if strings.Contains(command, "exclude5") {
		t.Errorf("Expected exclude5 to be excluded, got: %s", command)
	}

	// Should include first 3 exclusions
	if !strings.Contains(command, "exclude1") {
		t.Errorf("Expected exclude1 to be included, got: %s", command)
	}

	if !strings.Contains(command, "exclude2") {
		t.Errorf("Expected exclude2 to be included, got: %s", command)
	}

	if !strings.Contains(command, "exclude3") {
		t.Errorf("Expected exclude3 to be included, got: %s", command)
	}
}
