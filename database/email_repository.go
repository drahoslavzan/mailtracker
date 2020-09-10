package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type emailRepository struct {
	col *mongo.Collection
}

type EmailRepository interface {
	TrackSeen(id string)
}

func (m *Database) NewEmailRepository() EmailRepository {
	repo := emailRepository{
		col: m.db.Collection(m.cfg.EmailCol),
	}

	indexModel := repo.getIndexModel()
	_, err := repo.col.Indexes().CreateMany(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}

	return &repo
}

func (m *emailRepository) TrackSeen(id string) {
	ctxLogger := log.WithFields(log.Fields{
		"package":  "database",
		"location": "track_seen",
	})

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctxLogger.Warnf("invalid id: %v", id)
		return
	}

	filter := bson.M{"_id": oid}

	res := m.col.FindOneAndUpdate(context.Background(), filter, bson.M{
		"$inc": bson.M{
			"seen": 1,
		},
		"$set": bson.M{
			"lastSeen": time.Now(),
		},
	})

	err = res.Err()
	if res.Err() != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}
		ctxLogger.Warnf("id not found: %v", id)
		return
	}
}

func (m *emailRepository) getIndexModel() []mongo.IndexModel {
	models := []mongo.IndexModel{
		{
			Keys: bson.M{"seen": 1},
		},
		{
			Keys: bson.M{"lastSeen": 1},
		},
	}

	return models
}
