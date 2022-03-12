package user

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rkritchat/blog-user/internal/config"
	"github.com/rkritchat/blog-user/internal/repository"
	"net/http"
)

const (
	statusOK     = "SUCCESS"
	statusFailed = "FAILED"
)

type Service interface {
	CreateUser(event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	GetUser(event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
}

type service struct {
	userRepo repository.User
	env      config.Env
}

func NewService(userRepo repository.User, env config.Env) Service {
	return &service{
		userRepo: userRepo,
		env:      env,
	}
}

type CreateUserReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
type CommonResp struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (s service) CreateUser(event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req, err := validateReq(event)
	if err != nil {
		fmt.Printf("validateReq: %v", err)
		return s.toJson(CommonResp{Status: statusFailed}, http.StatusBadRequest)
	}
	fmt.Printf("firstname: %v", req.Firstname)
	return s.toJson(CommonResp{Status: statusOK}, http.StatusOK)
}

func (s service) GetUser(event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	//validate request
	email := event.QueryStringParameters["email"]
	if len(email) == 0 {
		return s.toJson(CommonResp{Status: statusFailed, Message: "email is required"}, http.StatusBadRequest)
	}

	fmt.Printf("email:%v\n", email)
	//get email by id
	entity, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return s.toJson(CommonResp{Status: statusFailed, Message: "internal server error"}, http.StatusInternalServerError)
	}

	if entity == nil {
		return s.toJson(CommonResp{Status: statusFailed, Message: "Email is not found"}, http.StatusBadRequest)
	}
	return s.toJson(entity, http.StatusOK)
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
