package dataset

import (
	"github.com/ChizhovVadim/counterdev/internal/ml/model"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type Feature768Provider struct{}

func NewFeature768Provider() *Feature768Provider {
	return &Feature768Provider{}
}

// На самом деле признаков 736, тк пешки не могут быть на крайних горизонталях. Переделать?
func (p *Feature768Provider) FeatureSize() int { return 768 }

func (p *Feature768Provider) ComputeFeatures(pos *common.Position) model.Input {
	var input = make([]model.FeatureInfo, 0, common.PopCount(pos.AllPieces()))
	for x := pos.AllPieces(); x != 0; x &= x - 1 {
		var sq = common.FirstOne(x)
		var pt, side = pos.GetPieceTypeAndSide(sq)
		var piece12 = pt - common.Pawn
		if !side {
			piece12 += 6
		}
		var index = int16(sq ^ piece12<<6)
		if !(index >= 0 && index < 768) {
			panic("feature out of range")
		}
		input = append(input, model.FeatureInfo{
			Index: index,
			Value: 1,
		})
	}
	return model.Input{
		Features: input,
	}
}
