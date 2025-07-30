package server

import (
	"audit-query-mcp-server/types"
)

// HandleMCPRequest handles incoming MCP requests
func (s *AuditQueryMCPServer) HandleMCPRequest(request types.MCPRequest) types.MCPResponse {
	s.logger.Infof("Handling MCP request: %s", request.Method)

	switch request.Method {
	case "tools/list":
		return s.handleListTools(request)
	case "tools/call":
		return s.handleToolCall(request)
	default:
		return types.MCPResponse{
			ID: request.ID,
			Error: &types.MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
			JSONRPC: "2.0",
		}
	}
}

// handleListTools handles the tools/list method
func (s *AuditQueryMCPServer) handleListTools(request types.MCPRequest) types.MCPResponse {
	tools := s.GetTools()
	return types.MCPResponse{
		ID:      request.ID,
		Result:  map[string]interface{}{"tools": tools},
		JSONRPC: "2.0",
	}
}

// handleToolCall handles the tools/call method
func (s *AuditQueryMCPServer) handleToolCall(request types.MCPRequest) types.MCPResponse {
	params, ok := request.Params["arguments"].(map[string]interface{})
	if !ok {
		return types.MCPResponse{
			ID: request.ID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
			JSONRPC: "2.0",
		}
	}

	toolName, ok := request.Params["name"].(string)
	if !ok {
		return types.MCPResponse{
			ID: request.ID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "Tool name required",
			},
			JSONRPC: "2.0",
		}
	}

	switch toolName {
	case "generate_audit_query_with_result":
		return s.handleGenerateAuditQueryWithResult(request.ID, params)
	case "execute_audit_query_with_result":
		return s.handleExecuteAuditQueryWithResult(request.ID, params)
	case "parse_audit_results_with_result":
		return s.handleParseAuditResultsWithResult(request.ID, params)
	case "execute_complete_audit_query":
		return s.handleExecuteCompleteAuditQuery(request.ID, params)
	case "get_cache_stats":
		return s.handleGetCacheStats(request.ID, params)
	case "clear_cache":
		return s.handleClearCache(request.ID, params)
	case "get_cached_result":
		return s.handleGetCachedResult(request.ID, params)
	case "delete_cached_result":
		return s.handleDeleteCachedResult(request.ID, params)
	case "get_server_stats":
		return s.handleGetServerStats(request.ID, params)
	default:
		return types.MCPResponse{
			ID: request.ID,
			Error: &types.MCPError{
				Code:    -32601,
				Message: "Tool not found",
			},
			JSONRPC: "2.0",
		}
	}
}

// handleGenerateAuditQueryWithResult handles the generate_audit_query tool with AuditResult
func (s *AuditQueryMCPServer) handleGenerateAuditQueryWithResult(requestID string, params map[string]interface{}) types.MCPResponse {
	structuredParams, ok := params["structured_params"].(map[string]interface{})
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "structured_params required",
			},
			JSONRPC: "2.0",
		}
	}

	// Convert to AuditQueryParams
	auditParams := types.AuditQueryParams{}
	if logSource, ok := structuredParams["log_source"].(string); ok {
		auditParams.LogSource = logSource
	}
	if patterns, ok := structuredParams["patterns"].([]interface{}); ok {
		for _, p := range patterns {
			if pattern, ok := p.(string); ok {
				auditParams.Patterns = append(auditParams.Patterns, pattern)
			}
		}
	}
	if timeframe, ok := structuredParams["timeframe"].(string); ok {
		auditParams.Timeframe = timeframe
	}
	if exclude, ok := structuredParams["exclude"].([]interface{}); ok {
		for _, e := range exclude {
			if ex, ok := e.(string); ok {
				auditParams.Exclude = append(auditParams.Exclude, ex)
			}
		}
	}
	if username, ok := structuredParams["username"].(string); ok {
		auditParams.Username = username
	}
	if resource, ok := structuredParams["resource"].(string); ok {
		auditParams.Resource = resource
	}
	if verb, ok := structuredParams["verb"].(string); ok {
		auditParams.Verb = verb
	}
	if namespace, ok := structuredParams["namespace"].(string); ok {
		auditParams.Namespace = namespace
	}

	result, err := s.GenerateAuditQueryWithResult(auditParams)
	if err != nil {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32000,
				Message: err.Error(),
			},
			JSONRPC: "2.0",
		}
	}

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"audit_result": result,
		},
		JSONRPC: "2.0",
	}
}

// handleExecuteAuditQueryWithResult handles the execute_audit_query tool with AuditResult
func (s *AuditQueryMCPServer) handleExecuteAuditQueryWithResult(requestID string, params map[string]interface{}) types.MCPResponse {
	command, ok := params["command"].(string)
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "command required",
			},
			JSONRPC: "2.0",
		}
	}

	queryID, ok := params["query_id"].(string)
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "query_id required",
			},
			JSONRPC: "2.0",
		}
	}

	result, err := s.ExecuteAuditQueryWithResult(command, queryID)
	if err != nil {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32000,
				Message: err.Error(),
			},
			JSONRPC: "2.0",
		}
	}

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"audit_result": result,
		},
		JSONRPC: "2.0",
	}
}

// handleParseAuditResultsWithResult handles the parse_audit_results tool with AuditResult
func (s *AuditQueryMCPServer) handleParseAuditResultsWithResult(requestID string, params map[string]interface{}) types.MCPResponse {
	rawOutput, ok := params["raw_output"].(string)
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "raw_output required",
			},
			JSONRPC: "2.0",
		}
	}

	queryContext, ok := params["query_context"].(map[string]interface{})
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "query_context required",
			},
			JSONRPC: "2.0",
		}
	}

	queryID, ok := params["query_id"].(string)
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "query_id required",
			},
			JSONRPC: "2.0",
		}
	}

	result, err := s.ParseAuditResultsWithResult(rawOutput, queryContext, queryID)
	if err != nil {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32000,
				Message: err.Error(),
			},
			JSONRPC: "2.0",
		}
	}

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"audit_result": result,
		},
		JSONRPC: "2.0",
	}
}

// handleExecuteCompleteAuditQuery handles the complete audit query pipeline
func (s *AuditQueryMCPServer) handleExecuteCompleteAuditQuery(requestID string, params map[string]interface{}) types.MCPResponse {
	structuredParams, ok := params["structured_params"].(map[string]interface{})
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "structured_params required",
			},
			JSONRPC: "2.0",
		}
	}

	// Convert to AuditQueryParams
	auditParams := types.AuditQueryParams{}
	if logSource, ok := structuredParams["log_source"].(string); ok {
		auditParams.LogSource = logSource
	}
	if patterns, ok := structuredParams["patterns"].([]interface{}); ok {
		for _, p := range patterns {
			if pattern, ok := p.(string); ok {
				auditParams.Patterns = append(auditParams.Patterns, pattern)
			}
		}
	}
	if timeframe, ok := structuredParams["timeframe"].(string); ok {
		auditParams.Timeframe = timeframe
	}
	if exclude, ok := structuredParams["exclude"].([]interface{}); ok {
		for _, e := range exclude {
			if ex, ok := e.(string); ok {
				auditParams.Exclude = append(auditParams.Exclude, ex)
			}
		}
	}
	if username, ok := structuredParams["username"].(string); ok {
		auditParams.Username = username
	}
	if resource, ok := structuredParams["resource"].(string); ok {
		auditParams.Resource = resource
	}
	if verb, ok := structuredParams["verb"].(string); ok {
		auditParams.Verb = verb
	}
	if namespace, ok := structuredParams["namespace"].(string); ok {
		auditParams.Namespace = namespace
	}

	result, err := s.ExecuteCompleteAuditQuery(auditParams)
	if err != nil {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32000,
				Message: err.Error(),
			},
			JSONRPC: "2.0",
		}
	}

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"audit_result": result,
		},
		JSONRPC: "2.0",
	}
}

// handleGetCacheStats handles the get_cache_stats tool
func (s *AuditQueryMCPServer) handleGetCacheStats(requestID string, params map[string]interface{}) types.MCPResponse {
	stats := s.GetCacheStats()

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"cache_stats": stats,
		},
		JSONRPC: "2.0",
	}
}

// handleClearCache handles the clear_cache tool
func (s *AuditQueryMCPServer) handleClearCache(requestID string, params map[string]interface{}) types.MCPResponse {
	s.ClearCache()

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"message": "Cache cleared successfully",
		},
		JSONRPC: "2.0",
	}
}

// handleGetCachedResult handles the get_cached_result tool
func (s *AuditQueryMCPServer) handleGetCachedResult(requestID string, params map[string]interface{}) types.MCPResponse {
	queryID, ok := params["query_id"].(string)
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "query_id required",
			},
			JSONRPC: "2.0",
		}
	}

	result, found := s.GetCachedResult(queryID)
	if !found {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32001,
				Message: "Cached result not found",
			},
			JSONRPC: "2.0",
		}
	}

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"audit_result": result,
		},
		JSONRPC: "2.0",
	}
}

// handleDeleteCachedResult handles the delete_cached_result tool
func (s *AuditQueryMCPServer) handleDeleteCachedResult(requestID string, params map[string]interface{}) types.MCPResponse {
	queryID, ok := params["query_id"].(string)
	if !ok {
		return types.MCPResponse{
			ID: requestID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "query_id required",
			},
			JSONRPC: "2.0",
		}
	}

	s.DeleteCachedResult(queryID)

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"message": "Cached result deleted successfully",
		},
		JSONRPC: "2.0",
	}
}

// handleGetServerStats handles the get_server_stats tool
func (s *AuditQueryMCPServer) handleGetServerStats(requestID string, params map[string]interface{}) types.MCPResponse {
	stats := s.GetServerStats()

	return types.MCPResponse{
		ID: requestID,
		Result: map[string]interface{}{
			"server_stats": stats,
		},
		JSONRPC: "2.0",
	}
}
