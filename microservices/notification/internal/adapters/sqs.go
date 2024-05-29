package adapters

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/julianojj/pdi/notification/internal/ports"
)

type SQS struct {
	client *sqs.Client
}

func NewSQS() ports.Queue {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal("error to load default config")
	}
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String("http://localstack:4566")
	})
	return &SQS{
		client,
	}
}

func (s *SQS) Consume(queueName string, callback func(args []byte) error) error {
	for {
		ouptut, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		for _, msg := range ouptut.Messages {
			body := []byte(*msg.Body)
			err := callback(body)
			if err == nil {
				s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueName),
					ReceiptHandle: msg.ReceiptHandle,
				})
				continue
			}
			fmt.Println(err)
		}
	}
}