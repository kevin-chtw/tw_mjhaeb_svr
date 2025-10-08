package mjhaeb

import (
	"errors"
	"time"

	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_proto/game/pbmj"
	"google.golang.org/protobuf/proto"
)

type StateDiscard struct {
	*State
	operates *mahjong.Operates
	handlers map[int32]func(tile mahjong.Tile)
}

func NewStateDiscard(game mahjong.IGame, args ...any) mahjong.IState {
	s := &StateDiscard{
		State:    NewState(game),
		handlers: make(map[int32]func(tile mahjong.Tile)),
	}
	s.handlers[mahjong.OperateDiscard] = s.discard
	s.handlers[mahjong.OperateTing] = s.ting
	s.handlers[mahjong.OperateKon] = s.kon
	s.handlers[mahjong.OperateHu] = s.hu
	return s
}

func (s *StateDiscard) OnEnter() {
	s.operates = s.game.play.FetchSelfOperates()
	s.game.sender.SendRequestAck(s.game.play.GetCurSeat(), s.operates)
	discardTime := s.game.GetRule().GetValue(RuleDiscardTime) + 1

	if s.game.GetPlayer(s.game.play.GetCurSeat()).IsTrusted() {
		s.discard(mahjong.TileNull)
		return
	}
	s.AsyncMsgTimer(s.OnMsg, time.Duration(discardTime)*time.Second, s.OnTimeout)

}

func (s *StateDiscard) OnMsg(seat int32, msg proto.Message) error {
	if seat != s.game.play.GetCurSeat() {
		return errors.New("not current seat")
	}

	optReq, ok := msg.(*pbmj.MJRequestReq)
	if !ok {
		return nil
	}
	if optReq == nil || optReq.Seat != seat || !s.game.IsRequestID(seat, optReq.Requestid) {
		return errors.New("msg error")
	}

	if !s.operates.HasOperate(optReq.RequestType) {
		return errors.New("no request type")
	}
	if handler, exists := s.handlers[optReq.RequestType]; exists {
		handler(mahjong.Tile(optReq.Tile))
	}
	return nil
}

func (s *StateDiscard) ting(tile mahjong.Tile) {
	if s.game.play.Ting(tile) {
		s.game.sender.SendTingAck(s.game.play.GetCurSeat(), tile)
		s.game.sender.sendBaoAck()
		s.game.SetNextState(NewStateWait)
	}
}

func (s *StateDiscard) discard(tile mahjong.Tile) {
	if s.game.play.Discard(tile) {
		s.game.sender.SendDiscardAck()
		s.game.SetNextState(NewStateWait)
	}
}

func (s *StateDiscard) kon(tile mahjong.Tile) {
	if s.game.play.TryKon(tile, mahjong.KonTypeBu) {
		s.game.sender.SendKonAck(s.game.play.GetCurSeat(), tile, mahjong.KonTypeBu)
		s.game.SetNextState(NewStateWait)
	} else if s.game.play.TryKon(tile, mahjong.KonTypeAn) {
		s.game.sender.SendKonAck(s.game.play.GetCurSeat(), tile, mahjong.KonTypeAn)
		s.game.SetNextState(NewStateDraw)
	}
}

func (s *StateDiscard) hu(tile mahjong.Tile) {
	s.game.SetNextState(NewStateZimo)
}

func (s *StateDiscard) OnTimeout() {
	if s.game.MatchType == "fdtable" {
		return
	}
	s.discard(mahjong.TileNull)
	s.game.sender.SendTrustAck(s.game.play.GetCurSeat(), true)
}
