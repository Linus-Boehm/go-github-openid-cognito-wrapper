package github_openid_wrapper

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

// AuthorizeHandler function description
func AuthorizeHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log and return result
	req := &AuthorizeRequest{}
	if request.HTTPMethod == http.MethodGet {
		err := req.UnmarshalMap(request.QueryStringParameters)
		if err != nil {
			fmt.Println("Error occured")
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
	}else{
		//POST Request
		err := req.Unmarshal([]byte(request.Body))
		if err != nil {
			fmt.Println("Error occured")
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
	}


	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Location": req.GetAuthorizeUrl(),
		},
		StatusCode: 301,
	}, nil
}

// TokenHandler function description
func TokenHandler(ctx *context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log and return result
	jsonItem, err := MakeIDToken(TokenResponse{
		ExpiresAt: 11111,
		AccessToken: "AccessToken",
	})
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