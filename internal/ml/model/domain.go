package model

import "github.com/ChizhovVadim/counterdev/internal/ml"

type IModel interface {
	Train(samples []ml.Sample)
	CalculateCost(samples []ml.Sample) float64
	SaveWeights(path string) error
}
