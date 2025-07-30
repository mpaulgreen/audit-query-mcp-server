package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"audit-query-mcp-server/types"
)

// AuditTrailEntry represents an audit trail entry
type AuditTrailEntry struct {
	Timestamp     string                 `json:"timestamp"`
	QueryID       string                 `json:"query_id"`
	UserID        string                 `json:"user_id,omitempty"`
	Action        string                 `json:"action"`
	Parameters    map[string]interface{} `json:"parameters"`
	Result        *types.AuditResult     `json:"result,omitempty"`
	Error         string                 `json:"error,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	ExecutionTime int64                  `json:"execution_time_ms"`
}

// AuditTrail provides audit logging functionality
type AuditTrail struct {
	filePath string
	mutex    sync.Mutex
	file     *os.File
	encoder  *json.Encoder
}

// NewAuditTrail creates a new audit trail instance
func NewAuditTrail(filePath string) (*AuditTrail, error) {
	// Validate file path
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audit trail directory: %w", err)
	}

	// Open or create the audit trail file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit trail file: %w", err)
	}

	trail := &AuditTrail{
		filePath: filePath,
		file:     file,
		encoder:  json.NewEncoder(file),
	}

	return trail, nil
}

// LogQuery logs an audit query execution
func (at *AuditTrail) LogQuery(entry AuditTrailEntry) error {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	// Ensure timestamp is set
	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format(time.RFC3339)
	}

	// Write the entry as a JSON line
	if err := at.encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode audit trail entry: %w", err)
	}

	// Flush to ensure data is written
	if err := at.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync audit trail file: %w", err)
	}

	return nil
}

// LogQueryGeneration logs a query generation event
func (at *AuditTrail) LogQueryGeneration(queryID string, params types.AuditQueryParams, result *types.AuditResult, userID, ipAddress, userAgent string) error {
	entry := AuditTrailEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		QueryID:       queryID,
		UserID:        userID,
		Action:        "query_generation",
		Parameters:    paramsToMap(params),
		Result:        result,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExecutionTime: result.ExecutionTime,
	}

	if result.Error != "" {
		entry.Error = result.Error
	}

	return at.LogQuery(entry)
}

// LogQueryExecution logs a query execution event
func (at *AuditTrail) LogQueryExecution(queryID string, command string, result *types.AuditResult, userID, ipAddress, userAgent string) error {
	entry := AuditTrailEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		QueryID:       queryID,
		UserID:        userID,
		Action:        "query_execution",
		Parameters:    map[string]interface{}{"command": command},
		Result:        result,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExecutionTime: result.ExecutionTime,
	}

	if result.Error != "" {
		entry.Error = result.Error
	}

	return at.LogQuery(entry)
}

// LogQueryParsing logs a query parsing event
func (at *AuditTrail) LogQueryParsing(queryID string, queryContext map[string]interface{}, result *types.AuditResult, userID, ipAddress, userAgent string) error {
	entry := AuditTrailEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		QueryID:       queryID,
		UserID:        userID,
		Action:        "query_parsing",
		Parameters:    queryContext,
		Result:        result,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExecutionTime: result.ExecutionTime,
	}

	if result.Error != "" {
		entry.Error = result.Error
	}

	return at.LogQuery(entry)
}

// LogCompleteQuery logs a complete query pipeline execution
func (at *AuditTrail) LogCompleteQuery(queryID string, params types.AuditQueryParams, result *types.AuditResult, userID, ipAddress, userAgent string) error {
	entry := AuditTrailEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		QueryID:       queryID,
		UserID:        userID,
		Action:        "complete_query",
		Parameters:    paramsToMap(params),
		Result:        result,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExecutionTime: result.ExecutionTime,
	}

	if result.Error != "" {
		entry.Error = result.Error
	}

	return at.LogQuery(entry)
}

// LogCacheAccess logs a cache access event
func (at *AuditTrail) LogCacheAccess(queryID string, action string, userID, ipAddress, userAgent string) error {
	entry := AuditTrailEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		QueryID:   queryID,
		UserID:    userID,
		Action:    "cache_" + action,
		Parameters: map[string]interface{}{
			"cache_action": action,
		},
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	return at.LogQuery(entry)
}

// Close closes the audit trail file
func (at *AuditTrail) Close() error {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	if at.file != nil {
		return at.file.Close()
	}
	return nil
}

// paramsToMap converts AuditQueryParams to a map for logging
func paramsToMap(params types.AuditQueryParams) map[string]interface{} {
	return map[string]interface{}{
		"log_source": params.LogSource,
		"patterns":   params.Patterns,
		"timeframe":  params.Timeframe,
		"exclude":    params.Exclude,
		"username":   params.Username,
		"resource":   params.Resource,
		"verb":       params.Verb,
		"namespace":  params.Namespace,
	}
}
