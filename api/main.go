package main

import (
	"log"
	"net/http"
	"github.com/shreshthkhilani/royal-squid/dinners"
)

func main() {
	http.HandleFunc("/dinners", dinners.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}