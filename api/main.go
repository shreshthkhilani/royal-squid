package main

import (
	"log"
	"net/http"
	"github.com/shreshthkhilani/royal-squid/reserve"
)

func main() {
	http.HandleFunc("/api/reserve/", reserve.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}