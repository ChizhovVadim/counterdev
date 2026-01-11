package enginemini

import (
	"math/rand/v2"
	"sort"

	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type line struct {
	Score int
	Moves []common.Move
}

// Дополнительную случайность обеспечивает TT (если ее не чистить) и таблица истории (если ее не чистить).
func (t *Engine) searchRootRandomness(depth int) {
	const pvNode = true
	const height = 0
	t.clearPV(height)

	var mo = newMoveOrderContext(t, height, t.searchResult.MainLine[0])
	var ml = mo.PrepareMoves()
	var movesSearched int
	var best = -valueInfinity
	var bestLines []line
	var child = &t.stack[height+1].position

	for i := range ml {
		var move = ml[i].Move
		if !t.MakeMove(move, height) {
			continue
		}
		movesSearched += 1

		var isNoisy = isCaptureOrPromotion(move)
		var reduction int
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

		var score, alpha int
		if movesSearched == 1 {
			alpha = -valueInfinity
			score = -t.alphaBeta(-valueInfinity, -alpha, depth-1, height+1)
		} else {
			alpha = max(-valueInfinity, best-t.randomness)
			score = -t.alphaBeta(-(alpha + 1), -alpha, depth-1-reduction, height+1)
			if reduction > 0 && score > alpha {
				score = -t.alphaBeta(-(alpha + 1), -alpha, depth-1, height+1)
			}
			if score > alpha {
				score = -t.alphaBeta(-valueInfinity, -alpha, depth-1, height+1)
			}
		}
		t.UnmakeMove()
		if score > alpha {
			best = max(best, score)
			t.assignPV(height, move)
			bestLines = append(bestLines, line{
				Score: score,
				Moves: t.stack[height].pv.toSlice(),
			})
		}
	}

	var bestLine = selectLRandomLine(bestLines, best, t.randomness)

	t.searchResult = common.SearchInfo{
		Depth:    depth,
		Nodes:    t.nodes,
		Score:    newUciScore(bestLine.Score),
		MainLine: bestLine.Moves,
	}
}

func selectLRandomLine(lines []line, best, randomness int) line {
	lines = filterLines(lines, best-randomness)
	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i].Score > lines[j].Score
	})
	var n = selectN(len(lines), myKernel)
	return lines[n]
}

func filterLines(lines []line, lowerScore int) []line {
	var res []line
	for _, line := range lines {
		if !(line.Score >= lowerScore) {
			continue
		}
		res = append(res, line)
	}
	return res
}

// не равномерный выбор, чтобы лучшие ходы имели больший вес
func selectN(
	n int,
	kernel func(float64) float64,
) int {
	var totalWeight float64
	for i := range n {
		totalWeight += kernel(float64(i))
	}
	var r = rand.Float64()
	var weight float64
	for i := range n {
		weight += kernel(float64(i)) / totalWeight
		if r <= weight {
			return i
		}
	}
	return n - 1
}

func myKernel(x float64) float64 {
	return 1.0 / (1.0 + x)
}
