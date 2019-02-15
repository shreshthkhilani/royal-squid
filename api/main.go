package main

import (
	"log"
	"net/http"
	"github.com/shreshthkhilani/royal-squid/times"
)

func main() {
	http.HandleFunc("/times", times.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}