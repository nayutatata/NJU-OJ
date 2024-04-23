package handlers

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburl  string = "mongodb://localhost:27017"
	dbname string = "nayuta"
)

type Handler struct {
	ctx        context.Context
	database   *mongo.Database
	usercoll   *mongo.Collection
	assigncoll *mongo.Collection
	subcoll    *mongo.Collection
	procoll    *mongo.Collection
	gracoll    *mongo.Collection
}

func GetHandler(ctx context.Context) *Handler {
	db := getdb()
	return &Handler{
		ctx:        ctx,
		database:   db,
		usercoll:   db.Collection("users"),
		assigncoll: db.Collection("assignments"),
		subcoll:    db.Collection("submissions"),
		procoll:    db.Collection("problems"),
		gracoll:    db.Collection("graders"),
	}
}
func getdb() *mongo.Database {
	clientOptions := options.Client().ApplyURI(dburl)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongo connected successfully!")
	return client.Database(dbname)
}
