package models

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Address string
	Count   int
}

type Event struct {
	Hash        string `db:"hash"`
	UserAddress string `db:"user_address"`
	Nftids      string `db:"nft_ids"`
	UserNonce   string `db:"user_nonce"`
	Time        uint64 `db:"time"`
}

func ShapeEvent(userAddress, nfts, nonce, hash string, time uint64) Event {
	return Event{
		UserAddress: userAddress,
		Nftids:      nfts,
		UserNonce:   nonce,
		Time:        time,
		Hash:        hash,
	}
}
func CreateEventsTable(db *sqlx.DB) {
	// Check if the table exists
	var exists bool
	err := db.Get(&exists, `SELECT EXISTS (
		SELECT 1
		FROM   information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name = 'events'
	)`)
	if err != nil {
		log.Fatalf("Error checking if table exists: %v", err)
	}
	//SQL IF NOT EXIST CLAUSE
	if !exists {
		// Create the table
		_, err := db.Exec(`CREATE TABLE events (
			hash TEXT,
			user_address TEXT,
			nft_ids TEXT,
			user_nonce TEXT,
			time BIGINT,
			CONSTRAINT hash UNIQUE (hash)
		)`)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
		fmt.Println("Table 'events' created successfully")
	}
}
