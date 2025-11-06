package mjhaeb

import (
	"errors"
	"time"

	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
	"github.com/kevin-chtw/tw_proto/game/pbmj"
	"google.golang.org/protobuf/proto"
)

type StateTing struct {
	*State
}

func NewStateTing(game mahjong.IGame, args ...any) mahjong.IState {
	return &StateTing{
		State: NewState(game),
	}
}

func (s *StateTing) OnEnter() {
	discardTime := s.game.GetRule().GetValue(RuleDiscardTime) + 1
	s.AsyncMsgTimer(s.OnMsg, time.Duration(discardTime)*time.Second, s.OnTimeout)
}

func (s *StateTing) OnMsg(seat int32, msg proto.Message) error {
	if seat != s.game.play.GetCurSeat() {
		return errors.New("not current seat")
	}

	optReq, ok := msg.(*pbmj.MJRequestReq)
	if !ok {
		return nil
	}
	if optReq == nil || optReq.Seat != seat || !s.game.sender.IsRequestID(seat, optReq.Requestid) {
		return errors.New("msg error")
	}

	if optReq.RequestType != int32(mahjong.OperateTing) {
		return errors.New("no request type")
	}

	s.ting(mahjong.Tile(optReq.Tile))
	return nil
}

func (s *StateTing) ting(tile mahjong.Tile) {
	if s.game.play.Ting(tile) {
		s.game.sender.SendTingAck(s.game.play.GetCurSeat(), tile)
		s.game.sender.sendBaoAck()
		s.game.SetNextState(NewStateWait)
	}
}

func (s *StateTing) OnTimeout() {
	if s.game.MatchType == "fdtable" {
		return
	}
	s.ting(mahjong.TileNull)
	s.game.sender.SendTrustAck(s.game.play.GetCurSeat(), true)
}
