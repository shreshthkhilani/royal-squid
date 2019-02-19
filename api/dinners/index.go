package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Dinner struct {
	ID int `bson:"id" json:"id"`
	DinnerTime time.Time `bson:"dinnerTime" json:"dinnerTime"`
	Available  int `bson:"available" json:"available"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	client, err := mongo.Connect(ctx, "mongodb://silentdinneruser:" +
		os.Getenv("MONGODB_PW") +
		"@silentdinnercluster-shard-00-00-m8j4f.mongodb.net:27017," +
		"silentdinnercluster-shard-00-01-m8j4f.mongodb.net:27017," +
		"silentdinnercluster-shard-00-02-m8j4f.mongodb.net:27017" +
		"/test?ssl=true&replicaSet=silentdinnercluster-shard-0&" +
		"authSource=admin&retryWrites=true")
	if err != nil {
		fmt.Fprintf(w, "Error connecting.")
		return
	}
	collection := client.Database("silentdinnerdb").Collection("times")
	ctx, _ = context.WithTimeout(context.Background(), 30 * time.Second)
	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		fmt.Fprintf(w, "Find")
	}
	defer cur.Close(ctx)
	var dinners []Dinner
	for cur.Next(ctx) {
		var result Dinner
		err := cur.Decode(&result)
		if err != nil {
			fmt.Fprintf(w, "Decode")
			return
		}
		dinners = append(dinners, result)
	}
	if err := cur.Err(); err != nil {
		fmt.Fprintf(w, "Err")
		return
	}
	response := make(map[string]interface{})
	response["dinners"] = dinners
	rw, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Fprintf(w, "MarshalIndent")
		return
	}
	fmt.Fprintf(w, string(rw))
	return
}