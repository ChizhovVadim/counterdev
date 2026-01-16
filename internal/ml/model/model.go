package model

import (
	"encoding/binary"
	"io"
	"math"
	"math/rand/v2"
	"os"
)

type Model struct {
	layer1 Layer
	layer2 Layer
	cost   SquareFn
}

func NewModel() *Model {
	var inputSize = 768
	var hiddenSize = 512
	return &Model{
		layer1: NewLayer(
			inputSize,
			make([]Neuron, hiddenSize),
			&ReLuFn{}),
		layer2: NewLayer(
			hiddenSize,
			make([]Neuron, 1),
			&SigmoidFn{}),
	}
}

func (m *Model) Clone() Model {
	return Model{
		layer1: m.layer1.Clone(),
		layer2: m.layer2.Clone(),
		cost:   m.cost,
	}
}

func (m *Model) InitWeights(rnd *rand.Rand) *Model {
	m.layer1.InitWeightsReLU(rnd) // ненулевых входных признаков не более 32 (кол во фигур на доске)
	m.layer2.InitWeightsSigmoid(rnd)
	return m
}

func (m *Model) LoadWeights(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	/*{
		var header = make([]byte, 24)
		_, err = io.ReadFull(f, header[:])
		if err != nil {
			return err
		}
	}*/

	var data = [...][]float64{
		m.layer1.weights.Data,
		m.layer1.biases.Data,
		m.layer2.weights.Data,
		m.layer2.biases.Data,
	}
	for i := range data {
		var err = readSlice(f, data[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) SaveWeights(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var data = [...][]float64{
		m.layer1.weights.Data,
		m.layer1.biases.Data,
		m.layer2.weights.Data,
		m.layer2.biases.Data,
	}
	for i := range data {
		var err = writeSlice(f, data[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) Train(samples []Sample) {
	for _, sample := range samples {
		m.trainSample(sample)
	}
	m.applyGradients()
}

func (m *Model) trainSample(sample Sample) {
	predicted := m.Forward(sample.Input)
	m.layer2.outputs[0].Error = m.cost.DerivativeFn(predicted - float64(sample.Target))
	// back propagation
	m.layer2.Backward(m.layer1.outputs)
	m.layer1.BackwardToInput(sample.Input)
}

func (m *Model) addGradients(src *Model) {
	m.layer1.addGradients(&src.layer1)
	m.layer2.addGradients(&src.layer2)
}

func (m *Model) applyGradients() {
	m.layer1.ApplyGradients()
	m.layer2.ApplyGradients()
}

func (m *Model) CalculateCost(samples []Sample) float64 {
	var totalCost float64
	for _, sample := range samples {
		var predict = m.Forward(sample.Input)
		var x = predict - float64(sample.Target)
		totalCost += m.cost.Fn(x)
	}
	return totalCost / float64(len(samples))
}

func (m *Model) Forward(input Input) float64 {
	m.layer1.ForwardFromInput(input)
	m.layer2.Forward(m.layer1.outputs)
	return m.layer2.outputs[0].Activation
}

func readSlice(f io.Reader, data []float64) error {
	var buf [4]byte
	for i := range data {
		_, err := io.ReadFull(f, buf[:])
		if err != nil {
			return err
		}
		var val = math.Float32frombits(binary.LittleEndian.Uint32(buf[:]))
		data[i] = float64(val)
	}
	return nil
}

func writeSlice(f io.Writer, data []float64) error {
	var buf [4]byte
	for i := range data {
		var val = float32(data[i])
		binary.LittleEndian.PutUint32(buf[:], math.Float32bits(val))
		_, err := f.Write(buf[:])
		if err != nil {
			return err
		}
	}
	return nil
}
