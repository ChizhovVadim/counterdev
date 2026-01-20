package main

import (
	"flag"

	"github.com/ChizhovVadim/counterdev/internal/ml/dataset"
)

// вычисляет ошибку оценочной функции на валидационном датасете
func qualityHandler(args []string) error {
	var (
		sigmoidScale = 3.5 / 512
		datasetPath  = mapPath("~/chess/tuner/quiet-labeled.epd")
		eval         = ""
		searchRatio  = 1.0
	)

	var flagset = flag.NewFlagSet("", flag.ExitOnError)
	flagset.Float64Var(&sigmoidScale, "sigmoidscale", sigmoidScale, "")
	flagset.StringVar(&datasetPath, "path", datasetPath, "")
	flagset.StringVar(&eval, "eval", eval, "")
	flagset.Float64Var(&searchRatio, "searchratio", searchRatio, "")
	flagset.Parse(args)

	var data = dataset.LoadValidationDataset(mapPath("~/chess/tuner/quiet-labeled.epd"))
	//var data = dataset.LoadDataset(mapPath("~/chess/dataset/arena-2026-01-03_16_32.pgn"), sigmoidScale, searchRatio)
	var pureEvaluator = buildEvaluator(eval)
	var evaluator, ok = pureEvaluator.(dataset.IProbEvaluator)
	if !ok {
		evaluator = dataset.ProbEvaluatorFromEvaluator(pureEvaluator, sigmoidScale)
	}

	/*var evaluator dataset.IProbEvaluator
	{
		var m = model.NewModel()
		var fn = mapPath("~/chess/net/2026-01-09_03_55/n-10-1052.nn")
		var err = m.LoadWeights(fn)
		if err != nil {
			return err
		}
		evaluator = dataset.ProbEvaluatorFromModel(dataset.NewFeature768Provider(), m)
	}*/

	return dataset.CheckEvalQuality(evaluator, data)
}
