package evalmix

import "github.com/ChizhovVadim/counterdev/pkg/common"

type EvaluationService struct {
	a common.IEvaluator
	b common.IEvaluator
}

func NewEvaluationService(
	a common.IEvaluator,
	b common.IEvaluator,
) *EvaluationService {
	return &EvaluationService{
		a: a,
		b: b,
	}
}

func (e *EvaluationService) Init(p *common.Position) {
	e.a.Init(p)
	e.b.Init(p)
}

func (e *EvaluationService) MakeMove(p *common.Position, m common.Move) {
	e.a.MakeMove(p, m)
	e.b.MakeMove(p, m)
}

func (e *EvaluationService) UnmakeMove() {
	e.a.UnmakeMove()
	e.b.UnmakeMove()
}

func (e *EvaluationService) EvaluateQuick(p *common.Position) int {
	return mix(e.a.EvaluateQuick(p), e.b.EvaluateQuick(p))
}

func mix(x, y int) int {
	return (x + y) / 2
}
