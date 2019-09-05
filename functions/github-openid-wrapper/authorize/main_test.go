package main

import (
	"fmt"
	"github.com/Linus-Boehm/go-github-openid-cognito-wrapper/functions/github-openid-wrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAuthorizeGet(t *testing.T) {
	queryParams := map[string]string{
		"client_id": "test",
		"scope": "openid user:email",
		"state": "ExampleState",
		"response_type": "authorization_code",
	}
	request := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
		QueryStringParameters: queryParams,
	}
	resp, err := github_openid_wrapper.AuthorizeHandler(request)
	assert.NoError(t, err)
	loc, ok := resp.Headers["Location"]
	fmt.Println(loc)
	assert.True(t,ok, "Location not found")

	urlObj, err := url.Parse(loc)
	assert.Equal(t, "github.com", urlObj.Host)
	urlValues := urlObj.Query()
	fmt.Println(urlValues)
	for k, v := range queryParams {
		val, ok := urlValues[k]
		assert.True(t, ok)
		assert.Equal(t,v,strings.Join(val," "))
	}



}
