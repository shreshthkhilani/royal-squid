package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Reservation struct {
	Slots int `bson:"slots" json:"slots"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Dietary string `bson:"dietary" json:"dietary"`
	Confirmed bool `bson:"confirmed" json:"confirmed"`
	DGAE bool `bson:"dgae" json:"dgae"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Dinner struct {
	ID int `bson:"id" json:"id"`
	DinnerTime time.Time `bson:"dinnerTime" json:"dinnerTime"`
	Available  int `bson:"available" json:"available"`
	Reservations []Reservation `bson:"reservations" json:"reservations"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 4 {
		http.Error(w, "Incorrect URI params.", 500)
		return
	}
	dinnerID64, err := strconv.ParseInt(path[3], 0, 64)
	if err != nil {
		http.Error(w, "Incorrect URI params.", 500)
		return
	}
	dinnerID := int(dinnerID64)
    fmt.Println(dinnerID)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var reservation Reservation
	err = json.Unmarshal(b, &reservation)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	reservation.Timestamp = time.Now()
	reservation.Confirmed = false
	reservation.DGAE = true
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
	collection := client.Database("silentdinnerdb").Collection("times")
	ctx, _ = context.WithTimeout(context.Background(), 30 * time.Second)
	cur, err := collection.Find(ctx, bson.D{{"id", dinnerID}})
	if err != nil {
		http.Error(w, "Find.", 500)
		return
	}
	defer cur.Close(ctx)
	var dinner Dinner
	cur.Next(ctx)
	err = cur.Decode(&dinner)
	if err != nil {
		http.Error(w, "Decode.", 500)
		return
	}
	if err = cur.Err(); err != nil {
		http.Error(w, "Err.", 500)
		return
	}
	if (dinner.Available < reservation.Slots) {
		http.Error(w, "This reservation isn't possible.", 500)
		return
	}
	dinner.Available = dinner.Available - reservation.Slots
	dinner.Reservations = append(dinner.Reservations, reservation)
	var newdinner Dinner
	err = collection.FindOneAndReplace(ctx, bson.D{{"id", dinnerID}}, dinner).Decode(&newdinner)
	if err != nil {
		http.Error(w, "Decode 2.", 500)
		return
	}
	response := make(map[string]interface{})
	response["reservation"] = reservation
	rw, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "MarshalIndent.", 500)
		return
	}
	fmt.Fprintf(w, string(rw))
	return
}