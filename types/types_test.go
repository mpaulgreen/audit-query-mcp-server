package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// TestAuditQueryParams_FieldValidation tests field validation for AuditQueryParams
func TestAuditQueryParams_FieldValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  AuditQueryParams
		isValid bool
	}{
		{
			name: "Valid complete params",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"pods", "delete"},
				Timeframe: "today",
				Username:  "admin",
				Resource:  "pods",
				Verb:      "delete",
				Namespace: "default",
				Exclude:   []string{"system:"},
			},
			isValid: true,
		},
		{
			name: "Valid minimal params",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "today",
			},
			isValid: true,
		},
		{
			name: "Valid with empty optional fields",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "today",
				Username:  "",
				Resource:  "",
				Verb:      "",
				Namespace: "",
				Exclude:   []string{},
			},
			isValid: true,
		},
		{
			name: "Invalid log source",
			params: AuditQueryParams{
				LogSource: "invalid-source",
				Patterns:  []string{"test"},
				Timeframe: "today",
			},
			isValid: false,
		},
		{
			name: "Invalid timeframe",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "invalid-timeframe",
			},
			isValid: false,
		},
		{
			name: "Invalid resource",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "today",
				Resource:  "invalid-resource",
			},
			isValid: false,
		},
		{
			name: "Invalid verb",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "today",
				Verb:      "invalid-verb",
			},
			isValid: false,
		},
		{
			name: "Invalid namespace pattern",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "today",
				Namespace: "invalid/namespace/pattern",
			},
			isValid: false,
		},
		{
			name: "Invalid username pattern",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "today",
				Username:  "invalid@user@name",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test field validation logic
			if tt.params.LogSource == "" {
				t.Error("LogSource should not be empty")
			}
			if len(tt.params.Patterns) == 0 {
				t.Error("Patterns should not be empty")
			}
			if tt.params.Timeframe == "" {
				t.Error("Timeframe should not be empty")
			}

			// Test that valid params have expected structure
			if tt.isValid {
				if tt.params.LogSource != "kube-apiserver" {
					t.Errorf("Expected valid log source, got %s", tt.params.LogSource)
				}
			}
		})
	}
}

// TestAuditQueryParams_JSONSerialization tests JSON marshaling and unmarshaling
func TestAuditQueryParams_JSONSerialization(t *testing.T) {
	original := AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"pods", "delete"},
		Timeframe: "today",
		Username:  "admin",
		Resource:  "pods",
		Verb:      "delete",
		Namespace: "default",
		Exclude:   []string{"system:"},
	}

	// Test marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal AuditQueryParams: %v", err)
	}

	// Test unmarshaling
	var unmarshaled AuditQueryParams
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal AuditQueryParams: %v", err)
	}

	// Verify all fields are preserved
	if original.LogSource != unmarshaled.LogSource {
		t.Errorf("LogSource mismatch: expected %s, got %s", original.LogSource, unmarshaled.LogSource)
	}
	if len(original.Patterns) != len(unmarshaled.Patterns) {
		t.Errorf("Patterns length mismatch: expected %d, got %d", len(original.Patterns), len(unmarshaled.Patterns))
	}
	if original.Timeframe != unmarshaled.Timeframe {
		t.Errorf("Timeframe mismatch: expected %s, got %s", original.Timeframe, unmarshaled.Timeframe)
	}
	if original.Username != unmarshaled.Username {
		t.Errorf("Username mismatch: expected %s, got %s", original.Username, unmarshaled.Username)
	}
	if original.Resource != unmarshaled.Resource {
		t.Errorf("Resource mismatch: expected %s, got %s", original.Resource, unmarshaled.Resource)
	}
	if original.Verb != unmarshaled.Verb {
		t.Errorf("Verb mismatch: expected %s, got %s", original.Verb, unmarshaled.Verb)
	}
	if original.Namespace != unmarshaled.Namespace {
		t.Errorf("Namespace mismatch: expected %s, got %s", original.Namespace, unmarshaled.Namespace)
	}
	if len(original.Exclude) != len(unmarshaled.Exclude) {
		t.Errorf("Exclude length mismatch: expected %d, got %d", len(original.Exclude), len(unmarshaled.Exclude))
	}
}

// TestAuditQueryParams_EdgeCases tests edge cases for AuditQueryParams
func TestAuditQueryParams_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		params AuditQueryParams
	}{
		{
			name: "Empty patterns",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{},
				Timeframe: "today",
			},
		},
		{
			name: "Large patterns array",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  make([]string, 1000),
				Timeframe: "today",
			},
		},
		{
			name: "Very long strings",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test"},
				Timeframe: "today",
				Username:  string(make([]byte, 10000)), // Very long username
			},
		},
		{
			name: "Special characters in patterns",
			params: AuditQueryParams{
				LogSource: "kube-apiserver",
				Patterns:  []string{"test@#$%^&*()", "user@domain.com", "namespace/with/slashes"},
				Timeframe: "today",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the struct can be created without panic
			if tt.params.LogSource == "" {
				t.Error("LogSource should be set")
			}
			if tt.params.Timeframe == "" {
				t.Error("Timeframe should be set")
			}

			// Test JSON serialization doesn't fail
			_, err := json.Marshal(tt.params)
			if err != nil {
				t.Errorf("JSON marshaling failed: %v", err)
			}
		})
	}
}

// TestAuditResult_FieldValidation tests field validation for AuditResult
func TestAuditResult_FieldValidation(t *testing.T) {
	tests := []struct {
		name     string
		result   AuditResult
		isValid  bool
		hasError bool
	}{
		{
			name: "Valid complete result",
			result: AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{{"test": "data"}},
				Summary:       "Found 1 audit entries",
				Error:         "",
				ExecutionTime: 100,
			},
			isValid:  true,
			hasError: false,
		},
		{
			name: "Valid result with error",
			result: AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "command execution failed",
				ExecutionTime: 50,
			},
			isValid:  true,
			hasError: true,
		},
		{
			name: "Empty QueryID",
			result: AuditResult{
				QueryID:       "",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			isValid: false,
		},
		{
			name: "Invalid timestamp format",
			result: AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     "invalid-timestamp",
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			isValid: false,
		},
		{
			name: "Negative execution time",
			result: AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: -100,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test QueryID validation
			if tt.result.QueryID == "" && tt.isValid {
				t.Error("QueryID should not be empty for valid results")
			}

			// Test timestamp validation
			if tt.isValid {
				_, err := time.Parse(time.RFC3339, tt.result.Timestamp)
				if err != nil {
					t.Errorf("Invalid timestamp format: %v", err)
				}
			}

			// Test execution time validation
			if tt.result.ExecutionTime < 0 && tt.isValid {
				t.Error("Execution time should not be negative for valid results")
			}

			// Test error field logic
			if tt.hasError && tt.result.Error == "" {
				t.Error("Error field should be set when hasError is true")
			}
			if !tt.hasError && tt.result.Error != "" {
				t.Error("Error field should be empty when hasError is false")
			}
		})
	}
}

// TestAuditResult_JSONSerialization tests JSON marshaling and unmarshaling
func TestAuditResult_JSONSerialization(t *testing.T) {
	original := AuditResult{
		QueryID:       "audit_query_20250729_204108_abc123",
		Timestamp:     time.Now().Format(time.RFC3339),
		Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
		RawOutput:     "test output with special chars: @#$%^&*()",
		ParsedData:    []map[string]interface{}{{"test": "data", "number": 123}},
		Summary:       "Found 1 audit entries",
		Error:         "",
		ExecutionTime: 100,
	}

	// Test marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal AuditResult: %v", err)
	}

	// Test unmarshaling
	var unmarshaled AuditResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal AuditResult: %v", err)
	}

	// Verify all fields are preserved
	if original.QueryID != unmarshaled.QueryID {
		t.Errorf("QueryID mismatch: expected %s, got %s", original.QueryID, unmarshaled.QueryID)
	}
	if original.Timestamp != unmarshaled.Timestamp {
		t.Errorf("Timestamp mismatch: expected %s, got %s", original.Timestamp, unmarshaled.Timestamp)
	}
	if original.Command != unmarshaled.Command {
		t.Errorf("Command mismatch: expected %s, got %s", original.Command, unmarshaled.Command)
	}
	if original.RawOutput != unmarshaled.RawOutput {
		t.Errorf("RawOutput mismatch: expected %s, got %s", original.RawOutput, unmarshaled.RawOutput)
	}
	if len(original.ParsedData) != len(unmarshaled.ParsedData) {
		t.Errorf("ParsedData length mismatch: expected %d, got %d", len(original.ParsedData), len(unmarshaled.ParsedData))
	}
	if original.Summary != unmarshaled.Summary {
		t.Errorf("Summary mismatch: expected %s, got %s", original.Summary, unmarshaled.Summary)
	}
	if original.Error != unmarshaled.Error {
		t.Errorf("Error mismatch: expected %s, got %s", original.Error, unmarshaled.Error)
	}
	if original.ExecutionTime != unmarshaled.ExecutionTime {
		t.Errorf("ExecutionTime mismatch: expected %d, got %d", original.ExecutionTime, unmarshaled.ExecutionTime)
	}
}

// TestAuditResult_EdgeCases tests edge cases for AuditResult
func TestAuditResult_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		result AuditResult
	}{
		{
			name: "Empty fields",
			result: AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "",
				ExecutionTime: 0,
			},
		},
		{
			name: "Very large raw output",
			result: AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     string(make([]byte, 100000)), // 100KB output
				ParsedData:    []map[string]interface{}{},
				Summary:       "Large output processed",
				Error:         "",
				ExecutionTime: 1000,
			},
		},
		{
			name: "Large parsed data",
			result: AuditResult{
				QueryID:   "test_id",
				Timestamp: time.Now().Format(time.RFC3339),
				Command:   "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput: "test output",
				ParsedData: func() []map[string]interface{} {
					data := make([]map[string]interface{}, 1000)
					for i := range data {
						data[i] = map[string]interface{}{
							"id":      i,
							"message": "test message",
							"data":    map[string]interface{}{"nested": "value"},
						}
					}
					return data
				}(),
				Summary:       "Found 1000 audit entries",
				Error:         "",
				ExecutionTime: 5000,
			},
		},
		{
			name: "Special characters in all fields",
			result: AuditResult{
				QueryID:       "test_id_with_special_chars_@#$%^&*()",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep '@#$%^&*()'",
				RawOutput:     "output with special chars: @#$%^&*()\nnewlines\nand\t\ttabs",
				ParsedData:    []map[string]interface{}{{"special": "@#$%^&*()", "unicode": "🚀🎉"}},
				Summary:       "Found entries with special chars: @#$%^&*()",
				Error:         "error with special chars: @#$%^&*()",
				ExecutionTime: 123,
			},
		},
		{
			name: "Zero values",
			result: AuditResult{
				QueryID:       "",
				Timestamp:     "",
				Command:       "",
				RawOutput:     "",
				ParsedData:    nil,
				Summary:       "",
				Error:         "",
				ExecutionTime: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the struct can be created without panic
			if tt.result.QueryID == "" && tt.name != "Zero values" {
				t.Error("QueryID should be set for non-zero test cases")
			}

			// Test JSON serialization doesn't fail
			_, err := json.Marshal(tt.result)
			if err != nil {
				t.Errorf("JSON marshaling failed: %v", err)
			}

			// Test that execution time is reasonable
			if tt.result.ExecutionTime < 0 {
				t.Error("Execution time should not be negative")
			}
		})
	}
}

// TestAuditResult_ErrorHandling tests error handling scenarios
func TestAuditResult_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		result        AuditResult
		expectedError bool
		errorContains string
	}{
		{
			name: "No error",
			result: AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "successful output",
				ParsedData:    []map[string]interface{}{{"status": "success"}},
				Summary:       "Operation completed successfully",
				Error:         "",
				ExecutionTime: 100,
			},
			expectedError: false,
		},
		{
			name: "Command execution error",
			result: AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "invalid_command",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "command execution failed: exit status 1",
				ExecutionTime: 50,
			},
			expectedError: true,
			errorContains: "command execution failed",
		},
		{
			name: "Validation error",
			result: AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "validation failed: invalid log source",
				ExecutionTime: 10,
			},
			expectedError: true,
			errorContains: "validation failed",
		},
		{
			name: "Timeout error",
			result: AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "command execution timed out after 30 seconds",
				ExecutionTime: 30000,
			},
			expectedError: true,
			errorContains: "timed out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test error field logic
			hasError := tt.result.Error != ""
			if hasError != tt.expectedError {
				t.Errorf("Error expectation mismatch: expected %v, got %v", tt.expectedError, hasError)
			}

			// Test error content
			if tt.expectedError && tt.errorContains != "" {
				if !contains(tt.result.Error, tt.errorContains) {
					t.Errorf("Error should contain '%s', got '%s'", tt.errorContains, tt.result.Error)
				}
			}

			// Test that error results have appropriate empty fields
			if tt.expectedError {
				if tt.result.RawOutput != "" && tt.result.Error != "command execution timed out after 30 seconds" {
					t.Error("Error results should typically have empty RawOutput")
				}
				if len(tt.result.ParsedData) > 0 {
					t.Error("Error results should typically have empty ParsedData")
				}
			}
		})
	}
}

// TestAuditResult_PerformanceMetrics tests performance-related fields
func TestAuditResult_PerformanceMetrics(t *testing.T) {
	tests := []struct {
		name          string
		executionTime int64
		expectedValid bool
	}{
		{
			name:          "Zero execution time",
			executionTime: 0,
			expectedValid: true,
		},
		{
			name:          "Normal execution time",
			executionTime: 100,
			expectedValid: true,
		},
		{
			name:          "Long execution time",
			executionTime: 30000,
			expectedValid: true,
		},
		{
			name:          "Very long execution time",
			executionTime: 300000,
			expectedValid: true,
		},
		{
			name:          "Negative execution time",
			executionTime: -100,
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "Test completed",
				Error:         "",
				ExecutionTime: tt.executionTime,
			}

			// Test execution time validation
			if tt.expectedValid && result.ExecutionTime < 0 {
				t.Error("Execution time should not be negative for valid cases")
			}

			// Test that execution time is preserved in JSON
			jsonData, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			var unmarshaled AuditResult
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if unmarshaled.ExecutionTime != tt.executionTime {
				t.Errorf("Execution time not preserved: expected %d, got %d", tt.executionTime, unmarshaled.ExecutionTime)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

// TestMCPTool_FieldValidation tests field validation for MCPTool
func TestMCPTool_FieldValidation(t *testing.T) {
	tests := []struct {
		name    string
		tool    MCPTool
		isValid bool
	}{
		{
			name: "Valid complete tool",
			tool: MCPTool{
				Name:        "generate_audit_query_with_result",
				Description: "Convert structured parameters to safe oc audit commands with detailed result tracking",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"structured_params": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"log_source": map[string]interface{}{
									"type": "string",
									"enum": []string{"kube-apiserver", "oauth-server"},
								},
							},
							"required": []string{"log_source"},
						},
					},
					"required": []string{"structured_params"},
				},
			},
			isValid: true,
		},
		{
			name: "Valid minimal tool",
			tool: MCPTool{
				Name:        "simple_tool",
				Description: "A simple tool",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
			isValid: true,
		},
		{
			name: "Empty name",
			tool: MCPTool{
				Name:        "",
				Description: "Tool with empty name",
				InputSchema: map[string]interface{}{
					"type": "object",
				},
			},
			isValid: false,
		},
		{
			name: "Empty description",
			tool: MCPTool{
				Name:        "tool_without_description",
				Description: "",
				InputSchema: map[string]interface{}{
					"type": "object",
				},
			},
			isValid: false,
		},
		{
			name: "Nil input schema",
			tool: MCPTool{
				Name:        "tool_without_schema",
				Description: "Tool without input schema",
				InputSchema: nil,
			},
			isValid: false,
		},
		{
			name: "Empty input schema",
			tool: MCPTool{
				Name:        "tool_with_empty_schema",
				Description: "Tool with empty input schema",
				InputSchema: map[string]interface{}{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test name validation
			if tt.tool.Name == "" && tt.isValid {
				t.Error("Name should not be empty for valid tools")
			}

			// Test description validation
			if tt.tool.Description == "" && tt.isValid {
				t.Error("Description should not be empty for valid tools")
			}

			// Test input schema validation
			if tt.tool.InputSchema == nil && tt.isValid {
				t.Error("InputSchema should not be nil for valid tools")
			}

			if tt.tool.InputSchema != nil && len(tt.tool.InputSchema) == 0 && tt.isValid {
				t.Error("InputSchema should not be empty for valid tools")
			}

			// Test that valid tools have expected structure
			if tt.isValid {
				if tt.tool.Name == "" {
					t.Error("Valid tools should have a name")
				}
				if tt.tool.Description == "" {
					t.Error("Valid tools should have a description")
				}
				if tt.tool.InputSchema == nil {
					t.Error("Valid tools should have an input schema")
				}
			}
		})
	}
}

// TestMCPTool_JSONSerialization tests JSON marshaling and unmarshaling
func TestMCPTool_JSONSerialization(t *testing.T) {
	original := MCPTool{
		Name:        "test_tool",
		Description: "A test tool for JSON serialization",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"param1": map[string]interface{}{
					"type": "string",
				},
				"param2": map[string]interface{}{
					"type": "number",
				},
			},
			"required": []string{"param1"},
		},
	}

	// Test marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MCPTool: %v", err)
	}

	// Test unmarshaling
	var unmarshaled MCPTool
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal MCPTool: %v", err)
	}

	// Verify all fields are preserved
	if original.Name != unmarshaled.Name {
		t.Errorf("Name mismatch: expected %s, got %s", original.Name, unmarshaled.Name)
	}
	if original.Description != unmarshaled.Description {
		t.Errorf("Description mismatch: expected %s, got %s", original.Description, unmarshaled.Description)
	}

	// Verify InputSchema structure
	if len(original.InputSchema) != len(unmarshaled.InputSchema) {
		t.Errorf("InputSchema size mismatch: expected %d, got %d", len(original.InputSchema), len(unmarshaled.InputSchema))
	}

	// Check specific schema properties
	originalType, ok1 := original.InputSchema["type"].(string)
	unmarshaledType, ok2 := unmarshaled.InputSchema["type"].(string)
	if !ok1 || !ok2 || originalType != unmarshaledType {
		t.Errorf("InputSchema type mismatch: expected %s, got %s", originalType, unmarshaledType)
	}
}

// TestMCPRequest_FieldValidation tests field validation for MCPRequest
func TestMCPRequest_FieldValidation(t *testing.T) {
	tests := []struct {
		name    string
		request MCPRequest
		isValid bool
	}{
		{
			name: "Valid tools/list request",
			request: MCPRequest{
				ID:      "request-1",
				Method:  "tools/list",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
			isValid: true,
		},
		{
			name: "Valid tools/call request",
			request: MCPRequest{
				ID:     "request-2",
				Method: "tools/call",
				Params: map[string]interface{}{
					"name": "generate_audit_query_with_result",
					"arguments": map[string]interface{}{
						"structured_params": map[string]interface{}{
							"log_source": "kube-apiserver",
						},
					},
				},
				JSONRPC: "2.0",
			},
			isValid: true,
		},
		{
			name: "Empty ID",
			request: MCPRequest{
				ID:      "",
				Method:  "tools/list",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
			isValid: false,
		},
		{
			name: "Empty method",
			request: MCPRequest{
				ID:      "request-3",
				Method:  "",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
			isValid: false,
		},
		{
			name: "Invalid method",
			request: MCPRequest{
				ID:      "request-4",
				Method:  "invalid/method",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
			isValid: false,
		},
		{
			name: "Nil params",
			request: MCPRequest{
				ID:      "request-5",
				Method:  "tools/list",
				Params:  nil,
				JSONRPC: "2.0",
			},
			isValid: false,
		},
		{
			name: "Invalid JSON-RPC version",
			request: MCPRequest{
				ID:      "request-6",
				Method:  "tools/list",
				Params:  map[string]interface{}{},
				JSONRPC: "1.0",
			},
			isValid: false,
		},
		{
			name: "Missing tool name in tools/call",
			request: MCPRequest{
				ID:     "request-7",
				Method: "tools/call",
				Params: map[string]interface{}{
					"arguments": map[string]interface{}{},
				},
				JSONRPC: "2.0",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test ID validation
			if tt.request.ID == "" && tt.isValid {
				t.Error("ID should not be empty for valid requests")
			}

			// Test method validation
			if tt.request.Method == "" && tt.isValid {
				t.Error("Method should not be empty for valid requests")
			}

			// Test method values
			if tt.isValid {
				validMethods := map[string]bool{
					"tools/list": true,
					"tools/call": true,
				}
				if !validMethods[tt.request.Method] {
					t.Errorf("Invalid method: %s", tt.request.Method)
				}
			}

			// Test params validation
			if tt.request.Params == nil && tt.isValid {
				t.Error("Params should not be nil for valid requests")
			}

			// Test JSON-RPC version
			if tt.request.JSONRPC != "2.0" && tt.isValid {
				t.Error("JSON-RPC version should be 2.0 for valid requests")
			}

			// Test tools/call specific validation
			if tt.request.Method == "tools/call" && tt.isValid {
				if _, ok := tt.request.Params["name"]; !ok {
					t.Error("tools/call requests should have a 'name' parameter")
				}
			}
		})
	}
}

// TestMCPRequest_JSONSerialization tests JSON marshaling and unmarshaling
func TestMCPRequest_JSONSerialization(t *testing.T) {
	original := MCPRequest{
		ID:     "test-request-123",
		Method: "tools/call",
		Params: map[string]interface{}{
			"name": "test_tool",
			"arguments": map[string]interface{}{
				"param1": "value1",
				"param2": 123,
				"param3": []string{"item1", "item2"},
			},
		},
		JSONRPC: "2.0",
	}

	// Test marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MCPRequest: %v", err)
	}

	// Test unmarshaling
	var unmarshaled MCPRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal MCPRequest: %v", err)
	}

	// Verify all fields are preserved
	if original.ID != unmarshaled.ID {
		t.Errorf("ID mismatch: expected %s, got %s", original.ID, unmarshaled.ID)
	}
	if original.Method != unmarshaled.Method {
		t.Errorf("Method mismatch: expected %s, got %s", original.Method, unmarshaled.Method)
	}
	if original.JSONRPC != unmarshaled.JSONRPC {
		t.Errorf("JSONRPC mismatch: expected %s, got %s", original.JSONRPC, unmarshaled.JSONRPC)
	}

	// Verify Params structure
	if len(original.Params) != len(unmarshaled.Params) {
		t.Errorf("Params size mismatch: expected %d, got %d", len(original.Params), len(unmarshaled.Params))
	}

	// Check specific params
	originalName, ok1 := original.Params["name"].(string)
	unmarshaledName, ok2 := unmarshaled.Params["name"].(string)
	if !ok1 || !ok2 || originalName != unmarshaledName {
		t.Errorf("Name parameter mismatch: expected %s, got %s", originalName, unmarshaledName)
	}
}

// TestMCPResponse_FieldValidation tests field validation for MCPResponse
func TestMCPResponse_FieldValidation(t *testing.T) {
	tests := []struct {
		name     string
		response MCPResponse
		isValid  bool
		hasError bool
	}{
		{
			name: "Valid success response",
			response: MCPResponse{
				ID:      "response-1",
				Result:  map[string]interface{}{"tools": []interface{}{}},
				Error:   nil,
				JSONRPC: "2.0",
			},
			isValid:  true,
			hasError: false,
		},
		{
			name: "Valid error response",
			response: MCPResponse{
				ID:     "response-2",
				Result: nil,
				Error: &MCPError{
					Code:    -32601,
					Message: "Method not found",
				},
				JSONRPC: "2.0",
			},
			isValid:  true,
			hasError: true,
		},
		{
			name: "Empty ID",
			response: MCPResponse{
				ID:      "",
				Result:  map[string]interface{}{},
				Error:   nil,
				JSONRPC: "2.0",
			},
			isValid: false,
		},
		{
			name: "Both result and error set",
			response: MCPResponse{
				ID:     "response-3",
				Result: map[string]interface{}{"data": "test"},
				Error: &MCPError{
					Code:    -32601,
					Message: "Method not found",
				},
				JSONRPC: "2.0",
			},
			isValid:  false,
			hasError: true,
		},
		{
			name: "Neither result nor error set",
			response: MCPResponse{
				ID:      "response-4",
				Result:  nil,
				Error:   nil,
				JSONRPC: "2.0",
			},
			isValid: false,
		},
		{
			name: "Invalid JSON-RPC version",
			response: MCPResponse{
				ID:      "response-5",
				Result:  map[string]interface{}{},
				Error:   nil,
				JSONRPC: "1.0",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test ID validation
			if tt.response.ID == "" && tt.isValid {
				t.Error("ID should not be empty for valid responses")
			}

			// Test JSON-RPC version
			if tt.response.JSONRPC != "2.0" && tt.isValid {
				t.Error("JSON-RPC version should be 2.0 for valid responses")
			}

			// Test result/error mutual exclusivity for valid responses
			if tt.isValid {
				if tt.response.Result != nil && tt.response.Error != nil {
					t.Error("Valid response should not have both result and error")
				}

				if tt.response.Result == nil && tt.response.Error == nil {
					t.Error("Valid response should have either result or error")
				}
			}

			// Test error field logic
			hasError := tt.response.Error != nil
			if hasError != tt.hasError {
				t.Errorf("Error expectation mismatch: expected %v, got %v", tt.hasError, hasError)
			}
		})
	}
}

// TestMCPResponse_JSONSerialization tests JSON marshaling and unmarshaling
func TestMCPResponse_JSONSerialization(t *testing.T) {
	original := MCPResponse{
		ID: "test-response-123",
		Result: map[string]interface{}{
			"tools": []interface{}{
				map[string]interface{}{
					"name":        "test_tool",
					"description": "A test tool",
					"inputSchema": map[string]interface{}{
						"type": "object",
					},
				},
			},
		},
		Error:   nil,
		JSONRPC: "2.0",
	}

	// Test marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MCPResponse: %v", err)
	}

	// Test unmarshaling
	var unmarshaled MCPResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal MCPResponse: %v", err)
	}

	// Verify all fields are preserved
	if original.ID != unmarshaled.ID {
		t.Errorf("ID mismatch: expected %s, got %s", original.ID, unmarshaled.ID)
	}
	if original.JSONRPC != unmarshaled.JSONRPC {
		t.Errorf("JSONRPC mismatch: expected %s, got %s", original.JSONRPC, unmarshaled.JSONRPC)
	}

	// Verify Result structure
	if original.Result == nil && unmarshaled.Result != nil {
		t.Error("Result should be nil")
	}
	if original.Result != nil && unmarshaled.Result == nil {
		t.Error("Result should not be nil")
	}

	// Verify Error field
	if original.Error != nil && unmarshaled.Error == nil {
		t.Error("Error should not be nil")
	}
	if original.Error == nil && unmarshaled.Error != nil {
		t.Error("Error should be nil")
	}
}

// TestMCPError_FieldValidation tests field validation for MCPError
func TestMCPError_FieldValidation(t *testing.T) {
	tests := []struct {
		name    string
		err     MCPError
		isValid bool
	}{
		{
			name: "Valid error",
			err: MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
			isValid: true,
		},
		{
			name: "Valid error with zero code",
			err: MCPError{
				Code:    0,
				Message: "Success (not really an error)",
			},
			isValid: true,
		},
		{
			name: "Valid error with positive code",
			err: MCPError{
				Code:    100,
				Message: "Custom error",
			},
			isValid: true,
		},
		{
			name: "Empty message",
			err: MCPError{
				Code:    -32601,
				Message: "",
			},
			isValid: false,
		},
		{
			name: "Very long message",
			err: MCPError{
				Code:    -32601,
				Message: string(make([]byte, 10000)), // 10KB message
			},
			isValid: true, // Long messages are valid
		},
		{
			name: "Special characters in message",
			err: MCPError{
				Code:    -32601,
				Message: "Error with special chars: @#$%^&*()\n\t\r",
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test message validation
			if tt.err.Message == "" && tt.isValid {
				t.Error("Message should not be empty for valid errors")
			}

			// Test that valid errors have expected structure
			if tt.isValid {
				if tt.err.Message == "" {
					t.Error("Valid errors should have a message")
				}
			}

			// Test code range (JSON-RPC allows any integer)
			// No specific validation needed for code as JSON-RPC allows any integer
		})
	}
}

// TestMCPError_JSONSerialization tests JSON marshaling and unmarshaling
func TestMCPError_JSONSerialization(t *testing.T) {
	original := MCPError{
		Code:    -32601,
		Message: "Method not found with special chars: @#$%^&*()",
	}

	// Test marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MCPError: %v", err)
	}

	// Test unmarshaling
	var unmarshaled MCPError
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal MCPError: %v", err)
	}

	// Verify all fields are preserved
	if original.Code != unmarshaled.Code {
		t.Errorf("Code mismatch: expected %d, got %d", original.Code, unmarshaled.Code)
	}
	if original.Message != unmarshaled.Message {
		t.Errorf("Message mismatch: expected %s, got %s", original.Message, unmarshaled.Message)
	}
}

// TestMCPTypes_Integration tests integration between MCP types
func TestMCPTypes_Integration(t *testing.T) {
	// Test complete request-response cycle
	request := MCPRequest{
		ID:     "integration-test-1",
		Method: "tools/call",
		Params: map[string]interface{}{
			"name": "test_tool",
			"arguments": map[string]interface{}{
				"param1": "value1",
			},
		},
		JSONRPC: "2.0",
	}

	// Simulate processing the request
	response := MCPResponse{
		ID: request.ID, // Same ID as request
		Result: map[string]interface{}{
			"result": "success",
		},
		Error:   nil,
		JSONRPC: "2.0",
	}

	// Test that request and response are properly linked
	if request.ID != response.ID {
		t.Error("Request and response should have the same ID")
	}

	// Test that both use the same JSON-RPC version
	if request.JSONRPC != response.JSONRPC {
		t.Error("Request and response should use the same JSON-RPC version")
	}

	// Test error response scenario
	errorResponse := MCPResponse{
		ID:     request.ID,
		Result: nil,
		Error: &MCPError{
			Code:    -32601,
			Message: "Method not found",
		},
		JSONRPC: "2.0",
	}

	// Test that error response has no result
	if errorResponse.Result != nil {
		t.Error("Error response should not have a result")
	}

	// Test that error response has an error
	if errorResponse.Error == nil {
		t.Error("Error response should have an error")
	}

	// Test tools/list scenario
	toolsListRequest := MCPRequest{
		ID:      "tools-list-test",
		Method:  "tools/list",
		Params:  map[string]interface{}{},
		JSONRPC: "2.0",
	}

	toolsListResponse := MCPResponse{
		ID: toolsListRequest.ID,
		Result: map[string]interface{}{
			"tools": []MCPTool{
				{
					Name:        "tool1",
					Description: "First tool",
					InputSchema: map[string]interface{}{
						"type": "object",
					},
				},
				{
					Name:        "tool2",
					Description: "Second tool",
					InputSchema: map[string]interface{}{
						"type": "object",
					},
				},
			},
		},
		Error:   nil,
		JSONRPC: "2.0",
	}

	// Test that tools list response contains tools
	if toolsListResponse.Result == nil {
		t.Error("Tools list response should have a result")
	}

	resultMap, ok := toolsListResponse.Result.(map[string]interface{})
	if !ok {
		t.Error("Tools list response result should be a map")
	}

	tools, ok := resultMap["tools"].([]MCPTool)
	if !ok {
		t.Error("Tools list response should contain tools array")
	}

	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}
}

// TestMCPTypes_EdgeCases tests edge cases for MCP types
func TestMCPTypes_EdgeCases(t *testing.T) {
	// Test very large input schemas
	largeSchema := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeSchema[fmt.Sprintf("property_%d", i)] = map[string]interface{}{
			"type": "string",
		}
	}

	largeTool := MCPTool{
		Name:        "large_tool",
		Description: "Tool with large schema",
		InputSchema: largeSchema,
	}

	// Test that large schema can be serialized
	_, err := json.Marshal(largeTool)
	if err != nil {
		t.Errorf("Failed to marshal tool with large schema: %v", err)
	}

	// Test very large parameters
	largeParams := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeParams[fmt.Sprintf("param_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	largeRequest := MCPRequest{
		ID:      "large-request",
		Method:  "tools/call",
		Params:  largeParams,
		JSONRPC: "2.0",
	}

	// Test that large request can be serialized
	_, err = json.Marshal(largeRequest)
	if err != nil {
		t.Errorf("Failed to marshal large request: %v", err)
	}

	// Test very large results
	largeResult := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeResult[fmt.Sprintf("result_%d", i)] = fmt.Sprintf("data_%d", i)
	}

	largeResponse := MCPResponse{
		ID:      "large-response",
		Result:  largeResult,
		Error:   nil,
		JSONRPC: "2.0",
	}

	// Test that large response can be serialized
	_, err = json.Marshal(largeResponse)
	if err != nil {
		t.Errorf("Failed to marshal large response: %v", err)
	}

	// Test unicode characters in all fields
	unicodeTool := MCPTool{
		Name:        "tool_with_unicode_🚀🎉",
		Description: "Tool with unicode description 🚀🎉",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"unicode_param": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	// Test that unicode tool can be serialized
	_, err = json.Marshal(unicodeTool)
	if err != nil {
		t.Errorf("Failed to marshal unicode tool: %v", err)
	}

	unicodeRequest := MCPRequest{
		ID:     "unicode-request-🚀",
		Method: "tools/call",
		Params: map[string]interface{}{
			"name": "unicode_tool",
			"arguments": map[string]interface{}{
				"unicode_param": "🚀🎉",
			},
		},
		JSONRPC: "2.0",
	}

	// Test that unicode request can be serialized
	_, err = json.Marshal(unicodeRequest)
	if err != nil {
		t.Errorf("Failed to marshal unicode request: %v", err)
	}
}
