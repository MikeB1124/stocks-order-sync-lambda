package controllers

import (
	"context"
	"log"

	stockslambdautils "github.com/MikeB1124/stocks-lambda-utils/v2"
	"github.com/MikeB1124/stocks-order-sync-lambda/configuration"
	"github.com/aws/aws-lambda-go/events"
)

func SyncAlpacaOrderWithDB(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get all open trades from DB
	openTrades, err := configuration.MongoClient.GetOpenTrades()
	if err != nil {
		log.Println("Unable to get open trades from DB")
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "Unable to get open trades from DB", StatusCode: 500})
	}

	var openTradeUpdates []stockslambdautils.AlpacaTrade
	for _, trade := range openTrades {
		alpacaOrder, err := configuration.AlpacaClient.GetAlpacaOrderByID(trade.Order.ID)
		if err != nil {
			log.Printf("Unable to get Alpaca trade with ID: %s\n", trade.Order.ID)
			continue
		}
		formattedOrder := stockslambdautils.FormatAlpacaOrderForDB(alpacaOrder)
		trade.Order = formattedOrder
		openTradeUpdates = append(openTradeUpdates, trade)
	}

	updateResult, err := configuration.MongoClient.BulkUpdateTrades(openTradeUpdates)
	if err != nil {
		log.Println("Unable to update trades in DB")
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "Unable to update trades in DB", StatusCode: 500})
	}
	log.Printf("Total trades updated: %d\n", updateResult.ModifiedCount)
	return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "OK", StatusCode: 200})
}
