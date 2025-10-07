package mjhaeb

import (
	"errors"
	"time"

	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_proto/game/pbmj"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
	"google.golang.org/protobuf/proto"
)

type ReqOperate struct {
	Operate int32        //操作
	tile    mahjong.Tile //吃牌吃最左的牌
	disTile mahjong.Tile //吃听、碰听时出的牌
}

type StateWait struct {
	*State
	operatesForSeats   []*mahjong.Operates   // 每个座位可执行的操作
	reqOperateForSeats map[int32]*ReqOperate // 每个座位已请求的操作
}

func NewStateWait(game mahjong.IGame, args ...any) mahjong.IState {
	s := &StateWait{
		State: NewState(game),
	}
	s.operatesForSeats = make([]*mahjong.Operates, s.game.GetPlayerCount())
	s.reqOperateForSeats = make(map[int32]*ReqOperate)
	return s
}

func (s *StateWait) OnEnter() {
	discardSeat := s.game.play.GetCurSeat()
	for i := int32(0); i < s.game.GetPlayerCount(); i++ {
		if i == discardSeat {
			continue
		}
		trusted := s.game.GetPlayer(i).IsTrusted()
		operates := s.game.play.FetchWaitOperates(i)
		s.operatesForSeats[i] = operates

		if operates.Value != mahjong.OperatePass && !trusted {
			s.game.sender.SendRequestAck(i, operates)
		} else {
			s.setReqOperate(i, s.getDefaultOperate(i), s.game.play.GetCurTile(), mahjong.TileNull)
		}
	}

	timeout := s.game.GetRule().GetValue(RuleWaitTime) + 1
	logger.Log.Infof("discardSeat:%d timeout:%d", discardSeat, timeout)
	s.AsyncMsgTimer(s.OnMsg, time.Second*time.Duration(timeout), s.Timeout)
	s.tryHandleAction()
}

func (s *StateWait) OnMsg(seat int32, msg proto.Message) error {
	optReq, ok := msg.(*pbmj.MJRequestReq)
	if !ok {
		return nil
	}
	if optReq == nil || optReq.Seat != seat || !s.game.IsRequestID(seat, optReq.Requestid) {
		return errors.New("invalid msg")
	}

	if !s.isValidOperate(seat, int(optReq.RequestType)) {
		return errors.New("invalid operate")
	}
	s.setReqOperate(seat, optReq.RequestType, mahjong.Tile(optReq.Tile), mahjong.Tile(optReq.DisTile))
	s.tryHandleAction()
	return nil
}

func (s *StateWait) Timeout() {
	logger.Log.Info("timeout", s.operatesForSeats)
	for i := int32(0); i < s.game.GetPlayerCount(); i++ {
		if i == s.game.play.GetCurSeat() {
			continue
		}
		if _, ok := s.reqOperateForSeats[i]; !ok {
			s.setReqOperate(i, s.getDefaultOperate(i), s.game.play.GetCurTile(), mahjong.TileNull)
		}
	}
	s.tryHandleAction()
}

func (s *StateWait) setReqOperate(seat, operate int32, tile, disTile mahjong.Tile) {
	if s.game.IsValidSeat(seat) {
		s.reqOperateForSeats[seat] = &ReqOperate{Operate: operate, tile: tile, disTile: disTile}
	}
}

func (s *StateWait) tryHandleAction() {
	curSeat := s.game.play.GetCurSeat()
	huSeats := make([]int32, 0)
	for i := int32(1); i < s.game.GetPlayerCount(); i++ {
		seat := mahjong.GetNextSeat(curSeat, i, s.game.GetPlayerCount())
		if operate, ok := s.reqOperateForSeats[seat]; ok {
			if operate.Operate == mahjong.OperateHu {
				huSeats = append(huSeats, seat)
			}
		} else if s.getMaxOperate(seat) == mahjong.OperateHu {
			return
		}
	}

	if len(huSeats) > 0 {
		s.excuteHu(huSeats)
		return
	}

	maxOper := &ReqOperate{Operate: mahjong.OperatePass, tile: s.game.play.GetCurTile()}
	maxOperSeat := mahjong.SeatNull
	isMaxReq := true
	for i := int32(1); i < s.game.GetPlayerCount(); i++ {
		seat := mahjong.GetNextSeat(curSeat, i, s.game.GetPlayerCount())
		if operate, ok := s.reqOperateForSeats[seat]; ok {
			if operate.Operate > maxOper.Operate {
				maxOper = operate
				maxOperSeat = seat
				isMaxReq = true
			}
		} else if operate := s.getMaxOperate(seat); operate > maxOper.Operate {
			maxOper = &ReqOperate{Operate: operate, tile: s.game.play.GetCurTile()}
			maxOperSeat = seat
			isMaxReq = false
		}
	}
	if isMaxReq {
		s.excuteOperate(maxOperSeat, maxOper)
	}
}

func (s *StateWait) excuteOperate(seat int32, operate *ReqOperate) {
	if operate.Operate == mahjong.OperatePonTing {
		ponTile := s.game.play.GetCurTile()
		s.game.play.PonTing(seat, operate.disTile)
		s.game.sender.SendPonAck(seat, ponTile)
		s.game.sender.SendTingAck(seat, operate.disTile)
		s.toWaitState(seat)
		return
	}

	if operate.Operate == mahjong.OperateChowTing {
		chowTile := s.game.play.GetCurTile()
		s.game.play.ChowTing(seat, operate.tile, operate.disTile)
		s.game.sender.SendChowAck(seat, chowTile, operate.tile)
		s.game.sender.SendTingAck(seat, operate.disTile)
		s.toWaitState(seat)
		return
	}
	if operate.Operate == mahjong.OperateKon {
		s.game.play.ZhiKon(seat)
		s.game.sender.SendKonAck(seat, s.game.play.GetCurTile(), mahjong.KonTypeZhi)
		s.toDrawState(seat)
		return
	}
	if operate.Operate == mahjong.OperatePon {
		s.game.play.Pon(seat)
		s.game.sender.SendPonAck(seat, s.game.play.GetCurTile())
		s.toDiscardState(seat)
		return
	}
	if operate.Operate == mahjong.OperateChow {
		s.game.play.Chow(seat, operate.tile)
		s.game.sender.SendChowAck(seat, s.game.play.GetCurTile(), operate.tile)
		s.toDiscardState(seat)
		return
	}
	s.toDrawState(mahjong.SeatNull)
}

func (s *StateWait) excuteHu(huSeats []int32) {
	s.game.SetNextState(NewStatePaohu, huSeats)
}

func (s *StateWait) toWaitState(seat int32) {
	s.game.play.DoSwitchSeat(seat)
	s.game.SetNextState(NewStateWait)
}

func (s *StateWait) toDrawState(seat int32) {
	s.game.play.DoSwitchSeat(seat)
	s.game.SetNextState(NewStateDraw)
}

func (s *StateWait) toDiscardState(seat int32) {
	s.game.play.DoSwitchSeat(seat)
	s.game.SetNextState(NewStateDiscard)
}

func (s *StateWait) isValidOperate(seat int32, operate int) bool {
	// 检查操作是否有效
	if !s.game.IsValidSeat(seat) {
		return false
	}
	if s.operatesForSeats[seat] == nil {
		return false
	}
	return s.operatesForSeats[seat].HasOperate(int32(operate))
}

func (s *StateWait) getMaxOperate(seat int32) int32 {
	if ops := s.operatesForSeats[seat]; ops != nil {
		if ops.HasOperate(mahjong.OperateHu) {
			return mahjong.OperateHu
		}
		if ops.HasOperate(mahjong.OperatePonTing) {
			return mahjong.OperatePonTing
		}
		if ops.HasOperate(mahjong.OperateChowTing) {
			return mahjong.OperateChowTing
		}
		if ops.HasOperate(mahjong.OperateKon) {
			return mahjong.OperateKon
		}
		if ops.HasOperate(mahjong.OperatePon) {
			return mahjong.OperatePon
		}
		if ops.HasOperate(mahjong.OperateChow) {
			return mahjong.OperateChow
		}
	}
	return mahjong.OperatePass
}

func (s *StateWait) getDefaultOperate(seat int32) int32 {
	ops := s.operatesForSeats[seat]
	if ops != nil && ops.HasOperate(mahjong.OperateHu) {
		return mahjong.OperateHu
	}
	return mahjong.OperatePass
}
