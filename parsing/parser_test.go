package parsing

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestParseAuditLogLine(t *testing.T) {
	config := DefaultParserConfig()

	tests := []struct {
		name    string
		line    string
		wantErr bool
		check   func(t *testing.T, entry AuditLogEntry)
	}{
		{
			name:    "Valid JSON with all fields",
			line:    `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin","uid":"123","groups":["admin","users"]},"verb":"delete","objectRef":{"resource":"pods","namespace":"default","name":"test-pod"},"responseStatus":{"code":200,"message":"OK"},"requestURI":"/api/v1/namespaces/default/pods/test-pod","userAgent":"kubectl/v1.24.0","sourceIPs":["192.168.1.100"],"annotations":{"key":"value"}}`,
			wantErr: false,
			check: func(t *testing.T, entry AuditLogEntry) {
				if entry.Timestamp != "2024-01-15T10:30:00Z" {
					t.Errorf("Expected timestamp '2024-01-15T10:30:00Z', got '%s'", entry.Timestamp)
				}
				if entry.Username != "admin" {
					t.Errorf("Expected username 'admin', got '%s'", entry.Username)
				}
				if entry.UID != "123" {
					t.Errorf("Expected UID '123', got '%s'", entry.UID)
				}
				if len(entry.Groups) != 2 || entry.Groups[0] != "admin" || entry.Groups[1] != "users" {
					t.Errorf("Expected groups ['admin', 'users'], got %v", entry.Groups)
				}
				if entry.Verb != "delete" {
					t.Errorf("Expected verb 'delete', got '%s'", entry.Verb)
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
				if entry.StatusCode != 200 {
					t.Errorf("Expected status code 200, got %d", entry.StatusCode)
				}
				if entry.StatusMessage != "OK" {
					t.Errorf("Expected status message 'OK', got '%s'", entry.StatusMessage)
				}
				if entry.RequestURI != "/api/v1/namespaces/default/pods/test-pod" {
					t.Errorf("Expected request URI '/api/v1/namespaces/default/pods/test-pod', got '%s'", entry.RequestURI)
				}
				if entry.UserAgent != "kubectl/v1.24.0" {
					t.Errorf("Expected user agent 'kubectl/v1.24.0', got '%s'", entry.UserAgent)
				}
				if len(entry.SourceIPs) != 1 || entry.SourceIPs[0] != "192.168.1.100" {
					t.Errorf("Expected source IPs ['192.168.1.100'], got %v", entry.SourceIPs)
				}
				if entry.Annotations["key"] != "value" {
					t.Errorf("Expected annotation key='value', got %v", entry.Annotations)
				}
			},
		},
		{
			name:    "Valid JSON with minimal fields",
			line:    `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"get"}`,
			wantErr: false,
			check: func(t *testing.T, entry AuditLogEntry) {
				if entry.Timestamp != "2024-01-15T10:30:00Z" {
					t.Errorf("Expected timestamp '2024-01-15T10:30:00Z', got '%s'", entry.Timestamp)
				}
				if entry.Username != "admin" {
					t.Errorf("Expected username 'admin', got '%s'", entry.Username)
				}
				if entry.Verb != "get" {
					t.Errorf("Expected verb 'get', got '%s'", entry.Verb)
				}
			},
		},
		{
			name:    "Invalid JSON",
			line:    `{"malformed": json}`,
			wantErr: true,
			check:   func(t *testing.T, entry AuditLogEntry) {},
		},
		{
			name:    "Valid JSON with authentication fields",
			line:    `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"delete","authentication.openshift.io/decision":"allow","authorization.k8s.io/decision":"allow","impersonatedUser":"system:admin"}`,
			wantErr: false,
			check: func(t *testing.T, entry AuditLogEntry) {
				if entry.AuthDecision != "allow" {
					t.Errorf("Expected auth decision 'allow', got '%s'", entry.AuthDecision)
				}
				if entry.AuthzDecision != "allow" {
					t.Errorf("Expected authz decision 'allow', got '%s'", entry.AuthzDecision)
				}
				if entry.ImpersonatedUser != "system:admin" {
					t.Errorf("Expected impersonated user 'system:admin', got '%s'", entry.ImpersonatedUser)
				}
			},
		},
		{
			name:    "Valid JSON with nested structures",
			line:    `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin","extra":{"scopes":["user:full"],"client_id":"console"}},"verb":"create","objectRef":{"resource":"pods","namespace":"default","apiGroup":"","apiVersion":"v1"},"responseStatus":{"code":201,"message":"Created","reason":"Success"}}`,
			wantErr: false,
			check: func(t *testing.T, entry AuditLogEntry) {
				if entry.APIGroup != "" {
					t.Errorf("Expected empty API group, got '%s'", entry.APIGroup)
				}
				if entry.APIVersion != "v1" {
					t.Errorf("Expected API version 'v1', got '%s'", entry.APIVersion)
				}
				if entry.StatusCode != 201 {
					t.Errorf("Expected status code 201, got %d", entry.StatusCode)
				}
				if entry.StatusReason != "Success" {
					t.Errorf("Expected status reason 'Success', got '%s'", entry.StatusReason)
				}
				if entry.Extra["scopes"] == nil {
					t.Errorf("Expected extra scopes field")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := ParseAuditLogLine(tt.line, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAuditLogLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				tt.check(t, entry)
			}
		})
	}
}

func TestParseAuditLogs(t *testing.T) {
	config := DefaultParserConfig()

	validLines := []string{
		`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"get","objectRef":{"resource":"pods","namespace":"default"},"responseStatus":{"code":200,"message":"OK"}}`,
		`{"requestReceivedTimestamp":"2024-01-15T10:31:00Z","user":{"username":"user1"},"verb":"create","objectRef":{"resource":"services","namespace":"default"},"responseStatus":{"code":201,"message":"Created"}}`,
		`{"requestReceivedTimestamp":"2024-01-15T10:32:00Z","user":{"username":"user2"},"verb":"delete","objectRef":{"resource":"pods","namespace":"kube-system"},"responseStatus":{"code":404,"message":"Not Found"}}`,
	}

	invalidLines := []string{
		`{"malformed": json}`,
		`{"malformed": again}`,
	}

	t.Run("Valid lines only", func(t *testing.T) {
		result := ParseAuditLogs(validLines, config)

		if result.TotalLines != 3 {
			t.Errorf("Expected total lines 3, got %d", result.TotalLines)
		}
		if result.ParsedLines != 3 {
			t.Errorf("Expected parsed lines 3, got %d", result.ParsedLines)
		}
		if result.ErrorLines != 0 {
			t.Errorf("Expected error lines 0, got %d", result.ErrorLines)
		}
		if len(result.Entries) != 3 {
			t.Errorf("Expected 3 entries, got %d", len(result.Entries))
		}
		if result.ParseTime <= 0 {
			t.Errorf("Expected positive parse time, got %v", result.ParseTime)
		}
		if result.Performance.LinesPerSecond <= 0 {
			t.Errorf("Expected positive lines per second, got %f", result.Performance.LinesPerSecond)
		}
	})

	t.Run("Mixed valid and invalid lines", func(t *testing.T) {
		mixedLines := append(validLines, invalidLines...)
		result := ParseAuditLogs(mixedLines, config)

		if result.TotalLines != 5 {
			t.Errorf("Expected total lines 5, got %d", result.TotalLines)
		}
		if result.ParsedLines != 3 {
			t.Errorf("Expected parsed lines 3, got %d", result.ParsedLines)
		}
		if result.ErrorLines != 2 {
			t.Errorf("Expected error lines 2, got %d", result.ErrorLines)
		}
		if len(result.ParseErrors) != 2 {
			t.Errorf("Expected 2 parse errors, got %d", len(result.ParseErrors))
		}
	})

	t.Run("Timeout configuration", func(t *testing.T) {
		timeoutConfig := ParserConfig{
			MaxLineLength:    100000,
			MaxParseErrors:   1000,
			Timeout:          1 * time.Microsecond, // Very short timeout
			EnableValidation: true,
			EnableMetrics:    true,
		}

		result := ParseAuditLogs(validLines, timeoutConfig)

		if len(result.ParseErrors) == 0 {
			t.Errorf("Expected timeout error, got none")
		}
		if !strings.Contains(result.ParseErrors[0], "timeout") {
			t.Errorf("Expected timeout error message, got '%s'", result.ParseErrors[0])
		}
	})

	t.Run("Line length limit", func(t *testing.T) {
		longLine := strings.Repeat("a", 100001) // Exceeds 100KB limit
		lines := []string{longLine}

		result := ParseAuditLogs(lines, config)

		if result.ErrorLines != 1 {
			t.Errorf("Expected 1 error line, got %d", result.ErrorLines)
		}
		if !strings.Contains(result.ParseErrors[0], "exceeds max length") {
			t.Errorf("Expected max length error message, got '%s'", result.ParseErrors[0])
		}
	})
}

func TestValidateEntry(t *testing.T) {
	tests := []struct {
		name    string
		entry   AuditLogEntry
		wantErr bool
	}{
		{
			name: "Valid entry",
			entry: AuditLogEntry{
				Timestamp:  "2024-01-15T10:30:00Z",
				Username:   "admin",
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "Invalid timestamp format",
			entry: AuditLogEntry{
				Timestamp: "invalid-timestamp",
				Username:  "admin",
			},
			wantErr: true,
		},
		{
			name: "Invalid status code",
			entry: AuditLogEntry{
				Timestamp:  "2024-01-15T10:30:00Z",
				Username:   "admin",
				StatusCode: 999, // Invalid status code
			},
			wantErr: true,
		},
		{
			name: "Missing user identification",
			entry: AuditLogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				// No username or UID
			},
			wantErr: true,
		},
		{
			name: "Valid entry with UID only",
			entry: AuditLogEntry{
				Timestamp:  "2024-01-15T10:30:00Z",
				UID:        "123",
				StatusCode: 200,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEntry(&tt.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateSummary(t *testing.T) {
	entries := []AuditLogEntry{
		{
			Username:   "admin",
			StatusCode: 200,
			Verb:       "get",
			Resource:   "pods",
		},
		{
			Username:   "admin",
			StatusCode: 201,
			Verb:       "create",
			Resource:   "services",
		},
		{
			Username:   "user1",
			StatusCode: 404,
			Verb:       "delete",
			Resource:   "pods",
		},
		{
			Username:   "user2",
			StatusCode: 200,
			Verb:       "get",
			Resource:   "configmaps",
		},
	}

	summary := GenerateSummary(entries, nil)

	// Check that summary contains expected information
	if !strings.Contains(summary, "Found 4 audit entries") {
		t.Errorf("Summary should contain 'Found 4 audit entries', got: %s", summary)
	}
	if !strings.Contains(summary, "admin (2)") {
		t.Errorf("Summary should contain 'admin (2)', got: %s", summary)
	}
	if !strings.Contains(summary, "user1 (1)") {
		t.Errorf("Summary should contain 'user1 (1)', got: %s", summary)
	}
	if !strings.Contains(summary, "200 (2)") {
		t.Errorf("Summary should contain '200 (2)', got: %s", summary)
	}
	if !strings.Contains(summary, "get (2)") {
		t.Errorf("Summary should contain 'get (2)', got: %s", summary)
	}
	if !strings.Contains(summary, "pods (2)") {
		t.Errorf("Summary should contain 'pods (2)', got: %s", summary)
	}
}

func TestParseStatusCodes(t *testing.T) {
	entries := []AuditLogEntry{
		{StatusCode: 200},
		{StatusCode: 201},
		{StatusCode: 400},
		{StatusCode: 401},
		{StatusCode: 500},
		{StatusCode: 502},
		{StatusCode: 0}, // Should be ignored
	}

	result := ParseStatusCodes(entries)

	expected := map[string]int{
		"success":      2, // 200, 201
		"client_error": 2, // 400, 401
		"server_error": 2, // 500, 502
	}

	for category, count := range expected {
		if result[category] != count {
			t.Errorf("Expected %d entries in category '%s', got %d", count, category, result[category])
		}
	}
}

func TestConvertLegacyEntries(t *testing.T) {
	legacyEntries := []map[string]interface{}{
		{
			"timestamp":         "2024-01-15T10:30:00Z",
			"username":          "admin",
			"verb":              "get",
			"resource":          "pods",
			"namespace":         "default",
			"status_code":       "200",
			"status_message":    "OK",
			"request_uri":       "/api/v1/namespaces/default/pods",
			"user_agent":        "kubectl/v1.24.0",
			"source_ips":        []string{"192.168.1.100"},
			"auth_decision":     "allow",
			"authz_decision":    "allow",
			"impersonated_user": "system:admin",
		},
	}

	entries := ConvertLegacyEntries(legacyEntries)

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.Timestamp != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected timestamp '2024-01-15T10:30:00Z', got '%s'", entry.Timestamp)
	}
	if entry.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", entry.Username)
	}
	if entry.Verb != "get" {
		t.Errorf("Expected verb 'get', got '%s'", entry.Verb)
	}
	if entry.Resource != "pods" {
		t.Errorf("Expected resource 'pods', got '%s'", entry.Resource)
	}
	if entry.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", entry.Namespace)
	}
	if entry.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", entry.StatusCode)
	}
	if entry.StatusMessage != "OK" {
		t.Errorf("Expected status message 'OK', got '%s'", entry.StatusMessage)
	}
	if entry.RequestURI != "/api/v1/namespaces/default/pods" {
		t.Errorf("Expected request URI '/api/v1/namespaces/default/pods', got '%s'", entry.RequestURI)
	}
	if entry.UserAgent != "kubectl/v1.24.0" {
		t.Errorf("Expected user agent 'kubectl/v1.24.0', got '%s'", entry.UserAgent)
	}
	if len(entry.SourceIPs) != 1 || entry.SourceIPs[0] != "192.168.1.100" {
		t.Errorf("Expected source IPs ['192.168.1.100'], got %v", entry.SourceIPs)
	}
	if entry.AuthDecision != "allow" {
		t.Errorf("Expected auth decision 'allow', got '%s'", entry.AuthDecision)
	}
	if entry.AuthzDecision != "allow" {
		t.Errorf("Expected authz decision 'allow', got '%s'", entry.AuthzDecision)
	}
	if entry.ImpersonatedUser != "system:admin" {
		t.Errorf("Expected impersonated user 'system:admin', got '%s'", entry.ImpersonatedUser)
	}
}

func TestParseAuditLogField(t *testing.T) {
	line := `{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"admin"},"verb":"get","objectRef":{"resource":"pods","namespace":"default"},"responseStatus":{"code":200,"message":"OK"}}`

	tests := []struct {
		fieldName string
		expected  string
		found     bool
	}{
		{
			fieldName: "RequestReceivedTimestamp",
			expected:  "2024-01-15T10:30:00Z",
			found:     true,
		},
		{
			fieldName: "Username",
			expected:  "admin",
			found:     true,
		},
		{
			fieldName: "Verb",
			expected:  "get",
			found:     true,
		},
		{
			fieldName: "Resource",
			expected:  "pods",
			found:     true,
		},
		{
			fieldName: "Namespace",
			expected:  "default",
			found:     true,
		},
		{
			fieldName: "Code",
			expected:  "200",
			found:     true,
		},
		{
			fieldName: "Message",
			expected:  "OK",
			found:     true,
		},
		{
			fieldName: "NonExistentField",
			expected:  "",
			found:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			value, found := ParseAuditLogField(line, tt.fieldName)
			if found != tt.found {
				t.Errorf("ParseAuditLogField() found = %v, want %v", found, tt.found)
			}
			if found && value != tt.expected {
				t.Errorf("ParseAuditLogField() value = %v, want %v", value, tt.expected)
			}
		})
	}
}

func TestParserPerformance(t *testing.T) {
	// Generate test data
	var lines []string
	for i := 0; i < 1000; i++ {
		line := fmt.Sprintf(`{"requestReceivedTimestamp":"2024-01-15T10:30:00Z","user":{"username":"user%d"},"verb":"get","objectRef":{"resource":"pods","namespace":"default"},"responseStatus":{"code":200,"message":"OK"}}`, i)
		lines = append(lines, line)
	}

	config := DefaultParserConfig()
	startTime := time.Now()
	result := ParseAuditLogs(lines, config)
	duration := time.Since(startTime)

	// Performance assertions
	if result.TotalLines != 1000 {
		t.Errorf("Expected 1000 total lines, got %d", result.TotalLines)
	}
	if result.ParsedLines != 1000 {
		t.Errorf("Expected 1000 parsed lines, got %d", result.ParsedLines)
	}
	if result.ErrorLines != 0 {
		t.Errorf("Expected 0 error lines, got %d", result.ErrorLines)
	}
	if duration > 5*time.Second {
		t.Errorf("Parsing took too long: %v", duration)
	}
	if result.Performance.LinesPerSecond < 100 {
		t.Errorf("Performance too slow: %f lines/second", result.Performance.LinesPerSecond)
	}
	if result.Performance.AverageLineSize <= 0 {
		t.Errorf("Expected positive average line size, got %d", result.Performance.AverageLineSize)
	}
}

func TestParserConfig(t *testing.T) {
	config := DefaultParserConfig()

	if config.MaxLineLength != 100000 {
		t.Errorf("Expected MaxLineLength 100000, got %d", config.MaxLineLength)
	}
	if config.MaxParseErrors != 1000 {
		t.Errorf("Expected MaxParseErrors 1000, got %d", config.MaxParseErrors)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected Timeout 30s, got %v", config.Timeout)
	}
	if !config.EnableValidation {
		t.Errorf("Expected EnableValidation true, got %v", config.EnableValidation)
	}
	if !config.EnableMetrics {
		t.Errorf("Expected EnableMetrics true, got %v", config.EnableMetrics)
	}
}

func TestJSONMarshalUnmarshal(t *testing.T) {
	entry := AuditLogEntry{
		Timestamp:        "2024-01-15T10:30:00Z",
		Username:         "admin",
		UID:              "123",
		Groups:           []string{"admin", "users"},
		Verb:             "delete",
		Resource:         "pods",
		Namespace:        "default",
		Name:             "test-pod",
		APIGroup:         "",
		APIVersion:       "v1",
		RequestURI:       "/api/v1/namespaces/default/pods/test-pod",
		UserAgent:        "kubectl/v1.24.0",
		SourceIPs:        []string{"192.168.1.100"},
		StatusCode:       200,
		StatusMessage:    "OK",
		StatusReason:     "Success",
		AuthDecision:     "allow",
		AuthzDecision:    "allow",
		ImpersonatedUser: "system:admin",
		Annotations:      map[string]interface{}{"key": "value"},
		Extra:            map[string]interface{}{"scopes": []string{"user:full"}},
		Headers:          map[string]interface{}{"Authorization": "Bearer token"},
		RawLine:          "original line",
		ParseErrors:      []string{"error1", "error2"},
		ParseTime:        time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal entry: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledEntry AuditLogEntry
	err = json.Unmarshal(jsonData, &unmarshaledEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal entry: %v", err)
	}

	// Verify key fields are preserved
	if unmarshaledEntry.Timestamp != entry.Timestamp {
		t.Errorf("Timestamp not preserved: expected %s, got %s", entry.Timestamp, unmarshaledEntry.Timestamp)
	}
	if unmarshaledEntry.Username != entry.Username {
		t.Errorf("Username not preserved: expected %s, got %s", entry.Username, unmarshaledEntry.Username)
	}
	if unmarshaledEntry.Verb != entry.Verb {
		t.Errorf("Verb not preserved: expected %s, got %s", entry.Verb, unmarshaledEntry.Verb)
	}
	if unmarshaledEntry.StatusCode != entry.StatusCode {
		t.Errorf("StatusCode not preserved: expected %d, got %d", entry.StatusCode, unmarshaledEntry.StatusCode)
	}
	if len(unmarshaledEntry.Groups) != len(entry.Groups) {
		t.Errorf("Groups length not preserved: expected %d, got %d", len(entry.Groups), len(unmarshaledEntry.Groups))
	}
	if len(unmarshaledEntry.SourceIPs) != len(entry.SourceIPs) {
		t.Errorf("SourceIPs length not preserved: expected %d, got %d", len(entry.SourceIPs), len(unmarshaledEntry.SourceIPs))
	}
}
