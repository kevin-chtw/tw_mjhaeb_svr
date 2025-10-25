package mjhaeb

import "github.com/kevin-chtw/tw_common/gamebase/mahjong"

type CheckerHu struct {
	play *Play
}

func NewCheckerHu(play *Play) mahjong.CheckerSelf {
	return &CheckerHu{play: play}
}

func (c *CheckerHu) Check(opt *mahjong.Operates, tips []int) []int {
	playData := c.play.GetPlayData(c.play.GetCurSeat())
	if !playData.IsTing() {
		return tips
	}

	huTypes := c.play.selfHuTypes()
	if len(huTypes) == 0 {
		return tips
	}
	result := &mahjong.HuResult{
		HuTypes:   huTypes,
		TotalMuti: totalMuti(huTypes),
	}

	opt.RemoveOperate(mahjong.OperateDiscard)
	c.play.AddHuOperate(opt, c.play.GetCurSeat(), result, true)
	return tips
}

// 听检查器
type CheckerTing struct {
	play *Play
}

func NewCheckerTing(play *Play) mahjong.CheckerSelf {
	return &CheckerTing{play: play}
}
func (c *CheckerTing) Check(opt *mahjong.Operates, tips []int) []int {
	if opt.IsMustHu {
		return tips
	}

	playData := c.play.GetPlayData(c.play.GetCurSeat())
	if playData.IsTing() {
		return tips
	}
	huData := mahjong.NewHuData(playData, true)
	callData := huData.CheckCall()
	if len(callData) <= 0 {
		return tips
	}

	if playData.IsMenQin() {
		tips = append(tips, mahjong.TipsMenQin)
		return tips
	}

	opt.AddOperate(mahjong.OperateTing)
	return tips
}
