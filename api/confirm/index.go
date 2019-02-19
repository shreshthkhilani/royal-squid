package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Reservation struct {
	ReservationID int `bson:"reservationId" json:"reservationId"`
	DinnerID int `bson:"dinnerId" json:"dinnerId"`
	Slots int `bson:"slots" json:"slots"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Dietary string `bson:"dietary" json:"dietary"`
	OTP string `bson:"otp" json:"otp"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Dinner struct {
	ID int `bson:"id" json:"id"`
	DinnerTime time.Time `bson:"dinnerTime" json:"dinnerTime"`
	Available  int `bson:"available" json:"available"`
	Reservations []Reservation `bson:"reservations" json:"reservations"`
}

func send(r Reservation, dt time.Time) error {
	from := "silentencounter@shreshthkhilani.com"
	to := r.Email
	pass := os.Getenv("SMTP_PW")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: silent&counter dinner confirmed\n\n" +
		"your dinner is confirmed––see you!\n\n" +
		dt.Format("Mon, Jan _2") + "\n" +
		strconv.Itoa(r.Slots) + "\n" + r.Name + "\n" + r.Email + "\n" + r.Dietary

	return smtp.SendMail("smtp.mailgun.org:587",
		smtp.PlainAuth("", from, pass, "smtp.mailgun.org"),
		from, []string{to}, []byte(msg))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var res Reservation
	err = json.Unmarshal(b, &res)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Connect to DB
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
	// Get dinner object
	ctx, _ = context.WithTimeout(context.Background(), 30 * time.Second)
	dinnersCollection := client.Database("silentdinnerdb").Collection("times")
	cur, err := dinnersCollection.Find(ctx, bson.D{{"id", res.DinnerID}})
	if err != nil {
		http.Error(w, "Find.", 500)
		fmt.Println(err)
		return
	}
	defer cur.Close(ctx)
	var dinner Dinner
	cur.Next(ctx)
	err = cur.Decode(&dinner)
	if err != nil {
		http.Error(w, "Decode.", 500)
		fmt.Println(err)
		return
	}
	if err = cur.Err(); err != nil {
		http.Error(w, "Err.", 500)
		fmt.Println(err)
		return
	}
	// Get reservation object
	ctx, _ = context.WithTimeout(context.Background(), 30 * time.Second)
	reservationsCollection := client.Database("silentdinnerdb").Collection("reservations")
	curr, err := reservationsCollection.Find(ctx, bson.D{{"reservationId", res.ReservationID}})
	if err != nil {
		http.Error(w, "Find.", 500)
		fmt.Println(err)
		return
	}
	defer curr.Close(ctx)
	var reservation Reservation
	curr.Next(ctx)
	err = curr.Decode(&reservation)
	if err != nil {
		http.Error(w, "Decode res.", 500)
		fmt.Println(err)
		return
	}
	if err = curr.Err(); err != nil {
		http.Error(w, "Err res.", 500)
		fmt.Println(err)
		return
	}
	// Check if reservation can be made
	if dinner.Available < reservation.Slots {
		http.Error(w, "Not enough availablity.", 500)
		return
	}
	if res.OTP != reservation.OTP {
		http.Error(w, "Wrong OTP.", 500)
		return
	}
	// Add reservation to dinner
	dinner.Reservations = append(dinner.Reservations, reservation)
	dinner.Available = dinner.Available - reservation.Slots
	err = send(reservation, dinner.DinnerTime)
	if err != nil {
		http.Error(w, "Unable to send email.", 500)
		return
	}
	var newdinner Dinner
	err = dinnersCollection.FindOneAndReplace(ctx, bson.D{{"id", reservation.DinnerID}}, dinner).Decode(&newdinner)
	if err != nil {
		http.Error(w, "Decode 2.", 500)
		return
	}
	// Return confirmed reservation
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