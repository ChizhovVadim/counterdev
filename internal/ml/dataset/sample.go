package dataset

import (
	"iter"

	"github.com/ChizhovVadim/counterdev/internal/ml/model"
	"github.com/ChizhovVadim/counterdev/pkg/common"
)

type IFeatureProvider interface {
	ComputeFeatures(pos *common.Position) model.Input
	FeatureSize() int
}

func LoadSamples(
	dataset iter.Seq2[DatasetItem, error],
	featureProvider IFeatureProvider,
	mirrorPos bool,
	maxSize int,
) ([]model.Sample, error) {
	// сразу задаем capacity, чтобы при переалокации не тратить память на 2 вектора.
	var res = make([]model.Sample, 0, maxSize) //var res []model.Sample
	for item, err := range dataset {
		if err != nil {
			return nil, err
		}
		var input = featureProvider.ComputeFeatures(&item.Position)
		res = append(res, model.Sample{
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
			res = append(res, model.Sample{
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
