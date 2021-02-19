package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBTool struct {
	DB     *mongo.Database
	Client *mongo.Client
	Ctx    context.Context
}

const DFTimeout time.Duration = 20 * time.Second

func Create(dburl, dbname string) (*DBTool, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dburl))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DFTimeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbname)

	return &DBTool{DB: db, Client: client, Ctx: ctx}, nil
}

func (dbt *DBTool) Close() {
	dbt.Client.Disconnect(dbt.Ctx)
}

func FindOptions() *options.FindOptions {
	df := options.Find()
	df.SetLimit(1000)
	return df
}

func InsertOptions() *options.InsertManyOptions {
	return options.InsertMany()
}

func UpdateOption() *options.UpdateOptions {
	return options.Update()
}
