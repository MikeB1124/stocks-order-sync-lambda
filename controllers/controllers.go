package controllers

import (
	"context"
	"log"

	stockslambdautils "github.com/MikeB1124/stocks-lambda-utils/v2"
	"github.com/MikeB1124/stocks-order-sync-lambda/configuration"
	"github.com/aws/aws-lambda-go/events"
)

func SyncAlpacaOrderWithDB(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Fetch all Alpaca orders
	allAlpacaOrders, err := configuration.AlpacaClient.GetAllAlpacaOrders()
	if err != nil {
		log.Println("Unable to get Alpaca orders")
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "Unable to get Alpaca orders", StatusCode: 500})
	}

	// Update DB with Alpaca orders
	totalUpdatedOrders := 0
	orderIDsNotInDB := []string{}
	for _, order := range allAlpacaOrders {
		formattedOrder := stockslambdautils.FormatAlpacaOrderForDB(&order)
		updateResult, err := configuration.MongoClient.UpateOrder(*formattedOrder)
		if err != nil {
			log.Println("Unable to update order in DB")
			return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "Unable to update order in DB", StatusCode: 500})
		}
		if updateResult.ModifiedCount == 0 {
			orderIDsNotInDB = append(orderIDsNotInDB, order.ID)
		} else {
			totalUpdatedOrders++
		}
	}
	log.Printf("Total alpaca orders %d and total DB updates %d\n", len(allAlpacaOrders), totalUpdatedOrders)
	if len(orderIDsNotInDB) > 0 {
		log.Println("Alpaca Order IDs not found in DB: %v", orderIDsNotInDB)
	}
	return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "OK", StatusCode: 200})
}
