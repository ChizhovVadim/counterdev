package model

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"
)

func Train(
	samples []Sample,
	epochs int,
	concurrency int,
	netFolderPath string,
) error {
	log.Println("Train started")
	defer log.Println("Train finished")

	netFolderPath, err := createOutputFolder(netFolderPath)
	if err != nil {
		return err
	}

	var validationSize = len(samples) / 20
	var validation = samples[:validationSize]
	var training = samples[validationSize:]

	var rnd = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	var model = NewModel()
	model.InitWeights(rnd)
	const BatchSize = 16384

	for epoch := 1; epoch <= epochs; epoch++ {
		shuffle(rnd, training)
		for i := 0; i+BatchSize <= len(training); i += BatchSize {
			var batch = training[i : i+BatchSize]
			model.Train(batch)
		}
		log.Printf("Finished Epoch %v\n", epoch)
		var validationCost = model.CalculateCost(validation)
		log.Printf("Current validation cost is: %f\n", validationCost)
		var err = saveModel(model, netFolderPath, epoch, validationCost)
		if err != nil {
			return err
		}
	}

	return nil
}

func shuffle(rnd *rand.Rand, training []Sample) {
	rnd.Shuffle(len(training), func(i, j int) {
		training[i], training[j] = training[j], training[i]
	})
}

func saveModel(
	model *Model,
	netFolderPath string,
	epoch int,
	validationCost float64,
) error {
	if epoch < 5 {
		return nil
	}
	var valCostInt = int(100_000 * validationCost)
	var filename = filepath.Join(netFolderPath, fmt.Sprintf("n-%2d-%v.nn", epoch, valCostInt))
	return model.SaveWeights(filename)
}

func createOutputFolder(baseFolder string) (string, error) {
	var res = filepath.Join(baseFolder, time.Now().Format("2006-01-02_15_04"))
	err := os.MkdirAll(res, os.ModePerm)
	if err != nil {
		return "", err
	}
	return res, nil
}
