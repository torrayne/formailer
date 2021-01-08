package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 503,
			Body:       "Something went wrong :(",
		}, nil
	}

	cc := lc.ClientContext

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, " + cc.Client.AppTitle + "\n" + request.Body,
	}, nil

	// return &events.APIGatewayProxyResponse{
	// 	StatusCode:        200,
	// 	Headers:           map[string]string{"Content-Type": "text/plain"},
	// 	MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
	// 	Body:              "Hello, World!",
	// 	IsBase64Encoded:   false,
	// }, nil
}

func main() {
	lambda.Start(handler)
}
