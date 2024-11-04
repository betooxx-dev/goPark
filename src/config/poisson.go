package config

import (
	"math"
	"math/rand"
	"time"
)

type PoissonProcess struct {
	lambda   float64 // tasa media de llegadas por unidad de tiempo
	lastTime time.Time
}

func NewPoissonProcess(lambda float64) *PoissonProcess {
	return &PoissonProcess{
		lambda:   lambda,
		lastTime: time.Now(),
	}
}

func (p *PoissonProcess) NextInterval() time.Duration {
	u := rand.Float64()
	// Próximo intervalo = -ln(U)/λ donde U es uniforme(0,1)
	interval := -math.Log(u) / p.lambda

	// Convertir a duración
	return time.Duration(interval * float64(time.Second))
}

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
