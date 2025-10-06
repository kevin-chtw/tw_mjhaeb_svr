package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/mahjong"
)

type StateLiuju struct {
	*State
}

func NewStateLiuju(game mahjong.IGame, args ...any) mahjong.IState {
	return &StateLiuju{
		State: NewState(game),
	}
}

func (s *StateLiuju) OnEnter() {
	s.game.sender.SendResult(true)
	s.WaitAni(s.game.OnGameOver)
}
