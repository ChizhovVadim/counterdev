package model

import "math/rand/v2"

type Neuron struct {
	Activation float64
	Error      float64
	Prime      float64
}

type Layer struct {
	activationFn IFunc
	outputs      []Neuron
	weights      Matrix
	biases       Matrix
	wGradients   Gradients
	bGradients   Gradients
}

func NewLayer(
	inputSize int,
	outputSize int,
	activationFn IFunc,
) Layer {
	return Layer{
		outputs:      make([]Neuron, outputSize),
		activationFn: activationFn,
		weights:      NewMatrix(outputSize, inputSize),
		biases:       NewMatrix(outputSize, 1),
		wGradients:   NewGradients(outputSize, inputSize),
		bGradients:   NewGradients(outputSize, 1),
	}
}

// веса общие, нейроны и градиенты раздельные
func (l *Layer) Clone() Layer {
	return Layer{
		activationFn: l.activationFn,
		outputs:      make([]Neuron, len(l.outputs)),
		weights:      l.weights,
		biases:       l.biases,
		wGradients:   NewGradients(l.wGradients.Rows, l.wGradients.Cols),
		bGradients:   NewGradients(l.bGradients.Rows, l.bGradients.Cols),
	}
}

func (layer *Layer) InitWeightsSigmoid(rnd *rand.Rand) {
	var outputSize = layer.weights.Rows
	var inputSize = layer.weights.Cols
	var variance = 2.0 / float64(inputSize+outputSize)
	initUniform(rnd, layer.weights.Data, variance)
}

func (layer *Layer) InitWeightsReLU(rnd *rand.Rand) {
	var inputSize = layer.weights.Cols
	var variance = 2.0 / float64(inputSize)
	initUniform(rnd, layer.weights.Data, variance)
}

func (layer *Layer) ForwardFromInput(input Input) {
	for outputIndex := range layer.outputs {
		var x = layer.biases.Data[outputIndex]
		for _, input := range input.Features {
			var inputIndex = int(input)
			const inputValue = 1.0
			x += layer.weights.Get(outputIndex, inputIndex) * inputValue
		}
		var n = &layer.outputs[outputIndex]
		n.Activation = layer.activationFn.Fn(x)
		n.Prime = layer.activationFn.DerivativeFn(x)
	}
}

func (layer *Layer) Forward(input []Neuron) {
	for outputIndex := range layer.outputs {
		var x = layer.biases.Data[outputIndex]
		for inputIndex := range input {
			var inputValue = input[inputIndex].Activation
			x += layer.weights.Get(outputIndex, inputIndex) * inputValue
		}
		var n = &layer.outputs[outputIndex]
		n.Activation = layer.activationFn.Fn(x)
		n.Prime = layer.activationFn.DerivativeFn(x)
	}
}

func (layer *Layer) Backward(input1 []Neuron) {
	for inputIndex := range input1 {
		input1[inputIndex].Error = 0
	}
	for outputIndex := range layer.outputs {
		var n = &layer.outputs[outputIndex]
		var x = n.Error * n.Prime

		for inputIndex := range input1 {
			input1[inputIndex].Error += layer.weights.Get(outputIndex, inputIndex) * x
		}

		layer.bGradients.Add(outputIndex, 0, x*1)
		for inputIndex := range input1 {
			var inputValue = input1[inputIndex].Activation
			layer.wGradients.Add(outputIndex, inputIndex, x*inputValue)
		}
	}
}

func (layer *Layer) BackwardToInput(input2 Input) {
	for outputIndex := range layer.outputs {
		var n = &layer.outputs[outputIndex]
		var x = n.Error * n.Prime
		layer.bGradients.Add(outputIndex, 0, x*1)

		for _, input := range input2.Features {
			var inputIndex = int(input)
			const inputValue = 1.0
			layer.wGradients.Add(outputIndex, inputIndex, x*inputValue)
		}
	}
}

func (l *Layer) addGradients(src *Layer) {
	src.wGradients.AddTo(&l.wGradients)
	src.bGradients.AddTo(&l.bGradients)
}

func (layer *Layer) ApplyGradients() {
	layer.wGradients.Apply(&layer.weights)
	layer.bGradients.Apply(&layer.biases)
}
