// Package weather consulta a temperatura atual de uma cidade usando a
// WeatherAPI (https://www.weatherapi.com).
package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// ErrCityNotFound indica que a cidade informada não foi encontrada pelo
// provedor de clima.
var ErrCityNotFound = errors.New("city not found")

// Client consulta a temperatura atual (em Celsius) de uma cidade.
type Client interface {
	CurrentTempC(ctx context.Context, city string) (float64, error)
}

// WeatherAPIClient implementa Client usando a WeatherAPI.
type WeatherAPIClient struct {
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string
}

// NewWeatherAPIClient cria um WeatherAPIClient pronto para uso.
func NewWeatherAPIClient(httpClient *http.Client, apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{
		HTTPClient: httpClient,
		BaseURL:    "https://api.weatherapi.com/v1",
		APIKey:     apiKey,
	}
}

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *WeatherAPIClient) CurrentTempC(ctx context.Context, city string) (float64, error) {
	endpoint := fmt.Sprintf("%s/current.json?key=%s&q=%s", c.BaseURL, url.QueryEscape(c.APIKey), url.QueryEscape(city))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return 0, ErrCityNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi: unexpected status code %d", resp.StatusCode)
	}

	var body weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, err
	}

	return body.Current.TempC, nil
}
