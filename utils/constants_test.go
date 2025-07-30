package utils

import (
	"testing"
)

// TestValidLogSources validates the log sources constants
func TestValidLogSources(t *testing.T) {
	expectedSources := []string{
		"kube-apiserver",
		"oauth-server",
		"node",
		"openshift-apiserver",
		"oauth-apiserver",
	}

	if len(ValidLogSources) != len(expectedSources) {
		t.Errorf("Expected %d log sources, got %d", len(expectedSources), len(ValidLogSources))
	}

	for _, expected := range expectedSources {
		found := false
		for _, actual := range ValidLogSources {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log source %s not found", expected)
		}
	}
}

// TestValidResources validates the resources constants
func TestValidResources(t *testing.T) {
	// Test that core resources are present
	coreResources := []string{"pods", "services", "deployments", "namespaces", "nodes"}
	for _, resource := range coreResources {
		found := false
		for _, validResource := range ValidResources {
			if validResource == resource {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected core resource %s not found", resource)
		}
	}

	// Test that OpenShift specific resources are present
	openshiftResources := []string{"projects", "routes", "builds", "deploymentconfigs"}
	for _, resource := range openshiftResources {
		found := false
		for _, validResource := range ValidResources {
			if validResource == resource {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected OpenShift resource %s not found", resource)
		}
	}
}

// TestValidVerbs validates the verbs constants
func TestValidVerbs(t *testing.T) {
	// Test that standard CRUD verbs are present
	crudVerbs := []string{"create", "get", "list", "watch", "update", "patch", "delete"}
	for _, verb := range crudVerbs {
		found := false
		for _, validVerb := range ValidVerbs {
			if validVerb == verb {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected CRUD verb %s not found", verb)
		}
	}
}

// TestResponseStatusCodes validates the response status codes
func TestResponseStatusCodes(t *testing.T) {
	// Test success codes
	expectedSuccessCodes := map[string]int{
		"OK":        200,
		"Created":   201,
		"NoContent": 204,
	}

	for name, expectedCode := range expectedSuccessCodes {
		if actualCode, exists := ResponseStatusCodes[name]; !exists {
			t.Errorf("Expected status code %s not found", name)
		} else if actualCode != expectedCode {
			t.Errorf("Expected status code %s to be %d, got %d", name, expectedCode, actualCode)
		}
	}

	// Test error codes
	expectedErrorCodes := map[string]int{
		"BadRequest":          400,
		"Unauthorized":        401,
		"Forbidden":           403,
		"NotFound":            404,
		"InternalServerError": 500,
	}

	for name, expectedCode := range expectedErrorCodes {
		if actualCode, exists := ResponseStatusCodes[name]; !exists {
			t.Errorf("Expected status code %s not found", name)
		} else if actualCode != expectedCode {
			t.Errorf("Expected status code %s to be %d, got %d", name, expectedCode, actualCode)
		}
	}
}

// TestStatusCodeRanges validates the status code ranges
func TestStatusCodeRanges(t *testing.T) {
	// Test success range
	if successCodes, exists := StatusCodeRanges["success"]; !exists {
		t.Error("Expected success status code range not found")
	} else {
		expectedSuccessCodes := []int{200, 201, 202, 204}
		for _, expected := range expectedSuccessCodes {
			found := false
			for _, actual := range successCodes {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected success code %d not found", expected)
			}
		}
	}

	// Test client error range
	if clientErrorCodes, exists := StatusCodeRanges["client_error"]; !exists {
		t.Error("Expected client_error status code range not found")
	} else {
		expectedClientErrorCodes := []int{400, 401, 403, 404, 409, 422, 429}
		for _, expected := range expectedClientErrorCodes {
			found := false
			for _, actual := range clientErrorCodes {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected client error code %d not found", expected)
			}
		}
	}
}

// TestAuditLogFields validates the audit log fields
func TestAuditLogFields(t *testing.T) {
	// Test core fields
	coreFields := []string{"RequestReceivedTimestamp", "ResponseStatus", "ObjectRef", "User", "Verb"}
	for _, field := range coreFields {
		if _, exists := AuditLogFields[field]; !exists {
			t.Errorf("Expected audit log field %s not found", field)
		}
	}

	// Test that field values are not empty
	for name, value := range AuditLogFields {
		if value == "" {
			t.Errorf("Audit log field %s has empty value", name)
		}
	}
}

// TestTimeFrameConstants validates the timeframe constants
func TestTimeFrameConstants(t *testing.T) {
	// Test common timeframes
	commonTimeframes := []string{"Today", "Yesterday", "LastHour", "Last24Hours", "Last7Days"}
	for _, timeframe := range commonTimeframes {
		if _, exists := TimeFrameConstants[timeframe]; !exists {
			t.Errorf("Expected timeframe constant %s not found", timeframe)
		}
	}

	// Test short forms
	shortForms := []string{"1m", "5m", "1h", "1d", "1w"}
	for _, shortForm := range shortForms {
		if _, exists := TimeFrameConstants[shortForm]; !exists {
			t.Errorf("Expected short form timeframe %s not found", shortForm)
		}
	}
}

// TestSystemUserGroups validates the system user groups
func TestSystemUserGroups(t *testing.T) {
	// Test system authenticated group
	if systemAuth, exists := SystemUserGroups["system_authenticated"]; !exists {
		t.Error("Expected system_authenticated group not found")
	} else {
		expectedAuth := []string{"system:authenticated", "system:authenticated:oauth"}
		for _, expected := range expectedAuth {
			found := false
			for _, actual := range systemAuth {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected system authenticated user %s not found", expected)
			}
		}
	}

	// Test system service accounts
	if serviceAccounts, exists := SystemUserGroups["system_service_accounts"]; !exists {
		t.Error("Expected system_service_accounts group not found")
	} else {
		expectedServiceAccounts := []string{"system:serviceaccount", "system:serviceaccount:kube-system"}
		for _, expected := range expectedServiceAccounts {
			found := false
			for _, actual := range serviceAccounts {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected system service account %s not found", expected)
			}
		}
	}
}

// TestSecurityPatterns validates the security patterns
func TestSecurityPatterns(t *testing.T) {
	// Test privilege escalation patterns
	if privilegeEscalation, exists := SecurityPatterns["privilege_escalation"]; !exists {
		t.Error("Expected privilege_escalation patterns not found")
	} else {
		expectedPatterns := []string{"clusterrole", "clusterrolebinding", "role", "rolebinding"}
		for _, expected := range expectedPatterns {
			found := false
			for _, actual := range privilegeEscalation {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected privilege escalation pattern %s not found", expected)
			}
		}
	}

	// Test authentication failures
	if authFailures, exists := SecurityPatterns["authentication_failures"]; !exists {
		t.Error("Expected authentication_failures patterns not found")
	} else {
		expectedPatterns := []string{"401", "403", "Unauthorized", "Forbidden"}
		for _, expected := range expectedPatterns {
			found := false
			for _, actual := range authFailures {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected authentication failure pattern %s not found", expected)
			}
		}
	}
}

// TestCacheConfig validates the cache configuration
func TestCacheConfig(t *testing.T) {
	// Test default TTL
	if defaultTTL, exists := CacheConfig["default_ttl"]; !exists {
		t.Error("Expected default_ttl not found in cache config")
	} else {
		if ttl, ok := defaultTTL.(int); !ok {
			t.Error("Expected default_ttl to be an integer")
		} else if ttl <= 0 {
			t.Error("Expected default_ttl to be positive")
		}
	}

	// Test max cache size
	if maxSize, exists := CacheConfig["max_cache_size"]; !exists {
		t.Error("Expected max_cache_size not found in cache config")
	} else {
		if size, ok := maxSize.(int); !ok {
			t.Error("Expected max_cache_size to be an integer")
		} else if size <= 0 {
			t.Error("Expected max_cache_size to be positive")
		}
	}
}

// TestExecutionLimits validates the execution limits
func TestExecutionLimits(t *testing.T) {
	// Test max execution time
	if maxTime, exists := ExecutionLimits["max_execution_time"]; !exists {
		t.Error("Expected max_execution_time not found in execution limits")
	} else {
		if time, ok := maxTime.(int); !ok {
			t.Error("Expected max_execution_time to be an integer")
		} else if time <= 0 {
			t.Error("Expected max_execution_time to be positive")
		}
	}

	// Test max output size
	if maxOutput, exists := ExecutionLimits["max_output_size"]; !exists {
		t.Error("Expected max_output_size not found in execution limits")
	} else {
		if size, ok := maxOutput.(int); !ok {
			t.Error("Expected max_output_size to be an integer")
		} else if size <= 0 {
			t.Error("Expected max_output_size to be positive")
		}
	}
}

// TestTimeoutValues validates the timeout values
func TestTimeoutValues(t *testing.T) {
	// Test query execution timeout
	if timeout, exists := TimeoutValues["query_execution"]; !exists {
		t.Error("Expected query_execution timeout not found")
	} else if timeout <= 0 {
		t.Error("Expected query_execution timeout to be positive")
	}

	// Test validation timeout
	if timeout, exists := TimeoutValues["validation"]; !exists {
		t.Error("Expected validation timeout not found")
	} else if timeout <= 0 {
		t.Error("Expected validation timeout to be positive")
	}
}

// TestErrorMessages validates the error messages
func TestErrorMessages(t *testing.T) {
	// Test that error messages are not empty
	for name, message := range ErrorMessages {
		if message == "" {
			t.Errorf("Error message for %s is empty", name)
		}
	}

	// Test specific error messages
	expectedErrors := []string{"invalid_log_source", "invalid_resource", "invalid_verb"}
	for _, expected := range expectedErrors {
		if _, exists := ErrorMessages[expected]; !exists {
			t.Errorf("Expected error message %s not found", expected)
		}
	}
}

// TestErrorCodes validates the error codes
func TestErrorCodes(t *testing.T) {
	// Test that error codes are valid HTTP status codes
	for name, code := range ErrorCodes {
		if code < 100 || code > 599 {
			t.Errorf("Error code %s has invalid HTTP status code %d", name, code)
		}
	}

	// Test specific error codes
	expectedCodes := map[string]int{
		"VALIDATION_ERROR":  400,
		"PERMISSION_DENIED": 403,
		"NOT_FOUND":         404,
		"UNKNOWN_ERROR":     500,
	}

	for name, expectedCode := range expectedCodes {
		if actualCode, exists := ErrorCodes[name]; !exists {
			t.Errorf("Expected error code %s not found", name)
		} else if actualCode != expectedCode {
			t.Errorf("Expected error code %s to be %d, got %d", name, expectedCode, actualCode)
		}
	}
}

// TestFilePathPatterns validates the file path patterns
func TestFilePathPatterns(t *testing.T) {
	// Test that file paths are not empty
	for name, path := range FilePathPatterns {
		if path == "" {
			t.Errorf("File path for %s is empty", name)
		}
	}

	// Test specific file paths
	expectedPaths := []string{"kube_apiserver_log", "oauth_server_log", "node_log"}
	for _, expected := range expectedPaths {
		if _, exists := FilePathPatterns[expected]; !exists {
			t.Errorf("Expected file path %s not found", expected)
		}
	}
}

// TestLogFilePatterns validates the log file patterns
func TestLogFilePatterns(t *testing.T) {
	// Test that log file patterns are not empty
	for _, pattern := range LogFilePatterns {
		if pattern == "" {
			t.Error("Log file pattern is empty")
		}
	}

	// Test that basic patterns are present
	expectedPatterns := []string{"audit.log", "audit.log.*"}
	for _, expected := range expectedPatterns {
		found := false
		for _, actual := range LogFilePatterns {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log file pattern %s not found", expected)
		}
	}
}

// TestDangerousPatterns validates the dangerous patterns
func TestDangerousPatterns(t *testing.T) {
	// Test that dangerous patterns are not empty
	for _, pattern := range DangerousPatterns {
		if pattern == "" {
			t.Error("Dangerous pattern is empty")
		}
	}

	// Test that critical dangerous patterns are present
	criticalPatterns := []string{"oc delete", "kubectl delete", "&&", "||"}
	for _, expected := range criticalPatterns {
		found := false
		for _, actual := range DangerousPatterns {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected dangerous pattern %s not found", expected)
		}
	}
}

// TestSafeDatePatterns validates the safe date patterns
func TestSafeDatePatterns(t *testing.T) {
	// Test that safe date patterns are not empty
	for _, pattern := range SafeDatePatterns {
		if pattern == "" {
			t.Error("Safe date pattern is empty")
		}
	}

	// Test that safe date patterns are present
	expectedPatterns := []string{"$(date", "$(date -d", "$(date -v"}
	for _, expected := range expectedPatterns {
		found := false
		for _, actual := range SafeDatePatterns {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected safe date pattern %s not found", expected)
		}
	}
}
