package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/mahjong"
)

type StatePaohu struct {
	*StateResult
}

func NewStatePaohu(game mahjong.IGame, args ...any) mahjong.IState {
	s := &StatePaohu{
		StateResult: NewStateResult(game),
	}
	s.huSeats = args[0].([]int32)
	return s
}

func (s *StatePaohu) OnEnter() {
	s.game.GetMessager().sendHuAck(s.huSeats, s.GetPlay().GetCurSeat())
	multiples := s.GetPlay().PaoHu(s.huSeats)
	s.game.GetScorelator().AddMultiple(mahjong.ScoreReasonHu, multiples)
	s.game.GetScorelator().Calculate()

	s.game.GetMessager().sendResult(false)
	s.handleOver()
}
