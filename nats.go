package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
)

func NewNATSClient(clientID, clusterID, url string) (stan.Conn, error) {
	return stan.Connect(clusterID, clientID, stan.NatsURL(url))
}

func SubscribeOrders(sc stan.Conn, db *DB, cache *Cache) error {
	_, err := sc.Subscribe("orders", func(msg *stan.Msg) {
		var order Order
		
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Invalid JSON in message: %v", err)
			return 
		}

		if order.OrderUID == "" {
			log.Printf("Validation FAILED: Message missing required 'order_uid'.")
			return
		}
		
		data, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to re-marshal order %s for DB: %v", order.OrderUID, err)
			return
		}

		_, err = db.conn.Exec(context.Background(),
			`INSERT INTO raw_orders (order_uid, data) 
			 VALUES ($1, $2)
			 ON CONFLICT (order_uid) DO UPDATE SET data = EXCLUDED.data`,
			order.OrderUID, data)
			
		if err != nil {
			log.Printf("Failed to save order %s to DB: %v", order.OrderUID, err)
			return 
		}
		
		cache.Set(order.OrderUID, order)
		
		if err := msg.Ack(); err != nil {
			log.Printf("Failed to acknowledge message for order %s: %v", order.OrderUID, err)
		}
	}, stan.DurableName("order-service-durable"), stan.SetManualAckMode(), stan.StartWithLastReceived())
	
	return err
}