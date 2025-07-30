package commands

import (
	"strings"
	"testing"
)

// TestBuildUsernameFilter tests username filter generation
func TestBuildUsernameFilter(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected []string
	}{
		{
			name:     "simple username",
			username: "testuser",
			expected: []string{
				"\"user\":{\"[^\"]*\":\"testuser\"",
				"\"user\":\"testuser\"",
				"\"userInfo\":{\"[^\"]*\":\"testuser\"",
				"\"impersonatedUser\":\"testuser\"",
				"\"requestUser\":\"testuser\"",
				"\"authentication.kubernetes.io/username\":\"testuser\"",
				"\"oauth_user\":\"testuser\"",
				"\"auth_user\":\"testuser\"",
				"\"user_agent\":\"testuser\"",
				"\"requestHeaders\":\"testuser\"",
			},
		},
		{
			name:     "service account username",
			username: "system:serviceaccount:default:myapp",
			expected: []string{
				"\"user\":{\"[^\"]*\":\"system:serviceaccount:default:myapp\"",
				"\"user\":\"system:serviceaccount:default:myapp\"",
				"\"userInfo\":{\"[^\"]*\":\"system:serviceaccount:default:myapp\"",
				"\"impersonatedUser\":\"system:serviceaccount:default:myapp\"",
				"\"requestUser\":\"system:serviceaccount:default:myapp\"",
			},
		},
		{
			name:     "username with special characters",
			username: "test[user]",
			expected: []string{
				"\"user\":{\"[^\"]*\":\"test\\[user\\]\"",
				"\"user\":\"test\\[user\\]\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := BuildUsernameFilter(tt.username)

			// Check that all expected patterns are present
			for _, expected := range tt.expected {
				if !strings.Contains(filter, expected) {
					t.Errorf("Username filter should contain pattern: %s", expected)
				}
			}

			// Check that the filter starts with a pipe
			if !strings.HasPrefix(filter, "| grep") {
				t.Errorf("Filter should start with '| grep', got: %s", filter[:10])
			}

			// Check that special characters are escaped
			if strings.Contains(tt.username, "[") && !strings.Contains(filter, "\\[") {
				t.Errorf("Special characters should be escaped in filter")
			}
		})
	}
}

// TestBuildResourceFilter tests resource filter generation
func TestBuildResourceFilter(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected []string
	}{
		{
			name:     "simple resource",
			resource: "pods",
			expected: []string{
				"\"objectRef\":{\"[^\"]*\":\"pods\"",
				"\"objectRef\":{\"[^\"]*\":\"[^\"]*pods\"",
				"\"requestObject\":{\"[^\"]*\":\"pods\"",
				"\"responseObject\":{\"[^\"]*\":\"pods\"",
				"\"requestURI\":\"[^\"]*pods[^\"]*\"",
				"\"annotations\":{\"[^\"]*\":\"pods\"",
				"\"labels\":{\"[^\"]*\":\"pods\"",
				"\"metadata\":{\"[^\"]*\":\"pods\"",
				"\"spec\":{\"[^\"]*\":\"pods\"",
				"\"status\":{\"[^\"]*\":\"pods\"",
			},
		},
		{
			name:     "api resource",
			resource: "v1/pods",
			expected: []string{
				"\"objectRef\":{\"[^\"]*\":\"v1/pods\"",
				"\"objectRef\":{\"[^\"]*\":\"[^\"]*v1/pods\"",
				"\"requestObject\":{\"[^\"]*\":\"v1/pods\"",
				"\"responseObject\":{\"[^\"]*\":\"v1/pods\"",
			},
		},
		{
			name:     "resource with special characters",
			resource: "test[resource]",
			expected: []string{
				"\"objectRef\":{\"[^\"]*\":\"test\\[resource\\]\"",
				"\"objectRef\":{\"[^\"]*\":\"[^\"]*test\\[resource\\]\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := BuildResourceFilter(tt.resource)

			// Check that all expected patterns are present
			for _, expected := range tt.expected {
				if !strings.Contains(filter, expected) {
					t.Errorf("Resource filter should contain pattern: %s", expected)
				}
			}

			// Check that the filter starts with a pipe
			if !strings.HasPrefix(filter, "| grep") {
				t.Errorf("Filter should start with '| grep', got: %s", filter[:10])
			}

			// Check that special characters are escaped
			if strings.Contains(tt.resource, "[") && !strings.Contains(filter, "\\[") {
				t.Errorf("Special characters should be escaped in filter")
			}
		})
	}
}

// TestBuildVerbFilter tests verb filter generation
func TestBuildVerbFilter(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		expected []string
	}{
		{
			name: "GET verb",
			verb: "GET",
			expected: []string{
				"\"verb\":\"GET\"",
				"\"method\":\"GET\"",
				"\"action\":\"GET\"",
				"\"operation\":\"GET\"",
				"\"requestURI\":\"[^\"]*GET[^\"]*\"",
				"\"requestMethod\":\"GET\"",
				"\"auditID\":\"GET\"",
				"\"stage\":\"GET\"",
				"\"responseStatus\":{\"[^\"]*\":\"GET\"",
				"\"httpMethod\":\"GET\"",
			},
		},
		{
			name: "CREATE verb",
			verb: "create",
			expected: []string{
				"\"verb\":\"create\"",
				"\"method\":\"create\"",
				"\"action\":\"create\"",
				"\"operation\":\"create\"",
				"\"requestURI\":\"[^\"]*create[^\"]*\"",
				"\"requestMethod\":\"create\"",
			},
		},
		{
			name: "verb with special characters",
			verb: "test[verb]",
			expected: []string{
				"\"verb\":\"test\\[verb\\]\"",
				"\"method\":\"test\\[verb\\]\"",
				"\"action\":\"test\\[verb\\]\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := BuildVerbFilter(tt.verb)

			// Check that all expected patterns are present
			for _, expected := range tt.expected {
				if !strings.Contains(filter, expected) {
					t.Errorf("Verb filter should contain pattern: %s", expected)
				}
			}

			// Check that the filter starts with a pipe
			if !strings.HasPrefix(filter, "| grep") {
				t.Errorf("Filter should start with '| grep', got: %s", filter[:10])
			}

			// Check that special characters are escaped
			if strings.Contains(tt.verb, "[") && !strings.Contains(filter, "\\[") {
				t.Errorf("Special characters should be escaped in filter")
			}
		})
	}
}

// TestBuildNamespaceFilter tests namespace filter generation
func TestBuildNamespaceFilter(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		expected  []string
	}{
		{
			name:      "simple namespace",
			namespace: "default",
			expected: []string{
				"\"objectRef\":{\"[^\"]*\":\"default\"",
				"\"requestObject\":{\"[^\"]*\":{\"[^\"]*\":\"default\"",
				"\"responseObject\":{\"[^\"]*\":{\"[^\"]*\":\"default\"",
				"\"requestURI\":\"[^\"]*default[^\"]*\"",
				"\"annotations\":{\"[^\"]*\":\"default\"",
				"\"labels\":{\"[^\"]*\":\"default\"",
				"\"metadata\":{\"[^\"]*\":\"default\"",
				"\"spec\":{\"[^\"]*\":\"default\"",
				"\"status\":{\"[^\"]*\":\"default\"",
				"\"user\":{\"[^\"]*\":\"[^\"]*default[^\"]*\"",
			},
		},
		{
			name:      "namespace with special characters",
			namespace: "test[namespace]",
			expected: []string{
				"\"objectRef\":{\"[^\"]*\":\"test\\[namespace\\]\"",
				"\"requestObject\":{\"[^\"]*\":{\"[^\"]*\":\"test\\[namespace\\]\"",
				"\"responseObject\":{\"[^\"]*\":{\"[^\"]*\":\"test\\[namespace\\]\"",
			},
		},
		{
			name:      "kube-system namespace",
			namespace: "kube-system",
			expected: []string{
				"\"objectRef\":{\"[^\"]*\":\"kube-system\"",
				"\"requestObject\":{\"[^\"]*\":{\"[^\"]*\":\"kube-system\"",
				"\"responseObject\":{\"[^\"]*\":{\"[^\"]*\":\"kube-system\"",
				"\"requestURI\":\"[^\"]*kube-system[^\"]*\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := BuildNamespaceFilter(tt.namespace)

			// Check that all expected patterns are present
			for _, expected := range tt.expected {
				if !strings.Contains(filter, expected) {
					t.Errorf("Namespace filter should contain pattern: %s", expected)
				}
			}

			// Check that the filter starts with a pipe
			if !strings.HasPrefix(filter, "| grep") {
				t.Errorf("Filter should start with '| grep', got: %s", filter[:10])
			}

			// Check that special characters are escaped
			if strings.Contains(tt.namespace, "[") && !strings.Contains(filter, "\\[") {
				t.Errorf("Special characters should be escaped in filter")
			}
		})
	}
}

// TestEscapeForGrep tests the escapeForGrep function
func TestEscapeForGrep(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "simple",
			expected: "simple",
		},
		{
			name:     "brackets",
			input:    "test[user]",
			expected: "test\\[user\\]",
		},
		{
			name:     "parentheses",
			input:    "test(user)",
			expected: "test\\(user\\)",
		},
		{
			name:     "dots",
			input:    "test.user",
			expected: "test\\.user",
		},
		{
			name:     "asterisks",
			input:    "test*user",
			expected: "test\\*user",
		},
		{
			name:     "plus signs",
			input:    "test+user",
			expected: "test\\+user",
		},
		{
			name:     "question marks",
			input:    "test?user",
			expected: "test\\?user",
		},
		{
			name:     "carets",
			input:    "test^user",
			expected: "test\\^user",
		},
		{
			name:     "dollar signs",
			input:    "test$user",
			expected: "test\\$user",
		},
		{
			name:     "braces",
			input:    "test{user}",
			expected: "test\\{user\\}",
		},
		{
			name:     "pipes",
			input:    "test|user",
			expected: "test\\|user",
		},
		{
			name:     "backslashes",
			input:    "test\\user",
			expected: "test\\\\user",
		},
		{
			name:     "single quotes",
			input:    "test'user",
			expected: "test\\'user",
		},
		{
			name:     "multiple special characters",
			input:    "test[user].*+?^${}|\\'",
			expected: "test\\[user\\]\\.\\*\\+\\?\\^\\$\\{\\}\\|\\\\\\'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeForGrep(tt.input)
			if result != tt.expected {
				t.Errorf("escapeForGrep(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFilterIntegration tests integration of all filters
func TestFilterIntegration(t *testing.T) {
	// Test that all filters work together without conflicts
	usernameFilter := BuildUsernameFilter("testuser")
	resourceFilter := BuildResourceFilter("pods")
	verbFilter := BuildVerbFilter("GET")
	namespaceFilter := BuildNamespaceFilter("default")

	// All filters should be non-empty
	if usernameFilter == "" {
		t.Error("Username filter should not be empty")
	}
	if resourceFilter == "" {
		t.Error("Resource filter should not be empty")
	}
	if verbFilter == "" {
		t.Error("Verb filter should not be empty")
	}
	if namespaceFilter == "" {
		t.Error("Namespace filter should not be empty")
	}

	// All filters should start with pipe
	filters := []string{usernameFilter, resourceFilter, verbFilter, namespaceFilter}
	for i, filter := range filters {
		if !strings.HasPrefix(filter, "| grep") {
			t.Errorf("Filter %d should start with '| grep', got: %s", i, filter[:10])
		}
	}

	// All filters should contain the expected content
	if !strings.Contains(usernameFilter, "testuser") {
		t.Error("Username filter should contain 'testuser'")
	}
	if !strings.Contains(resourceFilter, "pods") {
		t.Error("Resource filter should contain 'pods'")
	}
	if !strings.Contains(verbFilter, "GET") {
		t.Error("Verb filter should contain 'GET'")
	}
	if !strings.Contains(namespaceFilter, "default") {
		t.Error("Namespace filter should contain 'default'")
	}
}

// TestFilterEdgeCases tests edge cases for filters
func TestFilterEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(string) string
		input    string
	}{
		{
			name:     "empty username",
			testFunc: BuildUsernameFilter,
			input:    "",
		},
		{
			name:     "empty resource",
			testFunc: BuildResourceFilter,
			input:    "",
		},
		{
			name:     "empty verb",
			testFunc: BuildVerbFilter,
			input:    "",
		},
		{
			name:     "empty namespace",
			testFunc: BuildNamespaceFilter,
			input:    "",
		},
		{
			name:     "very long username",
			testFunc: BuildUsernameFilter,
			input:    "verylongusername" + strings.Repeat("a", 1000),
		},
		{
			name:     "unicode username",
			testFunc: BuildUsernameFilter,
			input:    "userñame",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := tt.testFunc(tt.input)

			// Filter should not be empty even for empty input
			if filter == "" {
				t.Error("Filter should not be empty even for empty input")
			}

			// Filter should start with pipe
			if !strings.HasPrefix(filter, "| grep") {
				t.Errorf("Filter should start with '| grep', got: %s", filter[:10])
			}

			// Filter should contain the input (if not empty)
			if tt.input != "" && !strings.Contains(filter, tt.input) {
				t.Errorf("Filter should contain input '%s'", tt.input)
			}
		})
	}
}

// BenchmarkFilterGeneration benchmarks filter generation performance
func BenchmarkBuildUsernameFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildUsernameFilter("testuser")
	}
}

func BenchmarkBuildResourceFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildResourceFilter("pods")
	}
}

func BenchmarkBuildVerbFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildVerbFilter("GET")
	}
}

func BenchmarkBuildNamespaceFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildNamespaceFilter("default")
	}
}

func BenchmarkEscapeForGrep(b *testing.B) {
	input := "test[user].*+?^${}|\\'"
	for i := 0; i < b.N; i++ {
		escapeForGrep(input)
	}
}
