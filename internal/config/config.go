package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type Env struct {
	AwsRegion       string `env:"AWS_TARGET_REGION"`
	DynamoTableName string `env:"DYNAMO_TABLE_NAME"`
}

type Conf struct {
	DB  *dynamodb.Client
	Env Env
}

func InitConfig() Conf {
	var localEnv Env
	err := env.Parse(&localEnv)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(localEnv.DynamoTableName)
	return Conf{
		DB:  initDBConn(localEnv),
		Env: localEnv,
	}
}

func initDBConn(env Env) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(options *config.LoadOptions) error {
		options.Region = env.AwsRegion
		return nil
	})
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	db := dynamodb.NewFromConfig(cfg)
	//create new table
	createTable(db, env)

	//check if table ready
	isTableReady(db, env)
	return db
}

func createTable(db *dynamodb.Client, env Env) {
	_, err := db.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(env.DynamoTableName),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func isTableReady(db *dynamodb.Client, env Env) {
	w := dynamodb.NewTableExistsWaiter(db)
	err := w.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(env.DynamoTableName),
	},
		2*time.Minute,
		func(options *dynamodb.TableExistsWaiterOptions) {
			options.MaxDelay = 5 * time.Second
			options.MinDelay = 5 * time.Second
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
func (c *Conf) Free() {

}
