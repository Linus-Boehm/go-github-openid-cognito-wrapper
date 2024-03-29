package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// UserinfoHandler function description
func UserinfoHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  	// Log and return result
	jsonItem, err := json.Marshal(map[string]string{"msg": "Userinfo invoked successfully"})
    if err != nil {
        fmt.Println("Error occured")
        return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil 
    }

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body:       string(jsonItem),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(UserinfoHandler)
}