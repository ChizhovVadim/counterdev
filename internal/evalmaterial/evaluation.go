package evalmaterial

import "github.com/ChizhovVadim/counterdev/pkg/common"

// Оценочная функция учитывает только материал. Удобно для:
// оценить performance движка без затрат на оценку.
// тактические тесты решает даже на такой оценочной функции.
// Сравнить ошибку разных оценочных функций на валидационном датасете.
type EvaluationService struct{}

func NewEvaluationService() *EvaluationService {
	return &EvaluationService{}
}

func (e *EvaluationService) Init(p *common.Position)                    {}
func (e *EvaluationService) MakeMove(p *common.Position, m common.Move) {}
func (e *EvaluationService) UnmakeMove()                                {}
func (e *EvaluationService) EvaluateQuick(p *common.Position) int {
	var val = 100*(common.PopCount(p.Pawns&p.White)-common.PopCount(p.Pawns&p.Black)) +
		400*(common.PopCount(p.Knights&p.White)-common.PopCount(p.Knights&p.Black)) +
		400*(common.PopCount(p.Bishops&p.White)-common.PopCount(p.Bishops&p.Black)) +
		600*(common.PopCount(p.Rooks&p.White)-common.PopCount(p.Rooks&p.Black)) +
		1200*(common.PopCount(p.Queens&p.White)-common.PopCount(p.Queens&p.Black))
	if !p.WhiteMove {
		val = -val
	}
	return val
}
