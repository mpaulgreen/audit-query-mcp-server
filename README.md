# OpenShift Audit Query MCP Server

A production-ready Model Context Protocol (MCP) server that provides comprehensive, structured access to OpenShift audit logs through intelligent query generation, execution, and result tracking.

## Overview

The OpenShift Audit Query MCP Server enables users to query OpenShift audit logs using structured parameters, which are then converted into safe `oc` commands. The server provides comprehensive audit query capabilities with detailed result tracking, caching, and compliance features.

### Core Capabilities

1. **Generate Audit Queries**: Convert structured parameters into safe `oc adm node-logs` commands with rolling log support
2. **Execute Audit Queries**: Safely run the generated commands and return detailed results
3. **Parse Audit Results**: Convert raw audit log output into structured, readable format with enhanced parsing
4. **Complete Pipeline**: Execute the full query pipeline (generate → execute → parse) in one operation
5. **Result Tracking**: Comprehensive tracking with query IDs, execution times, and error handling
6. **Caching**: Intelligent caching of query results for improved performance
7. **Audit Trail**: Complete audit trail logging for compliance and debugging

9. **Multi-file Log Support**: Intelligent handling of rolling log files across timeframes
10. **Comprehensive Testing**: Extensive test suite with multiple execution modes

## Features

- **Safe Command Generation**: All generated commands are validated for safety with complexity controls
- **Multiple Log Sources**: Support for kube-apiserver, oauth-server, node, openshift-apiserver, and oauth-apiserver
- **Flexible Filtering**: Filter by username, resource, verb, namespace, patterns, and timeframes
- **Comprehensive Validation**: Input validation for all parameters with enhanced security patterns
- **MCP Protocol Support**: Full Model Context Protocol implementation with 9 tools
- **Structured Output**: Parse JSON audit logs into readable summaries with performance metrics
- **Query Result Tracking**: Comprehensive tracking with unique query IDs and execution times
- **Intelligent Caching**: Cache query results for improved performance with TTL support
- **Audit Trail Logging**: Complete audit trail for compliance and debugging
- **Error Handling**: Detailed error information and recovery mechanisms

- **Rolling Log Support**: Intelligent handling of log files across different time periods
- **Performance Optimization**: Efficient parsing and command generation with complexity limits
- **Comprehensive Testing**: 50+ test functions across all packages with multiple execution modes

## Quick Start

### Prerequisites

- Go 1.21 or higher
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

4. Run setup and validation:
```bash
./audit-query-mcp-server setup
```

5. Run the server:
```bash
./audit-query-mcp-server serve
```

### Running Tests

The project includes a comprehensive test suite with multiple execution modes:

#### Basic Test Execution
```bash
# Run all tests
./audit-query-mcp-server test -all

# Run specific test categories
./audit-query-mcp-server test core integration

# Run fast tests only (skip slow integration tests)
./audit-query-mcp-server test -skip-slow

# Run with verbose output
./audit-query-mcp-server test -v -all

# Run with compact output
./audit-query-mcp-server test -compact -all
```

#### Available Test Categories
- **core**: command-builder, validation, caching, parser
- **integration**: mcp-protocol, integration, audit-trail
- **patterns**: nlp-patterns, nlp-simple, command-syntax
- **error**: error-handling
- **cluster**: real-cluster
- **fast**: All tests except slow integration tests
- **slow**: mcp-protocol, integration, nlp-patterns

#### Individual Test Options
```bash
# Test specific components
./audit-query-mcp-server test command-builder validation caching

# Test natural language patterns
./audit-query-mcp-server test nlp-patterns nlp-simple

# Test real cluster connectivity
./audit-query-mcp-server test real-cluster

# Show test help
./audit-query-mcp-server test -h
```

## Usage

### Server Modes

The server supports multiple operation modes:

#### 1. Setup Mode
```bash
./audit-query-mcp-server setup
```
Validates environment, dependencies, and OpenShift connectivity.

#### 2. Test Mode
```bash
./audit-query-mcp-server test [options] [test-names...]
```
Runs the comprehensive test suite with various options.

#### 3. HTTP Server Mode
```bash
./audit-query-mcp-server serve
```
Starts an HTTP server for testing and development (not for production).

### MCP Tools

The server provides 9 comprehensive MCP tools for audit query operations:

#### AuditResult-Based Tools

#### 1. `generate_audit_query_with_result`

Converts structured parameters into safe OpenShift audit commands with rolling log support and returns detailed AuditResult.

**Parameters:**
- `structured_params` (object): Query parameters including:
  - `log_source` (string): Audit log source (kube-apiserver, oauth-server, node, openshift-apiserver, oauth-apiserver)
  - `patterns` (array): Search patterns to filter logs (max 3 for complexity control)
  - `timeframe` (string): Time range for the query with rolling log support
  - `username` (string): Filter by specific username with pattern matching
  - `resource` (string): Filter by Kubernetes resource type
  - `verb` (string): Filter by API verb (create, get, list, delete, etc.)
  - `namespace` (string): Filter by namespace
  - `exclude` (array): Patterns to exclude from results (max 3 for complexity control)

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

Parses raw audit log output into structured format with enhanced parsing and detailed tracking.

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

**Returns:** Cache statistics including size, TTL, hit rates, and performance metrics

#### 6. `clear_cache`

Clears all cached audit results.

**Parameters:** None

**Returns:** Success message with cache statistics

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

**Returns:** Server statistics including version, features, tool counts, and performance metrics



## API Reference

### Enhanced AuditResult Structure

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

### Enhanced AuditQueryParams Structure

The `AuditQueryParams` structure defines the parameters for audit queries with rolling log support:

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

### Enhanced AuditLogEntry Structure

The `AuditLogEntry` structure provides detailed parsing of audit log entries:

```go
type AuditLogEntry struct {
    // Core fields
    Timestamp  string   `json:"timestamp,omitempty"`
    Username   string   `json:"username,omitempty"`
    UID        string   `json:"uid,omitempty"`
    Groups     []string `json:"groups,omitempty"`
    Verb       string   `json:"verb,omitempty"`
    Resource   string   `json:"resource,omitempty"`
    Namespace  string   `json:"namespace,omitempty"`
    Name       string   `json:"name,omitempty"`
    APIGroup   string   `json:"api_group,omitempty"`
    APIVersion string   `json:"api_version,omitempty"`
    RequestURI string   `json:"request_uri,omitempty"`
    UserAgent  string   `json:"user_agent,omitempty"`
    SourceIPs  []string `json:"source_ips,omitempty"`

    // Response fields
    StatusCode    int    `json:"status_code,omitempty"`
    StatusMessage string `json:"status_message,omitempty"`
    StatusReason  string `json:"status_reason,omitempty"`

    // Authentication fields
    AuthDecision     string `json:"auth_decision,omitempty"`
    AuthzDecision    string `json:"authz_decision,omitempty"`
    ImpersonatedUser string `json:"impersonated_user,omitempty"`

    // Additional fields
    Annotations map[string]interface{} `json:"annotations,omitempty"`
    Extra       map[string]interface{} `json:"extra,omitempty"`
    Headers     map[string]interface{} `json:"headers,omitempty"`

    // Metadata
    RawLine     string   `json:"raw_line,omitempty"`
    ParseErrors []string `json:"parse_errors,omitempty"`
    ParseTime   string   `json:"parse_time,omitempty"`
}
```

## Examples

### Basic Query Generation with Rolling Logs

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

### Complete Pipeline Execution with Enhanced Parsing

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



### Cache Management with Enhanced Statistics

```go
// Get comprehensive cache statistics
stats := server.GetCacheStats()
fmt.Printf("Cache size: %v\n", stats["size"])
fmt.Printf("Hit rate: %v\n", stats["hit_rate"])

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
- `PORT`: HTTP server port for testing mode (default: 3000)

### Logging

The server uses structured logging with the following levels:
- `INFO`: General operational information
- `WARN`: Warning messages
- `ERROR`: Error conditions

## Security

### Enhanced Command Validation

All generated commands are validated for safety:
- Only allows `oc adm node-logs` commands
- Validates log source parameters
- Prevents command injection attacks
- Enforces timeout limits
- Complexity controls for patterns and exclusions
- Rolling log file validation

### Enhanced Input Validation

All input parameters are validated:
- Log source must be from allowed list
- Timeframe must be valid format with rolling log support
- Username patterns are sanitized with comprehensive pattern matching
- Resource and verb parameters are validated
- Pattern and exclusion limits enforced

## Performance

### Enhanced Caching

The server implements intelligent caching:
- Query results are cached by query ID
- Configurable TTL for cache entries
- Cache statistics and monitoring
- Manual cache management tools
- Performance metrics tracking

### Optimization

Performance optimizations include:
- Query ID generation for tracking
- Execution time measurement
- Efficient JSON parsing with enhanced structures
- Structured data handling
- Rolling log file optimization
- Complexity controls for command generation

## Testing

### Comprehensive Test Suite

The project includes 50+ test functions across all packages:

#### Test Categories
- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end workflow testing
- **Performance Tests**: Performance and scalability testing
- **Security Tests**: Security validation and edge case testing

- **Real Cluster Tests**: Actual OpenShift cluster connectivity testing

#### Test Execution Modes
- **Fast Mode**: Skip slow integration tests
- **Verbose Mode**: Detailed test output
- **Compact Mode**: Minimal test output
- **Category Mode**: Run specific test categories
- **Individual Mode**: Run specific test functions

## Monitoring

### Enhanced Server Statistics

The server provides comprehensive statistics:
- Tool availability and counts (9 tools)
- Cache performance metrics
- Server version and features
- Execution time tracking
- Rolling log performance metrics

### Enhanced Audit Trail

Complete audit trail logging includes:
- Query generation events with rolling log support
- Query execution events with performance metrics
- Query parsing events with enhanced parsing details
- Cache access events with statistics
- Error conditions with detailed context


## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass with `./audit-query-mcp-server test -all`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions and support:
- Check the comprehensive test suite for usage examples
- Review the API documentation
- Run `./audit-query-mcp-server test -h` for test options
- Open an issue for bugs or feature requests 