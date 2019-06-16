package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"time"
)

// Login struct
type Login struct {
	Email string `json:"email"`
	Plate string `json:"plate"`
}

// Verify struct
type Verify struct {
	Code string `json:"code"`
}

// Message struct
type Message struct {
	Message string `json:"message"`
}

// Identity struct
type Identity struct {
	Ident struct {
		ID           string `json:"id"`
		Registations []struct {
			Plate string `json:"plate"`
		} `json:"registrations"`
	} `json:"identity"`
}

// LoginIdent struct
type LoginIdent struct {
	ID          string
	RequestTime time.Time
}

type GeneratedCode struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

// Handler what kind of request is it
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	message := Message{}
	if request.Resource == "/login" {
		resp, err := LoginService(request.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("login err: %v", err))
			return events.APIGatewayProxyResponse{}, err
		}

		r, err := json.Marshal(resp)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		message.Message = string(r)
	}

	if request.Resource == "/verify" {
		resp, err := VerifyService(request.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("verify err: %v", err))
			return events.APIGatewayProxyResponse{}, err
		}
		message.Message = resp
	}

	m, err := json.Marshal(message)
	if err != nil {
		fmt.Println(fmt.Sprintf("message marshall err: %v", err))
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(m),
		StatusCode: 200,
	}, nil
}

func StoreData(data []byte, id int) error {
	a := string(data)
	fmt.Println(fmt.Sprintf("Store Data: %s, ID: %v", a, id))

	return nil
}

func RetrieveData(id int) []byte {
	return []byte{}
}
