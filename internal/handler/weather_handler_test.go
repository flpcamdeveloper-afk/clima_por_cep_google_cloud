package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"climacep/internal/cep"
	"climacep/internal/handler"
	"climacep/internal/weather"
)

type fakeCEPClient struct {
	location cep.Location
	err      error
}

func (f fakeCEPClient) Find(ctx context.Context, cepCode string) (cep.Location, error) {
	return f.location, f.err
}

type fakeWeatherClient struct {
	tempC float64
	err   error
}

func (f fakeWeatherClient) CurrentTempC(ctx context.Context, city string) (float64, error) {
	return f.tempC, f.err
}

func newRequest(cepParam string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/?cep="+cepParam, nil)
	return req
}

func TestWeatherHandler_InvalidZipcodeFormat(t *testing.T) {
	h := &handler.WeatherHandler{
		CEPClient:     fakeCEPClient{},
		WeatherClient: fakeWeatherClient{},
	}

	cases := []string{"", "123", "1234567890", "abcdefgh", "0100100a"}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, newRequest(c))

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("cep %q: expected status 422, got %d", c, rec.Code)
		}
		if rec.Body.String() != "invalid zipcode" {
			t.Fatalf("cep %q: expected body %q, got %q", c, "invalid zipcode", rec.Body.String())
		}
	}
}

func TestWeatherHandler_ZipcodeNotFound(t *testing.T) {
	h := &handler.WeatherHandler{
		CEPClient:     fakeCEPClient{err: cep.ErrNotFound},
		WeatherClient: fakeWeatherClient{},
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, newRequest("99999999"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	if rec.Body.String() != "can not find zipcode" {
		t.Fatalf("expected body %q, got %q", "can not find zipcode", rec.Body.String())
	}
}

func TestWeatherHandler_WeatherCityNotFound(t *testing.T) {
	h := &handler.WeatherHandler{
		CEPClient:     fakeCEPClient{location: cep.Location{City: "Cidade Desconhecida"}},
		WeatherClient: fakeWeatherClient{err: weather.ErrCityNotFound},
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, newRequest("01001000"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	if rec.Body.String() != "can not find zipcode" {
		t.Fatalf("expected body %q, got %q", "can not find zipcode", rec.Body.String())
	}
}

func TestWeatherHandler_Success(t *testing.T) {
	h := &handler.WeatherHandler{
		CEPClient:     fakeCEPClient{location: cep.Location{City: "São Paulo"}},
		WeatherClient: fakeWeatherClient{tempC: 28.5},
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, newRequest("01001000"))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body struct {
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.TempC != 28.5 {
		t.Fatalf("expected temp_C 28.5, got %v", body.TempC)
	}
	if !almostEqual(body.TempF, 83.3) {
		t.Fatalf("expected temp_F 83.3, got %v", body.TempF)
	}
	if !almostEqual(body.TempK, 301.5) {
		t.Fatalf("expected temp_K 301.5, got %v", body.TempK)
	}
}

func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

func TestWeatherHandler_UnexpectedCEPError(t *testing.T) {
	h := &handler.WeatherHandler{
		CEPClient:     fakeCEPClient{err: errors.New("boom")},
		WeatherClient: fakeWeatherClient{},
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, newRequest("01001000"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}
