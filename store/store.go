package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/my_projects/ff-draft-dashboard-api/api"
)

// Store

type M api.M

type Store struct {
	ctx      context.Context
	client   *mongo.Client
	database *mongo.Database
	dbName   string
}

func NewStore(mongoDBName, mongoHost, mongoUser, mongoPwd string) (store *Store, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.%s.mongodb.net/?retryWrites=true&w=majority", mongoUser, mongoPwd, mongoHost)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var client *mongo.Client
	if client, err = mongo.Connect(ctx, opts); err != nil {
		return
	}
	store = &Store{
		ctx:      ctx,
		client:   client,
		database: client.Database(mongoDBName),
		dbName:   mongoDBName,
	}

	return
}

func (s *Store) C(colName string) *mongo.Collection {
	return s.database.Collection(colName)
}

func (s *Store) Close() {
	s.client.Disconnect(s.ctx)
}

func (s *Store) PlayerReportsCol() *mongo.Collection {
	return s.database.Collection("playerReports")
}
