package model

import (
	"math"
	"math/rand/v2"
)

func initUniform(rnd *rand.Rand, data []float64, variance float64) {
	var uniformVariance = 1.0 / 12
	var scale = math.Sqrt(variance / uniformVariance)
	for i := range data {
		data[i] = (rnd.Float64() - 0.5) * scale
	}
}
