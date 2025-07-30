# OpenShift Audit Query MCP Server

A production-ready Model Context Protocol (MCP) server that provides comprehensive, structured access to OpenShift audit logs through intelligent query generation, execution, and result tracking.

## Overview

The OpenShift Audit Query MCP Server enables users to query OpenShift audit logs using structured parameters, which are then converted into safe `oc` commands. The server provides comprehensive audit query capabilities with detailed result tracking, caching, and compliance features.

### Core Capabilities

1. **Generate Audit Queries**: Convert structured parameters into safe `oc adm node-logs` commands
2. **Execute Audit Queries**: Safely run the generated commands and return detailed results
3. **Parse Audit Results**: Convert raw audit log output into structured, readable format
4. **Complete Pipeline**: Execute the full query pipeline (generate → execute → parse) in one operation
5. **Result Tracking**: Comprehensive tracking with query IDs, execution times, and error handling
6. **Caching**: Intelligent caching of query results for improved performance
7. **Audit Trail**: Complete audit trail logging for compliance and debugging

## Features

- **Safe Command Generation**: All generated commands are validated for safety
- **Multiple Log Sources**: Support for kube-apiserver, oauth-server, node, openshift-apiserver, and oauth-apiserver
- **Flexible Filtering**: Filter by username, resource, verb, namespace, patterns, and timeframes
- **Comprehensive Validation**: Input validation for all parameters
- **MCP Protocol Support**: Full Model Context Protocol implementation
- **Structured Output**: Parse JSON audit logs into readable summaries
- **Query Result Tracking**: Comprehensive tracking with unique query IDs and execution times
- **Intelligent Caching**: Cache query results for improved performance
- **Audit Trail Logging**: Complete audit trail for compliance and debugging
- **Error Handling**: Detailed error information and recovery mechanisms

## Quick Start

### Prerequisites

- Go 1.19 or higher
- OpenShift CLI (`oc`) installed and configured
- Access to an OpenShift cluster with audit logging enabled

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd audit-query-mcp-server
```

2. Install dependencies:
```bash
go mod download
```

3. Build the server:
```bash
go build -o audit-query-mcp-server .
```

4. Run the server:
```bash
./audit-query-mcp-server
```

### Running Tests

To run the comprehensive test suite:
```bash
go run main.go test
```

This runs the complete test suite including:
- Comprehensive AuditResult integration tests
- MCP protocol tests with all tools
- Tool availability and categorization tests


## Usage

### MCP Tools

The server provides comprehensive MCP tools for audit query operations:

#### AuditResult-Based Tools

#### 1. `generate_audit_query_with_result`

Converts structured parameters into safe OpenShift audit commands and returns detailed AuditResult.

**Parameters:**
- `structured_params` (object): Query parameters including:
  - `log_source` (string): Audit log source (kube-apiserver, oauth-server, node, openshift-apiserver, oauth-apiserver)
  - `patterns` (array): Search patterns to filter logs
  - `timeframe` (string): Time range for the query
  - `username` (string): Filter by specific username
  - `resource` (string): Filter by Kubernetes resource type
  - `verb` (string): Filter by API verb (create, get, list, delete, etc.)
  - `namespace` (string): Filter by namespace
  - `exclude` (array): Patterns to exclude from results

**Returns:** AuditResult object with query ID, command, execution time, and error information

**Example:**
```json
{
  "structured_params": {
    "log_source": "kube-apiserver",
    "patterns": ["pods", "delete"],
    "username": "admin",
    "timeframe": "today",
    "exclude": ["system:"]
  }
}
```

#### 2. `execute_audit_query_with_result`

Safely executes the generated `oc` command and returns detailed AuditResult.

**Parameters:**
- `command` (string): The `oc` command to execute
- `query_id` (string): Unique query identifier for tracking

**Returns:** AuditResult object with raw output, execution time, and error information

#### 3. `parse_audit_results_with_result`

Parses raw audit log output into structured format with detailed tracking.

**Parameters:**
- `raw_output` (string): Raw audit log output
- `query_context` (object): Context information for parsing
- `query_id` (string): Unique query identifier for tracking

**Returns:** AuditResult object with parsed data, summary, and execution metrics

#### 4. `execute_complete_audit_query`

Executes the complete audit query pipeline (generate → execute → parse) in one operation.

**Parameters:** Same as `generate_audit_query_with_result`

**Returns:** Complete AuditResult object with all pipeline results

#### Cache Management Tools

#### 5. `get_cache_stats`

Retrieves cache statistics and performance metrics.

**Parameters:** None

**Returns:** Cache statistics including size, TTL, and hit rates

#### 6. `clear_cache`

Clears all cached audit results.

**Parameters:** None

**Returns:** Success message

#### 7. `get_cached_result`

Retrieves a cached audit result by query ID.

**Parameters:**
- `query_id` (string): Query identifier to retrieve

**Returns:** Cached AuditResult object or error if not found

#### 8. `delete_cached_result`

Deletes a specific cached audit result by query ID.

**Parameters:**
- `query_id` (string): Query identifier to delete

**Returns:** Success message

#### 9. `get_server_stats`

Retrieves comprehensive server statistics and feature information.

**Parameters:** None

**Returns:** Server statistics including version, features, and tool counts

## API Reference

### AuditResult Structure

The `AuditResult` structure provides comprehensive information about audit query operations:

```go
type AuditResult struct {
    QueryID       string                   `json:"query_id"`
    Timestamp     string                   `json:"timestamp"`
    Command       string                   `json:"command"`
    RawOutput     string                   `json:"raw_output"`
    ParsedData    []map[string]interface{} `json:"parsed_data"`
    Summary       string                   `json:"summary"`
    Error         string                   `json:"error,omitempty"`
    ExecutionTime int64                    `json:"execution_time_ms"`
}
```

### AuditQueryParams Structure

The `AuditQueryParams` structure defines the parameters for audit queries:

```go
type AuditQueryParams struct {
    LogSource string   `json:"log_source"`
    Patterns  []string `json:"patterns"`
    Timeframe string   `json:"timeframe"`
    Exclude   []string `json:"exclude"`
    Username  string   `json:"username,omitempty"`
    Resource  string   `json:"resource,omitempty"`
    Verb      string   `json:"verb,omitempty"`
    Namespace string   `json:"namespace,omitempty"`
}
```

## Examples

### Basic Query Generation

```go
params := types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"pods", "delete"},
    Timeframe: "today",
    Username:  "admin",
}

result, err := server.GenerateAuditQueryWithResult(params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Query ID: %s\n", result.QueryID)
fmt.Printf("Command: %s\n", result.Command)
fmt.Printf("Execution Time: %dms\n", result.ExecutionTime)
```

### Complete Pipeline Execution

```go
params := types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"authentication", "failed"},
    Timeframe: "last hour",
}

result, err := server.ExecuteCompleteAuditQuery(params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Found %d audit entries\n", len(result.ParsedData))
fmt.Printf("Summary: %s\n", result.Summary)
```

### Cache Management

```go
// Get cache statistics
stats := server.GetCacheStats()
fmt.Printf("Cache size: %v\n", stats["size"])

// Retrieve cached result
cachedResult, found := server.GetCachedResult("query_id_123")
if found {
    fmt.Printf("Retrieved cached result\n")
}
```

## Configuration

### Environment Variables

The server can be configured using environment variables:

- `OPENAI_API_KEY`: OpenAI API key for future LLM integration (optional)
- `CACHE_TTL`: Cache time-to-live duration (default: 1 hour)
- `AUDIT_TRAIL_PATH`: Path for audit trail logging (default: ./logs/audit_trail.json)

### Logging

The server uses structured logging with the following levels:
- `INFO`: General operational information
- `WARN`: Warning messages
- `ERROR`: Error conditions

## Security

### Command Validation

All generated commands are validated for safety:
- Only allows `oc adm node-logs` commands
- Validates log source parameters
- Prevents command injection attacks
- Enforces timeout limits

### Input Validation

All input parameters are validated:
- Log source must be from allowed list
- Timeframe must be valid format
- Username patterns are sanitized
- Resource and verb parameters are validated

## Performance

### Caching

The server implements intelligent caching:
- Query results are cached by query ID
- Configurable TTL for cache entries
- Cache statistics and monitoring
- Manual cache management tools

### Optimization

Performance optimizations include:
- Query ID generation for tracking
- Execution time measurement
- Efficient JSON parsing
- Structured data handling

## Monitoring

### Server Statistics

The server provides comprehensive statistics:
- Tool availability and counts
- Cache performance metrics
- Server version and features
- Execution time tracking

### Audit Trail

Complete audit trail logging includes:
- Query generation events
- Query execution events
- Query parsing events
- Cache access events
- Error conditions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions and support:
- Check the test suite for usage examples
- Review the API documentation
- Open an issue for bugs or feature requests 