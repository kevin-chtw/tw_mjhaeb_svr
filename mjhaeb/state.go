package mjhaeb

import (
	"time"

	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_proto/game/pbmj"
	"google.golang.org/protobuf/proto"
)

type State struct {
	*mahjong.State
	game  *Game
	aniFn func()
}

func NewState(game mahjong.IGame) *State {
	g := game.(*Game)
	return &State{
		State: mahjong.NewState(g.Game),
		game:  g,
		aniFn: nil,
	}
}

func (s *State) OnAniMsg(seat int32, msg proto.Message) error {
	aniReq, ok := msg.(*pbmj.MJAnimationReq)
	if !ok {
		return nil
	}
	if aniReq != nil && seat == aniReq.Seat && s.game.IsRequestID(seat, aniReq.Requestid) {
		s.aniFn()
	}
	return nil
}

func (s *State) WaitAni(reqFn func()) {
	s.game.sender.SendAnimationAck()
	s.aniFn = reqFn
	s.AsyncMsgTimer(s.OnAniMsg, time.Second*5, reqFn)
}
