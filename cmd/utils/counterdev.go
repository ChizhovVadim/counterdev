package main

import (
	"flag"
	"log"

	"github.com/ChizhovVadim/counterdev/internal/enginemini"
	"github.com/ChizhovVadim/counterdev/pkg/uci"
)

// dev версия движка
func counterdevHandler(args []string) error {
	var eval = ""

	var flagset = flag.NewFlagSet("", flag.ExitOnError)
	flagset.StringVar(&eval, "eval", eval, "")
	flagset.Parse(args)

	var eng = enginemini.NewEngine(buildEvaluator(eval), 128)
	var protocol = uci.New("Counter", "Vadim Chizhov", "dev", eng, []uci.Option{
		&uci.BoolOption{Name: "ExperimentSettings", Value: &eng.ExperimentSettings},
	})
	protocol.Run(log.Default())
	return nil
}
