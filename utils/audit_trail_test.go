package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"audit-query-mcp-server/types"
)

// TestNewAuditTrail tests the constructor function
func TestNewAuditTrail(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Valid file path",
			filePath: "testdata/audit_trail.json",
			wantErr:  false,
		},
		{
			name:     "Nested directory path",
			filePath: "testdata/nested/deep/audit_trail.json",
			wantErr:  false,
		},
		{
			name:     "Empty file path",
			filePath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up after test
			defer func() {
				if tt.filePath != "" {
					os.RemoveAll(filepath.Dir(tt.filePath))
				}
			}()

			trail, err := NewAuditTrail(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuditTrail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if trail == nil {
					t.Error("NewAuditTrail() returned nil when no error expected")
					return
				}

				// Verify file was created
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Errorf("Audit trail file was not created: %s", tt.filePath)
				}

				// Verify file can be written to
				if trail.file == nil {
					t.Error("Audit trail file is nil")
				}

				// Clean up
				trail.Close()
			}
		})
	}
}

// TestAuditTrail_LogQuery tests the core logging functionality
func TestAuditTrail_LogQuery(t *testing.T) {
	filePath := "testdata/log_query_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	// Test basic logging
	entry := AuditTrailEntry{
		QueryID:    "test_query_123",
		UserID:     "test_user",
		Action:     "test_action",
		Parameters: map[string]interface{}{"key": "value"},
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	}

	err = trail.LogQuery(entry)
	if err != nil {
		t.Errorf("LogQuery() failed: %v", err)
	}

	// Verify entry was written
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "test_query_123") {
		t.Error("Log entry was not written to file")
	}

	// Test auto-timestamp generation
	entryNoTimestamp := AuditTrailEntry{
		QueryID: "test_query_no_timestamp",
		Action:  "test_action",
	}

	err = trail.LogQuery(entryNoTimestamp)
	if err != nil {
		t.Errorf("LogQuery() with auto-timestamp failed: %v", err)
	}

	// Verify timestamp was added
	content, err = os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "test_query_no_timestamp") {
		t.Error("Entry with auto-timestamp was not written")
	}
}

// TestAuditTrail_LogQueryGeneration tests query generation logging
func TestAuditTrail_LogQueryGeneration(t *testing.T) {
	filePath := "testdata/log_generation_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"pods", "delete"},
		Timeframe: "today",
		Username:  "admin",
		Resource:  "pods",
		Verb:      "delete",
		Namespace: "default",
		Exclude:   []string{"system:"},
	}

	result := &types.AuditResult{
		QueryID:       "test_query_gen_123",
		Timestamp:     time.Now().Format(time.RFC3339),
		Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
		RawOutput:     "test output",
		ParsedData:    []map[string]interface{}{{"test": "data"}},
		Summary:       "Found 1 audit entries",
		Error:         "",
		ExecutionTime: 100,
	}

	err = trail.LogQueryGeneration("test_query_gen_123", params, result, "test_user", "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("LogQueryGeneration() failed: %v", err)
	}

	// Verify entry was written with correct action
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "query_generation") {
		t.Error("Query generation entry was not written with correct action")
	}

	if !strings.Contains(string(content), "kube-apiserver") {
		t.Error("Query parameters were not logged correctly")
	}
}

// TestAuditTrail_LogQueryExecution tests query execution logging
func TestAuditTrail_LogQueryExecution(t *testing.T) {
	filePath := "testdata/log_execution_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	command := "oc adm node-logs --role=master --path=kube-apiserver/audit.log | head -10"
	result := &types.AuditResult{
		QueryID:       "test_query_exec_123",
		Timestamp:     time.Now().Format(time.RFC3339),
		Command:       command,
		RawOutput:     "execution output",
		ParsedData:    []map[string]interface{}{},
		Summary:       "Executed successfully",
		Error:         "",
		ExecutionTime: 150,
	}

	err = trail.LogQueryExecution("test_query_exec_123", command, result, "test_user", "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("LogQueryExecution() failed: %v", err)
	}

	// Verify entry was written
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "query_execution") {
		t.Error("Query execution entry was not written with correct action")
	}

	if !strings.Contains(string(content), command) {
		t.Error("Command was not logged correctly")
	}
}

// TestAuditTrail_LogQueryParsing tests query parsing logging
func TestAuditTrail_LogQueryParsing(t *testing.T) {
	filePath := "testdata/log_parsing_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	queryContext := map[string]interface{}{
		"log_source": "kube-apiserver",
		"timeframe":  "today",
		"username":   "admin",
		"resource":   "pods",
		"verb":       "delete",
		"namespace":  "default",
	}

	result := &types.AuditResult{
		QueryID:       "test_query_parse_123",
		Timestamp:     time.Now().Format(time.RFC3339),
		Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
		RawOutput:     "raw output to parse",
		ParsedData:    []map[string]interface{}{{"parsed": "data"}},
		Summary:       "Parsed 1 entry",
		Error:         "",
		ExecutionTime: 75,
	}

	err = trail.LogQueryParsing("test_query_parse_123", queryContext, result, "test_user", "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("LogQueryParsing() failed: %v", err)
	}

	// Verify entry was written
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "query_parsing") {
		t.Error("Query parsing entry was not written with correct action")
	}

	if !strings.Contains(string(content), "kube-apiserver") {
		t.Error("Query context was not logged correctly")
	}
}

// TestAuditTrail_LogCompleteQuery tests complete query logging
func TestAuditTrail_LogCompleteQuery(t *testing.T) {
	filePath := "testdata/log_complete_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"pods", "delete"},
		Timeframe: "today",
		Username:  "admin",
		Resource:  "pods",
		Verb:      "delete",
		Namespace: "default",
		Exclude:   []string{"system:"},
	}

	result := &types.AuditResult{
		QueryID:       "test_query_complete_123",
		Timestamp:     time.Now().Format(time.RFC3339),
		Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
		RawOutput:     "complete pipeline output",
		ParsedData:    []map[string]interface{}{{"complete": "data"}},
		Summary:       "Complete pipeline executed successfully",
		Error:         "",
		ExecutionTime: 300,
	}

	err = trail.LogCompleteQuery("test_query_complete_123", params, result, "test_user", "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("LogCompleteQuery() failed: %v", err)
	}

	// Verify entry was written
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "complete_query") {
		t.Error("Complete query entry was not written with correct action")
	}

	if !strings.Contains(string(content), "kube-apiserver") {
		t.Error("Complete query parameters were not logged correctly")
	}
}

// TestAuditTrail_LogCacheAccess tests cache access logging
func TestAuditTrail_LogCacheAccess(t *testing.T) {
	filePath := "testdata/log_cache_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	// Test cache hit
	err = trail.LogCacheAccess("test_query_cache_123", "hit", "test_user", "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("LogCacheAccess() for hit failed: %v", err)
	}

	// Test cache miss
	err = trail.LogCacheAccess("test_query_cache_456", "miss", "test_user", "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("LogCacheAccess() for miss failed: %v", err)
	}

	// Verify entries were written
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "cache_hit") {
		t.Error("Cache hit entry was not written with correct action")
	}

	if !strings.Contains(string(content), "cache_miss") {
		t.Error("Cache miss entry was not written with correct action")
	}
}

// TestAuditTrail_ThreadSafety tests concurrent access
func TestAuditTrail_ThreadSafety(t *testing.T) {
	filePath := "testdata/thread_safety_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	// Test concurrent writes
	var wg sync.WaitGroup
	numGoroutines := 10
	entriesPerGoroutine := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < entriesPerGoroutine; j++ {
				entry := AuditTrailEntry{
					QueryID: fmt.Sprintf("query_%d_%d", goroutineID, j),
					Action:  "concurrent_test",
					Parameters: map[string]interface{}{
						"goroutine_id": goroutineID,
						"entry_id":     j,
					},
				}
				if err := trail.LogQuery(entry); err != nil {
					t.Errorf("Concurrent LogQuery failed: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify all entries were written
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	expectedLines := numGoroutines * entriesPerGoroutine
	if len(lines) != expectedLines {
		t.Errorf("Expected %d entries, got %d", expectedLines, len(lines))
	}
}

// TestAuditTrail_ErrorHandling tests error scenarios
func TestAuditTrail_ErrorHandling(t *testing.T) {
	// Test with invalid file path (read-only directory)
	filePath := "/tmp/readonly/audit_trail.json"

	// Create a read-only directory
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0444); err == nil {
		defer os.RemoveAll(dir)

		_, err := NewAuditTrail(filePath)
		if err == nil {
			t.Error("Expected error when creating audit trail in read-only directory")
		}
	}
}

// TestAuditTrail_Close tests the close functionality
func TestAuditTrail_Close(t *testing.T) {
	filePath := "testdata/close_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}

	// Write some entries
	entry := AuditTrailEntry{
		QueryID: "test_close_123",
		Action:  "close_test",
	}
	err = trail.LogQuery(entry)
	if err != nil {
		t.Errorf("LogQuery before close failed: %v", err)
	}

	// Close the trail
	err = trail.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Verify file is closed by trying to write again (should fail)
	err = trail.LogQuery(entry)
	if err == nil {
		t.Error("Expected error when writing to closed audit trail")
	}
}

// TestParamsToMap tests the parameter conversion utility
func TestParamsToMap(t *testing.T) {
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Patterns:  []string{"pods", "delete"},
		Timeframe: "today",
		Exclude:   []string{"system:"},
		Username:  "admin",
		Resource:  "pods",
		Verb:      "delete",
		Namespace: "default",
	}

	result := paramsToMap(params)

	// Verify all fields are present
	expectedFields := []string{"log_source", "patterns", "timeframe", "exclude", "username", "resource", "verb", "namespace"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Expected field %s not found in paramsToMap result", field)
		}
	}

	// Verify specific values
	if result["log_source"] != "kube-apiserver" {
		t.Errorf("Expected log_source to be 'kube-apiserver', got %v", result["log_source"])
	}

	if patterns, ok := result["patterns"].([]string); !ok || len(patterns) != 2 {
		t.Errorf("Expected patterns to be []string with 2 elements, got %v", result["patterns"])
	}
}

// TestAuditTrailEntry_JSONMarshaling tests JSON serialization
func TestAuditTrailEntry_JSONMarshaling(t *testing.T) {
	entry := AuditTrailEntry{
		Timestamp:     "2023-01-01T12:00:00Z",
		QueryID:       "test_json_123",
		UserID:        "test_user",
		Action:        "json_test",
		Parameters:    map[string]interface{}{"key": "value"},
		IPAddress:     "127.0.0.1",
		UserAgent:     "test-agent",
		ExecutionTime: 100,
	}

	// Test marshaling
	data, err := json.Marshal(entry)
	if err != nil {
		t.Errorf("JSON marshaling failed: %v", err)
	}

	// Test unmarshaling
	var unmarshaled AuditTrailEntry
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("JSON unmarshaling failed: %v", err)
	}

	// Verify fields
	if unmarshaled.QueryID != entry.QueryID {
		t.Errorf("QueryID mismatch: expected %s, got %s", entry.QueryID, unmarshaled.QueryID)
	}

	if unmarshaled.Action != entry.Action {
		t.Errorf("Action mismatch: expected %s, got %s", entry.Action, unmarshaled.Action)
	}
}

// TestAuditTrail_ErrorLogging tests logging with errors
func TestAuditTrail_ErrorLogging(t *testing.T) {
	filePath := "testdata/error_logging_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		t.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	// Test logging with error result
	result := &types.AuditResult{
		QueryID:       "test_error_123",
		Timestamp:     time.Now().Format(time.RFC3339),
		Command:       "invalid command",
		RawOutput:     "",
		ParsedData:    []map[string]interface{}{},
		Summary:       "",
		Error:         "command execution failed",
		ExecutionTime: 50,
	}

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
	}

	err = trail.LogQueryGeneration("test_error_123", params, result, "test_user", "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("LogQueryGeneration with error failed: %v", err)
	}

	// Verify error was logged
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	if !strings.Contains(string(content), "command execution failed") {
		t.Error("Error message was not logged correctly")
	}
}

// BenchmarkAuditTrail_LogQuery benchmarks the logging performance
func BenchmarkAuditTrail_LogQuery(b *testing.B) {
	filePath := "testdata/benchmark_test.json"
	defer os.RemoveAll(filepath.Dir(filePath))

	trail, err := NewAuditTrail(filePath)
	if err != nil {
		b.Fatalf("Failed to create audit trail: %v", err)
	}
	defer trail.Close()

	entry := AuditTrailEntry{
		QueryID:    "benchmark_query",
		Action:     "benchmark_test",
		Parameters: map[string]interface{}{"benchmark": true},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry.QueryID = fmt.Sprintf("benchmark_query_%d", i)
		if err := trail.LogQuery(entry); err != nil {
			b.Errorf("LogQuery failed: %v", err)
		}
	}
}
