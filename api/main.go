package main

import (
	"log"
	"net/http"
	"github.com/shreshthkhilani/royal-squid/reserve"
	"github.com/shreshthkhilani/royal-squid/confirm"
)

func main() {
	http.HandleFunc("/api/reserve/", reserve.Handler)
	http.HandleFunc("/api/confirm/", confirm.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}