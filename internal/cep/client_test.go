package cep_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"climacep/internal/cep"
)

func TestViaCEPClient_Find_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"cep":"01001-000","localidade":"São Paulo","uf":"SP","erro":false}`))
	}))
	defer server.Close()

	client := &cep.ViaCEPClient{HTTPClient: server.Client(), BaseURL: server.URL}

	loc, err := client.Find(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.City != "São Paulo" {
		t.Fatalf("expected city %q, got %q", "São Paulo", loc.City)
	}
}

func TestViaCEPClient_Find_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"erro":true}`))
	}))
	defer server.Close()

	client := &cep.ViaCEPClient{HTTPClient: server.Client(), BaseURL: server.URL}

	_, err := client.Find(context.Background(), "00000000")
	if err != cep.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
