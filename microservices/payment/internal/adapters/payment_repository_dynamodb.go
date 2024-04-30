package adapters

import (
	"context"
	"log"
	"pdi/payment/internal/core/domain"
	"pdi/payment/internal/ports"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PaymentRepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewPaymentRepositoryDynamoDB() ports.PaymentRepository {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal("error to load default config")
	}
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localstack:4566")
	})
	return &PaymentRepositoryDynamoDB{
		client,
	}
}

func (o *PaymentRepositoryDynamoDB) Save(payment *domain.Payment) error {
	item, err := attributevalue.MarshalMap(&payment)
	if err != nil {
		panic(err)
	}
	_, err = o.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("payment"), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}
