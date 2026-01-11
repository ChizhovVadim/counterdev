package enginemini

import "github.com/ChizhovVadim/counterdev/pkg/common"

type HistoryTable struct {
	mainHistory [8192]int16
}

func (t *HistoryTable) Clear() {
	for i := range t.mainHistory {
		t.mainHistory[i] = 0
	}
}

// Exponential moving average
func updateHistory(v *int16, bonus int, good bool) {
	const historyMax = 1 << 14

	var newVal int
	if good {
		newVal = historyMax
	} else {
		newVal = -historyMax
	}
	*v += int16((newVal - int(*v)) * bonus / 512)
}

func sideFromToIndex(side bool, move common.Move) int {
	var result = (move.From() << 6) | move.To()
	if side {
		result |= 1 << 12
	}
	return result
}
