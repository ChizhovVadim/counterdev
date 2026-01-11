package model

import (
	"sync"
	"sync/atomic"
)

// модель для многопоточного обучения
type ModelMt struct {
	models []Model
}

func (m *ModelMt) Train(smaples []Sample) {
	//TODO
}

func (m *ModelMt) CalculateCost(samples []Sample) float64 {
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
