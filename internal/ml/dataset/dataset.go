package dataset

import (
	"iter"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"github.com/ChizhovVadim/counterdev/internal/pgn"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type DatasetItem struct {
	Position common.Position
	Target   float64
}

func LoadDataset(
	fileNames []string,
	sigmoidScale float64,
	searchRatio float64,
) iter.Seq2[DatasetItem, error] {
	return func(yield func(DatasetItem, error) bool) {
		for _, path := range fileNames {
			for gameRaw, err := range pgn.LoadGames(path) {
				if err != nil {
					yield(DatasetItem{}, err)
					return
				}
				var g, err = pgn.ParseGame(gameRaw)
				if err != nil {
					yield(DatasetItem{}, err)
					return
				}
				var items = AnalyzeGame(&g, sigmoidScale, searchRatio)
				for i := range items {
					if !yield(items[i], nil) {
						return
					}
				}
			}
		}
	}
}

/*
важно решить:
результат поиска относим к позиции до или после сделанного хода
какие позиции фильтруем (дебют, лучший ход - взятие, повторы, под шахом, малая глубина перебора)
*/
func AnalyzeGame(
	g *game.Game,
	sigmoidScale float64,
	searchRatio float64,
) []DatasetItem {
	var gameRes, gameResOk = calcGameResult(g.Result)
	if !gameResOk {
		return nil
	}
	var replay, err = game.NewGame(g.StartFen)
	if err != nil {
		return nil
	}
	var result []DatasetItem
	for _, item := range g.Moves {
		if !(item.IsOpening ||
			item.Depth < 1 ||
			item.Score.Mate != 0 ||
			replay.Position.IsCheck() ||
			item.Move.CapturedPiece() != common.Empty ||
			item.Move.Promotion() != common.Empty ||
			replay.Position.Rule50 >= 90 ||
			replay.Repeats[replay.Position.Key] >= 2) {
			// mix game result and search result
			var target = (1-searchRatio)*gameRes +
				searchRatio*computeSearchTarget(sigmoidScale, item.Score, replay.Position.WhiteMove)
			result = append(result, DatasetItem{
				Position: replay.Position,
				Target:   target,
			})
		}
		if !replay.MakeMove(item) {
			panic("replay.MakeMove")
		}
	}
	return result
}

func calcGameResult(res int) (float64, bool) {
	switch res {
	case game.GameResultWhiteWins:
		return 1, true
	case game.GameResultBlackWins:
		return 0, true
	case game.GameResultDraw:
		return 0.5, true
	default:
		return 0, false
	}
}

func computeSearchTarget(
	sigmoidScale float64,
	searchScore common.UciScore,
	whiteMove bool,
) float64 {
	var res float64
	if searchScore.Mate != 0 {
		if searchScore.Mate > 0 {
			res = 1
		} else {
			res = 0
		}
	} else {
		res = sigmoid(sigmoidScale * float64(searchScore.Centipawns))
	}
	if !whiteMove {
		res = 1 - res
	}
	return res
}
