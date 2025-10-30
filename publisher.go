package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/nats-io/stan.go"
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

func main() {
	clusterID := "test-cluster"
	clientID := "order-publisher"
	natsURL := "localhost:4222" 

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer sc.Close()

	log.Println("Connected to NATS Streaming. Publishing orders...")

	for i := 1; i <= 1000; i++ {
		uid := fmt.Sprintf("test_uid_%d", i)
		
		order := Order{
			OrderUID: uid,
			TrackNumber: fmt.Sprintf("TRACK%d", i),
			Entry: "WB",
			Delivery: map[string]interface{}{"name": "Test Name"},
			Payment: map[string]interface{}{"currency": "USD", "amount": 100},
			Items: []map[string]interface{}{{"chrt_id": i, "name": "Test Item", "price": 50}},
			Locale: "en",
			CustomerID: "test_customer",
			DeliveryService: "meest",
			Shardkey: "20",
			SmID: 99,
			DateCreated: time.Now().Format(time.RFC3339),
			OofShard: "1",
		}

		data, _ := json.Marshal(order)

		err = sc.Publish("orders", data)
		if err != nil {
			log.Printf("Failed to publish order %s: %v", uid, err)
		}
	}
	
	log.Println("Finished publishing 1000 orders.")
}