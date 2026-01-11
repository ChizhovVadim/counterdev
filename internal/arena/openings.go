package arena

import (
	"context"
	"iter"
	"log"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

func GenerateRandomOpenings(
	count int,
	eng common.IEngine,
	randomness int,
	openingSize int,
	fixedNodes int,
) iter.Seq[game.Game] {
	return func(yield func(game.Game) bool) {
		var repeats = make(map[uint64]struct{})
		var filteredRepeats int
		defer func() {
			log.Println("GenerateRandomOpenings finished",
				"repeats", filteredRepeats)
		}()
		for index := 0; index < count; {
			var opening, ok = playRandomOpening(eng, randomness, openingSize, fixedNodes)
			if !ok {
				continue
			}
			if _, found := repeats[opening.Position.Key]; found {
				filteredRepeats += 1
				continue
			}
			repeats[opening.Position.Key] = struct{}{}
			index += 1
			if !yield(opening) {
				return
			}
		}
	}
}

func playRandomOpening(
	eng common.IEngine,
	randomness int,
	openingSize int,
	fixedNodes int,
) (game.Game, bool) {
	var g, _ = game.NewGame("")
	for range openingSize {
		var searchRes = eng.Search(context.Background(), common.SearchParams{
			Position:   g.Position,
			Repeats:    g.Repeats,
			Limits:     common.LimitsType{Nodes: fixedNodes},
			Randomness: randomness,
		})
		if !g.MakeMove(game.MoveItem{
			Move:      searchRes.MainLine[0],
			IsOpening: true,
		}) {
			return game.Game{}, false
		}
	}
	return g, true
}
