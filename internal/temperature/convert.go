// Package temperature converte temperaturas em Celsius para Fahrenheit e
// Kelvin, seguindo as fórmulas do enunciado do desafio.
package temperature

// CToF converte Celsius para Fahrenheit: F = C × 1.8 + 32.
func CToF(c float64) float64 {
	return c*1.8 + 32
}

// CToK converte Celsius para Kelvin: K = C + 273.
func CToK(c float64) float64 {
	return c + 273
}
