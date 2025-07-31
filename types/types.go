package types

import "time"

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

// AuditQueryConfig represents configuration for audit queries
type AuditQueryConfig struct {
	MaxRotatedFiles      int           `json:"max_rotated_files" default:"3"`
	CommandTimeout       time.Duration `json:"command_timeout" default:"30s"`
	UseJSONParsing       bool          `json:"use_json_parsing" default:"false"`
	EnableCompression    bool          `json:"enable_compression" default:"false"`
	ParallelProcessing   bool          `json:"parallel_processing" default:"false"`
	ForceSimple          bool          `json:"force_simple" default:"true"` // New default for reliability
	EnableFileDiscovery  bool          `json:"enable_file_discovery" default:"false"`
	MaxConcurrentQueries int           `json:"max_concurrent_queries" default:"5"`
}

// EnvironmentInfo represents information about the OpenShift environment
type EnvironmentInfo struct {
	OpenShiftVersion string   `json:"openshift_version"`
	JQAvailable      bool     `json:"jq_available"`
	LogFormats       []string `json:"log_formats"`
	AvailableFiles   []string `json:"available_files"`
	MaxConcurrentOC  int      `json:"max_concurrent_oc"`
	PermissionLevel  string   `json:"permission_level"`
}

// MigrationConfig represents configuration for the migration strategy
type MigrationConfig struct {
	// Feature flags for gradual rollout
	EnableNewBuilder    bool `json:"enable_new_builder" default:"false"`
	EnableFileDiscovery bool `json:"enable_file_discovery" default:"false"`
	EnableJSONParsing   bool `json:"enable_json_parsing" default:"false"`
	EnableParallelProc  bool `json:"enable_parallel_proc" default:"false"`

	// Safety limits
	MaxFiles             int           `json:"max_files" default:"3"`
	CommandTimeout       time.Duration `json:"command_timeout" default:"30s"`
	MaxConcurrentQueries int           `json:"max_concurrent_queries" default:"5"`

	// Backward compatibility
	PreserveOldBehavior bool `json:"preserve_old_behavior" default:"true"`
	FallbackOnError     bool `json:"fallback_on_error" default:"true"`
}

// CircuitBreaker represents a circuit breaker pattern for command execution
type CircuitBreaker struct {
	FailureThreshold int           `json:"failure_threshold"`
	ResetTimeout     time.Duration `json:"reset_timeout"`
	State            CircuitState  `json:"state"`
	FailureCount     int           `json:"failure_count"`
	LastFailureTime  time.Time     `json:"last_failure_time"`
}

// CircuitState represents the state of a circuit breaker
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"
	CircuitStateOpen     CircuitState = "open"
	CircuitStateHalfOpen CircuitState = "half_open"
)

// LogFileInfo represents information about a log file
type LogFileInfo struct {
	Path      string    `json:"path"`
	Date      time.Time `json:"date"`
	IsCurrent bool      `json:"is_current"`
	Exists    bool      `json:"exists"`
	Size      int64     `json:"size"`
}

// FileDiscoveryCache represents a cache for file discovery results
type FileDiscoveryCache struct {
	Cache     map[string][]string `json:"cache"`
	LastCheck time.Time           `json:"last_check"`
	TTL       time.Duration       `json:"ttl"`
}

// DiscoveryConfig represents configuration for file discovery
type DiscoveryConfig struct {
	EnableDiscovery bool          `json:"enable_discovery" default:"false"`
	CacheTTL        time.Duration `json:"cache_ttl" default:"5m"`
	MaxFilesToCheck int           `json:"max_files_to_check" default:"5"`
	FallbackFiles   []string      `json:"fallback_files"`
}

// DefaultAuditQueryConfig returns the default audit query configuration
func DefaultAuditQueryConfig() AuditQueryConfig {
	return AuditQueryConfig{
		MaxRotatedFiles:      3,
		CommandTimeout:       30 * time.Second,
		UseJSONParsing:       true, // Phase 2: Enable JSON parsing by default
		EnableCompression:    false,
		ParallelProcessing:   false,
		ForceSimple:          true, // Default to simple for reliability
		EnableFileDiscovery:  false,
		MaxConcurrentQueries: 5,
	}
}

// DefaultMigrationConfig returns the default migration configuration
func DefaultMigrationConfig() MigrationConfig {
	return MigrationConfig{
		EnableNewBuilder:     false,
		EnableFileDiscovery:  false,
		EnableJSONParsing:    false,
		EnableParallelProc:   false,
		MaxFiles:             3,
		CommandTimeout:       30 * time.Second,
		MaxConcurrentQueries: 5,
		PreserveOldBehavior:  true,
		FallbackOnError:      true,
	}
}

// DefaultDiscoveryConfig returns the default discovery configuration
func DefaultDiscoveryConfig() DiscoveryConfig {
	return DiscoveryConfig{
		EnableDiscovery: false,
		CacheTTL:        5 * time.Minute,
		MaxFilesToCheck: 5,
		FallbackFiles:   []string{"audit.log", "audit.log.1", "audit.log.2"},
	}
}
