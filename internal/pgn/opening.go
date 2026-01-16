package pgn

import (
	"bufio"
	"io"
	"strings"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

func LoadOpenings(r io.Reader) ([]game.Game, error) {
	var res []game.Game
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		opening, err := parseOpening(line)
		if err != nil {
			return nil, err
		}
		res = append(res, opening)
	}
	return res, nil
}

func parseOpening(strMoves string) (game.Game, error) {
	var g, err = game.NewGame("")
	if err != nil {
		return game.Game{}, err
	}
	var tokens = parsePgnBody(strMoves)
	for i := range tokens {
		var san = strings.TrimSpace(tokens[i].Value)
		var move = common.ParseMoveSAN(&g.Position, san)
		if move == common.MoveEmpty {
			break
		}
		if !g.MakeMove(game.MoveItem{Move: move, IsOpening: true}) {
			break
		}
	}
	return g, nil
}
