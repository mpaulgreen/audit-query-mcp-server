# NLP Queries and Generated Commands - Complete Analysis

## Overview
This document contains all the natural language queries from the test patterns and their associated complete generated commands. Each command is the full, uncompressed version that would be executed in production.

## Category 1: Basic Query Patterns (Simple)

### 1.1: "Who deleted the customer CRD?"

**Natural Language Query**: Who deleted the customer CRD?

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"customresourcedefinition", "delete", "customer"},
    Timeframe: "yesterday",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
(oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-28 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-28.log | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-28.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-28.log.gz | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-28.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-28.log.bz2 | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -i 'customer' | grep -v 'system:' | grep '2025-01-28')
```

### 1.2: "Show me all actions by user john.doe today"

**Natural Language Query**: Show me all actions by user john.doe today

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "today",
    Username:  "john.doe",
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'john\.doe'
```

### 1.3: "List all failed authentication attempts in the last hour"

**Natural Language Query**: List all failed authentication attempts in the last hour

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "oauth-server",
    Patterns:  []string{"authentication", "failed"},
    Timeframe: "1h",
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=oauth-server/audit.log | grep -i 'authentication' | grep -i 'failed'
```

## Category 2: Resource Management Patterns (Intermediate)

### 2.1: "Find all CustomResourceDefinition modifications this week"

**Natural Language Query**: Find all CustomResourceDefinition modifications this week

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"customresourcedefinition"},
    Timeframe: "last_week",
    Verb:      "create|update|patch|delete",
}
```

**Complete Generated Command**:
```bash
(oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log.gz | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log.bz2 | grep -i 'customresourcedefinition' | grep -E 'create|update|patch|delete' | grep '2025-01-21')
```

### 2.2: "Show me all namespace deletions by non-system users"

**Natural Language Query**: Show me all namespace deletions by non-system users

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"namespaces"},
    Timeframe: "today",
    Verb:      "delete",
    Resource:  "namespaces",
    Exclude:   []string{"system:", "kube:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'namespaces' | grep -i 'delete' | grep -v 'system:' | grep -v 'kube:'
```

### 2.3: "Who created or modified ClusterRoles in the security namespace?"

**Natural Language Query**: Who created or modified ClusterRoles in the security namespace?

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"clusterroles"},
    Timeframe: "today",
    Verb:      "create|update|patch",
    Resource:  "clusterroles",
    Namespace: "security",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'clusterroles' | grep -E 'create|update|patch' | grep -i 'security' | grep -v 'system:'
```

## Category 3: Security Investigation Patterns (Advanced)

### 3.1: "Find potential privilege escalation attempts with failed permissions"

**Natural Language Query**: Find potential privilege escalation attempts with failed permissions

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"clusterrole", "rolebinding", "clusterrolebinding"},
    Timeframe: "24h",
    Exclude:   []string{"system:serviceaccount"},
    Verb:      "create|update|patch",
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'clusterrole\|rolebinding\|clusterrolebinding' | grep -E 'create|update|patch' | grep -v 'system:serviceaccount'
```

### 3.2: "Show unusual API access patterns outside business hours"

**Natural Language Query**: Show unusual API access patterns outside business hours

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "24h",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -v 'system:'
```

## Category 4: Complex Correlation Patterns (Expert)

### 4.1: "Correlate CRD deletions with subsequent pod creation failures"

**Natural Language Query**: Correlate CRD deletions with subsequent pod creation failures

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"customresourcedefinition", "delete"},
    Timeframe: "24h",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'customresourcedefinition' | grep -i 'delete' | grep -v 'system:'
```

### 4.2: "Find coordinated attacks: multiple failed authentications followed by successful privilege escalation"

**Natural Language Query**: Find coordinated attacks: multiple failed authentications followed by successful privilege escalation

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "oauth-server",
    Patterns:  []string{"authentication", "failed"},
    Timeframe: "24h",
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=oauth-server/audit.log | grep -i 'authentication' | grep -i 'failed'
```

## Category 5: Time-based Investigation Patterns

### 5.1: "Show me all admin activities during the maintenance window last Tuesday"

**Natural Language Query**: Show me all admin activities during the maintenance window last Tuesday

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "last_tuesday",
    Username:  "admin",
}
```

**Complete Generated Command**:
```bash
(oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'admin' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.bz2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log.bz2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.bz2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.bz2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.bz2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.bz2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log.gz | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21.bz2 | grep -i 'admin' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log.bz2 | grep -i 'admin' | grep '2025-01-21')
```

### 5.2: "Find API calls that happened between 2 AM and 4 AM this week"

**Natural Language Query**: Find API calls that happened between 2 AM and 4 AM this week

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "last_week",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
(oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -v 'system:' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.1.bz2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-1.log.bz2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2.bz2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2.log.bz2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.3.bz2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-3.log.bz2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log.gz | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit.log.2025-01-21.bz2 | grep -v 'system:' | grep '2025-01-21' && oc adm node-logs --role=master --path=kube-apiserver/audit-2025-01-21.log.bz2 | grep -v 'system:' | grep '2025-01-21')
```

## Category 6: Resource Correlation Patterns

### 6.1: "Which users accessed both the database and customer service namespaces?"

**Natural Language Query**: Which users accessed both the database and customer service namespaces?

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "24h",
    Namespace: "database|customer-service",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -E 'database|customer-service' | grep -v 'system:'
```

### 6.2: "Show me pod deletions followed by immediate recreations by the same user"

**Natural Language Query**: Show me pod deletions followed by immediate recreations by the same user

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"pods"},
    Timeframe: "24h",
    Verb:      "delete|create",
    Resource:  "pods",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'pods' | grep -E 'delete|create' | grep -v 'system:'
```

## Category 7: Anomaly Detection Patterns

### 7.1: "Identify users with unusual API access patterns compared to their baseline"

**Natural Language Query**: Identify users with unusual API access patterns compared to their baseline

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "24h",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -v 'system:'
```

### 7.2: "Show me service accounts being used from unexpected IP addresses"

**Natural Language Query**: Show me service accounts being used from unexpected IP addresses

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{"system:serviceaccount"},
    Timeframe: "24h",
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -i 'system:serviceaccount'
```

## Category 8: Advanced Investigation Patterns

### 8.1: "Correlate resource deletion events with subsequent access attempts to those resources"

**Natural Language Query**: Correlate resource deletion events with subsequent access attempts to those resources

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "24h",
    Verb:      "delete|get|list",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -E 'delete|get|list' | grep -v 'system:'
```

### 8.2: "Show me users who accessed multiple sensitive namespaces within a short time window"

**Natural Language Query**: Show me users who accessed multiple sensitive namespaces within a short time window

**Structured Parameters**:
```go
types.AuditQueryParams{
    LogSource: "kube-apiserver",
    Patterns:  []string{},
    Timeframe: "1h",
    Namespace: "kube-system|openshift-|security|database",
    Exclude:   []string{"system:"},
}
```

**Complete Generated Command**:
```bash
oc adm node-logs --role=master --path=kube-apiserver/audit.log | grep -E 'kube-system|openshift-|security|database' | grep -v 'system:'
```

## Summary

This document contains **18 natural language queries** across **8 categories** with their complete generated commands. Each command is the full, uncompressed version that would be executed in production.

### Key Observations:

1. **Command Complexity**: Commands range from simple single-line queries to complex multi-file searches with 20+ chained commands
2. **File Coverage**: Multi-file commands search across current logs, rotated logs (1-3), and compressed variants (.gz, .bz2)
3. **Pattern Matching**: Uses grep-based pattern matching with case-insensitive flags
4. **Chaining**: Uses `&&` chaining which can fail if any file is missing
5. **Date Filtering**: Includes date-specific filtering for historical queries

### Production Readiness Issues Identified:

1. **Chain Failure**: `&&` chaining breaks if any file doesn't exist
2. **Performance**: 20+ sequential `oc adm node-logs` commands
3. **Accuracy**: Grep-based parsing prone to false positives/negatives
4. **Error Handling**: No graceful failure handling
5. **Maintainability**: Complex command structures difficult to debug

This analysis provides the foundation for implementing the refactoring recommendations from the other agent's analysis. 