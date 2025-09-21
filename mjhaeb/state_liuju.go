package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/mahjong"
)

type StateLiuju struct {
	*StateResult
}

func NewStateLiuju(game mahjong.IGame, args ...any) mahjong.IState {
	return &StateLiuju{
		StateResult: NewStateResult(game),
	}
}

func (s *StateLiuju) OnEnter() {
	s.game.GetMessager().sendResult(true)
	s.handleOver()
}
