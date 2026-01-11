package main

import (
	"flag"
	"os"

	"github.com/ChizhovVadim/counterdev/internal/arena"
	"github.com/ChizhovVadim/counterdev/internal/enginemini"
	"github.com/ChizhovVadim/counterdev/internal/pgn"
)

// генератор дебютов
// движок должен поддерживать randomness
func openingHandler(args []string) error {
	var (
		gamesCount  = 50
		randomness  = 50
		openingSize = 8
		fixedNodes  = 150_000
		eval        = ""
	)

	var flagset = flag.NewFlagSet("", flag.ExitOnError)
	flagset.IntVar(&gamesCount, "games_count", gamesCount, "")
	flagset.IntVar(&randomness, "randomness", randomness, "")
	flagset.IntVar(&openingSize, "size", openingSize, "")
	flagset.IntVar(&fixedNodes, "nodes", fixedNodes, "")
	flagset.StringVar(&eval, "eval", eval, "")
	flagset.Parse(args)

	var eng = enginemini.NewEngine(buildEvaluator(eval), 128)
	for opening := range arena.GenerateRandomOpenings(gamesCount, eng, randomness, openingSize, fixedNodes) {
		var err = pgn.WriteMoves(os.Stdout, &opening, false)
		if err != nil {
			return err
		}
	}
	return nil
}
