package server

import (
	"fmt"
	"testing"

	"audit-query-mcp-server/types"

	"github.com/stretchr/testify/assert"
)

// TestHandleMCPRequest_ValidMethods tests the main request handler with valid methods
func TestHandleMCPRequest_ValidMethods(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name     string
		request  types.MCPRequest
		expected string
	}{
		{
			name: "Tools list request",
			request: types.MCPRequest{
				ID:      "test-1",
				Method:  "tools/list",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
			expected: "tools/list",
		},
		{
			name: "Tool call request",
			request: types.MCPRequest{
				ID:     "test-2",
				Method: "tools/call",
				Params: map[string]interface{}{
					"name":      "get_cache_stats",
					"arguments": map[string]interface{}{},
				},
				JSONRPC: "2.0",
			},
			expected: "tools/call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.HandleMCPRequest(tt.request)

			assert.Equal(t, tt.request.ID, response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Nil(t, response.Error)
			assert.NotNil(t, response.Result)
		})
	}
}

// TestHandleMCPRequest_InvalidMethods tests the main request handler with invalid methods
func TestHandleMCPRequest_InvalidMethods(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name           string
		request        types.MCPRequest
		expectedCode   int
		expectedMethod string
	}{
		{
			name: "Invalid method",
			request: types.MCPRequest{
				ID:      "test-1",
				Method:  "invalid/method",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
			expectedCode:   -32601,
			expectedMethod: "Method not found",
		},
		{
			name: "Empty method",
			request: types.MCPRequest{
				ID:      "test-2",
				Method:  "",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
			expectedCode:   -32601,
			expectedMethod: "Method not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.HandleMCPRequest(tt.request)

			assert.Equal(t, tt.request.ID, response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)
			assert.NotNil(t, response.Error)
			assert.Equal(t, tt.expectedCode, response.Error.Code)
			assert.Equal(t, tt.expectedMethod, response.Error.Message)
			assert.Nil(t, response.Result)
		})
	}
}

// TestHandleListTools tests the tools listing handler
func TestHandleListTools(t *testing.T) {
	server := NewAuditQueryMCPServer()

	request := types.MCPRequest{
		ID:      "test-tools-list",
		Method:  "tools/list",
		Params:  map[string]interface{}{},
		JSONRPC: "2.0",
	}

	response := server.handleListTools(request)

	assert.Equal(t, request.ID, response.ID)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Result)

	// Check that result contains tools
	result, ok := response.Result.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, result, "tools")

	tools, ok := result["tools"].([]types.MCPTool)
	assert.True(t, ok)
	assert.Greater(t, len(tools), 0)

	// Check for expected tools
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

	for _, expectedTool := range expectedTools {
		assert.True(t, toolNames[expectedTool], "Expected tool %s not found", expectedTool)
	}
}

// TestHandleToolCall_ValidTools tests tool call handler with valid tools
func TestHandleToolCall_ValidTools(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name     string
		toolName string
		params   map[string]interface{}
	}{
		{
			name:     "Get cache stats",
			toolName: "get_cache_stats",
			params:   map[string]interface{}{},
		},
		{
			name:     "Clear cache",
			toolName: "clear_cache",
			params:   map[string]interface{}{},
		},
		{
			name:     "Get server stats",
			toolName: "get_server_stats",
			params:   map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := types.MCPRequest{
				ID:     "test-tool-call",
				Method: "tools/call",
				Params: map[string]interface{}{
					"name":      tt.toolName,
					"arguments": tt.params,
				},
				JSONRPC: "2.0",
			}

			response := server.handleToolCall(request)

			assert.Equal(t, request.ID, response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Nil(t, response.Error)
			assert.NotNil(t, response.Result)
		})
	}
}

// TestHandleToolCall_InvalidParams tests tool call handler with invalid parameters
func TestHandleToolCall_InvalidParams(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name          string
		params        map[string]interface{}
		expectedCode  int
		expectedError string
	}{
		{
			name: "Missing arguments",
			params: map[string]interface{}{
				"name": "get_cache_stats",
			},
			expectedCode:  -32602,
			expectedError: "Invalid params",
		},
		{
			name: "Missing tool name",
			params: map[string]interface{}{
				"arguments": map[string]interface{}{},
			},
			expectedCode:  -32602,
			expectedError: "Tool name required",
		},
		{
			name: "Invalid arguments type",
			params: map[string]interface{}{
				"name":      "get_cache_stats",
				"arguments": "not-a-map",
			},
			expectedCode:  -32602,
			expectedError: "Invalid params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := types.MCPRequest{
				ID:      "test-invalid-params",
				Method:  "tools/call",
				Params:  tt.params,
				JSONRPC: "2.0",
			}

			response := server.handleToolCall(request)

			assert.Equal(t, request.ID, response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)
			assert.NotNil(t, response.Error)
			assert.Equal(t, tt.expectedCode, response.Error.Code)
			assert.Equal(t, tt.expectedError, response.Error.Message)
			assert.Nil(t, response.Result)
		})
	}
}

// TestHandleToolCall_UnknownTool tests tool call handler with unknown tool
func TestHandleToolCall_UnknownTool(t *testing.T) {
	server := NewAuditQueryMCPServer()

	request := types.MCPRequest{
		ID:     "test-unknown-tool",
		Method: "tools/call",
		Params: map[string]interface{}{
			"name":      "unknown_tool",
			"arguments": map[string]interface{}{},
		},
		JSONRPC: "2.0",
	}

	response := server.handleToolCall(request)

	assert.Equal(t, request.ID, response.ID)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.NotNil(t, response.Error)
	assert.Equal(t, -32601, response.Error.Code)
	assert.Equal(t, "Tool not found", response.Error.Message)
	assert.Nil(t, response.Result)
}

// TestHandleGenerateAuditQueryWithResult tests the generate audit query handler
func TestHandleGenerateAuditQueryWithResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name          string
		params        map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Valid structured params",
			params: map[string]interface{}{
				"structured_params": map[string]interface{}{
					"log_source": "kube-apiserver",
					"patterns":   []string{"pods", "create"},
					"timeframe":  "today",
					"username":   "admin",
				},
			},
			expectSuccess: true,
		},
		{
			name: "Missing structured_params",
			params: map[string]interface{}{
				"other_param": "value",
			},
			expectSuccess: false,
		},
		{
			name: "Invalid structured_params type",
			params: map[string]interface{}{
				"structured_params": "not-a-map",
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.handleGenerateAuditQueryWithResult("test-id", tt.params)

			assert.Equal(t, "test-id", response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)

			if tt.expectSuccess {
				assert.Nil(t, response.Error)
				assert.NotNil(t, response.Result)

				result, ok := response.Result.(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, result, "audit_result")
			} else {
				assert.NotNil(t, response.Error)
				assert.Equal(t, -32602, response.Error.Code)
				assert.Equal(t, "structured_params required", response.Error.Message)
				assert.Nil(t, response.Result)
			}
		})
	}
}

// TestHandleExecuteAuditQueryWithResult tests the execute audit query handler
func TestHandleExecuteAuditQueryWithResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name          string
		params        map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Valid parameters",
			params: map[string]interface{}{
				"command":  "oc adm node-logs --since=1h",
				"query_id": "test-query-123",
			},
			expectSuccess: true,
		},
		{
			name: "Missing command",
			params: map[string]interface{}{
				"query_id": "test-query-123",
			},
			expectSuccess: false,
		},
		{
			name: "Missing query_id",
			params: map[string]interface{}{
				"command": "oc adm node-logs --since=1h",
			},
			expectSuccess: false,
		},
		{
			name: "Invalid command type",
			params: map[string]interface{}{
				"command":  123,
				"query_id": "test-query-123",
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.handleExecuteAuditQueryWithResult("test-id", tt.params)

			assert.Equal(t, "test-id", response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)

			if tt.expectSuccess {
				// Note: This might fail in test environment due to missing oc command
				// but the handler should still process the request
				if response.Error != nil {
					// If there's an error, it should be related to command execution, not parameter validation
					assert.NotEqual(t, -32602, response.Error.Code)
				}
			} else {
				assert.NotNil(t, response.Error)
				assert.Equal(t, -32602, response.Error.Code)
				assert.Nil(t, response.Result)
			}
		})
	}
}

// TestHandleParseAuditResultsWithResult tests the parse audit results handler
func TestHandleParseAuditResultsWithResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name          string
		params        map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Valid parameters",
			params: map[string]interface{}{
				"raw_output": "test raw output",
				"query_context": map[string]interface{}{
					"log_source": "kube-apiserver",
				},
				"query_id": "test-query-123",
			},
			expectSuccess: true,
		},
		{
			name: "Missing raw_output",
			params: map[string]interface{}{
				"query_context": map[string]interface{}{},
				"query_id":      "test-query-123",
			},
			expectSuccess: false,
		},
		{
			name: "Missing query_context",
			params: map[string]interface{}{
				"raw_output": "test raw output",
				"query_id":   "test-query-123",
			},
			expectSuccess: false,
		},
		{
			name: "Missing query_id",
			params: map[string]interface{}{
				"raw_output":    "test raw output",
				"query_context": map[string]interface{}{},
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.handleParseAuditResultsWithResult("test-id", tt.params)

			assert.Equal(t, "test-id", response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)

			if tt.expectSuccess {
				assert.Nil(t, response.Error)
				assert.NotNil(t, response.Result)

				result, ok := response.Result.(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, result, "audit_result")
			} else {
				assert.NotNil(t, response.Error)
				assert.Equal(t, -32602, response.Error.Code)
				assert.Nil(t, response.Result)
			}
		})
	}
}

// TestHandleExecuteCompleteAuditQuery tests the complete audit query handler
func TestHandleExecuteCompleteAuditQuery(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name          string
		params        map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Valid structured params",
			params: map[string]interface{}{
				"structured_params": map[string]interface{}{
					"log_source": "kube-apiserver",
					"patterns":   []string{"pods", "create"},
					"timeframe":  "today",
					"username":   "admin",
				},
			},
			expectSuccess: true,
		},
		{
			name: "Missing structured_params",
			params: map[string]interface{}{
				"other_param": "value",
			},
			expectSuccess: false,
		},
		{
			name: "Invalid structured_params type",
			params: map[string]interface{}{
				"structured_params": "not-a-map",
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.handleExecuteCompleteAuditQuery("test-id", tt.params)

			assert.Equal(t, "test-id", response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)

			if tt.expectSuccess {
				// Note: This might fail in test environment due to missing oc command
				// but the handler should still process the request
				if response.Error != nil {
					// If there's an error, it should be related to command execution, not parameter validation
					assert.NotEqual(t, -32602, response.Error.Code)
				}
			} else {
				assert.NotNil(t, response.Error)
				assert.Equal(t, -32602, response.Error.Code)
				assert.Equal(t, "structured_params required", response.Error.Message)
				assert.Nil(t, response.Result)
			}
		})
	}
}

// TestHandleGetCacheStats tests the get cache stats handler
func TestHandleGetCacheStats(t *testing.T) {
	server := NewAuditQueryMCPServer()

	response := server.handleGetCacheStats("test-id", map[string]interface{}{})

	assert.Equal(t, "test-id", response.ID)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Result)

	result, ok := response.Result.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, result, "cache_stats")
}

// TestHandleClearCache tests the clear cache handler
func TestHandleClearCache(t *testing.T) {
	server := NewAuditQueryMCPServer()

	response := server.handleClearCache("test-id", map[string]interface{}{})

	assert.Equal(t, "test-id", response.ID)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Result)

	result, ok := response.Result.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, result, "message")
	assert.Equal(t, "Cache cleared successfully", result["message"])
}

// TestHandleGetCachedResult tests the get cached result handler
func TestHandleGetCachedResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name          string
		params        map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Valid query_id",
			params: map[string]interface{}{
				"query_id": "test-query-123",
			},
			expectSuccess: false, // Should fail because cache is empty
		},
		{
			name: "Missing query_id",
			params: map[string]interface{}{
				"other_param": "value",
			},
			expectSuccess: false,
		},
		{
			name: "Invalid query_id type",
			params: map[string]interface{}{
				"query_id": 123,
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.handleGetCachedResult("test-id", tt.params)

			assert.Equal(t, "test-id", response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)

			if tt.expectSuccess {
				assert.Nil(t, response.Error)
				assert.NotNil(t, response.Result)
			} else {
				assert.NotNil(t, response.Error)
				assert.Nil(t, response.Result)
			}
		})
	}
}

// TestHandleDeleteCachedResult tests the delete cached result handler
func TestHandleDeleteCachedResult(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name          string
		params        map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Valid query_id",
			params: map[string]interface{}{
				"query_id": "test-query-123",
			},
			expectSuccess: true,
		},
		{
			name: "Missing query_id",
			params: map[string]interface{}{
				"other_param": "value",
			},
			expectSuccess: false,
		},
		{
			name: "Invalid query_id type",
			params: map[string]interface{}{
				"query_id": 123,
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.handleDeleteCachedResult("test-id", tt.params)

			assert.Equal(t, "test-id", response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)

			if tt.expectSuccess {
				assert.Nil(t, response.Error)
				assert.NotNil(t, response.Result)

				result, ok := response.Result.(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, result, "message")
				assert.Equal(t, "Cached result deleted successfully", result["message"])
			} else {
				assert.NotNil(t, response.Error)
				assert.Equal(t, -32602, response.Error.Code)
				assert.Equal(t, "query_id required", response.Error.Message)
				assert.Nil(t, response.Result)
			}
		})
	}
}

// TestHandleGetServerStats tests the get server stats handler
func TestHandleGetServerStats(t *testing.T) {
	server := NewAuditQueryMCPServer()

	response := server.handleGetServerStats("test-id", map[string]interface{}{})

	assert.Equal(t, "test-id", response.ID)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Result)

	result, ok := response.Result.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, result, "server_stats")
}

// TestParameterTypeConversion tests the conversion of interface{} parameters to structured types
func TestParameterTypeConversion(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Test conversion in generate audit query handler
	params := map[string]interface{}{
		"structured_params": map[string]interface{}{
			"log_source": "kube-apiserver",
			"patterns":   []interface{}{"pods", "create", "delete"},
			"timeframe":  "today",
			"exclude":    []interface{}{"system:", "kube-system"},
			"username":   "admin",
			"resource":   "pods",
			"verb":       "create",
			"namespace":  "default",
		},
	}

	response := server.handleGenerateAuditQueryWithResult("test-conversion", params)

	assert.Equal(t, "test-conversion", response.ID)
	assert.Equal(t, "2.0", response.JSONRPC)

	// Should succeed with valid parameter conversion
	if response.Error == nil {
		assert.NotNil(t, response.Result)
	} else {
		// If there's an error, it should not be related to parameter conversion
		assert.NotEqual(t, -32602, response.Error.Code)
	}
}

// TestEdgeCases tests various edge cases in the MCP handler
func TestEdgeCases(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name   string
		testFn func(t *testing.T)
	}{
		{
			name: "Empty request ID",
			testFn: func(t *testing.T) {
				request := types.MCPRequest{
					ID:      "",
					Method:  "tools/list",
					Params:  map[string]interface{}{},
					JSONRPC: "2.0",
				}
				response := server.HandleMCPRequest(request)
				assert.Equal(t, "", response.ID)
				assert.Nil(t, response.Error)
			},
		},
		{
			name: "Nil parameters",
			testFn: func(t *testing.T) {
				request := types.MCPRequest{
					ID:      "test-nil-params",
					Method:  "tools/call",
					Params:  nil,
					JSONRPC: "2.0",
				}
				response := server.HandleMCPRequest(request)
				assert.NotNil(t, response.Error)
				assert.Equal(t, -32602, response.Error.Code)
			},
		},
		{
			name: "Empty parameters map",
			testFn: func(t *testing.T) {
				request := types.MCPRequest{
					ID:      "test-empty-params",
					Method:  "tools/call",
					Params:  map[string]interface{}{},
					JSONRPC: "2.0",
				}
				response := server.HandleMCPRequest(request)
				assert.NotNil(t, response.Error)
				assert.Equal(t, -32602, response.Error.Code)
			},
		},
		{
			name: "Large parameter values",
			testFn: func(t *testing.T) {
				largeString := string(make([]byte, 10000)) // 10KB string
				params := map[string]interface{}{
					"structured_params": map[string]interface{}{
						"log_source": "kube-apiserver",
						"patterns":   []string{largeString},
						"timeframe":  "today",
					},
				}
				response := server.handleGenerateAuditQueryWithResult("test-large-params", params)
				assert.Equal(t, "test-large-params", response.ID)
				// Should handle large parameters gracefully
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFn)
	}
}

// TestConcurrentAccess tests that the MCP handler can handle concurrent requests
func TestConcurrentAccess(t *testing.T) {
	server := NewAuditQueryMCPServer()

	// Run multiple goroutines to test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			// Create a new request for each goroutine to avoid race conditions
			request := types.MCPRequest{
				ID:      fmt.Sprintf("concurrent-test-%d", id),
				Method:  "tools/list",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			}
			response := server.HandleMCPRequest(request)
			assert.Equal(t, request.ID, response.ID)
			assert.Nil(t, response.Error)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestJSONRPCCompliance tests that responses follow JSON-RPC 2.0 specification
func TestJSONRPCCompliance(t *testing.T) {
	server := NewAuditQueryMCPServer()

	tests := []struct {
		name    string
		request types.MCPRequest
	}{
		{
			name: "Tools list request",
			request: types.MCPRequest{
				ID:      "jsonrpc-test-1",
				Method:  "tools/list",
				Params:  map[string]interface{}{},
				JSONRPC: "2.0",
			},
		},
		{
			name: "Tool call request",
			request: types.MCPRequest{
				ID:     "jsonrpc-test-2",
				Method: "tools/call",
				Params: map[string]interface{}{
					"name":      "get_cache_stats",
					"arguments": map[string]interface{}{},
				},
				JSONRPC: "2.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := server.HandleMCPRequest(tt.request)

			// Check JSON-RPC 2.0 compliance
			assert.Equal(t, tt.request.ID, response.ID)
			assert.Equal(t, "2.0", response.JSONRPC)

			// Either result or error should be present, but not both
			if response.Error != nil {
				assert.Nil(t, response.Result)
				assert.NotEmpty(t, response.Error.Message)
				assert.NotZero(t, response.Error.Code)
			} else {
				assert.NotNil(t, response.Result)
			}
		})
	}
}
