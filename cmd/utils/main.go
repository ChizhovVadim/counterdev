// Утилиты, которые используются при разработке движка, но не нужны для запуска самого движка.
package main

import "log"

func main() {
	var app = &App{}
	app.AddCommand("counterdev", counterdevHandler)
	app.AddCommand("tactic", tacticHandler)
	app.AddCommand("arena", matchHandler)
	app.AddCommand("quality", qualityHandler)
	app.AddCommand("perft", perftHandler)
	app.AddCommand("opening", openingHandler)
	app.AddCommand("dataset", datasetHandler)
	app.AddCommand("train", trainHandler)
	var err = app.Run()
	if err != nil {
		log.Println("run failed",
			"error", err)
		return
	}
}
