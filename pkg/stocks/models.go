package stocks

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// StockTradeType BUY/SALE
type StockTradeType string

const (
	StockSold   StockTradeType = "SALE"
	StockBought                = "BUY"
)

type StockTrade struct {
	ID         string    `pg:"default:gen_random_uuid()"`
	CREATED_AT time.Time `pg:"default:now()"`
	SYMBOL     string
	TYPE       StockTradeType
	QUANTITY   int64
	PRICE      uint64
}

// CreateSchema creates database for models.
func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*StockTrade)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil

}
