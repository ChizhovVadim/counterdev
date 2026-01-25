package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/ChizhovVadim/counterdev/internal/arena"
	"github.com/ChizhovVadim/counterdev/internal/enginemini"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

// генерирует датасет для обучения нейронной сети
func datasetHandler(args []string) error {
	var date = time.Now()
	var (
		gamesCount        = 100
		openingRandomness = 75
		openingSize       = 9
	)

	var flagset = flag.NewFlagSet("", flag.ExitOnError)
	flagset.IntVar(&gamesCount, "games_count", gamesCount, "")
	flagset.IntVar(&openingRandomness, "opening_randomness", openingRandomness, "")
	flagset.IntVar(&openingSize, "opening_size", openingSize, "")
	flagset.Parse(args)

	return arena.GenerateDataset(context.Background(), gamesCount, openingRandomness, openingSize,
		func() arena.Player {
			return arena.Player{
				Name:      "countermini",
				TimeLimit: common.LimitsType{Nodes: 150_000},
				Engine:    enginemini.NewEngine(buildEvaluator("nnue"), 32),
			}
		},
		func() arena.Player {
			return arena.Player{
				Name:      "countermini",
				TimeLimit: common.LimitsType{Depth: 8},
				Engine:    enginemini.NewEngine(buildEvaluator("weiss"), 32),
			}
		},
		mapPath(fmt.Sprintf("~/chess/dataset/arena-%v.pgn", date.Format("2006-01-02_15_04"))))
}
