package main

import (
	"context"
	handlers "github.com/Linus-Boehm/go-github-openid-cognito-wrapper/functions/github-openid-wrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

func main(){
	r := gin.Default()
	r.GET("/authorize", AuhtorizeWrapHandle)
	r.POST("/authorize", AuhtorizeWrapHandle)
	r.GET("/token", TokenWrapHandle)
	r.POST("/token", TokenWrapHandle)
	log.Fatal(r.Run("0.0.0.0:3002"))
}

func mapQuery(values url.Values) map[string]string{
	params := make(map[string]string)
	for k, v := range values{
		params[k] = strings.Join(v, " ")
	}
	return params
}
func mapHeaders(resp events.APIGatewayProxyResponse, c *gin.Context){
	for k, v := range resp.Headers{
		c.Header(k, v)
	}
}

func TokenWrapHandle( c *gin.Context) {
	req := c.Request
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil{
		panic(err)
	}
	awsReq := events.APIGatewayProxyRequest{
		Body: string(bodyBytes),
		QueryStringParameters: mapQuery(req.URL.Query()),
		HTTPMethod: req.Method,

	}
	ctx := context.Background()
	awsResp, err := handlers.TokenHandler(&ctx, awsReq)
	if err != nil{
		panic(err)
	}
	mapHeaders(awsResp, c)
	c.JSON(awsResp.StatusCode, awsResp.Body)
}

func AuhtorizeWrapHandle( c *gin.Context) {
	req := c.Request
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil{
		panic(err)
	}
	awsReq := events.APIGatewayProxyRequest{
		Body: string(bodyBytes),
		QueryStringParameters: mapQuery(req.URL.Query()),
		HTTPMethod: req.Method,

	}
	awsResp, err := handlers.AuthorizeHandler(awsReq)
	if err != nil{
		panic(err)
	}
	c.Status(awsResp.StatusCode)
	mapHeaders(awsResp, c)

}