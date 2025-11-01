package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
)

type StateDeal struct {
	*State
}

func NewStateDeal(game mahjong.IGame, args ...any) mahjong.IState {
	return &StateDeal{
		State: NewState(game),
	}
}

func (s *StateDeal) OnEnter() {
	s.game.play.Deal()
	s.game.play.initBaoTile()
	s.game.sender.SendOpenDoorAck()
	for i := range s.game.GetPlayerCount() {
		s.game.play.FreshCallData(i)
		s.game.sender.SendCallDataAck(i)
	}
	s.WaitAni(func() { s.game.SetNextState(NewStateDiscard) })
}
