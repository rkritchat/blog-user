package user

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

const (
	statusOK     = "SUCCESS"
	statusFailed = "FAILED"
)

type Service interface {
	CreateUser(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

type CreateUserReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
type CreateUserResp struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (s service) CreateUser(event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req, err := validateReq(event)
	if err != nil {
		fmt.Printf("validateReq: %v", err)
		return s.toJson(CreateUserResp{Status: statusFailed}, http.StatusBadRequest)
	}
	fmt.Printf("firstname: %v", req.Firstname)
	return s.toJson(CreateUserResp{Status: statusOK}, http.StatusOK)
}

func validateReq(event events.APIGatewayProxyRequest) (*CreateUserReq, error) {
	var req CreateUserReq
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		fmt.Printf("invalid request json: %v", err)
		return nil, err
	}
	return &req, nil
}

func (s service) toJson(body interface{}, statusCode int) (*events.APIGatewayProxyResponse, error) {
	b, _ := json.Marshal(body)
	resp := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"content-type": "application/json"},
		Body:       string(b),
	}
	return &resp, nil
}
