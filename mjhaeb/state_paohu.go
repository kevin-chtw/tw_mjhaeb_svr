package mjhaeb

import (
	"time"

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
	// TODO 点炮胡牌
	//multiples := s.GetPlay().PaoHu()
	//s.game.GetScorelator().Calculate(multiples)
	s.game.GetMessager().sendResult(true, 0, 0)

	s.game.GetMessager().sendAnimationAck()
	s.AsyncMsgTimer(s.onMsg, time.Second*5, s.game.OnGameOver)
}
