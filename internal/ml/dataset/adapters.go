package dataset

import (
	"context"

	"github.com/ChizhovVadim/counterdev/internal/ml"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type EvalToProbabilityAdapter struct {
	evaluator    common.IEvaluator
	sigmoidScale float64
}

func ProbEvaluatorFromEvaluator(
	evaluator common.IEvaluator,
	sigmoidScale float64,
) *EvalToProbabilityAdapter {
	return &EvalToProbabilityAdapter{
		evaluator:    evaluator,
		sigmoidScale: sigmoidScale,
	}
}

func (a *EvalToProbabilityAdapter) EvaluateProb(p *common.Position) float64 {
	a.evaluator.Init(p)
	var staticEval = a.evaluator.EvaluateQuick(p)
	if !p.WhiteMove {
		staticEval = -staticEval
	}
	return sigmoid(a.sigmoidScale * float64(staticEval))
}

type IModel interface {
	Forward(input ml.Input) float64
}

type ModelToProbabilityAdapter struct {
	featureProvider IFeatureProvider
	model           IModel
}

func ProbEvaluatorFromModel(
	featureProvider IFeatureProvider,
	model IModel,
) *ModelToProbabilityAdapter {
	return &ModelToProbabilityAdapter{
		featureProvider: featureProvider,
		model:           model,
	}
}

func (a *ModelToProbabilityAdapter) EvaluateProb(p *common.Position) float64 {
	var input = a.featureProvider.ComputeFeatures(p)
	var output = a.model.Forward(input)
	return output
}

type EngineToProbabilityAdapter struct {
	engine       common.IEngine
	sigmoidScale float64
}

func ProbEvaluatorFromEngine(
	engine common.IEngine,
	sigmoidScale float64,
) *EngineToProbabilityAdapter {
	return &EngineToProbabilityAdapter{
		engine:       engine,
		sigmoidScale: sigmoidScale,
	}
}

func (a *EngineToProbabilityAdapter) EvaluateProb(p *common.Position) float64 {
	var res = a.engine.Search(context.Background(), common.SearchParams{
		Position: *p,
		Repeats:  map[uint64]int{p.Key: 1},
		Limits:   common.LimitsType{Depth: 8},
	})
	return computeSearchTarget(a.sigmoidScale, res.Score, p.WhiteMove)
}
