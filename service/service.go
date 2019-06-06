package service

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(fmt.Sprintf("Request: %v", request))

	return events.APIGatewayProxyResponse{
		Body: "",
		StatusCode: 200,
	}, nil
}