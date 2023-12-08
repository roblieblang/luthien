package db

import (
	"context"
	// "log"
	// "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string, ctx context.Context) (*mongo.Client, error) {
    client, err := mongo.NewClient(options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    // ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    // defer cancel()

    err = client.Connect(ctx)
    if err != nil {
		return nil, err
	}

    return client, nil
}