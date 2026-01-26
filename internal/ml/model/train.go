package model

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ChizhovVadim/counterdev/internal/ml"
)

func Train(
	samples []ml.Sample,
	epochs int,
	netFolderPath string,
) error {

	var validationSize = min(len(samples)/20, 500_000)
	var validation = samples[:validationSize]
	var training = samples[validationSize:]

	var concurrency = runtime.GOMAXPROCS(0)
	var model = NewModelMt(NewModel().InitWeights(), concurrency)

	return trainCycle(model, validation, training, epochs, netFolderPath)
}

func trainCycle(
	model IModel,
	validation, training []ml.Sample,
	epochs int,
	netFolderPath string,
) error {
	log.Println("Train started")
	defer log.Println("Train finished")

	err := os.MkdirAll(netFolderPath, os.ModePerm)
	if err != nil {
		return err
	}

	const BatchSize = 16384

	for epoch := 1; epoch <= epochs; epoch++ {
		shuffle(training)
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

func shuffle(training []ml.Sample) {
	rand.Shuffle(len(training), func(i, j int) {
		training[i], training[j] = training[j], training[i]
	})
}

func saveModel(
	model IModel,
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
