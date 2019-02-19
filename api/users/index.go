package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type User struct {
	Email string `bson:"email" json:"email"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var user User
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user.Timestamp = time.Now()
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	client, err := mongo.Connect(ctx, "mongodb://silentdinneruser:" +
		os.Getenv("MONGODB_PW") +
		"@silentdinnercluster-shard-00-00-m8j4f.mongodb.net:27017," +
		"silentdinnercluster-shard-00-01-m8j4f.mongodb.net:27017," +
		"silentdinnercluster-shard-00-02-m8j4f.mongodb.net:27017" +
		"/test?ssl=true&replicaSet=silentdinnercluster-shard-0&" +
		"authSource=admin&retryWrites=true")
	if err != nil {
		http.Error(w, "Error connecting.", 500)
		return
	}
	collection := client.Database("silentdinnerdb").Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Insert.", 500)
		return
	}
	response := make(map[string]interface{})
	response["user"] = user
	rw, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "MarshalIndent.", 500)
		return
	}
	fmt.Fprintf(w, string(rw))
	return
}