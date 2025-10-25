package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
)

type StateZimo struct {
	*State
}

func NewStateZimo(game mahjong.IGame, args ...any) mahjong.IState {
	return &StateZimo{
		State: NewState(game),
	}
}

func (s *StateZimo) OnEnter() {
	s.game.sender.SendHuAck([]int32{s.game.play.GetCurSeat()}, mahjong.SeatNull)
	multiples := s.game.play.Zimo()
	s.game.scorelator.AddMultiple(mahjong.ScoreReasonHu, multiples)
	s.game.scorelator.Calculate()
	s.game.sender.SendResult(false)

	s.WaitAni(s.game.OnGameOver)
}
