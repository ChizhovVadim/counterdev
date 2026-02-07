package evalpesto

import "github.com/ChizhovVadim/counterdev/pkg/common"

const (
	Add    = 1
	Remove = -Add
)

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

func calculateUpdates(
	p *common.Position,
	m common.Move,
	updates *Updates,
) {
	updates.Size = 0

	// Null move
	if m == common.MoveEmpty {
		return
	}

	var (
		from          = m.From()
		to            = m.To()
		movingPiece   = m.MovingPiece()
		capturedPiece = m.CapturedPiece()
		promotionPt   = m.Promotion()
		epCapSq       = common.SquareNone
		isCastling    = false
	)
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

	updates.Add(calculateNetInputIndex(p.WhiteMove, movingPiece, from), Remove)

	if capturedPiece != common.Empty {
		var capSq = to
		if epCapSq != common.SquareNone {
			capSq = epCapSq
		}
		updates.Add(calculateNetInputIndex(!p.WhiteMove, capturedPiece, capSq), Remove)
	}

	var pieceAfterMove = movingPiece
	if promotionPt != common.Empty {
		pieceAfterMove = promotionPt
	}
	updates.Add(calculateNetInputIndex(p.WhiteMove, pieceAfterMove, to), Add)

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

		updates.Add(calculateNetInputIndex(p.WhiteMove, common.Rook, rookRemoveSq), Remove)
		updates.Add(calculateNetInputIndex(p.WhiteMove, common.Rook, rookAddSq), Add)
	}
}

func calculateNetInputIndex(whiteSide bool, pieceType, square int) int16 {
	var piece12 = pieceType - common.Pawn
	if !whiteSide {
		piece12 += 6
	}
	return int16(square ^ piece12<<6)
}
