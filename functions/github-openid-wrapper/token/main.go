package main

import (
	"github.com/Linus-Boehm/go-github-openid-cognito-wrapper/functions/github-openid-wrapper"

	"github.com/aws/aws-lambda-go/lambda"
)



func main() {
	lambda.Start(github_openid_wrapper.TokenHandler)
}