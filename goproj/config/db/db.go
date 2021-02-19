package db

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID             string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string `json:"name,omitempty" bson:"name,omitempty"`
	Surname        string `json:"surname,omitempty" bson:"surname,omitempty"`
	Email          string `json:"email,omitempty" bson:"email,omitempty"`
	HashedPassword string `json:"hashedPassword,omitempty" bson:"hashedPassword,omitempty"`
	Created        string `json:"created,omitempty" bson:"created,omitempty"`
}

func GetUserDBCollection() (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("clientoptions type:", reflect.TypeOf(clientOptions), "\n")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}
	// Access a MongoDB collection through a database
	collection := client.Database("golang").Collection("users")
	return collection, nil
}
func GetSocketDBCollection() (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("clientoptions type:", reflect.TypeOf(clientOptions), "\n")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}
	// Access a MongoDB collection through a database
	collection := client.Database("golang").Collection("socket_data")
	return collection, nil
}
