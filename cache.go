package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
)

type Order struct {
	OrderUID          string                 `json:"order_uid"`
	TrackNumber       string                 `json:"track_number"`
	Entry             string                 `json:"entry"`
	Delivery          map[string]interface{} `json:"delivery"`
	Payment           map[string]interface{} `json:"payment"`
	Items             []map[string]interface{} `json:"items"`
	Locale            string                 `json:"locale"`
	InternalSignature string                 `json:"internal_signature"`
	CustomerID        string                 `json:"customer_id"`
	DeliveryService   string                 `json:"delivery_service"`
	Shardkey          string                 `json:"shardkey"`
	SmID              int                    `json:"sm_id"`
	DateCreated       string                 `json:"date_created"`
	OofShard          string                 `json:"oof_shard"`
}

type Cache struct {
	orders map[string]Order
	mu     sync.RWMutex 
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]Order),
	}
}

func (c *Cache) LoadFromDB(db *DB) error {
	log.Println("Loading orders from raw_orders table...")
	rows, err := db.conn.Query(context.Background(), "SELECT order_uid, data FROM raw_orders")
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	defer rows.Close()

	c.mu.Lock()
	defer c.mu.Unlock()
	
	count := 0
	for rows.Next() {
		var uid string
		var data []byte
		if err := rows.Scan(&uid, &data); err != nil {
			log.Printf("Skip order (scan error): %v", err)
			continue
		}
		var order Order
		if err := json.Unmarshal(data, &order); err != nil {
			log.Printf("Skip order %s (unmarshal error): %v", uid, err)
			continue
		}
		c.orders[uid] = order
		count++
	}
	log.Printf("Loaded %d orders into cache", count)
	return nil
}

func (c *Cache) Set(uid string, order Order) {
	c.mu.Lock() 
	defer c.mu.Unlock()
	c.orders[uid] = order
}

func (c *Cache) Get(uid string) (Order, bool) {
	c.mu.RLock() 
	defer c.mu.RUnlock()
	order, ok := c.orders[uid]
	return order, ok
}