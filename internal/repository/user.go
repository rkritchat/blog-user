package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserEntity struct {
	Id        string
	Firstname string
	Lastname  string
}

type User interface {
	GetUserByEmail(email string) (*UserEntity, error)
}

type user struct {
	tableName string
	db        *dynamodb.Client
}

func NewUser(db *dynamodb.Client, tableName string) User {
	return &user{
		db:        db,
		tableName: tableName,
	}
}
func (repo user) GetUserByEmail(email string) (*UserEntity, error) {
	out, err := repo.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: email},
		},
	})
	if len(out.Item) == 0 {
		//email is not found
		return nil, nil
	}

	var r UserEntity
	err = attributevalue.UnmarshalMap(out.Item, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
