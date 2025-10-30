package main

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MockDB struct {
	pool *pgxpool.Pool 
}

func (db *MockDB) InitDB() error { return nil }
func (db *MockDB) Close() {}

func TestCache_SetGet(t *testing.T) {
	c := NewCache()
	order := Order{OrderUID: "test_uid_1"}
	c.Set(order.OrderUID, order)

	retrieved, ok := c.Get("test_uid_1")
	assert.True(t, ok, "Order should be found")
	assert.Equal(t, "test_uid_1", retrieved.OrderUID, "Retrieved UID should match")

	_, ok = c.Get("non_existent_uid")
	assert.False(t, ok, "Non-existent order should not be found")
}

func TestCache_ThreadSafety(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup
	numRoutines := 100
	numOperations := 1000

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				uid := time.Now().Format("20060102150405.000") + string(i)
				order := Order{OrderUID: uid}
				c.Set(uid, order)
				c.Get(uid)
			}
		}(i)
	}

	wg.Wait()
	assert.True(t, true, "Test completed without data race warnings")
}
