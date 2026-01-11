package enginemini

import (
	"context"
	"time"

	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type timeManager struct {
	start  time.Time
	limits common.LimitsType
	done   <-chan struct{}
	cancel context.CancelFunc
}

func newTimeManager(ctx context.Context, start time.Time,
	limits common.LimitsType, _ *common.Position) timeManager {

	var tm = timeManager{
		start:  start,
		limits: limits,
	}

	var cancel context.CancelFunc
	if limits.MoveTime > 0 {
		var maximum = time.Duration(limits.MoveTime) * time.Millisecond
		ctx, cancel = context.WithDeadline(ctx, start.Add(maximum))
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	tm.done = ctx.Done()
	tm.cancel = cancel
	return tm
}

func (tm *timeManager) IsDone() bool {
	select {
	case <-tm.done:
		return true
	default:
		return false
	}
}

func (tm *timeManager) OnNodesChanged(nodes int) {
	if tm.limits.Nodes > 0 && nodes >= tm.limits.Nodes {
		tm.cancel()
	}
}

func (tm *timeManager) OnIterationComplete(line common.SearchInfo) {
	if tm.limits.Infinite {
		return
	}
	if tm.limits.Depth != 0 && line.Depth >= tm.limits.Depth {
		tm.cancel()
		return
	}
	/*if line.score >= winIn(line.depth-5) ||
		line.score <= lossIn(line.depth-5) {
		tm.cancel()
		return
	}*/
}

func (tm *timeManager) Close() {
	tm.cancel()
}
