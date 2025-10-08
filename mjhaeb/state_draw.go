package mjhaeb

import "github.com/kevin-chtw/tw_common/mahjong"

type StateDraw struct {
	*State
}

func NewStateDraw(game mahjong.IGame, args ...any) mahjong.IState {
	return &StateDraw{
		State: NewState(game),
	}
}

func (s *StateDraw) OnEnter() {
	tile := s.game.play.Draw()
	if tile == mahjong.TileNull {
		s.game.SetNextState(NewStateLiuju)
		return
	}
	if s.game.play.swapBaoTile() {
		if s.game.play.bao == mahjong.TileNull {
			s.game.SetNextState(NewStateLiuju)
			return
		}
		for i := range s.game.GetPlayerCount() {
			if s.game.play.GetPlayData(i).IsTing() {
				s.game.sender.sendBaoAck()
			}
		}
	}
	s.game.sender.SendDrawAck(tile)
	s.game.SetNextState(NewStateDiscard)
}
