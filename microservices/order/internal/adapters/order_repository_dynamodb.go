package adapters

import (
	"context"
	"log"
	"pdi/order/internal/core/domain"
	"pdi/order/internal/ports"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type OrderRepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewOrderRepositoryDynamoDB() ports.OrderRepository {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal("error to load default config")
	}
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localstack:4566")
	})
	return &OrderRepositoryDynamoDB{
		client,
	}
}

func (o *OrderRepositoryDynamoDB) SaveOrder(order *domain.Order) error {
	item, err := attributevalue.MarshalMap(&order)
	if err != nil {
		panic(err)
	}
	_, err = o.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("order"), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (o *OrderRepositoryDynamoDB) GetOrder(orderID string) (*domain.Order, error) {
	var order *domain.Order
	response, err := o.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: orderID,
			},
		},
		TableName: aws.String("order"),
	})
	if err != nil {
		return nil, err
	}
	if err := attributevalue.UnmarshalMap(response.Item, &order); err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderRepositoryDynamoDB) UpdateOrder(order *domain.Order) error {
	_, err := o.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: order.ID,
			},
		},
		TableName:        aws.String("order"),
		UpdateExpression: aws.String("SET #status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{
				Value: order.Status,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
