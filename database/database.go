package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
	db     *mongo.Database
	cfg    *Config
}

func NewDatabase(cfg Config) *Database {
	d := Database{
		cfg: &cfg,
	}

	var err error
	d.client, err = mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = d.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	d.db = d.client.Database(cfg.DbName)

	return &d
}

func (m *Database) Close() {
	m.client.Disconnect(context.Background())
}
