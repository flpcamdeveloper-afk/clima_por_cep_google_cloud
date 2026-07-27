// Package cep resolve um CEP brasileiro para a cidade correspondente,
// usando a API pública do ViaCEP.
package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrNotFound indica que o CEP tem formato válido mas não existe na base do
// provedor de localização.
var ErrNotFound = errors.New("cep not found")

// Location é a localização resolvida a partir de um CEP.
type Location struct {
	City string
}

// Client resolve um CEP para uma Location.
type Client interface {
	Find(ctx context.Context, cepCode string) (Location, error)
}

// ViaCEPClient implementa Client usando a API pública do ViaCEP
// (https://viacep.com.br).
type ViaCEPClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

// NewViaCEPClient cria um ViaCEPClient pronto para uso.
func NewViaCEPClient(httpClient *http.Client) *ViaCEPClient {
	return &ViaCEPClient{
		HTTPClient: httpClient,
		BaseURL:    "https://viacep.com.br/ws",
	}
}

type viaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

func (c *ViaCEPClient) Find(ctx context.Context, cepCode string) (Location, error) {
	url := fmt.Sprintf("%s/%s/json/", c.BaseURL, cepCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Location{}, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Location{}, fmt.Errorf("viacep: unexpected status code %d", resp.StatusCode)
	}

	var body viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return Location{}, err
	}

	if body.Erro || body.Localidade == "" {
		return Location{}, ErrNotFound
	}

	return Location{City: body.Localidade}, nil
}
