package main

import (
	"main/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.ListenAndServe(":8080", nil)
}
