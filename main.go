package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()
	log.Println("Connected to DB")

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize DB tables: %v", err)
	}
	log.Println("DB tables initialized")

	cache := NewCache()
	if err := cache.LoadFromDB(db); err != nil {
		log.Printf("Cache restore warning: %v", err)
	} else {
		log.Println("Cache restored from DB")
	}

	sc, err := NewNATSClient("order-service", "test-cluster", "localhost:4222") 
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer sc.Close()
	log.Println("Connected to NATS Streaming")

	if err := SubscribeOrders(sc, db, cache); err != nil {
		log.Fatalf("Failed to subscribe to orders channel: %v", err)
	}
	log.Println("Subscribed to 'orders' channel")

	router := SetupRoutes(cache)
	go func() {
		log.Println("HTTP server starting on :8080")
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Printf("HTTP server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down service...")
	}