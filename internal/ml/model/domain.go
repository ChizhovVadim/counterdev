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
