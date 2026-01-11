package main

import (
	"flag"
	"time"

	"github.com/ChizhovVadim/counterdev/internal/enginemini"
	"github.com/ChizhovVadim/counterdev/internal/tactic"
)

// решает тактические тесты
func tacticHandler(args []string) error {
	var (
		tacticTestsPath = mapPath("~/chess/tests/tests.epd")
		moveTime        = 3 * time.Second
		eval            = ""
	)

	var flagset = flag.NewFlagSet("", flag.ExitOnError)
	flagset.StringVar(&tacticTestsPath, "path", tacticTestsPath, "")
	flagset.DurationVar(&moveTime, "time", moveTime, "")
	flagset.StringVar(&eval, "eval", eval, "")
	flagset.Parse(args)

	var tests, err = tactic.Load(tacticTestsPath)
	if err != nil {
		return err
	}
	var eng = enginemini.NewEngine(buildEvaluator(eval), 128)
	return tactic.Solve(tests, eng, moveTime)
}
