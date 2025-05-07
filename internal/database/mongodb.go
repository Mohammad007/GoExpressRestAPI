package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func (m *MongoDB) Connect() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}

func (m *MongoDB) Ping() error {
	return m.client.Ping(context.Background(), nil)
}