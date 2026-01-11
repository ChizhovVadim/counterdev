package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/ChizhovVadim/counterdev/internal/arena"
	"github.com/ChizhovVadim/counterdev/internal/enginemini"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

// генерирует датасет для обучения нейронной сети
func datasetHandler(args []string) error {
	var date = time.Now()
	var (
		gamesCount        = 10
		openingRandomness = 50
		openingSize       = 8
		fixedNodes        = 50_000
		eval              = ""
		concurrency       = runtime.NumCPU()
	)

	var flagset = flag.NewFlagSet("", flag.ExitOnError)
	flagset.IntVar(&gamesCount, "games_count", gamesCount, "")
	flagset.IntVar(&openingRandomness, "opening_randomness", openingRandomness, "")
	flagset.IntVar(&openingSize, "opening_size", openingSize, "")
	flagset.IntVar(&fixedNodes, "nodes", fixedNodes, "")
	flagset.StringVar(&eval, "eval", eval, "")
	flagset.Parse(args)

	var openings = arena.GenerateRandomOpenings(gamesCount,
		enginemini.NewEngine(buildEvaluator(eval), 128),
		openingRandomness, openingSize, 3*fixedNodes)
	return arena.GenerateDataset(context.Background(), openings, concurrency,
		func() arena.Player {
			return arena.Player{
				Name:      "countermini",
				TimeLimit: common.LimitsType{Nodes: fixedNodes},
				Engine:    enginemini.NewEngine(buildEvaluator(eval), 128),
			}
		},
		mapPath(fmt.Sprintf("~/chess/dataset/arena-%v.pgn", date.Format("2006-01-02_15_04"))))
}
