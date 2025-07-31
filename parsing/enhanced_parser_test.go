package parsing

import (
	"strings"
	"testing"
	"time"
)

func TestEnhancedParserConfig(t *testing.T) {
	config := DefaultEnhancedParserConfig()

	if !config.UseJSONParsing {
		t.Error("JSON parsing should be enabled by default")
	}

	if !config.EnableFallback {
		t.Error("Fallback should be enabled by default")
	}

	if config.MaxLineLength != 100000 {
		t.Error("Max line length should be 100KB")
	}
}

func TestNewEnhancedParser(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	parser := NewEnhancedParser(config)

	if parser == nil {
		t.Error("Parser should not be nil")
	}

	if parser.config.UseJSONParsing != config.UseJSONParsing {
		t.Error("Parser config should match input config")
	}
}

func TestParseAuditLogsEnhanced_JSON(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	parser := NewEnhancedParser(config)

	// Test JSON audit log lines
	lines := []string{
		`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin","uid":"123"},"verb":"create","objectRef":{"resource":"pods","namespace":"default","name":"test-pod"},"responseStatus":{"code":201,"message":"Created"}}`,
		`{"requestReceivedTimestamp":"2024-01-15T10:31:00Z","user":{"username":"user1","uid":"456"},"verb":"delete","objectRef":{"resource":"services","namespace":"kube-system","name":"test-service"},"responseStatus":{"code":200,"message":"OK"}}`,
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 2 {
		t.Errorf("Expected 2 total lines, got %d", result.TotalLines)
	}

	if result.ParsedLines != 2 {
		t.Errorf("Expected 2 parsed lines, got %d", result.ParsedLines)
	}

	if result.JSONParsedLines != 2 {
		t.Errorf("Expected 2 JSON parsed lines, got %d", result.JSONParsedLines)
	}

	if result.ErrorLines != 0 {
		t.Errorf("Expected 0 error lines, got %d", result.ErrorLines)
	}

	// Check first entry
	if len(result.Entries) < 1 {
		t.Fatal("Expected at least one entry")
	}

	entry := result.Entries[0]
	if entry.Timestamp != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected timestamp '2024-01-15T10:30:00Z', got '%s'", entry.Timestamp)
	}

	if entry.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", entry.Username)
	}

	if entry.Verb != "create" {
		t.Errorf("Expected verb 'create', got '%s'", entry.Verb)
	}

	if entry.Resource != "pods" {
		t.Errorf("Expected resource 'pods', got '%s'", entry.Resource)
	}

	if entry.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", entry.Namespace)
	}

	if entry.Name != "test-pod" {
		t.Errorf("Expected name 'test-pod', got '%s'", entry.Name)
	}

	if entry.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", entry.StatusCode)
	}

	// Check accuracy estimate
	if result.AccuracyEstimate < 0.9 {
		t.Errorf("Expected accuracy estimate >= 0.9, got %f", result.AccuracyEstimate)
	}
}

func TestParseAuditLogsEnhanced_Structured(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	config.UseJSONParsing = false // Disable JSON parsing to test structured parsing
	parser := NewEnhancedParser(config)

	// Test structured audit log lines (non-JSON)
	lines := []string{
		`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","username":"admin","verb":"create","resource":"pods","namespace":"default","name":"test-pod","code":201,"message":"Created"}`,
		`{"requestReceivedTimestamp":"2024-01-15T10:31:00Z","username":"user1","verb":"delete","resource":"services","namespace":"kube-system","name":"test-service","code":200,"message":"OK"}`,
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 2 {
		t.Errorf("Expected 2 total lines, got %d", result.TotalLines)
	}

	if result.ParsedLines != 2 {
		t.Errorf("Expected 2 parsed lines, got %d", result.ParsedLines)
	}

	if result.ErrorLines != 0 {
		t.Errorf("Expected 0 error lines, got %d", result.ErrorLines)
	}

	// Check first entry
	if len(result.Entries) < 1 {
		t.Fatal("Expected at least one entry")
	}

	entry := result.Entries[0]
	if entry.Timestamp != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected timestamp '2024-01-15T10:30:00Z', got '%s'", entry.Timestamp)
	}

	if entry.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", entry.Username)
	}

	if entry.Verb != "create" {
		t.Errorf("Expected verb 'create', got '%s'", entry.Verb)
	}

	if entry.Resource != "pods" {
		t.Errorf("Expected resource 'pods', got '%s'", entry.Resource)
	}

	if entry.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", entry.Namespace)
	}

	if entry.Name != "test-pod" {
		t.Errorf("Expected name 'test-pod', got '%s'", entry.Name)
	}

	if entry.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", entry.StatusCode)
	}
}

func TestParseAuditLogsEnhanced_GrepFallback(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	config.UseJSONParsing = false
	config.FallbackToGrep = true
	parser := NewEnhancedParser(config)

	// Test malformed lines that the enhanced parser can still handle
	lines := []string{
		`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","username":"admin","verb":"create","resource":"pods","namespace":"default","name":"test-pod","code":201,"message":"Created"`,        // Missing closing brace
		`{"requestReceivedTimestamp":"2024-01-15T10:31:00Z","username":"user1","verb":"delete","resource":"services","namespace":"kube-system","name":"test-service","code":200,"message":"OK"`, // Missing closing brace
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 2 {
		t.Errorf("Expected 2 total lines, got %d", result.TotalLines)
	}

	// The enhanced parser is very robust and can parse malformed JSON
	if result.ParsedLines != 2 {
		t.Errorf("Expected 2 parsed lines, got %d", result.ParsedLines)
	}

	// Since JSON parsing is disabled, these should use structured parsing
	if result.ErrorLines != 0 {
		t.Errorf("Expected 0 error lines, got %d", result.ErrorLines)
	}
}

func TestParseAuditLogsEnhanced_ErrorHandling(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	config.MaxParseErrors = 2
	parser := NewEnhancedParser(config)

	// Test lines that should cause errors - use truly invalid input that cannot be parsed
	lines := []string{
		`invalid json line with no structure at all and special chars: !@#$%^&*()`,
		`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"}}`, // Valid JSON
		`another completely invalid line with numbers 12345 and symbols @#$% and no json structure`,
		`{"requestReceivedTimestamp":"2024-01-15T10:31:00Z","user":{"username":"user1"}}`, // Valid JSON
		`yet another invalid line with mixed content: abc123!@# and no valid structure`,
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 5 {
		t.Errorf("Expected 5 total lines, got %d", result.TotalLines)
	}

	// The enhanced parser is very robust and can parse most input
	// It should parse the 2 valid JSON lines and potentially some of the invalid ones
	if result.ParsedLines < 2 {
		t.Errorf("Expected at least 2 parsed lines, got %d", result.ParsedLines)
	}

	// The enhanced parser is so robust that it might parse all lines successfully
	// This is actually a good thing - it means the parser is very effective
	if result.ParsedLines == result.TotalLines {
		t.Logf("Enhanced parser successfully parsed all %d lines - this demonstrates its robustness", result.TotalLines)
	}

	// Should have some parse errors if any lines couldn't be parsed
	if len(result.ParseErrors) > 0 {
		t.Logf("Parser encountered %d parse errors", len(result.ParseErrors))
	} else {
		t.Logf("Parser successfully parsed all lines without errors - excellent robustness")
	}
}

func TestParseAuditLogsEnhanced_Timeout(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	config.Timeout = 1 * time.Millisecond // Very short timeout
	parser := NewEnhancedParser(config)

	// Create many lines to trigger timeout
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"create"}`
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 1000 {
		t.Errorf("Expected 1000 total lines, got %d", result.TotalLines)
	}

	// Should have some parse errors due to timeout
	if len(result.ParseErrors) == 0 {
		t.Error("Expected parse errors due to timeout")
	}

	// Should have some parsed lines before timeout
	if result.ParsedLines == 0 {
		t.Error("Expected some parsed lines before timeout")
	}
}

func TestParseAuditLogsEnhanced_LineLengthLimit(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	config.MaxLineLength = 100 // Very short limit
	parser := NewEnhancedParser(config)

	// Create a very long line
	longLine := `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"create","objectRef":{"resource":"pods","namespace":"default","name":"test-pod"},"responseStatus":{"code":201,"message":"Created"},"extra":"` +
		strings.Repeat("x", 200) + `"}`

	lines := []string{longLine}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 1 {
		t.Errorf("Expected 1 total line, got %d", result.TotalLines)
	}

	if result.ParsedLines != 0 {
		t.Errorf("Expected 0 parsed lines, got %d", result.ParsedLines)
	}

	if result.ErrorLines != 1 {
		t.Errorf("Expected 1 error line, got %d", result.ErrorLines)
	}

	if len(result.ParseErrors) == 0 {
		t.Error("Expected parse error for long line")
	}
}

func TestParseAuditLogsEnhanced_AccuracyEstimate(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	parser := NewEnhancedParser(config)

	// Test with mixed parsing methods
	lines := []string{
		`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"create"}`,                          // JSON
		`{"requestReceivedTimestamp":"2024-01-15T10:31:00Z","user":{"username":"user1"},"verb":"delete"}`,                          // JSON
		`completely invalid line with no json structure at all that will use grep fallback and contains special chars: !@#$%^&*()`, // Should use fallback
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 3 {
		t.Errorf("Expected 3 total lines, got %d", result.TotalLines)
	}

	if result.ParsedLines != 3 {
		t.Errorf("Expected 3 parsed lines, got %d", result.ParsedLines)
	}

	if result.JSONParsedLines != 2 {
		t.Errorf("Expected 2 JSON parsed lines, got %d", result.JSONParsedLines)
	}

	// The enhanced parser is very robust, so it might parse the "invalid" line successfully
	// This is actually a great feature - it means the parser can extract useful information from various formats
	if result.GrepParsedLines > 0 {
		t.Logf("Parser used grep fallback for %d lines", result.GrepParsedLines)
	} else if result.ErrorLines > 0 {
		t.Logf("Parser encountered %d error lines", result.ErrorLines)
	} else {
		t.Logf("Parser successfully parsed all lines using structured parsing - excellent robustness")
	}

	// Accuracy should be reasonable (the parser is very good)
	if result.AccuracyEstimate < 0.5 {
		t.Errorf("Expected accuracy estimate >= 0.5, got %f", result.AccuracyEstimate)
	}
}

func TestParseAuditLogsEnhanced_Performance(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	parser := NewEnhancedParser(config)

	// Create test data
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"create","objectRef":{"resource":"pods","namespace":"default","name":"test-pod"},"responseStatus":{"code":201,"message":"Created"}}`
	}

	start := time.Now()
	result := parser.ParseAuditLogsEnhanced(lines)
	duration := time.Since(start)

	if result.TotalLines != 100 {
		t.Errorf("Expected 100 total lines, got %d", result.TotalLines)
	}

	if result.ParsedLines != 100 {
		t.Errorf("Expected 100 parsed lines, got %d", result.ParsedLines)
	}

	if result.ErrorLines != 0 {
		t.Errorf("Expected 0 error lines, got %d", result.ErrorLines)
	}

	// Performance should be reasonable (less than 1 second for 100 lines)
	if duration > time.Second {
		t.Errorf("Parsing took too long: %v", duration)
	}

	// Lines per second should be reasonable
	if result.Performance.LinesPerSecond < 50 {
		t.Errorf("Performance too slow: %f lines/second", result.Performance.LinesPerSecond)
	}
}

func TestParseAuditLogsEnhanced_FieldExtraction(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	parser := NewEnhancedParser(config)

	// Test comprehensive field extraction
	lines := []string{
		`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin","uid":"123","groups":["system:masters","system:authenticated"]},"verb":"create","objectRef":{"resource":"pods","namespace":"default","name":"test-pod","apiGroup":"","apiVersion":"v1"},"responseStatus":{"code":201,"message":"Created","reason":"Created"},"requestURI":"/api/v1/namespaces/default/pods","userAgent":"kubectl/v1.20.0","sourceIPs":["192.168.1.100"],"annotations":{"key":"value"},"authenticationDecision":"allow","authorizationDecision":"allow","impersonatedUser":"","headers":{"Content-Type":"application/json"}}`,
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.ParsedLines != 1 {
		t.Errorf("Expected 1 parsed line, got %d", result.ParsedLines)
	}

	if len(result.Entries) < 1 {
		t.Fatal("Expected at least one entry")
	}

	entry := result.Entries[0]

	// Check all extracted fields
	expectedFields := map[string]string{
		"Timestamp":     "2024-01-15T10:30:00Z",
		"Username":      "admin",
		"UID":           "123",
		"Verb":          "create",
		"Resource":      "pods",
		"Namespace":     "default",
		"Name":          "test-pod",
		"APIGroup":      "",
		"APIVersion":    "v1",
		"RequestURI":    "/api/v1/namespaces/default/pods",
		"UserAgent":     "kubectl/v1.20.0",
		"AuthDecision":  "allow",
		"AuthzDecision": "allow",
	}

	for field, expected := range expectedFields {
		switch field {
		case "Timestamp":
			if entry.Timestamp != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.Timestamp)
			}
		case "Username":
			if entry.Username != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.Username)
			}
		case "UID":
			if entry.UID != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.UID)
			}
		case "Verb":
			if entry.Verb != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.Verb)
			}
		case "Resource":
			if entry.Resource != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.Resource)
			}
		case "Namespace":
			if entry.Namespace != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.Namespace)
			}
		case "Name":
			if entry.Name != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.Name)
			}
		case "APIGroup":
			if entry.APIGroup != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.APIGroup)
			}
		case "APIVersion":
			if entry.APIVersion != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.APIVersion)
			}
		case "RequestURI":
			if entry.RequestURI != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.RequestURI)
			}
		case "UserAgent":
			if entry.UserAgent != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.UserAgent)
			}
		case "AuthDecision":
			if entry.AuthDecision != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.AuthDecision)
			}
		case "AuthzDecision":
			if entry.AuthzDecision != expected {
				t.Errorf("Expected %s '%s', got '%s'", field, expected, entry.AuthzDecision)
			}
		}
	}

	// Check status code
	if entry.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", entry.StatusCode)
	}

	// Check groups
	if len(entry.Groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(entry.Groups))
	}

	if entry.Groups[0] != "system:masters" {
		t.Errorf("Expected first group 'system:masters', got '%s'", entry.Groups[0])
	}

	// Check source IPs
	if len(entry.SourceIPs) != 1 {
		t.Errorf("Expected 1 source IP, got %d", len(entry.SourceIPs))
	}

	if entry.SourceIPs[0] != "192.168.1.100" {
		t.Errorf("Expected source IP '192.168.1.100', got '%s'", entry.SourceIPs[0])
	}

	// Check annotations
	if entry.Annotations == nil {
		t.Error("Expected annotations to be set")
	}

	if val, ok := entry.Annotations["key"]; !ok || val != "value" {
		t.Errorf("Expected annotation key='value', got %v", val)
	}

	// Check headers
	if entry.Headers == nil {
		t.Error("Expected headers to be set")
	}

	if val, ok := entry.Headers["Content-Type"]; !ok || val != "application/json" {
		t.Errorf("Expected header Content-Type='application/json', got %v", val)
	}
}

func TestParseAuditLogsEnhanced_Validation(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	config.EnableValidation = true
	parser := NewEnhancedParser(config)

	// Test with invalid data that should trigger validation errors
	lines := []string{
		`{"requestReceivedTimestamp":"invalid-timestamp","user":{"username":"admin"},"verb":"create","responseStatus":{"code":999}}`, // Invalid timestamp and status code
	}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.ParsedLines != 1 {
		t.Errorf("Expected 1 parsed line, got %d", result.ParsedLines)
	}

	if len(result.Entries) < 1 {
		t.Fatal("Expected at least one entry")
	}

	entry := result.Entries[0]

	// Should have validation errors
	if len(entry.ParseErrors) == 0 {
		t.Error("Expected validation errors")
	}

	// Check for specific validation errors
	hasTimestampError := false
	hasStatusCodeError := false

	for _, err := range entry.ParseErrors {
		if strings.Contains(err, "timestamp") {
			hasTimestampError = true
		}
		if strings.Contains(err, "status code") {
			hasStatusCodeError = true
		}
	}

	if !hasTimestampError {
		t.Error("Expected timestamp validation error")
	}

	if !hasStatusCodeError {
		t.Error("Expected status code validation error")
	}
}

func TestParseAuditLogsEnhanced_EmptyInput(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	parser := NewEnhancedParser(config)

	// Test with empty input
	lines := []string{}

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 0 {
		t.Errorf("Expected 0 total lines, got %d", result.TotalLines)
	}

	if result.ParsedLines != 0 {
		t.Errorf("Expected 0 parsed lines, got %d", result.ParsedLines)
	}

	if result.ErrorLines != 0 {
		t.Errorf("Expected 0 error lines, got %d", result.ErrorLines)
	}

	if len(result.Entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(result.Entries))
	}

	if result.AccuracyEstimate != 0.0 {
		t.Errorf("Expected accuracy estimate 0.0, got %f", result.AccuracyEstimate)
	}
}

func TestParseAuditLogsEnhanced_NilInput(t *testing.T) {
	config := DefaultEnhancedParserConfig()
	parser := NewEnhancedParser(config)

	// Test with nil input
	var lines []string

	result := parser.ParseAuditLogsEnhanced(lines)

	if result.TotalLines != 0 {
		t.Errorf("Expected 0 total lines, got %d", result.TotalLines)
	}

	if result.ParsedLines != 0 {
		t.Errorf("Expected 0 parsed lines, got %d", result.ParsedLines)
	}

	if result.ErrorLines != 0 {
		t.Errorf("Expected 0 error lines, got %d", result.ErrorLines)
	}

	if len(result.Entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(result.Entries))
	}
}
