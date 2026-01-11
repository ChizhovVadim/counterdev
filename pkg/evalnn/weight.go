package evalnn

import (
	"encoding/binary"
	"io"
	"math"
)

const (
	InputSize  = 64 * 12
	HiddenSize = 512
)

type Weights struct {
	HiddenWeights [InputSize * HiddenSize]float32
	HiddenBiases  [HiddenSize]float32
	OutputWeights [HiddenSize]float32
	OutputBias    float32
}

func LoadWeights(f io.Reader, oldFormat bool) (*Weights, error) {
	var w = &Weights{}

	if oldFormat {
		var buf = make([]byte, 24)
		var _, err = io.ReadFull(f, buf)
		if err != nil {
			return nil, err
		}
	}

	if err := readSlice(f, w.HiddenWeights[:]); err != nil {
		return nil, err
	}
	if err := readSlice(f, w.HiddenBiases[:]); err != nil {
		return nil, err
	}
	if err := readSlice(f, w.OutputWeights[:]); err != nil {
		return nil, err
	}

	var buf = make([]byte, 4)
	var _, err = io.ReadFull(f, buf)
	if err != nil {
		return nil, err
	}
	w.OutputBias = math.Float32frombits(binary.LittleEndian.Uint32(buf))

	return w, nil
}

func readSlice(f io.Reader, data []float32) error {
	var buf [4]byte
	for i := range data {
		_, err := io.ReadFull(f, buf[:])
		if err != nil {
			return err
		}
		var val = math.Float32frombits(binary.LittleEndian.Uint32(buf[:]))
		data[i] = val
	}
	return nil
}
