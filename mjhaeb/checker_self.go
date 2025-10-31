package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
	"github.com/kevin-chtw/tw_proto/game/pbmj"
)

type CheckerHu struct {
	play *Play
}

func NewCheckerHu(play *Play) mahjong.CheckerSelf {
	return &CheckerHu{play: play}
}

func (c *CheckerHu) Check(opt *mahjong.Operates) {
	playData := c.play.GetPlayData(c.play.GetCurSeat())
	if !playData.IsTing() {
		return
	}

	huTypes := c.play.selfHuTypes()
	if len(huTypes) == 0 {
		return
	}
	result := &pbmj.MJHuData{
		HuTypes: huTypes,
		Multi:   totalMuti(huTypes),
	}

	opt.RemoveOperate(mahjong.OperateDiscard)
	c.play.AddHuOperate(opt, c.play.GetCurSeat(), result, true)
}

// 听检查器
type CheckerTing struct {
	play *Play
}

func NewCheckerTing(play *Play) mahjong.CheckerSelf {
	return &CheckerTing{play: play}
}
func (c *CheckerTing) Check(opt *mahjong.Operates) {
	if opt.IsMustHu {
		return
	}

	playData := c.play.GetPlayData(c.play.GetCurSeat())
	if playData.IsTing() {
		return
	}
	huData := mahjong.NewHuData(playData, true)
	callData := huData.CheckCall()
	if len(callData) <= 0 {
		return
	}

	if playData.IsMenQin() {
		opt.Tips = append(opt.Tips, mahjong.TipsMenQin)
		return
	}

	opt.AddOperate(mahjong.OperateTing)
}
