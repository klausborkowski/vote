package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jmoiron/sqlx"
	"github.com/klausborkowski/vote/contract"
	"github.com/klausborkowski/vote/models"
	"github.com/klausborkowski/vote/repository"
	_ "github.com/lib/pq"
)

var url = "https://scroll.blockpi.network/v1/rpc/public"
var adressContract = common.HexToAddress("0xDC3D8318Fbaec2de49281843f5bba22e78338146")

func main() {
	db, err := sqlx.Connect("postgres", "postgres://admin:admin@localhost/ruby?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	Repository := repository.NewRepository(db)
	models.CreateEventsTable(db)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		log.Fatalf("Cannot create cleint: %s", err.Error())
	}
	defer client.Close()
	c, _ := contract.NewContractFilterer(adressContract, client)
	endBlock, _ := client.BlockNumber(ctx)
	startBlock := endBlock - 1024
	f, err := c.FilterMinted(&bind.FilterOpts{Start: startBlock, End: &endBlock}, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	var events []models.Event
	for f.Next() {
		if f == nil {
			fmt.Println("breakdance")
			break
		}
		n := f.Event
		txHash := n.Raw.TxHash
		x, _, err := client.TransactionByHash(ctx, txHash)
		if err != nil {
			log.Fatal(err)
		}
		blockHash := n.Raw.BlockHash
		v, _ := client.BlockByHash(ctx, blockHash)
		// if v == nil {
		// 	fmt.Println("block is nil")
		// 	break
		// }
		fmt.Println("====")
		fmt.Println("hash:", txHash)
		fmt.Printf("time:%v\n", x.Time().UTC())
		//fmt.Println("Time:", time.Unix(int64(v.Time()), 0))

		fmt.Println(n.UserAddress)
		fmt.Println(n.NftIds)
		fmt.Println(n.UserNonce)

		var timeOf uint64
		if v != nil {
			timeOf = v.Time()
		}
		events = append(events, models.ShapeEvent(fmt.Sprint(n.UserAddress), fmt.Sprint(n.NftIds), fmt.Sprint(n.UserNonce), fmt.Sprint(txHash), timeOf))

	}
	err = Repository.InsertEvents(events)
	if err != nil {
		fmt.Errorf("Error occured while insertion:%v", err)
	}
}
