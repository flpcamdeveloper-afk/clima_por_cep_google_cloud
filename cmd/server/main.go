package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"climacep/internal/cep"
	"climacep/internal/handler"
	"climacep/internal/weather"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	h := &handler.WeatherHandler{
		CEPClient:     cep.NewViaCEPClient(httpClient),
		WeatherClient: weather.NewWeatherAPIClient(httpClient, apiKey),
	}

	mux := http.NewServeMux()
	mux.Handle("/", h)

	// Cloud Run injeta a porta a ser usada via a variável de ambiente PORT.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
