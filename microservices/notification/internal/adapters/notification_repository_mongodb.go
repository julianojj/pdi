package adapters

import (
	"context"

	"github.com/julianojj/pdi/notification/internal/ports"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepositoryMongoDB struct {
	client *mongo.Client
}

func NewNotificationMongoBD() ports.NotificationRepository {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://juliano:12345678@mongodb:27017/"))
	if err != nil {
		panic(err)
	}
	return &NotificationRepositoryMongoDB{
		client,
	}
}

func (n *NotificationRepositoryMongoDB) Save(notification map[string]any) error {
	collection := n.client.Database("ms").Collection("notification")
	_, err := collection.InsertOne(context.TODO(), notification)
	if err != nil {
		return err
	}
	return nil
}
