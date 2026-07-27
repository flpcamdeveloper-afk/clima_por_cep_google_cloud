package temperature_test

import (
	"testing"

	"climacep/internal/temperature"
)

func TestCToF(t *testing.T) {
	tests := []struct {
		name string
		c    float64
		want float64
	}{
		{"zero", 0, 32},
		{"example from the challenge spec", 28.5, 83.3},
		{"negative", -10, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := temperature.CToF(tt.c); !almostEqual(got, tt.want) {
				t.Fatalf("CToF(%v) = %v, want %v", tt.c, got, tt.want)
			}
		})
	}
}

func TestCToK(t *testing.T) {
	tests := []struct {
		name string
		c    float64
		want float64
	}{
		{"zero", 0, 273},
		{"positive", 28.5, 301.5},
		{"negative", -273, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := temperature.CToK(tt.c); !almostEqual(got, tt.want) {
				t.Fatalf("CToK(%v) = %v, want %v", tt.c, got, tt.want)
			}
		})
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
