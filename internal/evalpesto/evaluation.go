package evalpesto

import "github.com/ChizhovVadim/counterdev/pkg/common"

const MaxHeight = 128

type EvaluationService struct {
	weights            []Score
	accumulators       [MaxHeight]Score
	currentAccumulator int
	updates            Updates
}

func NewEvaluationService() *EvaluationService {
	return &EvaluationService{weights: loadWeights()}
}

func (e *EvaluationService) Init(p *common.Position) {
	var s Score

	for b := p.AllPieces(); b != 0; b &= b - 1 {
		var sq = common.FirstOne(b)
		piece, side := p.GetPieceTypeAndSide(sq)
		var index = calculateNetInputIndex(side, piece, sq)
		s += e.weights[index]
	}

	e.currentAccumulator = 0
	e.accumulators[e.currentAccumulator] = s
}

func (e *EvaluationService) MakeMove(p *common.Position, m common.Move) {
	calculateUpdates(p, m, &e.updates)
	var s = e.accumulators[e.currentAccumulator]
	for i := range e.updates.Size {
		var index = e.updates.Indices[i]
		var coeff = Score(e.updates.Coeffs[i])
		s += e.weights[index] * coeff
	}
	e.currentAccumulator += 1
	e.accumulators[e.currentAccumulator] = s
}

func (e *EvaluationService) UnmakeMove() {
	e.currentAccumulator -= 1
}

func (e *EvaluationService) EvaluateQuick(p *common.Position) int {
	var s = e.accumulators[e.currentAccumulator]

	//TODO BishopPair?

	const totalPhase = 24

	var phase = min(totalPhase,
		common.PopCount(p.Knights|p.Bishops)+
			2*common.PopCount(p.Rooks)+
			4*common.PopCount(p.Queens))

	var result = (int(s.Middle())*phase + int(s.End())*(totalPhase-phase)) / totalPhase
	if !p.WhiteMove {
		result = -result
	}

	const Tempo = 15
	return result + Tempo
}
