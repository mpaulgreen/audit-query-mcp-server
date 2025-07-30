package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"

	"audit-query-mcp-server/commands"
	"audit-query-mcp-server/parsing"
	"audit-query-mcp-server/types"
	"audit-query-mcp-server/utils"
	"audit-query-mcp-server/validation"
)

// AuditQueryMCPServer represents the MCP server for OpenShift audit queries
type AuditQueryMCPServer struct {
	client     *openai.Client
	logger     *logrus.Logger
	cache      *utils.Cache
	auditTrail *utils.AuditTrail
}

// NewAuditQueryMCPServer creates a new MCP server instance
func NewAuditQueryMCPServer() *AuditQueryMCPServer {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize OpenAI client (optional for current implementation)
	var client *openai.Client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" && apiKey != "dummy_key_for_testing" {
		client = openai.NewClient(apiKey)
		log.Println("OpenAI client initialized for future LLM integration")
	} else {
		log.Println("OpenAI API key not provided - LLM features will be disabled")
		log.Println("This is normal for the current rule-based implementation")
	}

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Initialize cache with 1 hour default TTL
	cache := utils.NewCache(1 * time.Hour)

	// Initialize audit trail
	auditTrail, err := utils.NewAuditTrail("./logs/audit_trail.json")
	if err != nil {
		log.Printf("Warning: Failed to initialize audit trail: %v", err)
		auditTrail = nil
	}

	return &AuditQueryMCPServer{
		client:     client,
		logger:     logger,
		cache:      cache,
		auditTrail: auditTrail,
	}
}

// GetTools returns the list of available MCP tools
func (s *AuditQueryMCPServer) GetTools() []types.MCPTool {
	return []types.MCPTool{
		// New AuditResult-based tools
		{
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
								"enum": []string{"kube-apiserver", "oauth-server", "node", "openshift-apiserver", "oauth-apiserver"},
							},
							"patterns": map[string]interface{}{
								"type": "array",
								"items": map[string]interface{}{
									"type": "string",
								},
							},
							"timeframe": map[string]interface{}{
								"type": "string",
							},
							"username": map[string]interface{}{
								"type": "string",
							},
							"resource": map[string]interface{}{
								"type": "string",
							},
							"verb": map[string]interface{}{
								"type": "string",
							},
							"namespace": map[string]interface{}{
								"type": "string",
							},
							"exclude": map[string]interface{}{
								"type": "array",
								"items": map[string]interface{}{
									"type": "string",
								},
							},
						},
						"required": []string{"log_source", "timeframe"},
					},
				},
				"required": []string{"structured_params"},
			},
		},
		{
			Name:        "execute_audit_query_with_result",
			Description: "Safely execute the oc command and return detailed AuditResult",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type": "string",
					},
					"query_id": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"command", "query_id"},
			},
		},
		{
			Name:        "parse_audit_results_with_result",
			Description: "Parse oc output into structured AuditResult format with detailed tracking",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"raw_output": map[string]interface{}{
						"type": "string",
					},
					"query_context": map[string]interface{}{
						"type": "object",
					},
					"query_id": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"raw_output", "query_context", "query_id"},
			},
		},
		{
			Name:        "execute_complete_audit_query",
			Description: "Execute the complete audit query pipeline (generate, execute, parse) and return comprehensive AuditResult",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"structured_params": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"log_source": map[string]interface{}{
								"type": "string",
								"enum": []string{"kube-apiserver", "oauth-server", "node", "openshift-apiserver", "oauth-apiserver"},
							},
							"patterns": map[string]interface{}{
								"type": "array",
								"items": map[string]interface{}{
									"type": "string",
								},
							},
							"timeframe": map[string]interface{}{
								"type": "string",
							},
							"username": map[string]interface{}{
								"type": "string",
							},
							"resource": map[string]interface{}{
								"type": "string",
							},
							"verb": map[string]interface{}{
								"type": "string",
							},
							"namespace": map[string]interface{}{
								"type": "string",
							},
							"exclude": map[string]interface{}{
								"type": "array",
								"items": map[string]interface{}{
									"type": "string",
								},
							},
						},
						"required": []string{"log_source", "timeframe"},
					},
				},
				"required": []string{"structured_params"},
			},
		},
		// Cache management tools
		{
			Name:        "get_cache_stats",
			Description: "Get cache statistics and performance metrics",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "clear_cache",
			Description: "Clear all cached audit results",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "get_cached_result",
			Description: "Retrieve a cached audit result by query ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query_id": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"query_id"},
			},
		},
		{
			Name:        "delete_cached_result",
			Description: "Delete a specific cached audit result by query ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query_id": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"query_id"},
			},
		},
		{
			Name:        "get_server_stats",
			Description: "Get comprehensive server statistics and feature information",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

// GetLogger returns the logger instance
func (s *AuditQueryMCPServer) GetLogger() *logrus.Logger {
	return s.logger
}

// GenerateAuditQueryWithResult converts JSON parameters to safe oc audit commands and returns AuditResult
func (s *AuditQueryMCPServer) GenerateAuditQueryWithResult(params types.AuditQueryParams) (*types.AuditResult, error) {
	s.logger.Info("Generating audit query from parameters with result tracking")

	startTime := time.Now()
	queryID := s.generateQueryID()

	result := &types.AuditResult{
		QueryID:   queryID,
		Timestamp: startTime.Format(time.RFC3339),
		Error:     "",
	}

	// Safety validation
	if err := validation.ValidateQueryParams(params); err != nil {
		result.Error = fmt.Sprintf("validation failed: %v", err)
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		return result, fmt.Errorf("validation failed: %w", err)
	}

	// Build the oc command based on parameters
	command := commands.BuildOcCommand(params)
	result.Command = command

	// Additional safety check
	if err := validation.ValidateGeneratedCommand(command); err != nil {
		result.Error = fmt.Sprintf("command validation failed: %v", err)
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		return result, fmt.Errorf("command validation failed: %w", err)
	}

	result.ExecutionTime = time.Since(startTime).Milliseconds()
	s.logger.Infof("Generated command: %s", command)
	return result, nil
}

// ExecuteAuditQueryWithResult safely executes the oc command and returns AuditResult
func (s *AuditQueryMCPServer) ExecuteAuditQueryWithResult(command string, queryID string) (*types.AuditResult, error) {
	s.logger.Info("Executing audit query with result tracking")

	startTime := time.Now()

	result := &types.AuditResult{
		QueryID:   queryID,
		Timestamp: startTime.Format(time.RFC3339),
		Command:   command,
		Error:     "",
	}

	// Final safety validation
	if err := validation.ValidateGeneratedCommand(command); err != nil {
		result.Error = fmt.Sprintf("command validation failed: %v", err)
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		return result, fmt.Errorf("command validation failed: %w", err)
	}

	// Execute with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		result.Error = "command execution timed out after 30 seconds"
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		return result, fmt.Errorf("command execution timed out after 30 seconds")
	}

	if err != nil {
		result.Error = fmt.Sprintf("command execution failed: %v, output: %s", err, string(output))
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		return result, fmt.Errorf("command execution failed: %w, output: %s", err, string(output))
	}

	result.RawOutput = string(output)
	result.ExecutionTime = time.Since(startTime).Milliseconds()
	s.logger.Infof("Command executed successfully, output length: %d", len(output))
	return result, nil
}

// ParseAuditResultsWithResult parses oc output into structured AuditResult format
func (s *AuditQueryMCPServer) ParseAuditResultsWithResult(rawOutput string, queryContext map[string]interface{}, queryID string) (*types.AuditResult, error) {
	s.logger.Info("Parsing audit results with enhanced parser")

	startTime := time.Now()

	result := &types.AuditResult{
		QueryID:   queryID,
		Timestamp: startTime.Format(time.RFC3339),
		RawOutput: rawOutput,
		Error:     "",
	}

	// Split output into lines
	lines := strings.Split(rawOutput, "\n")
	var validLines []string

	// Filter out empty lines
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			validLines = append(validLines, line)
		}
	}

	// Use enhanced parser
	config := parsing.DefaultParserConfig()
	parseResult := parsing.ParseAuditLogs(validLines, config)

	// Convert to legacy format for backward compatibility
	var parsedEntries []map[string]interface{}
	for _, entry := range parseResult.Entries {
		legacyEntry := map[string]interface{}{
			"timestamp":         entry.Timestamp,
			"username":          entry.Username,
			"uid":               entry.UID,
			"groups":            entry.Groups,
			"verb":              entry.Verb,
			"resource":          entry.Resource,
			"namespace":         entry.Namespace,
			"name":              entry.Name,
			"api_group":         entry.APIGroup,
			"api_version":       entry.APIVersion,
			"request_uri":       entry.RequestURI,
			"user_agent":        entry.UserAgent,
			"source_ips":        entry.SourceIPs,
			"status_code":       entry.StatusCode,
			"status_message":    entry.StatusMessage,
			"status_reason":     entry.StatusReason,
			"auth_decision":     entry.AuthDecision,
			"authz_decision":    entry.AuthzDecision,
			"impersonated_user": entry.ImpersonatedUser,
			"annotations":       entry.Annotations,
			"extra":             entry.Extra,
			"headers":           entry.Headers,
			"raw_line":          entry.RawLine,
			"parse_errors":      entry.ParseErrors,
			"parse_time":        entry.ParseTime.Format(time.RFC3339),
		}
		parsedEntries = append(parsedEntries, legacyEntry)
	}

	result.ParsedData = parsedEntries

	// Generate summary using enhanced parser
	result.Summary = parsing.GenerateSummary(parseResult.Entries, queryContext)
	result.ExecutionTime = time.Since(startTime).Milliseconds()

	s.logger.Infof("Parsed %d entries (enhanced parser)", len(parsedEntries))
	s.logger.Infof("Parse performance: %.2f lines/second", parseResult.Performance.LinesPerSecond)

	return result, nil
}

// ExecuteCompleteAuditQuery executes the full audit query pipeline and returns AuditResult
func (s *AuditQueryMCPServer) ExecuteCompleteAuditQuery(params types.AuditQueryParams) (*types.AuditResult, error) {
	s.logger.Info("Executing complete audit query pipeline")

	// Step 1: Generate query
	generateResult, err := s.GenerateAuditQueryWithResult(params)
	if err != nil {
		// Log audit trail for failed generation
		if s.auditTrail != nil {
			s.auditTrail.LogQueryGeneration(generateResult.QueryID, params, generateResult, "", "", "")
		}
		return generateResult, err
	}

	// Check cache for existing result
	if cachedResult, found := s.cache.Get(generateResult.QueryID); found {
		s.logger.Infof("Cache hit for query ID: %s", generateResult.QueryID)
		// Log cache access
		if s.auditTrail != nil {
			s.auditTrail.LogCacheAccess(generateResult.QueryID, "hit", "", "", "")
		}
		return cachedResult, nil
	}

	// Step 2: Execute query
	executeResult, err := s.ExecuteAuditQueryWithResult(generateResult.Command, generateResult.QueryID)
	if err != nil {
		// Merge error information
		generateResult.Error = executeResult.Error
		generateResult.ExecutionTime += executeResult.ExecutionTime
		// Log audit trail for failed execution
		if s.auditTrail != nil {
			s.auditTrail.LogQueryExecution(generateResult.QueryID, generateResult.Command, executeResult, "", "", "")
		}
		return generateResult, err
	}

	// Step 3: Parse results
	queryContext := map[string]interface{}{
		"log_source": params.LogSource,
		"timeframe":  params.Timeframe,
		"username":   params.Username,
		"resource":   params.Resource,
		"verb":       params.Verb,
		"namespace":  params.Namespace,
	}

	parseResult, err := s.ParseAuditResultsWithResult(executeResult.RawOutput, queryContext, generateResult.QueryID)
	if err != nil {
		// Merge error information
		executeResult.Error = parseResult.Error
		executeResult.ExecutionTime += parseResult.ExecutionTime
		// Log audit trail for failed parsing
		if s.auditTrail != nil {
			s.auditTrail.LogQueryParsing(generateResult.QueryID, queryContext, parseResult, "", "", "")
		}
		return executeResult, err
	}

	// Combine all results
	finalResult := &types.AuditResult{
		QueryID:       generateResult.QueryID,
		Timestamp:     generateResult.Timestamp,
		Command:       generateResult.Command,
		RawOutput:     executeResult.RawOutput,
		ParsedData:    parseResult.ParsedData,
		Summary:       parseResult.Summary,
		Error:         "",
		ExecutionTime: generateResult.ExecutionTime + executeResult.ExecutionTime + parseResult.ExecutionTime,
	}

	// Cache the result
	s.cache.Set(generateResult.QueryID, finalResult)
	s.logger.Infof("Cached result for query ID: %s", generateResult.QueryID)

	// Log complete query execution
	if s.auditTrail != nil {
		s.auditTrail.LogCompleteQuery(generateResult.QueryID, params, finalResult, "", "", "")
	}

	return finalResult, nil
}

// generateQueryID creates a unique query identifier
func (s *AuditQueryMCPServer) generateQueryID() string {
	return fmt.Sprintf("audit_query_%s_%s",
		time.Now().Format("20060102_150405"),
		s.generateRandomSuffix())
}

// generateRandomSuffix creates a random suffix for query IDs
func (s *AuditQueryMCPServer) generateRandomSuffix() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// GetCacheStats returns cache statistics
func (s *AuditQueryMCPServer) GetCacheStats() map[string]interface{} {
	return s.cache.GetStats()
}

// ClearCache clears all cached results
func (s *AuditQueryMCPServer) ClearCache() {
	s.cache.Clear()
	s.logger.Info("Cache cleared")
}

// GetCachedResult retrieves a cached result by query ID
func (s *AuditQueryMCPServer) GetCachedResult(queryID string) (*types.AuditResult, bool) {
	return s.cache.Get(queryID)
}

// DeleteCachedResult removes a specific cached result
func (s *AuditQueryMCPServer) DeleteCachedResult(queryID string) {
	s.cache.Delete(queryID)
	s.logger.Infof("Deleted cached result for query ID: %s", queryID)
}

// GetServerStats returns comprehensive server statistics
func (s *AuditQueryMCPServer) GetServerStats() map[string]interface{} {
	stats := map[string]interface{}{
		"server_info": map[string]interface{}{
			"version":      "1.0.0",
			"phase":        "2",
			"audit_result": true,
			"caching":      true,
			"audit_trail":  s.auditTrail != nil,
		},
		"cache_stats": s.GetCacheStats(),
		"tools": map[string]interface{}{
			"audit_result_tools": 4,
			"cache_tools":        5,
			"total_tools":        9,
		},
		"features": map[string]interface{}{
			"query_id_generation":    true,
			"execution_tracking":     true,
			"error_handling":         true,
			"performance_monitoring": true,
		},
	}

	return stats
}
