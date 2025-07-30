package validation

import "audit-query-mcp-server/utils"

// =============================================================================
// CORE PATTERN CONSTANTS - Single source of truth for basic patterns
// =============================================================================

// DNSLabelPattern represents the standard DNS label format used in Kubernetes
// RFC 1123 compliant: lowercase alphanumeric with hyphens, no consecutive hyphens
const DNSLabelPattern = `^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`

// UsernameBasePattern represents basic alphanumeric usernames
const UsernameBasePattern = `^[a-zA-Z0-9._-]+$`

// EmailPattern represents standard email addresses
const EmailPattern = `^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

// IPv4Pattern represents valid IPv4 addresses
const IPv4Pattern = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

// HTTPStatusCodePattern represents valid HTTP status codes (100-599)
const HTTPStatusCodePattern = `^[1-5]\d{2}$`

// =============================================================================
// CONSOLIDATED PATTERN VARIABLES - No duplication
// =============================================================================

// ValidLogSources contains all supported OpenShift audit log sources
var ValidLogSources = utils.ValidLogSources

// ValidResources contains all supported Kubernetes and OpenShift resources
var ValidResources = utils.ValidResources

// ValidVerbs contains all supported Kubernetes API verbs
var ValidVerbs = utils.ValidVerbs

// TimeframePatterns contains regex patterns for timeframe validation
var TimeframePatterns = utils.TimeFramePatterns

// NamespacePatterns contains the single, consolidated namespace pattern
var NamespacePatterns = []string{
	DNSLabelPattern, // Standard Kubernetes/OpenShift namespace pattern (1-63 chars)
}

// ResourceNamePatterns contains the single, consolidated resource name pattern
var ResourceNamePatterns = []string{
	DNSLabelPattern, // Standard Kubernetes resource name pattern (1-253 chars)
}

// UsernamePatterns contains consolidated username validation patterns
var UsernamePatterns = []string{
	UsernameBasePattern, // Basic alphanumeric usernames
	EmailPattern,        // Email addresses (mapul@redhat.com)
	`^system:serviceaccount:[a-zA-Z0-9._-]+:[a-zA-Z0-9._-]+$`, // Service accounts
	`^system:node:[a-zA-Z0-9._-]+$`,                           // Node usernames
	`^system:admin$`,                                          // System admin
	`^kube:admin$`,                                            // Kube admin
	`^system:anonymous$`,                                      // Anonymous user
	`^system:unauthenticated$`,                                // Unauthenticated user
	`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+$`,                        // Domain usernames (user@domain)
	`^[a-zA-Z0-9._-]+\\[a-zA-Z0-9._-]+$`,                      // Windows domain users (DOMAIN\user)
	`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`,                       // LDAP users (cn=user,dc=example,dc=com)
	`^[a-zA-Z0-9._-]+:[a-zA-Z0-9._-]+$`,                       // OAuth provider users (provider:user)
}

// StatusCodePatterns contains consolidated HTTP status code patterns
var StatusCodePatterns = []string{
	HTTPStatusCodePattern,             // Valid HTTP status codes (100-599)
	`^2\d{2}$`,                        // Success codes (200-299)
	`^4\d{2}$`,                        // Client error codes (400-499)
	`^5\d{2}$`,                        // Server error codes (500-599)
	`^(200|201|202|204)$`,             // Common success codes
	`^(400|401|403|404|409|422|429)$`, // Common client error codes
	`^(500|502|503|504)$`,             // Common server error codes
}

// IPAddressPatterns contains consolidated IP address patterns
var IPAddressPatterns = []string{
	IPv4Pattern, // IPv4
	`^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`, // IPv6 (simplified)
	`^127\.0\.0\.1$`,                 // Localhost
	`^10\.`,                          // Private network (10.x.x.x)
	`^172\.(1[6-9]|2[0-9]|3[0-1])\.`, // Private network (172.16-31.x.x)
	`^192\.168\.`,                    // Private network (192.168.x.x)
}

// APIGroupPatterns contains consolidated Kubernetes API group patterns
var APIGroupPatterns = []string{
	`^$`, // Core group (empty string)
	`^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)*$`, // API group format
	`^apps$`,                                    // Apps API group
	`^batch$`,                                   // Batch API group
	`^extensions$`,                              // Extensions API group
	`^networking\.k8s\.io$`,                     // Networking API group
	`^rbac\.authorization\.k8s\.io$`,            // RBAC API group
	`^storage\.k8s\.io$`,                        // Storage API group
	`^autoscaling$`,                             // Autoscaling API group
	`^policy$`,                                  // Policy API group
	`^scheduling\.k8s\.io$`,                     // Scheduling API group
	`^coordination\.k8s\.io$`,                   // Coordination API group
	`^node\.k8s\.io$`,                           // Node API group
	`^discovery\.k8s\.io$`,                      // Discovery API group
	`^flowcontrol\.apiserver\.k8s\.io$`,         // Flow control API group
	`^admissionregistration\.k8s\.io$`,          // Admission registration API group
	`^apiextensions\.k8s\.io$`,                  // API extensions API group
	`^apiregistration\.k8s\.io$`,                // API registration API group
	`^config\.openshift\.io$`,                   // OpenShift config API group
	`^operator\.openshift\.io$`,                 // OpenShift operator API group
	`^user\.openshift\.io$`,                     // OpenShift user API group
	`^oauth\.openshift\.io$`,                    // OpenShift OAuth API group
	`^project\.openshift\.io$`,                  // OpenShift project API group
	`^route\.openshift\.io$`,                    // OpenShift route API group
	`^template\.openshift\.io$`,                 // OpenShift template API group
	`^build\.openshift\.io$`,                    // OpenShift build API group
	`^image\.openshift\.io$`,                    // OpenShift image API group
	`^deploymentconfig\.apps\.openshift\.io$`,   // OpenShift deployment config API group
	`^security\.openshift\.io$`,                 // OpenShift security API group
	`^quota\.openshift\.io$`,                    // OpenShift quota API group
	`^authorization\.openshift\.io$`,            // OpenShift authorization API group
	`^network\.openshift\.io$`,                  // OpenShift network API group
	`^monitoring\.coreos\.com$`,                 // Prometheus monitoring API group
	`^alertmanager\.monitoring\.coreos\.com$`,   // Alertmanager API group
	`^prometheus\.monitoring\.coreos\.com$`,     // Prometheus API group
	`^servicemonitor\.monitoring\.coreos\.com$`, // ServiceMonitor API group
	`^podmonitor\.monitoring\.coreos\.com$`,     // PodMonitor API group
	`^prometheusrule\.monitoring\.coreos\.com$`, // PrometheusRule API group
	`^thanosruler\.monitoring\.coreos\.com$`,    // ThanosRuler API group
}

// APIVersionPatterns contains consolidated Kubernetes API version patterns
var APIVersionPatterns = []string{
	`^v1$`,                                         // Core v1
	`^v1beta1$`,                                    // Beta v1
	`^v1alpha1$`,                                   // Alpha v1
	`^v2$`,                                         // Core v2
	`^v2beta1$`,                                    // Beta v2
	`^v2alpha1$`,                                   // Alpha v2
	`^v3$`,                                         // Core v3
	`^v3beta1$`,                                    // Beta v3
	`^v3alpha1$`,                                   // Alpha v3
	`^apps/v1$`,                                    // Apps v1
	`^apps/v1beta1$`,                               // Apps beta v1
	`^batch/v1$`,                                   // Batch v1
	`^batch/v1beta1$`,                              // Batch beta v1
	`^extensions/v1beta1$`,                         // Extensions beta v1
	`^networking\.k8s\.io/v1$`,                     // Networking v1
	`^rbac\.authorization\.k8s\.io/v1$`,            // RBAC v1
	`^storage\.k8s\.io/v1$`,                        // Storage v1
	`^autoscaling/v1$`,                             // Autoscaling v1
	`^policy/v1$`,                                  // Policy v1
	`^scheduling\.k8s\.io/v1$`,                     // Scheduling v1
	`^coordination\.k8s\.io/v1$`,                   // Coordination v1
	`^node\.k8s\.io/v1$`,                           // Node v1
	`^discovery\.k8s\.io/v1$`,                      // Discovery v1
	`^flowcontrol\.apiserver\.k8s\.io/v1$`,         // Flow control v1
	`^admissionregistration\.k8s\.io/v1$`,          // Admission registration v1
	`^apiextensions\.k8s\.io/v1$`,                  // API extensions v1
	`^apiregistration\.k8s\.io/v1$`,                // API registration v1
	`^config\.openshift\.io/v1$`,                   // OpenShift config v1
	`^operator\.openshift\.io/v1$`,                 // OpenShift operator v1
	`^user\.openshift\.io/v1$`,                     // OpenShift user v1
	`^oauth\.openshift\.io/v1$`,                    // OpenShift OAuth v1
	`^project\.openshift\.io/v1$`,                  // OpenShift project v1
	`^route\.openshift\.io/v1$`,                    // OpenShift route v1
	`^template\.openshift\.io/v1$`,                 // OpenShift template v1
	`^build\.openshift\.io/v1$`,                    // OpenShift build v1
	`^image\.openshift\.io/v1$`,                    // OpenShift image v1
	`^deploymentconfig\.apps\.openshift\.io/v1$`,   // OpenShift deployment config v1
	`^security\.openshift\.io/v1$`,                 // OpenShift security v1
	`^quota\.openshift\.io/v1$`,                    // OpenShift quota v1
	`^authorization\.openshift\.io/v1$`,            // OpenShift authorization v1
	`^network\.openshift\.io/v1$`,                  // OpenShift network v1
	`^monitoring\.coreos\.com/v1$`,                 // Prometheus monitoring v1
	`^alertmanager\.monitoring\.coreos\.com/v1$`,   // Alertmanager v1
	`^prometheus\.monitoring\.coreos\.com/v1$`,     // Prometheus v1
	`^servicemonitor\.monitoring\.coreos\.com/v1$`, // ServiceMonitor v1
	`^podmonitor\.monitoring\.coreos\.com/v1$`,     // PodMonitor v1
	`^prometheusrule\.monitoring\.coreos\.com/v1$`, // PrometheusRule v1
	`^thanosruler\.monitoring\.coreos\.com/v1$`,    // ThanosRuler v1
}
