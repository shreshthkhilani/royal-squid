package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/badoux/checkmail"
)

type Reservation struct {
	DinnerID int `bson:"dinnerId" json:"dinnerId"`
	Slots int `bson:"slots" json:"slots"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Dietary string `bson:"dietary" json:"dietary"`
	Confirmed bool `bson:"confirmed" json:"confirmed"`
	DGAE bool `bson:"dgae" json:"dgae"`
	OTP string `bson:"otp" json:"otp"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Dinner struct {
	ID int `bson:"id" json:"id"`
	DinnerTime time.Time `bson:"dinnerTime" json:"dinnerTime"`
	Available  int `bson:"available" json:"available"`
	Reservations []Reservation `bson:"reservations" json:"reservations"`
}

func send(to string, otp string) error {
	from := "silentencounter@shreshthkhilani.com"
	pass := os.Getenv("SMTP_PW")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: silent&counter confirmation code\n\n" +
		"use this code to confirm your dinner: " + otp

	return smtp.SendMail("smtp.mailgun.org:587",
		smtp.PlainAuth("", from, pass, "smtp.mailgun.org"),
		from, []string{to}, []byte(msg))
}

func getOTP(n int) string {
	var letterRunes = []rune("0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Handler(w http.ResponseWriter, r *http.Request) {
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
	reservation.DGAE = false
	err = checkmail.ValidateFormat(reservation.Email)
	if err != nil {
		http.Error(w, "Invalid Email.", 500)
		return
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
		http.Error(w, "Error connecting.", 500)
		return
	}
	collection := client.Database("silentdinnerdb").Collection("times")
	ctx, _ = context.WithTimeout(context.Background(), 30 * time.Second)
	cur, err := collection.Find(ctx, bson.D{{"id", reservation.DinnerID}})
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
	reservation.OTP = getOTP(4)
	err = send(reservation.Email, reservation.OTP)
	if err != nil {
		http.Error(w, "Unable to send email.", 500)
		log.Fatal(err)
		return
	}
	dinner.Available = dinner.Available - reservation.Slots
	dinner.Reservations = append(dinner.Reservations, reservation)
	var newdinner Dinner
	err = collection.FindOneAndReplace(ctx, bson.D{{"id", reservation.DinnerID}}, dinner).Decode(&newdinner)
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