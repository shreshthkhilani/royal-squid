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
	ReservationID int `bson:"reservationId" json:"reservationId"`
	DinnerID int `bson:"dinnerId" json:"dinnerId"`
	Slots int `bson:"slots" json:"slots"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Dietary string `bson:"dietary" json:"dietary"`
	OTP string `bson:"otp" json:"otp"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
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
	collection := client.Database("silentdinnerdb").Collection("reservations")
	count64, err := collection.CountDocuments(ctx, bson.D{{}});
	if err != nil {
		http.Error(w, "Count.", 500)
		log.Fatal(err)
		return
	}
	reservation.ReservationID = int(count64) + 1
	reservation.OTP = getOTP(4)
	reservation.Timestamp = time.Now()
	err = send(reservation.Email, reservation.OTP)
	if err != nil {
		http.Error(w, "Unable to send email.", 500)
		log.Fatal(err)
		return
	}
	_, err = collection.InsertOne(ctx, reservation)
	if err != nil {
		http.Error(w, "Insert.", 500)
		log.Fatal(err)
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