package utils

import (
	"sync"
	"time"
)

// IdempotencyRecord stores the result of an idempotent request
type IdempotencyRecord struct {
	StatusCode   int
	ResponseBody interface{}
	CreatedAt    time.Time
}

// IdempotencyStore manages idempotency keys and their responses
type IdempotencyStore struct {
	mu      sync.RWMutex
	records map[string]*IdempotencyRecord
	ttl     time.Duration
}

var (
	idempotencyStore     *IdempotencyStore
	idempotencyStoreOnce sync.Once
)

// GetIdempotencyStore returns the singleton idempotency store
func GetIdempotencyStore() *IdempotencyStore {
	idempotencyStoreOnce.Do(func() {
		idempotencyStore = &IdempotencyStore{
			records: make(map[string]*IdempotencyRecord),
			ttl:     10 * time.Minute, // Store idempotency records for 10 minutes
		}
		// Start cleanup goroutine
		go idempotencyStore.cleanupExpired()
	})
	return idempotencyStore
}

// Get retrieves a cached response for an idempotency key
func (s *IdempotencyStore) Get(key string) (*IdempotencyRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.records[key]
	if !exists {
		return nil, false
	}

	// Check if record has expired
	if time.Since(record.CreatedAt) > s.ttl {
		return nil, false
	}

	return record, true
}

// Set stores a response for an idempotency key
func (s *IdempotencyStore) Set(key string, statusCode int, responseBody interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[key] = &IdempotencyRecord{
		StatusCode:   statusCode,
		ResponseBody: responseBody,
		CreatedAt:    time.Now(),
	}
}

// Delete removes an idempotency key from the store
func (s *IdempotencyStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.records, key)
}

// cleanupExpired periodically removes expired idempotency records
func (s *IdempotencyStore) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, record := range s.records {
			if now.Sub(record.CreatedAt) > s.ttl {
				delete(s.records, key)
			}
		}
		s.mu.Unlock()
	}
}
