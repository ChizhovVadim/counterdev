package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/ChizhovVadim/counterdev/internal/ml/dataset"
	"github.com/ChizhovVadim/counterdev/internal/ml/model"
)

// Обучение нейросети
func trainHandler(args []string) error {
	var datasetFiles = findFiles(mapPath("~/chess/dataset/arena-2026-01-08_12_28.pgn"))
	log.Println("LoadDataset",
		"fileCount", len(datasetFiles))
	var data = dataset.LoadDataset(datasetFiles, &dataset.GameMarker{
		SigmoidScale: 3.5 / 512,
		SearchRatio:  1.0,
	})
	var samples, err = dataset.LoadSamples(data, dataset.NewFeature768Provider(), true, 10_000_000)
	if err != nil {
		return err
	}
	log.Println("Dataset loaded", "size", len(samples))

	var (
		epochs        = 10
		netFolderPath = filepath.Join(mapPath("~/chess/net"), time.Now().Format("2006-01-02_15_04"))
	)
	return model.Train(samples, epochs, netFolderPath)
}
