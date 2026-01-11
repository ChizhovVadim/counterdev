package enginemini

import (
	"log"

	"github.com/ChizhovVadim/counterdev/pkg/common"
)

const (
	boundLower = 1 << iota
	boundUpper
)

const boundExact = boundLower | boundUpper

func roundPowerOfTwo(size int) int {
	var x = 1
	for (x << 1) <= size {
		x <<= 1
	}
	return x
}

// 16 bytes
type transEntry struct {
	key32 uint32
	move  common.Move
	score int16
	depth int8
	bound uint8
}

type transTable struct {
	megabytes int
	entries   []transEntry
	date      uint16
	mask      uint32
}

// good test: position fen 8/k7/3p4/p2P1p2/P2P1P2/8/8/K7 w - - 0 1
// good test: position fen 8/pp6/2p5/P1P5/1P3k2/3K4/8/8 w - - 5 47
func newTransTable(megabytes int) transTable {
	log.Println("Init trans table", "size", megabytes)
	var size = roundPowerOfTwo(1024 * 1024 * megabytes / 16)
	return transTable{
		megabytes: megabytes,
		entries:   make([]transEntry, size),
		mask:      uint32(size - 1),
	}
}

func (tt *transTable) Size() int {
	return tt.megabytes
}

func (tt *transTable) IncDate() {
	tt.date += 1
}

func (tt *transTable) Clear() {
	tt.date = 0
	for i := range tt.entries {
		tt.entries[i] = transEntry{}
	}
}

func (tt *transTable) Read(key uint64) (depth, score, bound int, move common.Move, ok bool) {
	var entry = &tt.entries[uint32(key)&tt.mask]
	if entry.key32 == uint32(key>>32) {
		depth = int(entry.depth)
		score = int(entry.score)
		bound = int(entry.bound)
		move = entry.move
		ok = true
	}
	return
}

func (tt *transTable) Update(key uint64, depth, score, bound int, move common.Move) {
	var entry = &tt.entries[uint32(key)&tt.mask]
	var found = entry.key32 == uint32(key>>32)
	entry.key32 = uint32(key >> 32)
	entry.depth = int8(depth)
	entry.score = int16(score)
	entry.bound = uint8(bound)
	if !(found && move == common.MoveEmpty) {
		entry.move = move
	}
}
