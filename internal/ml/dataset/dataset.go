package dataset

import (
	"iter"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"github.com/ChizhovVadim/counterdev/internal/ml"
	"github.com/ChizhovVadim/counterdev/internal/pgn"
)

type IGameMarker interface {
	MarkGame(g *game.Game) []ml.DatasetItem
}

func LoadDataset(
	fileNames []string,
	gameMarker IGameMarker,
) iter.Seq2[ml.DatasetItem, error] {
	return func(yield func(ml.DatasetItem, error) bool) {
		for _, path := range fileNames {
			for gameRaw, err := range pgn.LoadGames(path) {
				if err != nil {
					yield(ml.DatasetItem{}, err)
					return
				}
				var g, err = pgn.ParseGame(gameRaw)
				if err != nil {
					yield(ml.DatasetItem{}, err)
					return
				}
				var items = gameMarker.MarkGame(&g)
				for i := range items {
					if !yield(items[i], nil) {
						return
					}
				}
			}
		}
	}
}
