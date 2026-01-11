package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/ChizhovVadim/counterdev/internal/arena"
	"github.com/ChizhovVadim/counterdev/internal/enginemini"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

// играет матч между двумя движками
func matchHandler(args []string) error {
	var date = time.Now()
	var (
		openingsPath   = mapPath("~/chess/openings.txt")
		outputGamePath = mapPath(fmt.Sprintf("~/chess/games/match-%v.pgn", date.Format("2006-01-02_15_04")))
		concurrency    = runtime.NumCPU()
		timeLimit      = common.LimitsType{Nodes: 1_000_000}
	)
	return arena.PlayMatch(context.Background(), concurrency, openingsPath, outputGamePath,
		func() arena.Player {
			return arena.Player{
				Name:      "Main",
				TimeLimit: timeLimit,
				Engine:    enginemini.NewEngine(buildEvaluator(""), 128),
			}
		},
		func() arena.Player {
			return arena.Player{
				Name:      "Exp",
				TimeLimit: timeLimit,
				Engine:    enginemini.NewEngine(buildEvaluator(""), 128),
			}
		})
}
