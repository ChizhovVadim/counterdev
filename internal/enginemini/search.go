package enginemini

import (
	"errors"
	"time"

	"github.com/ChizhovVadim/counterdev/pkg/common"
)

var errSearchTimeout = errors.New("search timeout")

func (e *Engine) iterativeDeepening() {
	defer func() {
		if r := recover(); r != nil {
			if r == errSearchTimeout {
				return
			}
			panic(r)
		}
	}()

	const height = 0

	//select random legal move
	var legalMoves = e.stack[height].position.GenerateLegalMoves()
	if len(legalMoves) == 0 {
		return
	}
	e.searchResult = common.SearchInfo{
		Depth:    0,
		Nodes:    0,
		Time:     0,
		Score:    common.UciScore{},
		MainLine: []common.Move{legalMoves[0]},
	}

	for depth := 1; depth <= maxHeight; depth += 1 {
		if e.timeManager.IsDone() {
			break
		}
		e.searchRoot(depth)
		assertGoodSearchResult(&e.searchResult)
		if e.progress != nil && e.searchResult.Nodes >= int64(e.ProgressMinNodes) {
			e.searchResult.Time = time.Since(e.start)
			e.progress(e.searchResult)
		}
		e.timeManager.OnIterationComplete(e.searchResult)
	}
}

func assertGoodSearchResult(sr *common.SearchInfo) {
	if len(sr.MainLine) == 0 {
		panic("empty best lines")
	}
}

func (t *Engine) searchRoot(depth int) {
	if t.randomness != 0 {
		t.searchRootRandomness(depth)
		return
	}
	const height = 0
	var score = t.alphaBeta(-valueInfinity, valueInfinity, depth, height)
	t.searchResult = common.SearchInfo{
		Depth:    depth,
		Nodes:    t.nodes,
		Score:    newUciScore(score),
		MainLine: t.stack[height].pv.toSlice(),
	}
}

// main search method
func (t *Engine) alphaBeta(alpha, beta, depth, height int) int {
	if depth <= 0 {
		return t.quiescence(alpha, beta, height)
	}
	t.clearPV(height)
	var rootNode = height == 0
	var pvNode = beta != alpha+1
	var position = &t.stack[height].position
	var isCheck = position.IsCheck()

	if !rootNode {
		if height >= maxHeight {
			return t.evaluator.EvaluateQuick(position)
		}
		if t.isRepeat(height) {
			return valueDraw
		}
		if isDraw(position) {
			return valueDraw
		}
		// mate distance pruning
		alpha = max(alpha, lossIn(height))
		beta = min(beta, winIn(height+1))
		if alpha >= beta {
			return alpha
		}
	}

	// transposition table
	var ttDepth, ttValue, ttBound, ttMove, ttHit = t.transTable.Read(position.Key)
	if ttHit {
		ttValue = valueFromTT(ttValue, height)
		if ttDepth >= depth && !pvNode {
			if ttValue >= beta && (ttBound&boundLower) != 0 {
				return ttValue
			}
			if ttValue <= alpha && (ttBound&boundUpper) != 0 {
				return ttValue
			}
		}
	}

	var child = &t.stack[height+1].position
	var staticEval = t.evaluator.EvaluateQuick(position)

	if !rootNode {
		// reverse futility pruning
		if !pvNode && !isCheck && depth <= 6 {
			var score = staticEval - 100*depth
			if score >= beta {
				return score
			}
		}
	}

	var mo = newMoveOrderContext(t, height, ttMove)
	var ml = mo.PrepareMoves()

	var (
		movesSearched int
		hasLegalMove  bool
		bestMove      common.Move
		oldAlpha      = alpha
	)

	for i := range ml {
		var move = ml[i].Move
		var isNoisy = isCaptureOrPromotion(move)

		if alpha > valueLoss && hasLegalMove && !rootNode {
			// futility pruning
			if !pvNode && !isCheck &&
				!isNoisy && depth <= 6 && staticEval+100*depth <= alpha {
				continue
			}
		}

		if !t.MakeMove(move, height) {
			continue
		}
		hasLegalMove = true
		movesSearched += 1

		var reduction, extension int
		if depth >= 3 && movesSearched > 1 && !isNoisy {
			reduction = t.lmr(depth, movesSearched)
			if pvNode {
				reduction -= 1
			}
			if child.IsCheck() {
				reduction -= 1
			}
			reduction = max(0, min(depth-2, reduction))
		}

		if !isNoisy {
			mo.AddQuiet(move)
		}
		var newDepth = depth - 1 + extension

		var score = alpha + 1
		// LMR
		if reduction > 0 {
			score = -t.alphaBeta(-(alpha + 1), -alpha, newDepth-reduction, height+1)
		}
		// PVS
		if score > alpha && beta != alpha+1 && movesSearched > 1 /*&& depth-1 > 0*/ {
			score = -t.alphaBeta(-(alpha + 1), -alpha, newDepth, height+1)
		}
		// full search
		if score > alpha {
			score = -t.alphaBeta(-beta, -alpha, newDepth, height+1)
		}

		t.UnmakeMove()

		if score > alpha {
			alpha = score
			bestMove = move
			t.assignPV(height, move)
			if alpha >= beta {
				break
			}
		}
	}

	if !hasLegalMove {
		if !isCheck /*&& skipMove == 0*/ {
			return valueDraw
		}
		return lossIn(height)
	}

	if alpha > oldAlpha {
		mo.UpdateStatistics(bestMove, depth)
	}

	{
		ttBound = 0
		if alpha > oldAlpha {
			ttBound |= boundLower
		}
		if alpha < beta {
			ttBound |= boundUpper
		}
		t.transTable.Update(position.Key, depth, valueToTT(alpha, height), ttBound, bestMove)
	}

	return alpha
}

func (t *Engine) quiescence(alpha, beta, height int) int {
	t.clearPV(height)
	var position = &t.stack[height].position
	if height >= maxHeight {
		return t.evaluator.EvaluateQuick(position)
	}
	if isDraw(position) {
		return valueDraw
	}
	if t.isRepeat(height) {
		return valueDraw
	}
	var isCheck = position.IsCheck()
	alpha = max(alpha, lossIn(height))
	beta = min(beta, winIn(height+1))
	if alpha >= beta {
		return alpha
	}
	if !isCheck {
		var staticEval = t.evaluator.EvaluateQuick(position)
		if staticEval > alpha {
			alpha = staticEval
			if alpha >= beta {
				return alpha
			}
		}
	}
	var mo = newMoveOrderContext(t, height, common.MoveEmpty)
	var ml = mo.PrepareNoisyMoves()
	for i := range ml {
		var move = ml[i].Move
		if alpha > valueLoss && !isCheck && !seeGEZero(position, move) {
			continue
		}
		if !t.MakeMove(move, height) {
			continue
		}
		var score = -t.quiescence(-beta, -alpha, height+1)
		t.UnmakeMove()
		if score > alpha {
			alpha = score
			if alpha >= beta {
				break
			}
		}
	}
	return alpha
}

func (t *Engine) MakeMove(move common.Move, height int) bool {
	var pos = &t.stack[height].position
	var child = &t.stack[height+1].position
	if move == common.MoveEmpty {
		pos.MakeNullMove(child)
	} else {
		if !pos.MakeMove(move, child) {
			return false
		}
	}
	t.evaluator.MakeMove(pos, move)
	t.incNodes()
	return true
}

func (t *Engine) UnmakeMove() {
	t.evaluator.UnmakeMove()
}

func (t *Engine) incNodes() {
	t.nodes += 1
	if t.nodes&255 == 0 {
		t.timeManager.OnNodesChanged(int(t.nodes))
		if t.timeManager.IsDone() {
			panic(errSearchTimeout)
		}
	}
}

func (e *Engine) lmr(d, m int) int {
	return int(e.reductions[min(d, 63)][min(m, 63)])
}

func isDraw(p *common.Position) bool {
	if p.Rule50 > 100 {
		return true
	}

	if (p.Pawns|p.Rooks|p.Queens) == 0 &&
		!common.MoreThanOne(p.Knights|p.Bishops) {
		return true
	}

	return false
}

func (t *Engine) isRepeat(height int) bool {
	var p = &t.stack[height].position

	if p.Rule50 == 0 || p.LastMove == common.MoveEmpty {
		return false
	}
	for i := height - 1; i >= 0; i-- {
		var temp = &t.stack[i].position
		if temp.Key == p.Key {
			return true
		}
		if temp.Rule50 == 0 || temp.LastMove == common.MoveEmpty {
			return false
		}
	}

	return t.historyKeys[p.Key] >= 2
}
