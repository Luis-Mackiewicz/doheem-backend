package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Doheem server is running! 🚀")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)

	fmt.Println("Doheem server is running! 🚀")

	http.ListenAndServe(":8080", mux)
}
