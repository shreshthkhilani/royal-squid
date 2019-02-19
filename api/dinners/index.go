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
	filter := false
	if val, ok := r.URL.Query()["filter"]; ok {
		if len(val) == 1 {
			if val[0] == "1" {
				filter = true
			}
		}
	}
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
	var filteredDinners []Dinner
	for i := 0; i < len(dinners); i++ {
		if dinners[i].DinnerTime.After(time.Now().Add((time.Hour * 6))) && dinners[i].Available != 0 {
			filteredDinners = append(filteredDinners, dinners[i])
		}
	}
	response := make(map[string]interface{})
	if filter {
		response["dinners"] = filteredDinners
	} else {
		response["dinners"] = dinners
	}
	rw, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Fprintf(w, "MarshalIndent")
		return
	}
	fmt.Fprintf(w, string(rw))
	return
}