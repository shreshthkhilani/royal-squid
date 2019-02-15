package main

import (
	"log"
	"net/http"
	"github.com/shreshthkhilani/royal-squid/dinners"
)

func main() {
	http.HandleFunc("/api/dinners", dinners.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}