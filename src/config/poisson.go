package config

import (
	"math"
	"math/rand"
	"time"
)

// PoissonProcess genera eventos siguiendo una distribución de Poisson
type PoissonProcess struct {
	lambda   float64 // tasa media de llegadas por unidad de tiempo
	lastTime time.Time
}

// NewPoissonProcess crea un nuevo proceso de Poisson con la tasa lambda dada
func NewPoissonProcess(lambda float64) *PoissonProcess {
	return &PoissonProcess{
		lambda:   lambda,
		lastTime: time.Now(),
	}
}

// NextInterval retorna el siguiente intervalo de tiempo según la distribución de Poisson
func (p *PoissonProcess) NextInterval() time.Duration {
	// La distribución exponencial modela el tiempo entre eventos de Poisson
	u := rand.Float64()
	// Próximo intervalo = -ln(U)/λ donde U es uniforme(0,1)
	interval := -math.Log(u) / p.lambda

	// Convertir a duración
	return time.Duration(interval * float64(time.Second))
}

// GeneratePoisson genera un número según la distribución de Poisson
func GeneratePoisson(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}
	return k - 1
}
