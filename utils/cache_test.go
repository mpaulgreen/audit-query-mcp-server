package utils

import (
	"fmt"
	"testing"
	"time"

	"audit-query-mcp-server/types"
)

// MockAuditResult creates a simple audit result for testing
func MockAuditResult(id string) *types.AuditResult {
	return &types.AuditResult{
		QueryID:       id,
		Timestamp:     time.Now().Format(time.RFC3339),
		Command:       "test-command",
		RawOutput:     "test-output",
		ParsedData:    []map[string]interface{}{{"test": "data"}},
		Summary:       "test-summary",
		ExecutionTime: 100,
	}
}

func TestNewCache(t *testing.T) {
	ttl := 1 * time.Hour
	cache := NewCache(ttl)

	if cache == nil {
		t.Fatal("NewCache returned nil")
	}

	if cache.ttl != ttl {
		t.Errorf("Expected TTL %v, got %v", ttl, cache.ttl)
	}

	if cache.entries == nil {
		t.Error("Cache entries map should be initialized")
	}

	if len(cache.entries) != 0 {
		t.Error("New cache should be empty")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	queryID := "test-query-123"
	result := MockAuditResult(queryID)

	// Test Set
	cache.Set(queryID, result)

	// Test Get
	retrieved, found := cache.Get(queryID)
	if !found {
		t.Error("Expected to find cached result")
	}

	if retrieved == nil {
		t.Fatal("Retrieved result should not be nil")
	}

	if retrieved.QueryID != queryID {
		t.Errorf("Expected QueryID %s, got %s", queryID, retrieved.QueryID)
	}
}

func TestCache_SetWithTTL(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	queryID := "test-query-ttl"
	result := MockAuditResult(queryID)
	customTTL := 100 * time.Millisecond

	// Test SetWithTTL
	cache.SetWithTTL(queryID, result, customTTL)

	// Should be found immediately
	retrieved, found := cache.Get(queryID)
	if !found {
		t.Error("Expected to find cached result immediately after setting")
	}

	if retrieved.QueryID != queryID {
		t.Errorf("Expected QueryID %s, got %s", queryID, retrieved.QueryID)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not be found after expiration
	retrieved, found = cache.Get(queryID)
	if found {
		t.Error("Expected cached result to be expired")
	}

	if retrieved != nil {
		t.Error("Retrieved result should be nil after expiration")
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	queryID := "non-existent-query"

	result, found := cache.Get(queryID)
	if found {
		t.Error("Expected not to find non-existent query")
	}

	if result != nil {
		t.Error("Result should be nil for non-existent query")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	queryID := "test-query-delete"
	result := MockAuditResult(queryID)

	// Set and verify
	cache.Set(queryID, result)
	_, found := cache.Get(queryID)
	if !found {
		t.Error("Expected to find result before deletion")
	}

	// Delete
	cache.Delete(queryID)

	// Verify deletion
	_, found = cache.Get(queryID)
	if found {
		t.Error("Expected not to find result after deletion")
	}
}

func TestCache_DeleteNonExistent(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	queryID := "non-existent-query"

	// Should not panic
	cache.Delete(queryID)
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	// Add multiple entries
	cache.Set("query1", MockAuditResult("query1"))
	cache.Set("query2", MockAuditResult("query2"))
	cache.Set("query3", MockAuditResult("query3"))

	// Verify entries exist
	if cache.Size() != 3 {
		t.Errorf("Expected 3 entries, got %d", cache.Size())
	}

	// Clear cache
	cache.Clear()

	// Verify cache is empty
	if cache.Size() != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", cache.Size())
	}

	// Verify individual entries are gone
	_, found := cache.Get("query1")
	if found {
		t.Error("Expected query1 to be removed after clear")
	}
}

func TestCache_Size(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	// Empty cache
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 for empty cache, got %d", cache.Size())
	}

	// Add entries
	cache.Set("query1", MockAuditResult("query1"))
	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}

	cache.Set("query2", MockAuditResult("query2"))
	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}

	// Delete entry
	cache.Delete("query1")
	if cache.Size() != 1 {
		t.Errorf("Expected size 1 after deletion, got %d", cache.Size())
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)
	queryID := "test-query-expiration"
	result := MockAuditResult(queryID)

	// Set entry
	cache.Set(queryID, result)

	// Should be found immediately
	_, found := cache.Get(queryID)
	if !found {
		t.Error("Expected to find result immediately")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should not be found after expiration
	_, found = cache.Get(queryID)
	if found {
		t.Error("Expected result to be expired")
	}

	// Cache size should be 0 after expiration
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after expiration, got %d", cache.Size())
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	done := make(chan bool, 10)

	// Start multiple goroutines reading and writing
	for i := 0; i < 5; i++ {
		go func(id int) {
			queryID := fmt.Sprintf("query-%d", id)
			result := MockAuditResult(queryID)

			// Write
			cache.Set(queryID, result)

			// Read
			_, found := cache.Get(queryID)
			if !found {
				t.Errorf("Goroutine %d: Expected to find result", id)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify final state
	if cache.Size() != 5 {
		t.Errorf("Expected 5 entries after concurrent access, got %d", cache.Size())
	}
}

func TestCache_GetStats(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	// Empty cache stats
	stats := cache.GetStats()
	if stats["size"] != 0 {
		t.Errorf("Expected size 0, got %v", stats["size"])
	}

	if stats["default_ttl"] != "1h0m0s" {
		t.Errorf("Expected default_ttl '1h0m0s', got %v", stats["default_ttl"])
	}

	// Add entries with different ages
	cache.Set("recent", MockAuditResult("recent"))
	time.Sleep(10 * time.Millisecond)

	cache.Set("older", MockAuditResult("older"))
	time.Sleep(10 * time.Millisecond)

	// Get stats
	stats = cache.GetStats()
	if stats["size"] != 2 {
		t.Errorf("Expected size 2, got %v", stats["size"])
	}

	ageDist, ok := stats["age_distribution"].(map[string]int)
	if !ok {
		t.Fatal("Expected age_distribution to be map[string]int")
	}

	// Should have entries in age categories
	if ageDist["<1m"] == 0 {
		t.Error("Expected some entries in <1m category")
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)

	// Add multiple entries
	cache.Set("query1", MockAuditResult("query1"))
	cache.Set("query2", MockAuditResult("query2"))
	cache.Set("query3", MockAuditResult("query3"))

	// Verify initial size
	if cache.Size() != 3 {
		t.Errorf("Expected 3 entries initially, got %d", cache.Size())
	}

	// Wait for expiration (cleanup runs every 5 minutes, but entries expire on access)
	time.Sleep(100 * time.Millisecond)

	// Try to get each entry - they should be expired and removed
	_, found := cache.Get("query1")
	if found {
		t.Error("Expected query1 to be expired")
	}

	_, found = cache.Get("query2")
	if found {
		t.Error("Expected query2 to be expired")
	}

	_, found = cache.Get("query3")
	if found {
		t.Error("Expected query3 to be expired")
	}

	// Cache should be empty after all expired entries are accessed
	if cache.Size() != 0 {
		t.Errorf("Expected 0 entries after expiration, got %d", cache.Size())
	}
}

func TestCache_Overwrite(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	queryID := "test-query-overwrite"

	// Set initial result
	result1 := MockAuditResult(queryID)
	result1.Summary = "initial"
	cache.Set(queryID, result1)

	// Verify initial result
	retrieved, found := cache.Get(queryID)
	if !found {
		t.Error("Expected to find initial result")
	}
	if retrieved.Summary != "initial" {
		t.Errorf("Expected summary 'initial', got %s", retrieved.Summary)
	}

	// Overwrite with new result
	result2 := MockAuditResult(queryID)
	result2.Summary = "updated"
	cache.Set(queryID, result2)

	// Verify updated result
	retrieved, found = cache.Get(queryID)
	if !found {
		t.Error("Expected to find updated result")
	}
	if retrieved.Summary != "updated" {
		t.Errorf("Expected summary 'updated', got %s", retrieved.Summary)
	}

	// Should still have only one entry
	if cache.Size() != 1 {
		t.Errorf("Expected 1 entry after overwrite, got %d", cache.Size())
	}
}

func TestCache_NilResult(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	queryID := "test-query-nil"

	// Set nil result
	cache.Set(queryID, nil)

	// Get nil result
	retrieved, found := cache.Get(queryID)
	if !found {
		t.Error("Expected to find nil result")
	}

	if retrieved != nil {
		t.Error("Expected retrieved result to be nil")
	}
}

func BenchmarkCache_Set(b *testing.B) {
	cache := NewCache(1 * time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queryID := fmt.Sprintf("query-%d", i)
		result := MockAuditResult(queryID)
		cache.Set(queryID, result)
	}
}

func BenchmarkCache_Get(b *testing.B) {
	cache := NewCache(1 * time.Hour)
	queryID := "benchmark-query"
	result := MockAuditResult(queryID)
	cache.Set(queryID, result)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(queryID)
	}
}

func BenchmarkCache_ConcurrentSet(b *testing.B) {
	cache := NewCache(1 * time.Hour)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			queryID := fmt.Sprintf("query-%d", i)
			result := MockAuditResult(queryID)
			cache.Set(queryID, result)
			i++
		}
	})
}
