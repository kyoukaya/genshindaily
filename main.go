package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kyoukaya/genshindaily/internal/genshindaily"
)

func main() {
	lambda.Start(genshindaily.HandleMessage)
}
