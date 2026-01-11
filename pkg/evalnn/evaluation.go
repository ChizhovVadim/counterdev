package evalnn

import (
	"math"

	"github.com/ChizhovVadim/counterdev/pkg/common"
)

const (
	Add    = 1
	Remove = -Add
)

const MaxHeight = 128

type EvaluationService struct {
	*Weights
	scale         float32
	updates       Updates
	hiddenOutputs [MaxHeight][HiddenSize]float32
	currentHidden int
	hiddenRelu    [HiddenSize]float32
}

type Updates struct {
	Indices [8]int16
	Coeffs  [8]int8
	Size    int
}

func (u *Updates) Add(index int16, coeff int8) {
	u.Indices[u.Size] = index
	u.Coeffs[u.Size] = coeff
	u.Size++
}

func NewEvaluationService(weights *Weights, scale float32) *EvaluationService {
	return &EvaluationService{
		Weights: weights,
		scale:   scale,
	}
}

func (e *EvaluationService) Init(p *common.Position) {
	input := make([]int, 0, 32)

	for b := p.AllPieces(); b != 0; b &= b - 1 {
		var sq = common.FirstOne(b)
		piece, side := p.GetPieceTypeAndSide(sq)
		if piece != common.Empty {
			input = append(input, int(calculateNetInputIndex(side, piece, sq)))
		}
	}

	e.currentHidden = 0
	hiddenOutputs := e.hiddenOutputs[e.currentHidden][:]

	for i := range hiddenOutputs {
		hiddenOutputs[i] = e.HiddenBiases[i]
	}

	for _, i := range input {
		for j := range hiddenOutputs {
			hiddenOutputs[j] += e.HiddenWeights[i*HiddenSize+j]
		}
	}
}

func (e *EvaluationService) MakeMove(p *common.Position, m common.Move) {
	e.updates.Size = 0

	// MakeNullMove
	if m == common.MoveEmpty {
		e.updateHidden()
		return
	}

	var from, to, movingPiece, capturedPiece, epCapSq, promotionPt, isCastling = unpackMove(p, m)

	e.updates.Add(calculateNetInputIndex(p.WhiteMove, movingPiece, from), Remove)

	if capturedPiece != common.Empty {
		var capSq = to
		if epCapSq != common.SquareNone {
			capSq = epCapSq
		}
		e.updates.Add(calculateNetInputIndex(!p.WhiteMove, capturedPiece, capSq), Remove)
	}

	var pieceAfterMove = movingPiece
	if promotionPt != common.Empty {
		pieceAfterMove = promotionPt
	}
	e.updates.Add(calculateNetInputIndex(p.WhiteMove, pieceAfterMove, to), Add)

	if isCastling {
		var rookRemoveSq, rookAddSq int
		if p.WhiteMove {
			if to == common.SquareG1 {
				rookRemoveSq = common.SquareH1
				rookAddSq = common.SquareF1
			} else {
				rookRemoveSq = common.SquareA1
				rookAddSq = common.SquareD1
			}
		} else {
			if to == common.SquareG8 {
				rookRemoveSq = common.SquareH8
				rookAddSq = common.SquareF8
			} else {
				rookRemoveSq = common.SquareA8
				rookAddSq = common.SquareD8
			}
		}

		e.updates.Add(calculateNetInputIndex(p.WhiteMove, common.Rook, rookRemoveSq), Remove)
		e.updates.Add(calculateNetInputIndex(p.WhiteMove, common.Rook, rookAddSq), Add)
	}

	e.updateHidden()
}

func (e *EvaluationService) UnmakeMove() {
	e.currentHidden -= 1
}

func (e *EvaluationService) EvaluateQuick(p *common.Position) int {
	var output = int(e.scale * e.quickFeed())
	const MaxEval = 15_000
	output = max(-MaxEval, min(MaxEval, output))

	/*var npMaterial = 4*common.PopCount(p.Knights|p.Bishops) + 6*common.PopCount(p.Rooks) + 12*common.PopCount(p.Queens)
	output = output * (160 + npMaterial) / 160
	output = output * (200 - p.Rule50) / 200*/

	if !p.WhiteMove {
		output = -output
	}
	return output + 15
}

/*func (e *EvaluationService) EvaluateProb(p *common.Position) float64 {
	e.Init(p)
	var output = float64(e.QuickFeed())
	return sigmoid(output)
}*/

func (e *EvaluationService) quickFeed() float32 {
	reluNEON(e.hiddenRelu[:], e.hiddenOutputs[e.currentHidden][:])
	return dotProductNEON(e.hiddenRelu[:], e.OutputWeights[:]) + e.OutputBias

	/*var output float32
	// zip(hiddenOutputs, e.OutputWeights)
	for i, x := range e.hiddenOutputs[e.currentHidden][:] {
		output += max(x, 0) * e.OutputWeights[i]
	}
	return output + e.OutputBias()*/
}

func (e *EvaluationService) updateHidden() {
	e.currentHidden += 1
	hiddenOutputs := e.hiddenOutputs[e.currentHidden][:]
	copy(hiddenOutputs, e.hiddenOutputs[e.currentHidden-1][:])

	for i := 0; i < e.updates.Size; i++ {
		var index = int(e.updates.Indices[i]) * HiddenSize
		// zip(hiddenOutputs, e.HiddenWeights)
		if e.updates.Coeffs[i] == Add {
			/*for j := range hiddenOutputs {
				hiddenOutputs[j] += e.HiddenWeights[index+j]
			}*/
			addNEON(hiddenOutputs[:], hiddenOutputs[:], e.HiddenWeights[index:index+HiddenSize])
		} else {
			/*for j := range hiddenOutputs {
				hiddenOutputs[j] -= e.HiddenWeights[index+j]
			}*/
			subNEON(hiddenOutputs[:], hiddenOutputs[:], e.HiddenWeights[index:index+HiddenSize])
		}
	}
}

func unpackMove(p *common.Position, m common.Move) (from, to, movingPiece, capturedPiece, epCapSq, promotionPt int, isCastling bool) {
	from = m.From()
	to = m.To()
	movingPiece = m.MovingPiece()
	capturedPiece = m.CapturedPiece()
	promotionPt = m.Promotion()
	epCapSq = common.SquareNone
	if movingPiece == common.King {
		if p.WhiteMove {
			if from == common.SquareE1 && (to == common.SquareG1 || to == common.SquareC1) {
				isCastling = true
			}
		} else {
			if from == common.SquareE8 && (to == common.SquareG8 || to == common.SquareC8) {
				isCastling = true
			}
		}
	} else if movingPiece == common.Pawn {
		if to == p.EpSquare {
			if p.WhiteMove {
				epCapSq = to - 8
			} else {
				epCapSq = to + 8
			}
		}
	}
	return
}

func calculateNetInputIndex(whiteSide bool, pieceType, square int) int16 {
	var piece12 = pieceType - common.Pawn
	if !whiteSide {
		piece12 += 6
	}
	return int16(square ^ piece12<<6)
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
