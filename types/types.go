package types

// AuditQueryParams represents the structured parameters for audit queries
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

// AuditResult represents the parsed audit query result
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

// AuditLogEntry represents a structured audit log entry (interface version)
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

// ParseResult represents the result of parsing audit logs (interface version)
type ParseResult struct {
	Entries     []AuditLogEntry  `json:"entries"`
	TotalLines  int              `json:"total_lines"`
	ParsedLines int              `json:"parsed_lines"`
	ErrorLines  int              `json:"error_lines"`
	ParseErrors []string         `json:"parse_errors"`
	ParseTime   string           `json:"parse_time"`
	Performance ParsePerformance `json:"performance"`
}

// ParsePerformance tracks parsing performance metrics (interface version)
type ParsePerformance struct {
	LinesPerSecond  float64 `json:"lines_per_second"`
	AverageLineSize int     `json:"average_line_size"`
	MemoryUsage     int64   `json:"memory_usage_bytes"`
}

// ParserConfig holds configuration for the parser (interface version)
type ParserConfig struct {
	MaxLineLength    int    `json:"max_line_length"`
	MaxParseErrors   int    `json:"max_parse_errors"`
	Timeout          string `json:"timeout"`
	EnableValidation bool   `json:"enable_validation"`
	EnableMetrics    bool   `json:"enable_metrics"`
}

// EnhancedAuditResult represents the enhanced audit query result with structured parsing
type EnhancedAuditResult struct {
	QueryID       string          `json:"query_id"`
	Timestamp     string          `json:"timestamp"`
	Command       string          `json:"command"`
	RawOutput     string          `json:"raw_output"`
	ParsedEntries []AuditLogEntry `json:"parsed_entries"`
	ParseResult   ParseResult     `json:"parse_result"`
	Summary       string          `json:"summary"`
	Error         string          `json:"error,omitempty"`
	ExecutionTime int64           `json:"execution_time_ms"`
}

// ParserConfiguration represents the configuration for the enhanced parser
type ParserConfiguration struct {
	Config        ParserConfig `json:"config"`
	EnableLegacy  bool         `json:"enable_legacy"`
	EnableMetrics bool         `json:"enable_metrics"`
	MaxBatchSize  int          `json:"max_batch_size"`
}

// DefaultParserConfiguration returns the default parser configuration
func DefaultParserConfiguration() ParserConfiguration {
	return ParserConfiguration{
		Config: ParserConfig{
			MaxLineLength:    100000, // 100KB
			MaxParseErrors:   1000,
			Timeout:          "30s",
			EnableValidation: true,
			EnableMetrics:    true,
		},
		EnableLegacy:  true,
		EnableMetrics: true,
		MaxBatchSize:  10000,
	}
}
