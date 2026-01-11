package enginemini

import (
	"math"

	. "github.com/ChizhovVadim/counterdev/pkg/common"
)

const (
	stackSize     = 128
	maxHeight     = stackSize - 1
	valueDraw     = 0
	valueMate     = 30000
	valueInfinity = valueMate + 1
	valueWin      = valueMate - 2*maxHeight
	valueLoss     = -valueWin
)

func winIn(height int) int {
	return valueMate - height
}

func lossIn(height int) int {
	return -valueMate + height
}

func valueToTT(v, height int) int {
	if v >= valueWin {
		return v + height
	}

	if v <= valueLoss {
		return v - height
	}

	return v
}

func valueFromTT(v, height int) int {
	if v >= valueWin {
		return v - height
	}

	if v <= valueLoss {
		return v + height
	}

	return v
}

func newUciScore(v int) UciScore {
	if v >= valueWin {
		return UciScore{Mate: (valueMate - v + 1) / 2}
	} else if v <= valueLoss {
		return UciScore{Mate: (-valueMate - v) / 2}
	} else {
		return UciScore{Centipawns: v}
	}
}

func isCaptureOrPromotion(move Move) bool {
	return move.CapturedPiece() != Empty ||
		move.Promotion() != Empty
}

func initLmr(reductions *[64][64]int8,
	f func(d, m float64) float64) {
	for d := 1; d < 64; d++ {
		for m := 1; m < 64; m++ {
			var r = f(float64(d), float64(m))
			reductions[d][m] = int8(r)
		}
	}
}

func lmrMain(d, m float64) float64 {
	return 0.75 + math.Log(d)*math.Log(m)/2.25
}
