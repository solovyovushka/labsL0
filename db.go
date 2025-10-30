package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB() (*DB, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:qwerty11@localhost:5432/labsL0")
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}

func (db *DB) InitDB() error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS raw_orders (
		order_uid VARCHAR(255) PRIMARY KEY,
		data JSONB NOT NULL
	);`
	log.Println("Initializing raw_orders table...")
	_, err := db.conn.Exec(context.Background(), createTableSQL)
	return err
}

func (db *DB) Close() {
	db.conn.Close(context.Background())
}