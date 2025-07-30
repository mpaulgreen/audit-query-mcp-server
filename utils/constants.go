package utils

// Valid log sources for OpenShift audit logs
var ValidLogSources = []string{
	"kube-apiserver",
	"oauth-server",
	"node",
	"openshift-apiserver",
	"oauth-apiserver",
}

// Valid Kubernetes/OpenShift resources
var ValidResources = []string{
	// Core Kubernetes Resources
	"pods", "services", "deployments", "replicasets", "statefulsets", "daemonsets",
	"namespaces", "nodes", "configmaps", "secrets", "persistentvolumes", "persistentvolumeclaims",
	"endpoints", "events", "limitranges", "resourcequotas", "serviceaccounts",

	// RBAC Resources
	"roles", "rolebindings", "clusterroles", "clusterrolebindings",

	// Network Resources
	"networkpolicies", "ingresses", "ingressclasses", "services",

	// Storage Resources
	"storageclasses", "persistentvolumes", "persistentvolumeclaims", "volumeattachments",

	// Custom Resources
	"customresourcedefinitions", "customresourcedefinition",

	// OpenShift Specific Resources
	"projects", "builds", "buildconfigs", "deploymentconfigs", "routes", "imagestreams",
	"imagestreamtags", "imagestreamimages", "templates", "templateinstances",
	"securitycontextconstraints", "groups", "identities", "oauthclients",
	"oauthaccesstokens", "oauthauthorizetokens", "oauthclientauthorizations",
	"clusternetworks", "hostsubnets", "netnamespaces", "egressnetworkpolicies",
	"clusterresourcequotas", "appliedclusterresourcequotas", "resourceaccessreviews",
	"localresourceaccessreviews", "subjectaccessreviews", "localsubjectaccessreviews",
	"selfsubjectaccessreviews", "selfsubjectrulesreviews", "subjectrulesreviews",
	"clusterroles", "clusterrolebindings", "roles", "rolebindings",

	// Monitoring and Metrics
	"prometheusrules", "servicemonitors", "podmonitors", "alertmanagers",

	// Security Resources
	"podsecuritypolicies", "poddisruptionbudgets", "validatingwebhookconfigurations",
	"mutatingwebhookconfigurations", "certificatesigningrequests",

	// API Resources
	"apiservices", "flowschemas", "prioritylevelconfigurations",

	// Short forms and aliases
	"pod", "service", "deployment", "replicaset", "statefulset", "daemonset",
	"namespace", "node", "configmap", "secret", "pv", "pvc", "endpoint", "event",
	"limitrange", "resourcequota", "serviceaccount", "role", "rolebinding",
	"clusterrole", "clusterrolebinding", "networkpolicy", "ingress", "ingressclass",
	"storageclass", "volumeattachment", "crd", "project", "build", "buildconfig",
	"deploymentconfig", "route", "imagestream", "imagestreamtag", "imagestreamimage",
	"template", "templateinstance", "scc", "group", "identity", "oauthclient",
	"oauthaccesstoken", "oauthauthorizetoken", "oauthclientauthorization",
	"clusternetwork", "hostsubnet", "netnamespace", "egressnetworkpolicy",
	"clusterresourcequota", "appliedclusterresourcequota", "resourceaccessreview",
	"localresourceaccessreview", "subjectaccessreview", "localsubjectaccessreview",
	"selfsubjectaccessreview", "selfsubjectrulesreview", "subjectrulesreview",
	"prometheusrule", "servicemonitor", "podmonitor", "alertmanager",
	"podsecuritypolicy", "poddisruptionbudget", "validatingwebhookconfiguration",
	"mutatingwebhookconfiguration", "certificatesigningrequest", "apiservice",
	"flowschema", "prioritylevelconfiguration",
}

// Valid Kubernetes API verbs
var ValidVerbs = []string{
	// Standard CRUD operations
	"create", "get", "list", "watch", "update", "patch", "delete", "deletecollection",

	// Additional operations
	"apply", "replace", "patch", "connect", "proxy", "redirect", "head", "options",

	// OpenShift specific verbs
	"impersonate", "escalate", "bind", "approve", "deny", "escalate", "impersonate",

	// Custom resource operations
	"custom", "scale", "rollback", "restart", "pause", "resume", "attach", "detach",

	// Short forms and aliases
	"created", "updated", "deleted", "patched", "applied", "replaced", "connected",
	"proxied", "redirected", "impersonated", "escalated", "bound", "approved", "denied",
	"scaled", "rolledback", "restarted", "paused", "resumed", "attached", "detached",
}

// HTTP Response Status Codes for audit log analysis
var ResponseStatusCodes = map[string]int{
	// Success codes
	"OK":        200,
	"Created":   201,
	"NoContent": 204,

	// Client error codes
	"BadRequest":          400,
	"Unauthorized":        401,
	"Forbidden":           403,
	"NotFound":            404,
	"Conflict":            409,
	"UnprocessableEntity": 422,

	// Server error codes
	"InternalServerError": 500,
	"BadGateway":          502,
	"ServiceUnavailable":  503,
}

// Common HTTP status code ranges for filtering
var StatusCodeRanges = map[string][]int{
	"success":      {200, 201, 202, 204},
	"client_error": {400, 401, 403, 404, 409, 422, 429},
	"server_error": {500, 502, 503, 504},
	"auth_error":   {401, 403},
	"not_found":    {404},
	"conflict":     {409},
}

// Audit Log Field Names - JSON fields commonly found in OpenShift audit logs
var AuditLogFields = map[string]string{
	// Core fields
	"RequestReceivedTimestamp": "requestReceivedTimestamp",
	"ResponseStatus":           "responseStatus",
	"ObjectRef":                "objectRef",
	"User":                     "user",
	"Verb":                     "verb",
	"RequestURI":               "requestURI",
	"SourceIPs":                "sourceIPs",
	"UserAgent":                "userAgent",
	"Annotations":              "annotations",

	// User-related fields
	"Username": "username",
	"UID":      "uid",
	"Groups":   "groups",
	"Extra":    "extra",

	// Object reference fields
	"Resource":   "resource",
	"Namespace":  "namespace",
	"Name":       "name",
	"APIGroup":   "apiGroup",
	"APIVersion": "apiVersion",

	// Response status fields
	"Code":    "code",
	"Message": "message",
	"Reason":  "reason",

	// Request fields
	"Method":  "method",
	"Path":    "path",
	"Query":   "query",
	"Headers": "headers",

	// Authentication fields
	"AuthenticationDecision": "authentication.openshift.io/decision",
	"AuthorizationDecision":  "authorization.k8s.io/decision",
	"ImpersonatedUser":       "impersonatedUser",
	"RequestUser":            "requestUser",
}

// Time Frame Constants for consistent timeframe handling
var TimeFrameConstants = map[string]string{
	// Common timeframes
	"Today":       "today",
	"Yesterday":   "yesterday",
	"ThisWeek":    "this week",
	"LastHour":    "last hour",
	"Last24Hours": "24h",
	"Last7Days":   "7d",
	"LastWeek":    "last week",
	"ThisMonth":   "this month",
	"LastMonth":   "last month",
	"Last30Days":  "last 30 days",

	// Underscore variants for compatibility
	"Last_24_Hours": "last_24_hours",
	"Last_7_Days":   "last_7_days",
	"Last_Week":     "last_week",
	"Last_Month":    "last_month",
	"Last_30_Days":  "last_30_days",
	"Last_Hour":     "last_hour",

	// Short forms
	"1m":  "1m",
	"5m":  "5m",
	"15m": "15m",
	"30m": "30m",
	"1h":  "1h",
	"2h":  "2h",
	"6h":  "6h",
	"12h": "12h",
	"1d":  "1d",
	"2d":  "2d",
	"3d":  "3d",
	"1w":  "1w",
	"2w":  "2w",
	"1y":  "1y",
}

// Time Frame Patterns for validation and parsing
var TimeFramePatterns = []string{
	// Basic patterns
	"^today$",
	"^yesterday$",
	"^this week$",
	"^last hour$",
	"^last_?24_?hours?$",
	"^last_?7_?days?$",
	"^last_?week$",
	"^this month$",
	"^last_?month$",
	"^last_?30_?days?$",
	"^last_?hour$",

	// Dynamic patterns
	"^last_?\\d+_?minute(s)?$",
	"^last_?\\d+_?hour(s)?$",
	"^last_?\\d+_?day(s)?$",
	"^last_?\\d+_?week(s)?$",
	"^last_?\\d+_?month(s)?$",
	"^last_?\\d+_?year(s)?$",

	// Short form patterns
	"^\\d+[mhdwy]$",
	"^\\d+[mhdwy] ago$",

	// Date patterns
	"^since \\d{4}-\\d{2}-\\d{2}$",
	"^since \\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}$",
}

// Dangerous command patterns that should be blocked
var DangerousPatterns = []string{
	"oc delete",
	"oc apply",
	"oc create",
	"oc patch",
	"oc replace",
	"kubectl delete",
	"kubectl apply",
	"kubectl create",
	"kubectl patch",
	"kubectl replace",
	";",
	"&&",
	"||",
	"`",
}

// Safe date command patterns that are allowed
var SafeDatePatterns = []string{
	"$(date",
	"$(date -d",
	"$(date -v",
}

// System User Groups and Users for filtering and security analysis
var SystemUserGroups = map[string][]string{
	"system_authenticated": {
		"system:authenticated",
		"system:authenticated:oauth",
	},
	"system_service_accounts": {
		"system:serviceaccount",
		"system:serviceaccount:kube-system",
		"system:serviceaccount:openshift-*",
		"system:serviceaccount:default",
	},
	"system_admins": {
		"system:admin",
		"kube:admin",
		"system:masters",
	},
	"system_anonymous": {
		"system:anonymous",
		"system:unauthenticated",
	},
	"system_nodes": {
		"system:node",
		"system:node:*",
	},
	"system_controllers": {
		"system:controller",
		"system:controller:*",
	},
	"system_bootstrap": {
		"system:bootstrap",
		"system:bootstrap:*",
	},
}

// Common system users that should be excluded from human user analysis
var SystemUsers = []string{
	"system:admin",
	"kube:admin",
	"system:anonymous",
	"system:unauthenticated",
	"system:serviceaccount",
	"system:node",
	"system:controller",
	"system:bootstrap",
	"system:masters",
	"system:authenticated",
	"system:authenticated:oauth",
}

// Audit Log Levels/Profiles for OpenShift audit configuration
var AuditLogLevels = map[string]string{
	"None":               "None",
	"Metadata":           "Metadata",
	"Request":            "Request",
	"RequestResponse":    "RequestResponse",
	"Default":            "Default",
	"WriteRequestBodies": "WriteRequestBodies",
	"AllRequestBodies":   "AllRequestBodies",
}

// Common Security Patterns for threat detection and investigation
var SecurityPatterns = map[string][]string{
	"privilege_escalation": {
		"clusterrole",
		"clusterrolebinding",
		"role",
		"rolebinding",
		"securitycontextconstraints",
		"impersonate",
		"escalate",
		"bind",
	},
	"resource_deletion": {
		"delete",
		"deletecollection",
		"namespaces",
		"customresourcedefinitions",
		"clusterroles",
		"clusterrolebindings",
	},
	"authentication_failures": {
		"401",
		"403",
		"Unauthorized",
		"Forbidden",
		"authentication.openshift.io/decision\":\"error",
		"authorization.k8s.io/decision\":\"forbid",
	},
	"unusual_access": {
		"system:serviceaccount",
		"system:node",
		"system:controller",
		"system:bootstrap",
	},
	"after_hours": {
		"18:00:00",
		"19:00:00",
		"20:00:00",
		"21:00:00",
		"22:00:00",
		"23:00:00",
		"00:00:00",
		"01:00:00",
		"02:00:00",
		"03:00:00",
		"04:00:00",
		"05:00:00",
		"06:00:00",
		"07:00:00",
	},
	"weekend_access": {
		"Saturday",
		"Sunday",
		"Sat",
		"Sun",
	},
}

// Common Namespace Patterns for filtering and analysis
var NamespacePatterns = map[string][]string{
	"system_namespaces": {
		"kube-system",
		"kube-public",
		"kube-node-lease",
		"default",
	},
	"openshift_system_namespaces": {
		"openshift-*",
		"openshift-apiserver",
		"openshift-authentication",
		"openshift-authorization",
		"openshift-config",
		"openshift-console",
		"openshift-controller-manager",
		"openshift-dns",
		"openshift-etcd",
		"openshift-image-registry",
		"openshift-ingress",
		"openshift-kube-apiserver",
		"openshift-kube-controller-manager",
		"openshift-kube-scheduler",
		"openshift-machine-api",
		"openshift-machine-config-operator",
		"openshift-marketplace",
		"openshift-monitoring",
		"openshift-multus",
		"openshift-network-diagnostics",
		"openshift-network-operator",
		"openshift-node",
		"openshift-oauth-apiserver",
		"openshift-operator-lifecycle-manager",
		"openshift-operators",
		"openshift-service-ca",
		"openshift-service-catalog-apiserver",
		"openshift-service-catalog-controller-manager",
		"openshift-user-workload-monitoring",
	},
	"sensitive_namespaces": {
		"security",
		"compliance",
		"audit",
		"admin",
		"management",
		"infrastructure",
		"monitoring",
		"logging",
	},
}

// Error Messages and Codes for standardized error handling
var ErrorMessages = map[string]string{
	"invalid_log_source":  "Invalid log source specified",
	"invalid_resource":    "Invalid resource type specified",
	"invalid_verb":        "Invalid API verb specified",
	"invalid_timeframe":   "Invalid timeframe specified",
	"invalid_username":    "Invalid username pattern specified",
	"invalid_namespace":   "Invalid namespace pattern specified",
	"dangerous_command":   "Command contains dangerous patterns",
	"timeout_error":       "Query execution timed out",
	"permission_denied":   "Permission denied for audit log access",
	"cluster_unavailable": "OpenShift cluster is unavailable",
	"invalid_response":    "Invalid response from audit log query",
	"parsing_error":       "Error parsing audit log results",
	"cache_error":         "Cache operation failed",
	"validation_error":    "Input validation failed",
}

// Error Codes for programmatic error handling
var ErrorCodes = map[string]int{
	"VALIDATION_ERROR":    400,
	"PERMISSION_DENIED":   403,
	"NOT_FOUND":           404,
	"TIMEOUT_ERROR":       408,
	"CLUSTER_UNAVAILABLE": 503,
	"PARSING_ERROR":       422,
	"CACHE_ERROR":         500,
	"UNKNOWN_ERROR":       500,
}

// File Path Patterns for audit log file handling
var FilePathPatterns = map[string]string{
	"kube_apiserver_log":      "kube-apiserver/audit.log",
	"oauth_server_log":        "oauth-server/audit.log",
	"node_log":                "audit/audit.log",
	"openshift_apiserver_log": "openshift-apiserver/audit.log",
	"oauth_apiserver_log":     "oauth-apiserver/audit.log",
	"audit_log_dir":           "/var/log/audit/",
	"node_logs_dir":           "/var/log/",
}

// Log File Naming Conventions
var LogFilePatterns = []string{
	"audit.log",
	"audit.log.*",
	"audit.log.*.gz",
	"audit.log.*.bz2",
	"audit-*.log",
	"audit-*.log.*",
	"audit-*.log.*.gz",
	"audit-*.log.*.bz2",
}

// File Rotation Patterns
var FileRotationPatterns = []string{
	"audit.log.1",
	"audit.log.2",
	"audit.log.3",
	"audit.log.4",
	"audit.log.5",
	"audit.log.6",
	"audit.log.7",
	"audit.log.8",
	"audit.log.9",
	"audit.log.10",
}

// Cache Configuration Constants for performance optimization
var CacheConfig = map[string]interface{}{
	"default_ttl":           3600,    // 1 hour in seconds
	"max_ttl":               86400,   // 24 hours in seconds
	"min_ttl":               60,      // 1 minute in seconds
	"max_cache_size":        1000,    // Maximum number of cached entries
	"cleanup_interval":      300,     // 5 minutes in seconds
	"eviction_policy":       "lru",   // Least Recently Used
	"compression_threshold": 1024,    // Compress entries larger than 1KB
	"max_entry_size":        1048576, // 1MB in bytes
}

// Command Execution Limits for safety and performance
var ExecutionLimits = map[string]interface{}{
	"max_command_length":    10000,    // Maximum command length in characters
	"max_execution_time":    300,      // Maximum execution time in seconds
	"max_output_size":       10485760, // 10MB in bytes
	"max_result_entries":    100000,   // Maximum number of result entries
	"rate_limit_per_minute": 60,       // Maximum queries per minute
	"rate_limit_per_hour":   1000,     // Maximum queries per hour
	"concurrent_queries":    10,       // Maximum concurrent queries
	"timeout_buffer":        30,       // Additional timeout buffer in seconds
}

// Timeout Values for different operations
var TimeoutValues = map[string]int{
	"query_generation":   5,   // 5 seconds
	"query_execution":    300, // 5 minutes
	"query_parsing":      30,  // 30 seconds
	"cache_operation":    1,   // 1 second
	"validation":         2,   // 2 seconds
	"file_operation":     10,  // 10 seconds
	"network_operation":  30,  // 30 seconds
	"database_operation": 60,  // 1 minute
}

// Rate Limiting Configuration
var RateLimitConfig = map[string]interface{}{
	"enabled":                 true,
	"window_size":             60, // 1 minute window
	"max_requests_per_window": 60, // 60 requests per minute
	"burst_limit":             10, // Allow burst of 10 requests
	"retry_after":             60, // Retry after 60 seconds when rate limited
}

// Performance Monitoring Thresholds
var PerformanceThresholds = map[string]interface{}{
	"slow_query_threshold":     10, // Queries taking >10 seconds are slow
	"memory_usage_threshold":   80, // Alert when memory usage >80%
	"cpu_usage_threshold":      70, // Alert when CPU usage >70%
	"cache_hit_rate_threshold": 50, // Alert when cache hit rate <50%
	"error_rate_threshold":     5,  // Alert when error rate >5%
}

// Logging Configuration
var LoggingConfig = map[string]interface{}{
	"log_level":          "info",
	"log_format":         "json",
	"log_file":           "./logs/audit_query.log",
	"max_log_size":       10485760, // 10MB
	"max_log_files":      5,
	"log_retention_days": 30,
	"enable_audit_trail": true,
	"audit_trail_file":   "./logs/audit_trail.json",
}

// Security Configuration
var SecurityConfig = map[string]interface{}{
	"enable_input_sanitization": true,
	"enable_command_validation": true,
	"enable_rate_limiting":      true,
	"enable_audit_logging":      true,
	"max_input_length":          1000,
	"allowed_command_patterns": []string{
		"oc adm node-logs",
		"grep",
		"jq",
		"date",
	},
	"blocked_command_patterns": []string{
		"rm",
		"delete",
		"drop",
		"truncate",
		"format",
		"mkfs",
		"dd",
		"shred",
	},
}

// Database Configuration (for future use)
var DatabaseConfig = map[string]interface{}{
	"connection_string":      "",
	"max_connections":        10,
	"connection_timeout":     30,
	"query_timeout":          300,
	"enable_connection_pool": true,
	"pool_size":              5,
	"pool_max_idle":          2,
	"pool_max_lifetime":      3600,
}

// API Configuration
var APIConfig = map[string]interface{}{
	"port":              8080,
	"host":              "localhost",
	"read_timeout":      30,
	"write_timeout":     30,
	"idle_timeout":      60,
	"max_request_size":  1048576, // 1MB
	"enable_cors":       true,
	"cors_origins":      []string{"*"},
	"enable_metrics":    true,
	"metrics_path":      "/metrics",
	"health_check_path": "/health",
}

// Metrics Configuration
var MetricsConfig = map[string]interface{}{
	"enabled":           true,
	"port":              9090,
	"path":              "/metrics",
	"collect_interval":  15,    // 15 seconds
	"retention_period":  86400, // 24 hours
	"enable_histograms": true,
	"enable_summaries":  true,
	"enable_counters":   true,
	"enable_gauges":     true,
}

// Alerting Configuration
var AlertingConfig = map[string]interface{}{
	"enabled":               false,
	"alert_webhook_url":     "",
	"alert_webhook_timeout": 10,
	"alert_retry_count":     3,
	"alert_retry_delay":     5,
	"critical_threshold":    90,
	"warning_threshold":     70,
	"info_threshold":        50,
}

// Backup and Recovery Configuration
var BackupConfig = map[string]interface{}{
	"enabled":               false,
	"backup_interval":       86400, // 24 hours
	"backup_retention_days": 7,
	"backup_path":           "./backups",
	"backup_compression":    true,
	"backup_encryption":     false,
	"restore_enabled":       false,
}

// Maintenance Configuration
var MaintenanceConfig = map[string]interface{}{
	"maintenance_window_start": "02:00",
	"maintenance_window_end":   "04:00",
	"maintenance_timezone":     "UTC",
	"enable_auto_cleanup":      true,
	"cleanup_interval":         3600, // 1 hour
	"max_log_age_days":         30,
	"max_cache_age_hours":      24,
	"enable_health_checks":     true,
	"health_check_interval":    300, // 5 minutes
}
