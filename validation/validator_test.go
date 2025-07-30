package validation

import (
	"testing"
	"time"

	"audit-query-mcp-server/types"
)

// TestValidateAuditResult tests the ValidateAuditResult function
func TestValidateAuditResult(t *testing.T) {
	tests := []struct {
		name    string
		result  types.AuditResult
		wantErr bool
	}{
		{
			name: "Valid complete result",
			result: types.AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{{"test": "data"}},
				Summary:       "Found 1 audit entries",
				Error:         "",
				ExecutionTime: 100,
			},
			wantErr: false,
		},
		{
			name: "Valid result with error",
			result: types.AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "command execution failed",
				ExecutionTime: 50,
			},
			wantErr: false,
		},
		{
			name: "Missing QueryID",
			result: types.AuditResult{
				QueryID:       "",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "Missing Timestamp",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     "",
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "Invalid Timestamp format",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     "invalid-timestamp",
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "Missing Command when no error",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "Command too long",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       string(make([]byte, 10001)),
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "RawOutput too large",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     string(make([]byte, 1000001)),
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "Nil ParsedData",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    nil,
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "ParsedData too large",
			result: types.AuditResult{
				QueryID:   "test_id",
				Timestamp: time.Now().Format(time.RFC3339),
				Command:   "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput: "test output",
				ParsedData: func() []map[string]interface{} {
					data := make([]map[string]interface{}, 100001)
					for i := range data {
						data[i] = map[string]interface{}{"id": i}
					}
					return data
				}(),
				Summary:       "Too many entries",
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "Summary too long",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       string(make([]byte, 10001)),
				Error:         "",
				ExecutionTime: 0,
			},
			wantErr: true,
		},
		{
			name: "Error too long",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         string(make([]byte, 5001)),
				ExecutionTime: 50,
			},
			wantErr: true,
		},
		{
			name: "Negative ExecutionTime",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: -100,
			},
			wantErr: true,
		},
		{
			name: "ExecutionTime too high",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "No entries found",
				Error:         "",
				ExecutionTime: 3600001,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuditResult(tt.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAuditResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateAuditResultStrict tests the ValidateAuditResultStrict function
func TestValidateAuditResultStrict(t *testing.T) {
	tests := []struct {
		name    string
		result  types.AuditResult
		wantErr bool
	}{
		{
			name: "Valid complete result",
			result: types.AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{{"test": "data"}},
				Summary:       "Found 1 audit entries",
				Error:         "",
				ExecutionTime: 100,
			},
			wantErr: false,
		},
		{
			name: "Valid result with error",
			result: types.AuditResult{
				QueryID:       "audit_query_20250729_204108_abc123",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "command execution failed",
				ExecutionTime: 50,
			},
			wantErr: false,
		},
		{
			name: "Error with RawOutput (non-timeout)",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "some output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "command execution failed",
				ExecutionTime: 50,
			},
			wantErr: true,
		},
		{
			name: "Error with RawOutput (timeout)",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "partial output",
				ParsedData:    []map[string]interface{}{},
				Summary:       "",
				Error:         "command execution timed out after 30 seconds",
				ExecutionTime: 30000,
			},
			wantErr: false,
		},
		{
			name: "Error with ParsedData",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "",
				ParsedData:    []map[string]interface{}{{"partial": "data"}},
				Summary:       "",
				Error:         "command execution failed",
				ExecutionTime: 50,
			},
			wantErr: true,
		},
		{
			name: "No error but missing Command",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{{"test": "data"}},
				Summary:       "Found 1 audit entries",
				Error:         "",
				ExecutionTime: 100,
			},
			wantErr: true,
		},
		{
			name: "No error, has ParsedData, missing Summary",
			result: types.AuditResult{
				QueryID:       "test_id",
				Timestamp:     time.Now().Format(time.RFC3339),
				Command:       "oc adm node-logs --role=master --path=kube-apiserver/audit.log",
				RawOutput:     "test output",
				ParsedData:    []map[string]interface{}{{"test": "data"}},
				Summary:       "",
				Error:         "",
				ExecutionTime: 100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuditResultStrict(tt.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAuditResultStrict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateStatusCode tests the ValidateStatusCode function
func TestValidateStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode string
		want       bool
	}{
		{
			name:       "Valid success code",
			statusCode: "200",
			want:       true,
		},
		{
			name:       "Valid created code",
			statusCode: "201",
			want:       true,
		},
		{
			name:       "Valid no content code",
			statusCode: "204",
			want:       true,
		},
		{
			name:       "Valid client error code",
			statusCode: "400",
			want:       true,
		},
		{
			name:       "Valid unauthorized code",
			statusCode: "401",
			want:       true,
		},
		{
			name:       "Valid forbidden code",
			statusCode: "403",
			want:       true,
		},
		{
			name:       "Valid not found code",
			statusCode: "404",
			want:       true,
		},
		{
			name:       "Valid server error code",
			statusCode: "500",
			want:       true,
		},
		{
			name:       "Invalid code too low",
			statusCode: "99",
			want:       false,
		},
		{
			name:       "Invalid code too high",
			statusCode: "600",
			want:       false,
		},
		{
			name:       "Invalid non-numeric",
			statusCode: "abc",
			want:       false,
		},
		{
			name:       "Empty string",
			statusCode: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateStatusCode(tt.statusCode)
			if got != tt.want {
				t.Errorf("ValidateStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateStatusCodeRange tests the ValidateStatusCodeRange function
func TestValidateStatusCodeRange(t *testing.T) {
	tests := []struct {
		name       string
		statusCode string
		rangeName  string
		want       bool
	}{
		{
			name:       "Success code in success range",
			statusCode: "200",
			rangeName:  "success",
			want:       true,
		},
		{
			name:       "Created code in success range",
			statusCode: "201",
			rangeName:  "success",
			want:       true,
		},
		{
			name:       "Client error in client_error range",
			statusCode: "400",
			rangeName:  "client_error",
			want:       true,
		},
		{
			name:       "Unauthorized in auth_error range",
			statusCode: "401",
			rangeName:  "auth_error",
			want:       true,
		},
		{
			name:       "Forbidden in auth_error range",
			statusCode: "403",
			rangeName:  "auth_error",
			want:       true,
		},
		{
			name:       "Not found in not_found range",
			statusCode: "404",
			rangeName:  "not_found",
			want:       true,
		},
		{
			name:       "Conflict in conflict range",
			statusCode: "409",
			rangeName:  "conflict",
			want:       true,
		},
		{
			name:       "Server error in server_error range",
			statusCode: "500",
			rangeName:  "server_error",
			want:       true,
		},
		{
			name:       "Success code not in client_error range",
			statusCode: "200",
			rangeName:  "client_error",
			want:       false,
		},
		{
			name:       "Invalid range name",
			statusCode: "200",
			rangeName:  "invalid_range",
			want:       false,
		},
		{
			name:       "Invalid status code",
			statusCode: "999",
			rangeName:  "success",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateStatusCodeRange(tt.statusCode, tt.rangeName)
			if got != tt.want {
				t.Errorf("ValidateStatusCodeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateIPAddress tests the ValidateIPAddress function
func TestValidateIPAddress(t *testing.T) {
	tests := []struct {
		name      string
		ipAddress string
		want      bool
	}{
		{
			name:      "Valid IPv4 address",
			ipAddress: "192.168.1.1",
			want:      true,
		},
		{
			name:      "Valid localhost",
			ipAddress: "127.0.0.1",
			want:      true,
		},
		{
			name:      "Valid private network 10.x.x.x",
			ipAddress: "10.0.0.1",
			want:      true,
		},
		{
			name:      "Valid private network 172.16.x.x",
			ipAddress: "172.16.0.1",
			want:      true,
		},
		{
			name:      "Valid private network 192.168.x.x",
			ipAddress: "192.168.0.1",
			want:      true,
		},
		{
			name:      "Invalid IP address",
			ipAddress: "256.256.256.256",
			want:      false,
		},
		{
			name:      "Invalid format - malformed IP",
			ipAddress: "999.999.999.999",
			want:      false,
		},
		{
			name:      "Empty string",
			ipAddress: "",
			want:      false,
		},
		{
			name:      "Non-IP string",
			ipAddress: "not-an-ip",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateIPAddress(tt.ipAddress)
			if got != tt.want {
				t.Errorf("ValidateIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateResourceName tests the ValidateResourceName function
func TestValidateResourceName(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		want         bool
	}{
		{
			name:         "Valid resource name",
			resourceName: "my-resource",
			want:         true,
		},
		{
			name:         "Valid resource name with numbers",
			resourceName: "resource-123",
			want:         true,
		},
		{
			name:         "Valid single character",
			resourceName: "a",
			want:         true,
		},
		{
			name:         "Valid long name",
			resourceName: "my-very-long-resource-name-that-is-still-valid",
			want:         true,
		},
		{
			name:         "Empty string",
			resourceName: "",
			want:         false,
		},
		{
			name:         "Too long",
			resourceName: string(make([]byte, 254)),
			want:         false,
		},
		{
			name:         "Starts with hyphen",
			resourceName: "-resource",
			want:         false,
		},
		{
			name:         "Ends with hyphen",
			resourceName: "resource-",
			want:         false,
		},
		{
			name:         "Contains uppercase",
			resourceName: "Resource",
			want:         false,
		},
		{
			name:         "Contains special characters",
			resourceName: "resource@name",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateResourceName(tt.resourceName)
			if got != tt.want {
				t.Errorf("ValidateResourceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateAPIGroup tests the ValidateAPIGroup function
func TestValidateAPIGroup(t *testing.T) {
	tests := []struct {
		name     string
		apiGroup string
		want     bool
	}{
		{
			name:     "Empty core group",
			apiGroup: "",
			want:     true,
		},
		{
			name:     "Valid apps group",
			apiGroup: "apps",
			want:     true,
		},
		{
			name:     "Valid batch group",
			apiGroup: "batch",
			want:     true,
		},
		{
			name:     "Valid networking group",
			apiGroup: "networking.k8s.io",
			want:     true,
		},
		{
			name:     "Valid RBAC group",
			apiGroup: "rbac.authorization.k8s.io",
			want:     true,
		},
		{
			name:     "Valid OpenShift config group",
			apiGroup: "config.openshift.io",
			want:     true,
		},
		{
			name:     "Valid OpenShift user group",
			apiGroup: "user.openshift.io",
			want:     true,
		},
		{
			name:     "Invalid format with uppercase",
			apiGroup: "Apps",
			want:     false,
		},
		{
			name:     "Invalid format with special characters",
			apiGroup: "apps@k8s.io",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAPIGroup(tt.apiGroup)
			if got != tt.want {
				t.Errorf("ValidateAPIGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateAPIVersion tests the ValidateAPIVersion function
func TestValidateAPIVersion(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		want       bool
	}{
		{
			name:       "Valid core v1",
			apiVersion: "v1",
			want:       true,
		},
		{
			name:       "Valid beta version",
			apiVersion: "v1beta1",
			want:       true,
		},
		{
			name:       "Valid alpha version",
			apiVersion: "v1alpha1",
			want:       true,
		},
		{
			name:       "Valid apps v1",
			apiVersion: "apps/v1",
			want:       true,
		},
		{
			name:       "Valid batch v1",
			apiVersion: "batch/v1",
			want:       true,
		},
		{
			name:       "Valid networking v1",
			apiVersion: "networking.k8s.io/v1",
			want:       true,
		},
		{
			name:       "Valid OpenShift config v1",
			apiVersion: "config.openshift.io/v1",
			want:       true,
		},
		{
			name:       "Invalid format",
			apiVersion: "v1.0",
			want:       false,
		},
		{
			name:       "Invalid with special characters",
			apiVersion: "v1@beta1",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAPIVersion(tt.apiVersion)
			if got != tt.want {
				t.Errorf("ValidateAPIVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateAuditLogField tests the ValidateAuditLogField function
func TestValidateAuditLogField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		want      bool
	}{
		{
			name:      "Valid RequestReceivedTimestamp",
			fieldName: "requestReceivedTimestamp",
			want:      true,
		},
		{
			name:      "Valid ResponseStatus",
			fieldName: "responseStatus",
			want:      true,
		},
		{
			name:      "Valid ObjectRef",
			fieldName: "objectRef",
			want:      true,
		},
		{
			name:      "Valid User",
			fieldName: "user",
			want:      true,
		},
		{
			name:      "Valid Verb",
			fieldName: "verb",
			want:      true,
		},
		{
			name:      "Valid RequestURI",
			fieldName: "requestURI",
			want:      true,
		},
		{
			name:      "Valid SourceIPs",
			fieldName: "sourceIPs",
			want:      true,
		},
		{
			name:      "Valid UserAgent",
			fieldName: "userAgent",
			want:      true,
		},
		{
			name:      "Valid Annotations",
			fieldName: "annotations",
			want:      true,
		},
		{
			name:      "Valid Username",
			fieldName: "username",
			want:      true,
		},
		{
			name:      "Valid UID",
			fieldName: "uid",
			want:      true,
		},
		{
			name:      "Valid Groups",
			fieldName: "groups",
			want:      true,
		},
		{
			name:      "Valid Extra",
			fieldName: "extra",
			want:      true,
		},
		{
			name:      "Valid Resource",
			fieldName: "resource",
			want:      true,
		},
		{
			name:      "Valid Namespace",
			fieldName: "namespace",
			want:      true,
		},
		{
			name:      "Valid Name",
			fieldName: "name",
			want:      true,
		},
		{
			name:      "Valid APIGroup",
			fieldName: "apiGroup",
			want:      true,
		},
		{
			name:      "Valid APIVersion",
			fieldName: "apiVersion",
			want:      true,
		},
		{
			name:      "Valid Code",
			fieldName: "code",
			want:      true,
		},
		{
			name:      "Valid Message",
			fieldName: "message",
			want:      true,
		},
		{
			name:      "Valid Reason",
			fieldName: "reason",
			want:      true,
		},
		{
			name:      "Valid Method",
			fieldName: "method",
			want:      true,
		},
		{
			name:      "Valid Path",
			fieldName: "path",
			want:      true,
		},
		{
			name:      "Valid Query",
			fieldName: "query",
			want:      true,
		},
		{
			name:      "Valid Headers",
			fieldName: "headers",
			want:      true,
		},
		{
			name:      "Valid AuthenticationDecision",
			fieldName: "authentication.openshift.io/decision",
			want:      true,
		},
		{
			name:      "Valid AuthorizationDecision",
			fieldName: "authorization.k8s.io/decision",
			want:      true,
		},
		{
			name:      "Valid ImpersonatedUser",
			fieldName: "impersonatedUser",
			want:      true,
		},
		{
			name:      "Valid RequestUser",
			fieldName: "requestUser",
			want:      true,
		},
		{
			name:      "Invalid field name",
			fieldName: "InvalidField",
			want:      false,
		},
		{
			name:      "Empty field name",
			fieldName: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAuditLogField(tt.fieldName)
			if got != tt.want {
				t.Errorf("ValidateAuditLogField() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateTimeFrameConstant tests the ValidateTimeFrameConstant function
func TestValidateTimeFrameConstant(t *testing.T) {
	tests := []struct {
		name      string
		timeframe string
		want      bool
	}{
		{
			name:      "Valid Today",
			timeframe: "today",
			want:      true,
		},
		{
			name:      "Valid Yesterday",
			timeframe: "yesterday",
			want:      true,
		},
		{
			name:      "Valid ThisWeek",
			timeframe: "this week",
			want:      true,
		},
		{
			name:      "Valid LastHour",
			timeframe: "last hour",
			want:      true,
		},
		{
			name:      "Valid Last24Hours",
			timeframe: "24h",
			want:      true,
		},
		{
			name:      "Valid Last7Days",
			timeframe: "7d",
			want:      true,
		},
		{
			name:      "Valid LastWeek",
			timeframe: "last week",
			want:      true,
		},
		{
			name:      "Valid ThisMonth",
			timeframe: "this month",
			want:      true,
		},
		{
			name:      "Valid LastMonth",
			timeframe: "last month",
			want:      true,
		},
		{
			name:      "Valid Last30Days",
			timeframe: "last 30 days",
			want:      true,
		},
		{
			name:      "Valid short form 1m",
			timeframe: "1m",
			want:      true,
		},
		{
			name:      "Valid short form 5m",
			timeframe: "5m",
			want:      true,
		},
		{
			name:      "Valid short form 1h",
			timeframe: "1h",
			want:      true,
		},
		{
			name:      "Valid short form 1d",
			timeframe: "1d",
			want:      true,
		},
		{
			name:      "Valid short form 1w",
			timeframe: "1w",
			want:      true,
		},
		{
			name:      "Valid short form 1y",
			timeframe: "1y",
			want:      true,
		},
		{
			name:      "Invalid timeframe",
			timeframe: "invalid",
			want:      false,
		},
		{
			name:      "Empty timeframe",
			timeframe: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateTimeFrameConstant(tt.timeframe)
			if got != tt.want {
				t.Errorf("ValidateTimeFrameConstant() = %v, want %v", got, tt.want)
			}
		})
	}
}
