package main

import (
	"log/slog"
	"net/http"

	adapterhttp "doheem-backend/internal/adapter/http"
)

func main() {
	router := adapterhttp.NewRouter()

	slog.Info("Doheem server is running!")

	http.ListenAndServe(":8080", router)
}
