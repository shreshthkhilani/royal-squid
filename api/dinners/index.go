package dinners

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Reservation struct {
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Dinner struct {
	DinnerTime time.Time `bson:"dinnerTime" json:"dinnerTime"`
	Available  int `bson:"available" json:"available"`
	Reservations []Reservation `bson:"reservations" json:"reservations"`
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
	}
	collection := client.Database("silentdinnerdb").Collection("times")
	ctx, _ = context.WithTimeout(context.Background(), 30 * time.Second)
	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Print("Find")
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	var dinners []Dinner
	for cur.Next(ctx) {
		var result Dinner
		err := cur.Decode(&result)
		if err != nil {
			log.Print("Decode")
			log.Fatal(err)
		}
		dinners = append(dinners, result)
	}
	if err := cur.Err(); err != nil {
		log.Print("Err")
		log.Fatal(err)
	}
	response := make(map[string]interface{})
	response["dinners"] = dinners
	rw, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Print("MarshalIndent")
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(rw))
}