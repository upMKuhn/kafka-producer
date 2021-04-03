package stocks

import (
	"fmt"

	"github.com/upmkuhn/kafka-trades-demo/pkg/schema"
	proto_stocks "github.com/upmkuhn/kafka-trades-demo/proto/build/go/stocks"
	"google.golang.org/protobuf/proto"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Topic written to
const Topic string = "demo-stocks-trades"

var producer *kafka.Producer = nil

// CreateStockTradeProducer created producer
func CreateStockTradeProducer(bostrapServers string, clientID string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bostrapServers,
		"client.id":         clientID,
		"acks":              "all",
	})
	if err != nil {
		panic(err)
	}
	producer = p
}

// ConnectToKafka connectes to kafka
func ConnectToKafka() {
	defer producer.Close()
	defer producerEventLoop()
}

// PublishStockTrade Serializes to Protobuf and publishes Trades
func PublishStockTrade(trade *StockTrade) {
	topic := Topic
	key, value := toMessage(trade)
	//delivery_chan := make(chan kafka.Event, 10000)
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            marshalMessage(key, true),
		Value:          marshalMessage(value, false),
	}, nil)
	fmt.Printf("Published trade for symbol %s\n", trade.SYMBOL)
}

// RegisterSchemas updates schema registry
func RegisterSchemas() {
	schema.RegisterSchemaFromFile(Topic, false, "proto/stocks/stock_trades.proto")
	schema.RegisterSchemaFromFile(Topic, true, "proto/stocks/stock_trades.key.proto")
}

// Loads schemas for producer
func LoadSchemas() {
	schema.LoadSchemas([]string{Topic})
}

func producerEventLoop() {
	// Delivery report handler for produced messages
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

func toMessage(trade *StockTrade) (*proto_stocks.Key, *proto_stocks.StockTrade) {

	messageKey := proto_stocks.Key{
		Id: trade.SYMBOL + "-" + trade.CREATED_AT.UTC().Format("02-01-2006"),
	}
	messageValue := proto_stocks.StockTrade{
		Id:        trade.ID,
		CreatedAt: trade.CREATED_AT.UnixNano(),
		Symbol:    trade.SYMBOL,
		Type:      string(trade.TYPE),
		Quantity:  trade.QUANTITY,
		Price:     trade.PRICE,
	}

	return &messageKey, &messageValue
}

func marshalMessage(message proto.Message, isKey bool) []byte {
	data, err := schema.MarschalProtobuf(Topic, isKey, message)
	if err != nil {
		panic(err)
	}
	return data
}
