package ml

import "github.com/ChizhovVadim/counterdev/pkg/common"

type DatasetItem struct {
	Position common.Position
	Target   float64
}

type Input struct {
	Features []int16
}

type Sample struct {
	Input  Input
	Target float32
}
