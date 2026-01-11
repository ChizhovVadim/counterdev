package enginemini

import "github.com/ChizhovVadim/counterdev/pkg/common"

type MoveOrderContext struct {
	engine         *Engine
	height         int
	transMove      common.Move
	sideToMove     bool
	quietsSearched []common.Move
}

func newMoveOrderContext(engine *Engine, height int, transMove common.Move) MoveOrderContext {
	var stack = &engine.stack[height]
	return MoveOrderContext{
		engine:         engine,
		height:         height,
		transMove:      transMove,
		sideToMove:     stack.position.WhiteMove,
		quietsSearched: stack.quietsSearched[:0],
	}
}

const sortTableKeyImportant = 100_000

func (mo *MoveOrderContext) PrepareMoves() []common.OrderedMove {
	var stack = &mo.engine.stack[mo.height]
	var pos = &stack.position
	var ml = pos.GenerateMoves(stack.moveList[:])
	for i := range ml {
		var m = ml[i].Move
		var score int
		if m == mo.transMove {
			score = sortTableKeyImportant + 2_000
		} else if isCaptureOrPromotion(m) {
			if seeGEZero(pos, m) {
				score = sortTableKeyImportant + 1_000 + mvvlva(m)
			} else {
				score = -100_000 + mvvlva(m)
			}
		} else {
			score = int(mo.engine.history.mainHistory[sideFromToIndex(mo.sideToMove, m)])
		}
		ml[i].Key = int32(score)
	}
	sortMoves(ml)
	return ml
}

func (mo *MoveOrderContext) PrepareNoisyMoves() []common.OrderedMove {
	var stack = &mo.engine.stack[mo.height]
	var pos = &stack.position
	if pos.IsCheck() {
		return mo.PrepareMoves()
	}
	var ml = pos.GenerateCaptures(stack.moveList[:])
	for i := range ml {
		var m = ml[i].Move
		var score = mvvlva(m)
		ml[i].Key = int32(score)
	}
	sortMoves(ml)
	return ml
}

func (mo *MoveOrderContext) AddQuiet(mv common.Move) {
	mo.quietsSearched = append(mo.quietsSearched, mv)
}

func (mo *MoveOrderContext) UpdateStatistics(bestMove common.Move, depth int) {
	if isCaptureOrPromotion(bestMove) {
		return
	}

	var bonus = min(depth*depth, 400)
	var history = &mo.engine.history
	var sideToMove = mo.sideToMove

	for _, m := range mo.quietsSearched {
		var good = m == bestMove

		var fromToIndex = sideFromToIndex(sideToMove, m)
		updateHistory(&history.mainHistory[fromToIndex], bonus, good)

		if good {
			break
		}
	}
}

var sortPieceValues = [...]int{common.Empty: 0, common.Pawn: 1, common.Knight: 2, common.Bishop: 3, common.Rook: 4, common.Queen: 5, common.King: 6}

func mvvlva(move common.Move) int {
	return 8*(sortPieceValues[move.CapturedPiece()]+
		sortPieceValues[move.Promotion()]) -
		sortPieceValues[move.MovingPiece()]
}

func sortMoves(moves []common.OrderedMove) {
	for i := 1; i < len(moves); i++ {
		j, t := i, moves[i]
		for ; j > 0 && moves[j-1].Key < t.Key; j-- {
			moves[j] = moves[j-1]
		}
		moves[j] = t
	}
}
