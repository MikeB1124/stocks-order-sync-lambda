package controllers

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func HelloWorld(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("HELLO WORLD!!")

	return createResponse(Response{Message: "HELLO WORLD!!!", StatusCode: 200})
}
