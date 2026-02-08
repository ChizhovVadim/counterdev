package uciadapter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type StockfishEvaluationService struct {
	r *bufio.Scanner
	w io.Writer
}

func NewStockfishEvaluationService(uci *UciProcess) *StockfishEvaluationService {
	return &StockfishEvaluationService{
		r: bufio.NewScanner(uci.out),
		w: uci.in,
	}
}

func (e *StockfishEvaluationService) Evaluate(p *common.Position) int {
	fmt.Fprintln(e.w, "position fen", p.String())
	fmt.Fprintln(e.w, "eval")
	var res int
	for e.r.Scan() {
		var msg = e.r.Text()
		if strings.HasPrefix(msg, "NNUE evaluation") {
			//NNUE evaluation        -0.18 (white side)
			var s = strings.TrimPrefix(msg, "NNUE evaluation")
			s = strings.TrimSuffix(s, "(white side)")
			s = strings.TrimSpace(s)
			var f, err = strconv.ParseFloat(s, 64)
			if err == nil {
				res = int(100 * f)
				if !p.WhiteMove {
					res = -res
				}
			}
		}
		if strings.HasPrefix(msg, "Final evaluation") {
			break
		}
	}
	return res
}
