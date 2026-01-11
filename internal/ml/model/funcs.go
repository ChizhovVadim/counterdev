package model

import "math"

type IFunc interface {
	Fn(x float64) float64
	DerivativeFn(x float64) float64
}

type SquareFn struct{}

func (*SquareFn) Fn(x float64) float64 {
	return x * x
}

func (*SquareFn) DerivativeFn(x float64) float64 {
	return 2 * x
}

type ReLuFn struct{}

func (*ReLuFn) Fn(x float64) float64 {
	if x > 0 {
		return x
	}
	return 0
}

func (*ReLuFn) DerivativeFn(x float64) float64 {
	if x > 0 {
		return 1
	}
	return 0
}

type SigmoidFn struct{}

func (s *SigmoidFn) Fn(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func (s *SigmoidFn) DerivativeFn(x float64) float64 {
	var y = s.Fn(x)
	return y * (1 - y)
}
