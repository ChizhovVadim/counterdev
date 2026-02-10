package uciadapter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type Engine struct {
	r *bufio.Scanner
	w io.Writer
}

func NewEngine(uci *UciProcess) *Engine {
	return &Engine{
		r: bufio.NewScanner(uci.out),
		w: uci.in,
	}
}

func (e *Engine) Uci() EngineInfo {
	fmt.Fprintln(e.w, "uci")
	for e.r.Scan() {
		var msg = e.r.Text()
		if msg == "uciok" {
			return EngineInfo{}
		}
	}
	// TODO err
	return EngineInfo{}
}

func (e *Engine) SetOption(option Option) error {
	fmt.Fprintf(e.w, "setoption name %v value %v\n", option.Name, option.Value)
	return nil
}

func (e *Engine) UciNewgame() {
	fmt.Fprintln(e.w, "ucinewgame")
}

func (e *Engine) IsReady() {
	fmt.Fprintln(e.w, "isready")
	for e.r.Scan() {
		var msg = e.r.Text()
		if msg == "readyok" {
			return
		}
	}
}

func (e *Engine) Stop() {
	fmt.Fprintln(e.w, "stop")
}

func (e *Engine) Position(game *game.Game) {
	if game.StartFen == "" || game.StartFen == common.InitialPositionFen {
		fmt.Fprintf(e.w, "position startpos")
	} else {
		fmt.Fprintf(e.w, "position fen %v", game.StartFen)
	}
	if len(game.Moves) > 0 {
		fmt.Fprintf(e.w, " moves")
		for _, m := range game.Moves {
			fmt.Fprintf(e.w, " %v", m)
		}
	}
	fmt.Fprintln(e.w)
}

func (e *Engine) Go(
	tc common.LimitsType,
	progress func(common.SearchInfo),
) (common.Move, error) {
	fmt.Fprintf(e.w, "go ")
	if tc.MoveTime != 0 {
		fmt.Fprintf(e.w, "movetime %v", tc.MoveTime)
	} else if tc.Nodes != 0 {
		fmt.Fprintf(e.w, "nodes %v", tc.Nodes)
	}
	fmt.Fprintln(e.w)

	for e.r.Scan() {
		var msg = e.r.Text()
		if msg == "" {
			continue
		}
		var tokens = Tokens{items: strings.Fields(msg)}
		if tokens.Scan() {
			if tokens.Text() == "info" {
				var si = parseSearchInfo(tokens)
				if si.Depth != 0 {
					if progress != nil {
						progress(si)
					}
				}
			} else if tokens.Text() == "bestmove" {
				//TODO
				var bestMove = common.MoveEmpty
				return bestMove, nil
			}
		}
	}

	return common.MoveEmpty, fmt.Errorf("uci protocol")
}

func parseSearchInfo(tokens Tokens) common.SearchInfo {
	var res common.SearchInfo
	for tokens.Scan() {
		var name = tokens.Text()
		if name == "depth" {
			if tokens.Scan() {
				res.Depth, _ = strconv.Atoi(tokens.Text())
			}
		} else if name == "nodes" {
			if tokens.Scan() {
				res.Nodes, _ = strconv.ParseInt(tokens.Text(), 10, 64)
			}
		} else if name == "time" {
			if tokens.Scan() {
				var t, _ = strconv.Atoi(tokens.Text())
				res.Time = time.Duration(t) * time.Millisecond
			}
		} else if name == "score" {
			if tokens.Scan() {
				if tokens.Text() == "cp" {
					if tokens.Scan() {
						var v, _ = strconv.Atoi(tokens.Text())
						res.Score = common.UciScore{Centipawns: v}
					}
				} else if tokens.Text() == "mate" {
					if tokens.Scan() {
						var v, _ = strconv.Atoi(tokens.Text())
						res.Score = common.UciScore{Mate: v}
					}
				}
			}
		} else if name == "pv" {
			/*for tokens.Scan() {
				common.ParseMoveLAN() tokens.Text()
			}*/
		}
	}
	return res
}
