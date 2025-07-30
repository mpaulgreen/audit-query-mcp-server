package commands

import (
	"fmt"
	"strings"
)

// BuildUsernameFilter creates a comprehensive username filter for audit logs
func BuildUsernameFilter(username string) string {
	// Escape special characters for grep
	escapedUsername := escapeForGrep(username)

	// Build multiple grep patterns to catch different username formats in audit logs
	var patterns []string

	// Pattern 1: Username in user.username field (most common)
	patterns = append(patterns, fmt.Sprintf("| grep '\"user\":{\"[^\"]*\":\"%s\"'", escapedUsername))

	// Pattern 2: Username in user field directly
	patterns = append(patterns, fmt.Sprintf("| grep '\"user\":\"%s\"'", escapedUsername))

	// Pattern 3: Username in userInfo.username field
	patterns = append(patterns, fmt.Sprintf("| grep '\"userInfo\":{\"[^\"]*\":\"%s\"'", escapedUsername))

	// Pattern 4: Username in impersonatedUser field
	patterns = append(patterns, fmt.Sprintf("| grep '\"impersonatedUser\":\"%s\"'", escapedUsername))

	// Pattern 5: Username in requestUser field
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestUser\":\"%s\"'", escapedUsername))

	// Pattern 6: Username in annotations (for service accounts)
	patterns = append(patterns, fmt.Sprintf("| grep '\"authentication.kubernetes.io/username\":\"%s\"'", escapedUsername))

	// Pattern 7: Username in OAuth context
	patterns = append(patterns, fmt.Sprintf("| grep '\"oauth_user\":\"%s\"'", escapedUsername))

	// Pattern 8: Username in authentication context
	patterns = append(patterns, fmt.Sprintf("| grep '\"auth_user\":\"%s\"'", escapedUsername))

	// Pattern 9: Username in user-agent (for some OAuth flows)
	patterns = append(patterns, fmt.Sprintf("| grep '\"user_agent\":\"%s\"'", escapedUsername))

	// Pattern 10: Username in request headers
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestHeaders\":\"%s\"'", escapedUsername))

	return strings.Join(patterns, " ")
}

// BuildResourceFilter creates a comprehensive resource filter for audit logs
func BuildResourceFilter(resource string) string {
	// Escape special characters for grep
	escapedResource := escapeForGrep(resource)

	// Build multiple grep patterns to catch different resource formats in audit logs
	var patterns []string

	// Pattern 1: Resource in objectRef.resource field (most common)
	patterns = append(patterns, fmt.Sprintf("| grep '\"objectRef\":{\"[^\"]*\":\"%s\"'", escapedResource))

	// Pattern 2: Resource in objectRef.apiVersion field
	patterns = append(patterns, fmt.Sprintf("| grep '\"objectRef\":{\"[^\"]*\":\"[^\"]*%s\"'", escapedResource))

	// Pattern 3: Resource in requestObject.kind field
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestObject\":{\"[^\"]*\":\"%s\"'", escapedResource))

	// Pattern 4: Resource in responseObject.kind field
	patterns = append(patterns, fmt.Sprintf("| grep '\"responseObject\":{\"[^\"]*\":\"%s\"'", escapedResource))

	// Pattern 5: Resource in requestURI path
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestURI\":\"[^\"]*%s[^\"]*\"'", escapedResource))

	// Pattern 6: Resource in annotations
	patterns = append(patterns, fmt.Sprintf("| grep '\"annotations\":{\"[^\"]*\":\"%s\"'", escapedResource))

	// Pattern 7: Resource in labels
	patterns = append(patterns, fmt.Sprintf("| grep '\"labels\":{\"[^\"]*\":\"%s\"'", escapedResource))

	// Pattern 8: Resource in metadata.name field
	patterns = append(patterns, fmt.Sprintf("| grep '\"metadata\":{\"[^\"]*\":\"%s\"'", escapedResource))

	// Pattern 9: Resource in spec field
	patterns = append(patterns, fmt.Sprintf("| grep '\"spec\":{\"[^\"]*\":\"%s\"'", escapedResource))

	// Pattern 10: Resource in status field
	patterns = append(patterns, fmt.Sprintf("| grep '\"status\":{\"[^\"]*\":\"%s\"'", escapedResource))

	return strings.Join(patterns, " ")
}

// BuildVerbFilter creates a comprehensive verb filter for audit logs
func BuildVerbFilter(verb string) string {
	// Escape special characters for grep
	escapedVerb := escapeForGrep(verb)

	// Build multiple grep patterns to catch different verb formats in audit logs
	var patterns []string

	// Pattern 1: Verb in verb field (most common)
	patterns = append(patterns, fmt.Sprintf("| grep '\"verb\":\"%s\"'", escapedVerb))

	// Pattern 2: Verb in method field
	patterns = append(patterns, fmt.Sprintf("| grep '\"method\":\"%s\"'", escapedVerb))

	// Pattern 3: Verb in action field
	patterns = append(patterns, fmt.Sprintf("| grep '\"action\":\"%s\"'", escapedVerb))

	// Pattern 4: Verb in operation field
	patterns = append(patterns, fmt.Sprintf("| grep '\"operation\":\"%s\"'", escapedVerb))

	// Pattern 5: Verb in requestURI path
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestURI\":\"[^\"]*%s[^\"]*\"'", escapedVerb))

	// Pattern 6: Verb in requestMethod field
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestMethod\":\"%s\"'", escapedVerb))

	// Pattern 7: Verb in auditID field (less common)
	patterns = append(patterns, fmt.Sprintf("| grep '\"auditID\":\"%s\"'", escapedVerb))

	// Pattern 8: Verb in stage field
	patterns = append(patterns, fmt.Sprintf("| grep '\"stage\":\"%s\"'", escapedVerb))

	// Pattern 9: Verb in responseStatus field
	patterns = append(patterns, fmt.Sprintf("| grep '\"responseStatus\":{\"[^\"]*\":\"%s\"'", escapedVerb))

	// Pattern 10: Verb in HTTP method context
	patterns = append(patterns, fmt.Sprintf("| grep '\"httpMethod\":\"%s\"'", escapedVerb))

	return strings.Join(patterns, " ")
}

// BuildNamespaceFilter creates a comprehensive namespace filter for audit logs
func BuildNamespaceFilter(namespace string) string {
	// Escape special characters for grep
	escapedNamespace := escapeForGrep(namespace)

	// Build multiple grep patterns to catch different namespace formats in audit logs
	var patterns []string

	// Pattern 1: Namespace in objectRef.namespace field (most common)
	patterns = append(patterns, fmt.Sprintf("| grep '\"objectRef\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 2: Namespace in requestObject.metadata.namespace field
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestObject\":{\"[^\"]*\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 3: Namespace in responseObject.metadata.namespace field
	patterns = append(patterns, fmt.Sprintf("| grep '\"responseObject\":{\"[^\"]*\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 4: Namespace in requestURI path
	patterns = append(patterns, fmt.Sprintf("| grep '\"requestURI\":\"[^\"]*%s[^\"]*\"'", escapedNamespace))

	// Pattern 5: Namespace in annotations
	patterns = append(patterns, fmt.Sprintf("| grep '\"annotations\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 6: Namespace in labels
	patterns = append(patterns, fmt.Sprintf("| grep '\"labels\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 7: Namespace in metadata.namespace field
	patterns = append(patterns, fmt.Sprintf("| grep '\"metadata\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 8: Namespace in spec field
	patterns = append(patterns, fmt.Sprintf("| grep '\"spec\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 9: Namespace in status field
	patterns = append(patterns, fmt.Sprintf("| grep '\"status\":{\"[^\"]*\":\"%s\"'", escapedNamespace))

	// Pattern 10: Namespace in user context (for service accounts)
	patterns = append(patterns, fmt.Sprintf("| grep '\"user\":{\"[^\"]*\":\"[^\"]*%s[^\"]*\"'", escapedNamespace))

	return strings.Join(patterns, " ")
}

// escapeForGrep escapes special characters for safe grep usage
func escapeForGrep(input string) string {
	// Escape special grep characters: [ ] ( ) . * + ? ^ $ { } | \
	escaped := input
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "[", "\\[")
	escaped = strings.ReplaceAll(escaped, "]", "\\]")
	escaped = strings.ReplaceAll(escaped, "(", "\\(")
	escaped = strings.ReplaceAll(escaped, ")", "\\)")
	escaped = strings.ReplaceAll(escaped, ".", "\\.")
	escaped = strings.ReplaceAll(escaped, "*", "\\*")
	escaped = strings.ReplaceAll(escaped, "+", "\\+")
	escaped = strings.ReplaceAll(escaped, "?", "\\?")
	escaped = strings.ReplaceAll(escaped, "^", "\\^")
	escaped = strings.ReplaceAll(escaped, "$", "\\$")
	escaped = strings.ReplaceAll(escaped, "{", "\\{")
	escaped = strings.ReplaceAll(escaped, "}", "\\}")
	escaped = strings.ReplaceAll(escaped, "|", "\\|")
	escaped = strings.ReplaceAll(escaped, "'", "\\'")

	return escaped
}
