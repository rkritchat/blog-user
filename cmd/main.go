package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rkritchat/blog-user/internal/config"
	"github.com/rkritchat/blog-user/internal/repository"
	"github.com/rkritchat/blog-user/internal/user"
	"github.com/rkritchat/jsonmask"
	"net/http"
)

var userService user.Service

func main() {
	//init config
	cfg := config.InitConfig()
	defer cfg.Free()

	//init repository
	userRepo := repository.NewUser(cfg.DB, cfg.Env.DynamoTableName)

	userService = user.NewService(userRepo, cfg.Env)
	lambda.Start(handler)
}

func handler(_ context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqLogger(req)
	switch req.HTTPMethod {
	case http.MethodGet:
		return userService.GetUser(req)
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
