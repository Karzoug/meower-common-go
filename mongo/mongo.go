package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func New(ctx context.Context, cfg Config, appName string) (client *mongo.Client, closeFn func(context.Context) error, err error) {
	client, err = mongo.Connect(
		options.Client().
			ApplyURI(cfg.URI).
			SetAppName(appName))
	// TODO: otelmongo does not support mongo-driver v2 still
	// options.Client().SetMonitor(otelmongo.NewMonitor())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	return client, client.Disconnect, nil
}
