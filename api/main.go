package main

import (
	"log"
	"net/http"
	"github.com/shreshthkhilani/royal-squid/dinners"
	"github.com/shreshthkhilani/royal-squid/reserve"
	"github.com/shreshthkhilani/royal-squid/confirm"
	"github.com/shreshthkhilani/royal-squid/users"
)

func main() {
	http.HandleFunc("/api/dinners/", dinners.Handler)
	http.HandleFunc("/api/reserve/", reserve.Handler)
	http.HandleFunc("/api/confirm/", confirm.Handler)
	http.HandleFunc("/api/users/", users.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}