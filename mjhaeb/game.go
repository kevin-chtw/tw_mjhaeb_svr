package mjhaeb

import (
	"errors"

	"github.com/kevin-chtw/tw_common/game"
	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_common/utils"
	"github.com/kevin-chtw/tw_proto/game/pbhaeb"
	"github.com/kevin-chtw/tw_proto/game/pbmj"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
)

type Game struct {
	*mahjong.Game
	play       *Play
	sender     *Sender
	scorelator mahjong.Scorelator
}

func NewGame(t *game.Table, id int32) game.IGame {
	g := &Game{}
	g.Game = mahjong.NewGame(g, t, id)
	g.play = NewPlay(g)
	g.sender = NewSender(g)
	g.scorelator = mahjong.NewScorelatorOnce(g.Game, mahjong.ScoreType(g.GetRule().GetValue(RuleScoreType)))

	return g
}

func (g *Game) OnStart() {
	g.Game.SetNextState(NewStateInit)
}

func (g *Game) OnReqMsg(player *game.Player, data []byte) error {
	msg := &pbhaeb.HAEBReq{}
	if err := utils.Unmarshal(player.Ctx, data, msg); err != nil {
		return err
	}
	logger.Log.Info(msg)

	req, err := msg.Req.UnmarshalNew()
	if err != nil {
		return err
	}

	trust, ok := req.(*pbmj.MJTrustReq)
	if ok && !trust.GetTrust() {
		g.sender.SendTrustAck(player.GetSeat(), trust.GetTrust())
		return nil
	}

	if g.CurState == nil {
		return errors.New("invalid state")
	}
	return g.CurState.OnPlayerMsg(player.GetSeat(), req)
}
