package mjhaeb

import (
	"time"

	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_proto/haebpb"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
	"google.golang.org/protobuf/proto"
)

type StateResult struct {
	*State
	huSeats []int32
}

func NewStateResult(game mahjong.IGame) *StateResult {
	return &StateResult{
		State:   NewState(game),
		huSeats: make([]int32, 0),
	}
}

func (s *StateResult) onMsg(seat int32, msg proto.Message) error {
	req := msg.(*haebpb.HAEBReq)
	aniReq := req.GetHaebAnimationReq()
	if aniReq != nil && seat == aniReq.Seat && s.game.IsRequestID(seat, aniReq.Requestid) {
		s.game.OnGameOver()
	}
	return nil
}

func (s *StateResult) handleOver() {
	for _, seat := range s.huSeats {
		tiles := s.GetPlay().GetPlayData(seat).GetHandTiles()
		logger.Log.Info(tiles)
	}

	s.game.GetMessager().sendAnimationAck()
	s.AsyncMsgTimer(s.onMsg, time.Second*5, s.game.OnGameOver)
}
