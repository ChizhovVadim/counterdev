package model

type Input struct {
	Features []int16
}

type Sample struct {
	Input  Input
	Target float32
}

type IModel interface {
	Train(samples []Sample)
	CalculateCost(samples []Sample) float64
	SaveWeights(path string) error
}
