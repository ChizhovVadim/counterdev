package main

import (
	"log"
	"runtime"

	"github.com/ChizhovVadim/counterdev/internal/ml/dataset"
	"github.com/ChizhovVadim/counterdev/internal/ml/model"
)

// Обучение нейросети
func trainHandler(args []string) error {
	var netFolderPath = mapPath("~/chess/net")
	var epochs int = 10
	var sigmoidScale = 3.5 / 512
	var featureService = dataset.NewFeature768Provider()
	var data = dataset.LoadDataset(mapPath("~/chess/dataset/arena-2026-01-08_12_28.pgn"), sigmoidScale, 1.0)
	var samples, err = dataset.LoadSamples(data, featureService, true, 10_000_000)
	if err != nil {
		return err
	}
	log.Println("Dataset loaded", "size", len(samples))
	return model.Train(samples, epochs, runtime.NumCPU(), netFolderPath)
}
