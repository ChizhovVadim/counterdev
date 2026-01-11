package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ChizhovVadim/counterdev/internal/evalmaterial"
	"github.com/ChizhovVadim/counterdev/internal/evalweiss"
	"github.com/ChizhovVadim/counterdev/pkg/common"
	"github.com/ChizhovVadim/counterdev/pkg/evalnn"
)

func buildEvaluator(key string) common.IEvaluator {
	if key == "weiss" {
		return evalweiss.NewEvaluationService()
	}
	if key == "material" {
		return evalmaterial.NewEvaluationService()
	}
	if key == "" || key == "nnue" {
		var fn = mapPath("~/chess/n-30-5268.nn")
		var eval, err = buildEvalNN2(fn, true, 1)
		if err != nil {
			panic(err)
		}
		return eval
	}
	if key == "nnue2" {
		var fn = mapPath("~/chess/net/2026-01-09_03_55/n-10-1052.nn")
		var eval, err = buildEvalNN2(fn, false, 146)
		if err != nil {
			panic(err)
		}
		return eval
	}
	panic(fmt.Errorf("bad eval key %v", key))
}

func buildEvalNN2(path string, oldFormat bool, scale float32) (*evalnn.EvaluationService, error) {
	// TODO кешировать веса
	var f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	weights, err := evalnn.LoadWeights(f, oldFormat)
	if err != nil {
		return nil, err
	}
	log.Println("Loaded nnue weights", "path", path)
	return evalnn.NewEvaluationService(weights, scale), nil
}
