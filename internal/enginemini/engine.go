/*
Упрощенная версия движка.
Можно использовать как tutorial, для генерации дебютов, для генерации датасета для обучения NN eval.
  - код однопоточный
  - Time manager simple (fixed time, fixed nodes)
  - TT: replace always
  - search: fail hard PVS, TT, rev futility, LMR.
*/
package enginemini

import (
	"context"
	"log"
	"time"

	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type Engine struct {
	ExperimentSettings bool
	ProgressMinNodes   int
	timeManager        timeManager
	transTable         transTable
	historyKeys        map[uint64]int
	randomness         int
	progress           func(common.SearchInfo)
	searchResult       common.SearchInfo
	start              time.Time
	evaluator          common.IEvaluator
	nodes              int64
	stack              [stackSize]struct {
		position       common.Position
		moveList       [common.MaxMoves]common.OrderedMove
		quietsSearched [common.MaxMoves]common.Move
		pv             pv
	}
	reductions [64][64]int8
	history    HistoryTable
}

func NewEngine(evaluator common.IEvaluator, hash int) *Engine {
	var e = &Engine{
		transTable:       newTransTable(hash),
		evaluator:        evaluator,
		ProgressMinNodes: 100_000,
	}
	initLmr(&e.reductions, lmrMain)
	return e
}

type pv struct {
	items [stackSize]common.Move
	size  int
}

func (e *Engine) Prepare() {}

func (e *Engine) Clear() {
	e.transTable.Clear()
	e.history.Clear()
}

func (e *Engine) Search(ctx context.Context, searchParams common.SearchParams) common.SearchInfo {
	e.start = time.Now()
	e.Prepare()
	var p = &searchParams.Position
	e.timeManager = newTimeManager(ctx, e.start, searchParams.Limits, p)
	defer e.timeManager.Close()
	e.transTable.IncDate()
	e.historyKeys = searchParams.Repeats
	e.nodes = 0
	e.stack[0].position = *p
	e.evaluator.Init(&e.stack[0].position)
	e.randomness = searchParams.Randomness
	e.progress = searchParams.Progress
	e.iterativeDeepening()
	e.searchResult.Nodes = e.nodes
	e.searchResult.Time = time.Since(e.start)
	if e.searchResult.Depth < 1 {
		log.Println("not finished search with depth=1")
	}
	return e.searchResult
}

func (t *Engine) clearPV(height int) {
	t.stack[height].pv.size = 0
}

func (t *Engine) assignPV(height int, m common.Move) {
	var pv = &t.stack[height].pv
	var child = &t.stack[height+1].pv
	pv.size = 1
	pv.items[0] = m
	if child.size > 0 {
		pv.size += child.size
		copy(pv.items[1:], child.items[:child.size])
	}
}

func (pv *pv) toSlice() []common.Move {
	var result = make([]common.Move, pv.size)
	copy(result, pv.items[:pv.size])
	return result
}
