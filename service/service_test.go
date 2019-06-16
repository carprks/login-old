package service_test

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	if os.Getenv("IDENTITY_URI") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	login, err := json.Marshal(service.Login{
		Email: os.Getenv("TEST_EMAIL"),
		Plate: os.Getenv("TEST_PLATE"),
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("marshal login err: %v", err))
	}

	verify, err := json.Marshal(service.Verify{
		Code: "1234",
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("marshal verify err: %v", err))
	}

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/login",
				Body:     string(login),
			},
			expect: `{"message":"check email for a code"}`,
			err:    nil,
		},
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/verify",
				Body:     string(verify),
			},
			expect: `{"message":"login failed"}`,
			err:    nil,
		},
	}

	for _, test := range tests {
		response, err := service.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
