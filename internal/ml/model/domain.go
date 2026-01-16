package model

type FeatureInfo struct {
	Index int16
	Value int16
}

type Input struct {
	Features []FeatureInfo
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
