package weather_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flpcamdeveloper-afk/clima_por_cep_google_cloud/internal/weather"
)

func TestWeatherAPIClient_CurrentTempC_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	client := &weather.WeatherAPIClient{HTTPClient: server.Client(), BaseURL: server.URL, APIKey: "test-key"}

	temp, err := client.CurrentTempC(context.Background(), "São Paulo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if temp != 28.5 {
		t.Fatalf("expected temp_c 28.5, got %v", temp)
	}
}

func TestWeatherAPIClient_CurrentTempC_CityNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := &weather.WeatherAPIClient{HTTPClient: server.Client(), BaseURL: server.URL, APIKey: "test-key"}

	_, err := client.CurrentTempC(context.Background(), "Cidade Inexistente")
	if err != weather.ErrCityNotFound {
		t.Fatalf("expected ErrCityNotFound, got %v", err)
	}
}
