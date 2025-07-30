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

## Testing

The OpenShift Audit Query MCP Server includes a comprehensive testing framework with multiple execution modes, unit tests, and integration tests. The testing suite is designed to validate all aspects of the system including command generation, validation, caching, parsing, and MCP protocol compliance.

### Test Suite Overview

The project includes **50+ test functions** across all packages with multiple execution modes:

- **Unit Tests**: Individual component testing for each package
- **Integration Tests**: End-to-end workflow testing
- **Performance Tests**: Performance and scalability testing
- **Security Tests**: Security validation and edge case testing
- **Real Cluster Tests**: Actual OpenShift cluster connectivity testing

### Running Tests

#### Using the Custom Test Runner

The server provides a custom test runner with multiple execution modes:

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

# Show test help
./audit-query-mcp-server test -h
```

#### Test Command Options

The test command supports the following options:

- `-all`: Run all available tests
- `-v`: Verbose output with detailed test information
- `-skip-slow`: Skip slow tests (integration, mcp-protocol)
- `-skip-integration`: Skip integration tests
- `-compact`: Compact output (less verbose)
- `-h`: Show detailed help information

#### Available Test Categories

- **core**: Core functionality (command-builder, validation, caching, parser)
- **integration**: Integration tests (mcp-protocol, integration, audit-trail)
- **patterns**: Pattern matching (nlp-patterns, nlp-simple, command-syntax)
- **error**: Error handling (error-handling)
- **cluster**: Cluster connectivity tests (real-cluster)
- **fast**: Fast tests only (excludes slow tests)
- **slow**: Slow tests only (mcp-protocol, integration, nlp-patterns)

#### Individual Test Components

- **command-builder**: Enhanced command builder functionality
- **validation**: Robust validation patterns
- **caching**: Improved caching mechanisms
- **audit-trail**: Audit trail functionality
- **parser**: Enhanced parser capabilities
- **mcp-protocol**: Comprehensive MCP protocol (slow)
- **integration**: Integration scenarios (slow)
- **error-handling**: Error handling and recovery
- **nlp-patterns**: Natural language patterns (comprehensive)
- **nlp-simple**: Natural language patterns (simple)
- **nlp-compact**: Natural language patterns (compact)
- **command-syntax**: Command syntax validation
- **real-cluster**: Real cluster connectivity test

### Unit Testing with Go

#### Running Go Unit Tests

The project includes comprehensive unit tests for each package. You can run these using standard Go testing:

```bash
# Run all unit tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./commands
go test ./validation
go test ./parsing
go test ./utils
go test ./server
go test ./types

# Run tests with race detection
go test -race ./...

# Run tests with benchmarks
go test -bench=. ./...
```

#### Available Unit Test Files

The project includes unit tests for all major components:

- `commands/builder_test.go` - Command builder functionality tests
- `commands/filters_test.go` - Filter functionality tests
- `validation/validator_test.go` - Input validation tests
- `parsing/parser_test.go` - Audit log parsing tests
- `utils/cache_test.go` - Caching mechanism tests
- `utils/audit_trail_test.go` - Audit trail functionality tests
- `utils/constants_test.go` - Constants and configuration tests
- `server/mcp_handler_test.go` - MCP protocol handler tests
- `server/server_test.go` - Server functionality tests
- `types/types_test.go` - Data structure tests

#### Test Examples

```bash
# Test command builder with specific test
go test -v ./commands -run TestBuildAuditQuery

# Test validation with coverage
go test -cover ./validation

# Test parsing with benchmarks
go test -bench=. ./parsing

# Test caching with race detection
go test -race ./utils

# Test server with specific test
go test -v ./server -run TestNewAuditQueryMCPServer
```

### Integration Testing

#### Real Cluster Testing

The server includes integration tests that connect to actual OpenShift clusters:

```bash
# Test real cluster connectivity
./audit-query-mcp-server test real-cluster

# Run integration tests
./audit-query-mcp-server test integration

# Run MCP protocol tests
./audit-query-mcp-server test mcp-protocol
```

#### Integration Test Requirements

- OpenShift CLI (`oc`) installed and configured
- Access to an OpenShift cluster with audit logging enabled
- Proper authentication and permissions

### Test Execution Modes

#### Fast Mode
```bash
# Run only fast tests (excludes slow integration tests)
./audit-query-mcp-server test -skip-slow
go test ./commands ./validation ./utils ./types
```

#### Verbose Mode
```bash
# Run with detailed output
./audit-query-mcp-server test -v -all
go test -v ./...
```

#### Compact Mode
```bash
# Run with minimal output
./audit-query-mcp-server test -compact -all
```

#### Coverage Mode
```bash
# Run with coverage analysis
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Help Output

Running `./audit-query-mcp-server test -h` provides detailed information:

```
🧪 Audit Query MCP Server Test Suite
=====================================

Usage:
  go run . test [options] [test-names...]

Options:
  -all              Run all tests
  -v                Verbose output
  -skip-slow        Skip slow tests (integration, mcp-protocol)
  -skip-integration Skip integration tests
  -compact          Compact output (less verbose)
  -h                Show this help

Test Categories:
  core             - Core functionality (command-builder, validation, caching, parser)
  integration      - Integration tests (mcp-protocol, integration, audit-trail)
  patterns         - Pattern matching (nlp-patterns, nlp-simple, command-syntax)
  error            - Error handling (error-handling)
  cluster          - Cluster connectivity tests (real-cluster)
  fast             - Fast tests only (excludes slow tests)
  slow             - Slow tests only (mcp-protocol, integration, nlp-patterns)

Available Tests:
  command-builder   - Enhanced command builder functionality
  validation        - Robust validation patterns
  caching           - Improved caching mechanisms
  audit-trail       - Audit trail functionality
  parser            - Enhanced parser capabilities
  mcp-protocol      - Comprehensive MCP protocol (slow)
  integration       - Integration scenarios (slow)
  error-handling    - Error handling and recovery
  nlp-patterns      - Natural language patterns (comprehensive)
  nlp-simple        - Natural language patterns (simple)
  nlp-compact       - Natural language patterns (compact)
  command-syntax    - Command syntax validation
  real-cluster      - Real cluster connectivity test

Examples:
  go run . test -all                    # Run all tests
  go run . test command-builder         # Run specific test
  go run . test validation caching      # Run multiple tests
  go run . test -skip-slow              # Run fast tests only
  go run . test core                    # Run core tests
  go run . test -v command-builder      # Verbose output
  go run . test -compact command-builder # Compact output
```

### Test Development

#### Adding New Tests

When adding new functionality, include corresponding tests:

1. **Unit Tests**: Add tests in the appropriate `*_test.go` file
2. **Integration Tests**: Add to the custom test runner in `test_client.go`
3. **Test Categories**: Update the `testCategories` map in `test_client.go`

#### Test Best Practices

- Write tests for all public functions
- Include both positive and negative test cases
- Test edge cases and error conditions
- Use descriptive test names
- Include benchmarks for performance-critical code
- Maintain test coverage above 80%

#### Running Tests in CI/CD

```bash
# Run all tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run tests with race detection
go test -race ./...

# Run integration tests
./audit-query-mcp-server test integration

# Run security tests
./audit-query-mcp-server test validation
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