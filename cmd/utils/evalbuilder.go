package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ChizhovVadim/counterdev/internal/evalmaterial"
	"github.com/ChizhovVadim/counterdev/internal/evalpesto"
	"github.com/ChizhovVadim/counterdev/internal/evalweiss"
	"github.com/ChizhovVadim/counterdev/pkg/common"
	"github.com/ChizhovVadim/counterdev/pkg/evalnn"
)

var networkWeightsCounter55 = weightsLoader("~/chess/n-30-5268.nn", true)
var networkWeightsExp2 = weightsLoader("~/chess/net/2026-01-09_03_55/n-10-1052.nn", false)

func buildEvaluator(key string) common.IEvaluator {
	if key == "material" {
		return evalmaterial.NewEvaluationService()
	}
	if key == "pesto" {
		return evalpesto.NewEvaluationService()
	}
	if key == "" || key == "weiss" {
		return evalweiss.NewEvaluationService()
	}
	if key == "nnue" {
		return evalnn.NewEvaluationService(networkWeightsCounter55(), 1)
	}
	if key == "nnue2" {
		return evalnn.NewEvaluationService(networkWeightsExp2(), 146)
	}
	panic(fmt.Errorf("bad eval key %v", key))
}

func weightsLoader(path string, oldFormat bool) func() *evalnn.Weights {
	var (
		once    sync.Once
		weights *evalnn.Weights
	)
	return func() *evalnn.Weights {
		once.Do(func() {
			path = mapPath(path)
			var f, err = os.Open(path)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			weights, err = evalnn.LoadWeights(f, oldFormat)
			if err != nil {
				panic(err)
			}
			log.Println("Loaded nnue weights", "path", path)
		})
		return weights
	}
}
