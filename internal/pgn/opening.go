package pgn

import (
	"bufio"
	"os"
	"strings"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

// Каждая строка в файле - дебют в SAN формате:
// 1. e4 e6 2. Qe2
// 1. e4 e6 2. d3 d5 3. Nd2 Nf6 4. Ngf3
func LoadOpenings(openingsPath string) ([]game.Game, error) {
	var file, err = os.Open(openingsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var res []game.Game
	scanner := bufio.NewScanner(file)
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
