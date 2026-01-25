package arena

import (
	"context"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

func playRandomOpening(
	g *game.Game,
	eng common.IEngine,
	randomness int,
	openingSize int,
	limits common.LimitsType,
) bool {
	for range openingSize {
		var searchRes = eng.Search(context.Background(), common.SearchParams{
			Position:   g.Position,
			Repeats:    g.Repeats,
			Limits:     limits,
			Randomness: randomness,
		})
		if !g.MakeMove(game.MoveItem{
			Move:      searchRes.MainLine[0],
			IsOpening: true,
		}) {
			return false
		}
	}
	return true
}
