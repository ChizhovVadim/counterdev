package dataset

import (
	"iter"

	"github.com/ChizhovVadim/counterdev/internal/ml"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type IFeatureProvider interface {
	ComputeFeatures(pos *common.Position) ml.Input
	FeatureSize() int
}

func LoadSamples(
	dataset iter.Seq2[ml.DatasetItem, error],
	featureProvider IFeatureProvider,
	mirrorPos bool,
	maxSize int,
) ([]ml.Sample, error) {
	// сразу задаем capacity, чтобы при переалокации не тратить память на 2 вектора.
	var res = make([]ml.Sample, 0, maxSize) //var res []model.Sample
	for item, err := range dataset {
		if err != nil {
			return nil, err
		}
		var input = featureProvider.ComputeFeatures(&item.Position)
		res = append(res, ml.Sample{
			Input:  input,
			Target: float32(item.Target),
		})
		if len(res) >= maxSize {
			break
		}
		if mirrorPos {
			var mirrorPos = common.MirrorPosition(&item.Position)
			var mirrorInput = featureProvider.ComputeFeatures(&mirrorPos)
			var mirrorTarget = 1 - item.Target
			res = append(res, ml.Sample{
				Input:  mirrorInput,
				Target: float32(mirrorTarget),
			})
			if len(res) >= maxSize {
				break
			}
		}
	}
	return res, nil
}
