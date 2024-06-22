package main

import (
	"github.com/MikeB1124/stocks-order-sync-lambda/controllers"
	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/lambda"
)

var router *lmdrouter.Router

func init() {
	router = lmdrouter.NewRouter("")
	router.Route("POST", "/sync/orders", controllers.HelloWorld)
}

func main() {
	lambda.Start(router.Handler)
}
