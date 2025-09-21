package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/mahjong"
)

type StateZimo struct {
	*StateResult
}

func NewStateZimo(game mahjong.IGame, args ...any) mahjong.IState {
	return &StateZimo{
		StateResult: NewStateResult(game),
	}
}

func (s *StateZimo) OnEnter() {
	s.huSeats = append(s.huSeats, s.GetPlay().GetCurSeat())
	s.game.GetMessager().sendHuAck(s.huSeats, mahjong.SeatNull)

	multiples := s.GetPlay().Zimo()
	s.game.GetScorelator().AddMultiple(mahjong.ScoreReasonHu, multiples)
	s.game.GetScorelator().Calculate()

	s.game.GetMessager().sendResult(false)
	s.handleOver()
}
