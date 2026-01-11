package common

import (
	"context"
	"time"
)

type LimitsType struct {
	Ponder         bool
	Infinite       bool
	WhiteTime      int
	BlackTime      int
	WhiteIncrement int
	BlackIncrement int
	MoveTime       int
	MovesToGo      int
	Depth          int
	Nodes          int
	Mate           int
}

type SearchParams struct {
	Position   Position
	Repeats    map[uint64]int //включая текущую позицию
	Limits     LimitsType
	Randomness int
	Progress   func(si SearchInfo)
}

type SearchInfo struct {
	Depth    int
	Nodes    int64
	Time     time.Duration
	Score    UciScore
	MainLine []Move
}

type UciScore struct {
	Centipawns int
	Mate       int
}

type IEvaluator interface {
	Init(p *Position)
	MakeMove(p *Position, m Move)
	UnmakeMove()
	EvaluateQuick(p *Position) int
}

type IEngine interface {
	Clear()
	Search(ctx context.Context, searchParams SearchParams) SearchInfo
}
