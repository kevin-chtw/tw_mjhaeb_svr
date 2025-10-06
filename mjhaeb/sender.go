package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/game"
	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_proto/game/pbhaeb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Sender struct {
	*mahjong.Sender
	play *Play
}

func NewSender(game *Game) *Sender {
	s := &Sender{
		play: game.play,
	}
	s.Sender = mahjong.NewSender(game.Game, game.play.Play, s)
	return s
}

func (s *Sender) PackMsg(msg proto.Message) (proto.Message, error) {
	data, err := anypb.New(msg)
	if err != nil {
		return nil, err
	}
	ack := &pbhaeb.HAEBAck{Ack: data}
	return ack, nil
}

func (s *Sender) sendBaoAck() {
	ack := &pbhaeb.HaebBaoAck{
		Tile: s.play.bao.ToInt32(),
	}
	s.SendMsg(ack, game.SeatAll)
}
