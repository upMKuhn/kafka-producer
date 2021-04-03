package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/upmkuhn/kafka-trades-demo/pkg/schema"
	"github.com/upmkuhn/kafka-trades-demo/pkg/stocks"
)

var demoDb *pg.DB

func init() {
	demoDb = pg.Connect(&pg.Options{
		User:     "user1",
		Password: "password",
		Database: "kafka-demo",
		Addr:     "127.0.0.1:5432",
	})
	stocks.CreateStockTradeProducer("localhost:9094", "kafka-trades-demo")
}

func main() {

	fmt.Println("Registering Schema")
	schema.ConnectSchemaRegistry("http://localhost:8081")
	stocks.RegisterSchemas()
	stocks.LoadSchemas()

	fmt.Println("Booting Server")
	//defer demoDb.Close()
	//error := stocks.CreateSchema(demoDb)
	//if error != nil {
	//	panic(error)
	//}
	fmt.Println("Created schema")
	defer stocks.ConnectToKafka()
	produceFakeTrades()
}

func produceFakeTrades() {
	fmt.Println("Starting Fake messages")
	for true {
		stocks.PublishStockTrade(makeFakeTrade())
		time.Sleep(time.Duration(int(time.Second) * rand.Intn(10)))
	}
}

func populateDb() {
	//demoDb.Model(trade).Insert()
}

func makeFakeTrade() *stocks.StockTrade {
	SYMBOLS := []string{"TSLA", "APPL", "GOOG"}
	TradeType := []stocks.StockTradeType{stocks.StockBought, stocks.StockSold}
	trade := &stocks.StockTrade{
		ID:         uuid.New().String(),
		SYMBOL:     SYMBOLS[rand.Intn(len(SYMBOLS))],
		QUANTITY:   int64(rand.Intn(70000) + 100),
		TYPE:       TradeType[rand.Intn(2)],
		PRICE:      uint64(rand.Intn(99999999)),
		CREATED_AT: time.Now(),
	}
	return trade
}
