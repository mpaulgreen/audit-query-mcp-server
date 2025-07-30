package types

// MCPTool represents an MCP tool definition
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPRequest represents an MCP tool call request
type MCPRequest struct {
	ID      string                 `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	JSONRPC string                 `json:"jsonrpc"`
}

// MCPResponse represents an MCP tool call response
type MCPResponse struct {
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   *MCPError   `json:"error,omitempty"`
	JSONRPC string      `json:"jsonrpc"`
}

// MCPError represents an MCP error response
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
