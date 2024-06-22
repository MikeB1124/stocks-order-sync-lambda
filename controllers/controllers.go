package controllers

import (
	"context"
	"log"

	"github.com/MikeB1124/stocks-order-sync-lambda/alpaca"
	"github.com/MikeB1124/stocks-order-sync-lambda/db"
	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func SyncAlpacaOrderWithDB(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Fetch all Alpaca orders
	allAlpacaOrders, err := alpaca.GetAllAlpacaOrders()
	if err != nil {
		log.Println("Unable to get Alpaca orders")
		return createResponse(Response{Message: "Unable to get Alpaca orders", StatusCode: 500})
	}

	// Update DB with Alpaca orders
	totalUpdatedOrders := 0
	orderIDsNotInDB := []string{}
	for _, order := range allAlpacaOrders {
		updateResult, err := db.UpateOrder(order)
		if err != nil {
			log.Println("Unable to update order in DB")
			return createResponse(Response{Message: "Unable to update order in DB", StatusCode: 500})
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
	return createResponse(Response{Message: "OK", StatusCode: 200})
}
