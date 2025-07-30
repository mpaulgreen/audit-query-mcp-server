package parsing

import (
	"time"

	"audit-query-mcp-server/types"
)

// ConvertToTypesAuditLogEntry converts parsing.AuditLogEntry to types.AuditLogEntry
func ConvertToTypesAuditLogEntry(entry AuditLogEntry) types.AuditLogEntry {
	return types.AuditLogEntry{
		Timestamp:        entry.Timestamp,
		Username:         entry.Username,
		UID:              entry.UID,
		Groups:           entry.Groups,
		Verb:             entry.Verb,
		Resource:         entry.Resource,
		Namespace:        entry.Namespace,
		Name:             entry.Name,
		APIGroup:         entry.APIGroup,
		APIVersion:       entry.APIVersion,
		RequestURI:       entry.RequestURI,
		UserAgent:        entry.UserAgent,
		SourceIPs:        entry.SourceIPs,
		StatusCode:       entry.StatusCode,
		StatusMessage:    entry.StatusMessage,
		StatusReason:     entry.StatusReason,
		AuthDecision:     entry.AuthDecision,
		AuthzDecision:    entry.AuthzDecision,
		ImpersonatedUser: entry.ImpersonatedUser,
		Annotations:      entry.Annotations,
		Extra:            entry.Extra,
		Headers:          entry.Headers,
		RawLine:          entry.RawLine,
		ParseErrors:      entry.ParseErrors,
		ParseTime:        entry.ParseTime.Format(time.RFC3339),
	}
}

// ConvertToTypesParseResult converts parsing.ParseResult to types.ParseResult
func ConvertToTypesParseResult(result ParseResult) types.ParseResult {
	entries := make([]types.AuditLogEntry, len(result.Entries))
	for i, entry := range result.Entries {
		entries[i] = ConvertToTypesAuditLogEntry(entry)
	}

	return types.ParseResult{
		Entries:     entries,
		TotalLines:  result.TotalLines,
		ParsedLines: result.ParsedLines,
		ErrorLines:  result.ErrorLines,
		ParseErrors: result.ParseErrors,
		ParseTime:   result.ParseTime.String(),
		Performance: types.ParsePerformance(result.Performance),
	}
}

// ConvertToTypesParserConfig converts parsing.ParserConfig to types.ParserConfig
func ConvertToTypesParserConfig(config ParserConfig) types.ParserConfig {
	return types.ParserConfig{
		MaxLineLength:    config.MaxLineLength,
		MaxParseErrors:   config.MaxParseErrors,
		Timeout:          config.Timeout.String(),
		EnableValidation: config.EnableValidation,
		EnableMetrics:    config.EnableMetrics,
	}
}

// ConvertFromTypesParserConfig converts types.ParserConfig to parsing.ParserConfig
func ConvertFromTypesParserConfig(config types.ParserConfig) ParserConfig {
	timeout, _ := time.ParseDuration(config.Timeout)
	return ParserConfig{
		MaxLineLength:    config.MaxLineLength,
		MaxParseErrors:   config.MaxParseErrors,
		Timeout:          timeout,
		EnableValidation: config.EnableValidation,
		EnableMetrics:    config.EnableMetrics,
	}
}

// ConvertToTypesEnhancedAuditResult converts parsing results to types.EnhancedAuditResult
func ConvertToTypesEnhancedAuditResult(queryID, timestamp, command, rawOutput, summary string, executionTime int64, parseResult ParseResult) types.EnhancedAuditResult {
	entries := make([]types.AuditLogEntry, len(parseResult.Entries))
	for i, entry := range parseResult.Entries {
		entries[i] = ConvertToTypesAuditLogEntry(entry)
	}

	return types.EnhancedAuditResult{
		QueryID:       queryID,
		Timestamp:     timestamp,
		Command:       command,
		RawOutput:     rawOutput,
		ParsedEntries: entries,
		ParseResult:   ConvertToTypesParseResult(parseResult),
		Summary:       summary,
		ExecutionTime: executionTime,
	}
}
