#!/bin/bash

docker stop stocks-order-sync-lambda
docker rm stocks-order-sync-lambda
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
docker build -t stocks-lambda-image .
docker run --name stocks-order-sync-lambda -p 9000:8080 --env-file .env stocks-lambda-image