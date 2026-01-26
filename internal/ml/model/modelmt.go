package model

import (
	"sync"
	"sync/atomic"

	"github.com/ChizhovVadim/counterdev/internal/ml"
)

// модель для многопоточного обучения
type ModelMt struct {
	models []Model
}

func NewModelMt(model *Model, concurrency int) *ModelMt {
	var models = make([]Model, concurrency)
	for i := range models {
		models[i] = model.Clone()
	}
	return &ModelMt{models: models}
}

func (m *ModelMt) Train(samples []ml.Sample) {
	var index int32 = -1
	var wg = &sync.WaitGroup{}
	for modelIndex := range m.models {
		wg.Add(1)
		go func(m *Model) {
			defer wg.Done()
			for {
				var i = int(atomic.AddInt32(&index, 1))
				if i >= len(samples) {
					break
				}
				m.trainSample(samples[i])
			}
		}(&m.models[modelIndex])
	}
	wg.Wait()
	for i := 1; i < len(m.models); i += 1 {
		m.models[0].addGradients(&m.models[i])
	}
	m.models[0].applyGradients()
}

func (m *ModelMt) CalculateCost(samples []ml.Sample) float64 {
	var index int32 = -1
	var wg = &sync.WaitGroup{}
	var totalCost float64
	var mu = &sync.Mutex{}
	for modelIndex := range m.models {
		wg.Add(1)
		go func(m *Model) {
			defer wg.Done()
			var localCost float64
			for {
				var i = int(atomic.AddInt32(&index, 1))
				if i >= len(samples) {
					break
				}
				var sample = samples[i]
				var predict = m.Forward(sample.Input)
				var x = predict - float64(sample.Target)
				localCost += m.cost.Fn(x)
			}
			mu.Lock()
			totalCost += localCost
			mu.Unlock()
		}(&m.models[modelIndex])
	}
	wg.Wait()
	averageCost := totalCost / float64(len(samples))
	return averageCost
}

func (m *ModelMt) SaveWeights(path string) error {
	return m.models[0].SaveWeights(path)
}
