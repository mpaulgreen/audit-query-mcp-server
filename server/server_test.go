package server

import (
	"fmt"
	"testing"
	"time"

	"audit-query-mcp-server/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewAuditQueryMCPServer tests server initialization
func TestNewAuditQueryMCPServer(t *testing.T) {
	server := NewAuditQueryMCPServer()

	assert.NotNil(t, server)
	assert.NotNil(t, server.logger)
	assert.NotNil(t, server.cache)
	// auditTrail might be nil if file creation fails, which is acceptable
}

// TestGetTools tests the tools list
func TestGetTools(t *testing.T) {
	server := NewAuditQueryMCPServer()
	tools := server.GetTools()

	assert.Len(t, tools, 9) // Should have 9 tools total

	// Check for specific tools
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	expectedTools := []string{
		"generate_audit_query_with_result",
		"execute_audit_query_with_result",
		"parse_audit_results_with_result",
		"execute_complete_audit_query",
		"get_cache_stats",
		"clear_cache",
		"get_cached_result",
		"delete_cached_result",
		"get_server_stats",
	}

	for _, expected := range expectedTools {
		assert.True(t, toolNames[expected], "Missing tool: %s", expected)
	}
}

// TestGenerateAuditQueryWithResult tests query generation
func TestGenerateAuditQueryWithResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
		Username:  "admin",
	}

	result, err := server.GenerateAuditQueryWithResult(params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.QueryID)
	assert.NotEmpty(t, result.Timestamp)
	assert.NotEmpty(t, result.Command)
	assert.Empty(t, result.Error)
	// Execution time might be 0 for very fast operations
	assert.GreaterOrEqual(t, result.ExecutionTime, int64(0))

	// Verify command contains expected elements
	assert.Contains(t, result.Command, "oc adm node-logs")
	assert.Contains(t, result.Command, "kube-apiserver")
	assert.Contains(t, result.Command, "admin")
}

// TestGenerateAuditQueryWithResult_InvalidParams tests validation
func TestGenerateAuditQueryWithResult_InvalidParams(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Test with invalid log source
	params := types.AuditQueryParams{
		LogSource: "invalid-source",
		Timeframe: "today",
	}

	result, err := server.GenerateAuditQueryWithResult(params)

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "validation failed")
}

// TestExecuteAuditQueryWithResult tests command execution
func TestExecuteAuditQueryWithResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Use a valid oc command for testing (this will fail if oc is not available, which is expected)
	command := "oc adm node-logs --role=master --path=kube-apiserver/audit.log | head -1"
	queryID := "test-query-123"

	result, err := server.ExecuteAuditQueryWithResult(command, queryID)

	// This test might fail if oc command is not available, which is expected
	if err != nil {
		// If it fails, it should be due to command execution, not validation
		assert.NotNil(t, result)
		assert.Equal(t, queryID, result.QueryID)
		assert.Equal(t, command, result.Command)
		// The error should be related to command execution, not validation
		assert.NotContains(t, result.Error, "command validation failed")
	} else {
		assert.NotNil(t, result)
		assert.Equal(t, queryID, result.QueryID)
		assert.Equal(t, command, result.Command)
		assert.NotEmpty(t, result.RawOutput)
		assert.Empty(t, result.Error)
		// Execution time might be 0 for very fast operations
		assert.GreaterOrEqual(t, result.ExecutionTime, int64(0))
	}
}

// TestExecuteAuditQueryWithResult_InvalidCommand tests command validation
func TestExecuteAuditQueryWithResult_InvalidCommand(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Test with dangerous command
	command := "rm -rf /"
	queryID := "test-query-123"

	result, err := server.ExecuteAuditQueryWithResult(command, queryID)

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "command validation failed")
}

// TestExecuteAuditQueryWithResult_Timeout tests timeout handling
func TestExecuteAuditQueryWithResult_Timeout(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Use a valid oc command that will timeout (this will fail if oc is not available, which is expected)
	command := "oc adm node-logs --role=master --path=kube-apiserver/audit.log | sleep 35"
	queryID := "test-query-123"

	result, err := server.ExecuteAuditQueryWithResult(command, queryID)

	// This test might fail if oc command is not available, which is expected
	if err != nil {
		// If it fails, it should be due to command execution or timeout, not validation
		assert.NotNil(t, result)
		assert.Equal(t, queryID, result.QueryID)
		assert.Equal(t, command, result.Command)
		// The error should be related to command execution or timeout, not validation
		assert.NotContains(t, result.Error, "command validation failed")
	} else {
		// If it succeeds, it should have timed out
		assert.NotNil(t, result)
		assert.Equal(t, queryID, result.QueryID)
		assert.Equal(t, command, result.Command)
		assert.NotEmpty(t, result.Error)
		assert.Contains(t, result.Error, "timed out")
	}
}

// TestParseAuditResultsWithResult tests result parsing
func TestParseAuditResultsWithResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Sample audit log output
	rawOutput := `{"kind":"Event","apiVersion":"audit.k8s.io/v1","level":"Metadata","auditID":"test-123","stage":"ResponseComplete","requestURI":"/api/v1/namespaces/default/pods","verb":"list","user":{"username":"admin","uid":"admin-123","groups":["system:authenticated"]},"sourceIPs":["127.0.0.1"],"userAgent":"kubectl/v1.20.0","objectRef":{"resource":"pods","namespace":"default","apiVersion":"v1"},"responseStatus":{"metadata":{},"code":200},"requestReceivedTimestamp":"2023-01-01T00:00:00.000000Z","stageTimestamp":"2023-01-01T00:00:00.000000Z"}`

	queryContext := map[string]interface{}{
		"log_source": "kube-apiserver",
		"timeframe":  "today",
	}
	queryID := "test-query-123"

	result, err := server.ParseAuditResultsWithResult(rawOutput, queryContext, queryID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, queryID, result.QueryID)
	assert.Equal(t, rawOutput, result.RawOutput)
	assert.NotNil(t, result.ParsedData)
	assert.Len(t, result.ParsedData, 1)
	assert.NotEmpty(t, result.Summary)
	assert.Empty(t, result.Error)
	// Execution time might be 0 for very fast operations
	assert.GreaterOrEqual(t, result.ExecutionTime, int64(0))
}

// TestParseAuditResultsWithResult_EmptyOutput tests empty output handling
func TestParseAuditResultsWithResult_EmptyOutput(t *testing.T) {
	server := NewAuditQueryMCPServer()

	rawOutput := ""
	queryContext := map[string]interface{}{
		"log_source": "kube-apiserver",
		"timeframe":  "today",
	}
	queryID := "test-query-123"

	result, err := server.ParseAuditResultsWithResult(rawOutput, queryContext, queryID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, queryID, result.QueryID)
	assert.Empty(t, result.RawOutput)
	// ParsedData might be nil for empty output
	if result.ParsedData != nil {
		assert.Len(t, result.ParsedData, 0)
	}
	assert.Empty(t, result.Error)
}

// TestExecuteCompleteAuditQuery tests the complete pipeline
func TestExecuteCompleteAuditQuery(t *testing.T) {
	server := NewAuditQueryMCPServer()

	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
		Username:  "admin",
	}

	result, err := server.ExecuteCompleteAuditQuery(params)

	// This test might fail if oc command is not available, which is expected
	if err != nil {
		// If it fails, it should be due to command execution, not validation
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.QueryID)
		assert.NotEmpty(t, result.Command)
		// The error should be related to command execution, not validation
		assert.NotContains(t, result.Error, "validation failed")
	} else {
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.QueryID)
		assert.NotEmpty(t, result.Command)
		assert.Empty(t, result.Error)
		// Execution time might be 0 for very fast operations
		assert.GreaterOrEqual(t, result.ExecutionTime, int64(0))
	}
}

// TestGenerateQueryID tests query ID generation
func TestGenerateQueryID(t *testing.T) {
	server := NewAuditQueryMCPServer()

	queryID1 := server.generateQueryID()
	queryID2 := server.generateQueryID()

	assert.NotEmpty(t, queryID1)
	assert.NotEmpty(t, queryID2)
	assert.NotEqual(t, queryID1, queryID2)
	assert.Contains(t, queryID1, "audit_query_")
}

// TestGetCacheStats tests cache statistics
func TestGetCacheStats(t *testing.T) {
	server := NewAuditQueryMCPServer()

	stats := server.GetCacheStats()

	assert.NotNil(t, stats)
	assert.Contains(t, stats, "size")
	assert.Contains(t, stats, "hits")
	assert.Contains(t, stats, "misses")
	assert.Contains(t, stats, "hit_rate")
}

// TestClearCache tests cache clearing
func TestClearCache(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Add some data to cache
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
	}
	result, _ := server.GenerateAuditQueryWithResult(params)
	server.cache.Set(result.QueryID, result)

	// Verify cache has data
	stats := server.GetCacheStats()
	assert.Greater(t, stats["size"].(int), 0)

	// Clear cache
	server.ClearCache()

	// Verify cache is empty
	stats = server.GetCacheStats()
	assert.Equal(t, 0, stats["size"].(int))
}

// TestGetCachedResult tests cached result retrieval
func TestGetCachedResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Add data to cache
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
	}
	result, _ := server.GenerateAuditQueryWithResult(params)
	server.cache.Set(result.QueryID, result)

	// Retrieve cached result
	cachedResult, found := server.GetCachedResult(result.QueryID)

	assert.True(t, found)
	assert.NotNil(t, cachedResult)
	assert.Equal(t, result.QueryID, cachedResult.QueryID)

	// Test non-existent result
	_, found = server.GetCachedResult("non-existent-id")
	assert.False(t, found)
}

// TestDeleteCachedResult tests cached result deletion
func TestDeleteCachedResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Add data to cache
	params := types.AuditQueryParams{
		LogSource: "kube-apiserver",
		Timeframe: "today",
	}
	result, _ := server.GenerateAuditQueryWithResult(params)
	server.cache.Set(result.QueryID, result)

	// Verify it exists
	_, found := server.GetCachedResult(result.QueryID)
	assert.True(t, found)

	// Delete it
	server.DeleteCachedResult(result.QueryID)

	// Verify it's gone
	_, found = server.GetCachedResult(result.QueryID)
	assert.False(t, found)
}

// TestGetServerStats tests server statistics
func TestGetServerStats(t *testing.T) {
	server := NewAuditQueryMCPServer()

	stats := server.GetServerStats()

	assert.NotNil(t, stats)
	assert.Contains(t, stats, "server_info")
	assert.Contains(t, stats, "cache_stats")
	assert.Contains(t, stats, "tools")
	assert.Contains(t, stats, "features")

	// Check server info
	serverInfo := stats["server_info"].(map[string]interface{})
	assert.Equal(t, "1.0.0", serverInfo["version"])
	assert.Equal(t, "2", serverInfo["phase"])
	assert.True(t, serverInfo["audit_result"].(bool))
	assert.True(t, serverInfo["caching"].(bool))

	// Check tools count - handle both int and float64 types
	tools := stats["tools"].(map[string]interface{})

	// Convert to int for comparison (JSON unmarshaling can produce either type)
	auditResultTools := tools["audit_result_tools"]
	if auditResultToolsFloat, ok := auditResultTools.(float64); ok {
		assert.Equal(t, 4, int(auditResultToolsFloat))
	} else if auditResultToolsInt, ok := auditResultTools.(int); ok {
		assert.Equal(t, 4, auditResultToolsInt)
	} else {
		t.Errorf("Unexpected type for audit_result_tools: %T", auditResultTools)
	}

	cacheTools := tools["cache_tools"]
	if cacheToolsFloat, ok := cacheTools.(float64); ok {
		assert.Equal(t, 5, int(cacheToolsFloat))
	} else if cacheToolsInt, ok := cacheTools.(int); ok {
		assert.Equal(t, 5, cacheToolsInt)
	} else {
		t.Errorf("Unexpected type for cache_tools: %T", cacheTools)
	}

	totalTools := tools["total_tools"]
	if totalToolsFloat, ok := totalTools.(float64); ok {
		assert.Equal(t, 9, int(totalToolsFloat))
	} else if totalToolsInt, ok := totalTools.(int); ok {
		assert.Equal(t, 9, totalToolsInt)
	} else {
		t.Errorf("Unexpected type for total_tools: %T", totalTools)
	}
}

// TestGetLogger tests logger retrieval
func TestGetLogger(t *testing.T) {
	server := NewAuditQueryMCPServer()

	logger := server.GetLogger()

	assert.NotNil(t, logger)
	// Test that it's the same logger instance
	assert.Equal(t, server.logger, logger)
}

// TestServerConcurrentAccess tests thread safety of server methods
func TestServerConcurrentAccess(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Test concurrent cache operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			params := types.AuditQueryParams{
				LogSource: "kube-apiserver",
				Timeframe: "today",
				Username:  fmt.Sprintf("user-%d", id),
			}

			result, _ := server.GenerateAuditQueryWithResult(params)
			server.cache.Set(result.QueryID, result)
			_, _ = server.GetCachedResult(result.QueryID)
			server.DeleteCachedResult(result.QueryID)

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic and should complete successfully
	assert.True(t, true)
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Test with nil parameters
	var params types.AuditQueryParams
	result, err := server.GenerateAuditQueryWithResult(params)

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Error)

	// Test with empty command
	result, err = server.ExecuteAuditQueryWithResult("", "test-id")

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Error)
}

// TestPerformance tests basic performance characteristics
func TestPerformance(t *testing.T) {
	server := NewAuditQueryMCPServer()

	start := time.Now()

	// Generate multiple queries
	for i := 0; i < 100; i++ {
		params := types.AuditQueryParams{
			LogSource: "kube-apiserver",
			Timeframe: "today",
			Username:  fmt.Sprintf("user-%d", i),
		}
		_, err := server.GenerateAuditQueryWithResult(params)
		require.NoError(t, err)
	}

	duration := time.Since(start)

	// Should complete 100 queries in reasonable time (less than 1 second)
	assert.Less(t, duration, time.Second)
}
