package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rkritchat/blog-user/internal/user"
	"github.com/rkritchat/jsonmask"
	"net/http"
)

var userService user.Service

func main() {
	userService = user.NewService()
	lambda.Start(handler)
}

func handler(_ context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqLogger(req)
	switch req.HTTPMethod {
	case http.MethodPost:
		return userService.CreateUser(req)
	default:
		fmt.Printf("method: %v, is not supppport", req.HTTPMethod)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, nil
	}
}

func reqLogger(req events.APIGatewayProxyRequest) {
	m := jsonmask.Init()
	r, err := m.Json([]byte(req.Body))
	if err == nil && r != nil {
		fmt.Println(*r)
	}
}
