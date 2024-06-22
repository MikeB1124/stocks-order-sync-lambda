package alpaca

import (
	"log"

	"github.com/MikeB1124/stocks-order-sync-lambda/configuration"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

var client *alpaca.Client

func init() {
	log.Println("Initializing Alpaca client...")
	configuration := configuration.GetConfig()
	client = alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    configuration.Alpaca.ApiKey,
		APISecret: configuration.Alpaca.ApiSecret,
		BaseURL:   configuration.Alpaca.PaperApiUrl,
	})
}

func GetAllAlpacaOrders() ([]alpaca.Order, error) {
	orders, err := client.GetOrders(alpaca.GetOrdersRequest{
		Nested: true,
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}
