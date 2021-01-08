package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var body map[string]interface{}
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       http.StatusText(http.StatusBadRequest),
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("Hello, World!\n%+v", body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
