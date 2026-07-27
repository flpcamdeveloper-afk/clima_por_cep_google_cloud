// Package handler expõe a lógica de temperatura por CEP como um handler
// HTTP, desacoplado dos clientes concretos de CEP e clima (injetados via
// interface).
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"climacep/internal/cep"
	"climacep/internal/temperature"
	"climacep/internal/weather"
)

const (
	msgInvalidZipcode  = "invalid zipcode"
	msgZipcodeNotFound = "can not find zipcode"
)

var cepFormat = regexp.MustCompile(`^\d{8}$`)

// WeatherHandler resolve um CEP recebido via query string (?cep=) para a
// temperatura atual da cidade correspondente.
type WeatherHandler struct {
	CEPClient     cep.Client
	WeatherClient weather.Client
}

type weatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cepParam := r.URL.Query().Get("cep")
	if !cepFormat.MatchString(cepParam) {
		writeError(w, http.StatusUnprocessableEntity, msgInvalidZipcode)
		return
	}

	location, err := h.CEPClient.Find(r.Context(), cepParam)
	if err != nil {
		if errors.Is(err, cep.ErrNotFound) {
			writeError(w, http.StatusNotFound, msgZipcodeNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	tempC, err := h.WeatherClient.CurrentTempC(r.Context(), location.City)
	if err != nil {
		if errors.Is(err, weather.ErrCityNotFound) {
			writeError(w, http.StatusNotFound, msgZipcodeNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(weatherResponse{
		TempC: tempC,
		TempF: temperature.CToF(tempC),
		TempK: temperature.CToK(tempC),
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}
