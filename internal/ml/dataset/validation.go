package dataset

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"strings"

	"github.com/ChizhovVadim/counterdev/internal/ml"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

func LoadValidationDataset(path string) iter.Seq2[ml.DatasetItem, error] {
	return func(yield func(ml.DatasetItem, error) bool) {
		file, err := os.Open(path)
		if err != nil {
			if !yield(ml.DatasetItem{}, err) {
				return
			}
		}
		defer file.Close()

		var scanner = bufio.NewScanner(file)
		for scanner.Scan() {
			var s = scanner.Text()
			var item, err = parseItem(s)
			if !yield(item, err) {
				break
			}
		}
	}
}

func parseItem(s string) (ml.DatasetItem, error) {
	var index = strings.Index(s, "\"")
	if index < 0 {
		return ml.DatasetItem{}, fmt.Errorf("zurichessParser failed %v", s)
	}

	var fen = s[:index]
	var pos, err = common.NewPositionFromFEN(fen)
	if err != nil {
		return ml.DatasetItem{}, err
	}

	var strScore = s[index+1:]

	var prob float64
	if strings.HasPrefix(strScore, "1/2-1/2") {
		prob = 0.5
	} else if strings.HasPrefix(strScore, "1-0") {
		prob = 1.0
	} else if strings.HasPrefix(strScore, "0-1") {
		prob = 0.0
	} else {
		return ml.DatasetItem{}, fmt.Errorf("zurichessParser failed %v", s)
	}
	return ml.DatasetItem{Position: pos, Target: prob}, nil
}
