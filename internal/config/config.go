package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	//check if table exist
	isTableExist(db, env)
	return db
}

func isTableExist(db *dynamodb.Client, env Env) {
	w := dynamodb.NewTableExistsWaiter(db)
	err := w.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(env.DynamoTableName),
	},
		30*time.Second,
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
