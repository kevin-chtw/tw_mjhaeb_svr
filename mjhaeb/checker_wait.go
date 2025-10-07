package mjhaeb

import "github.com/kevin-chtw/tw_common/mahjong"

type CheckerPao struct{ play *Play } // 点炮检查器
func NewCheckerPao(play *Play) mahjong.CheckerWait {
	return &CheckerPao{play: play}
}
func (c *CheckerPao) Check(seat int32, opt *mahjong.Operates, tips []int) []int {
	playData := c.play.GetPlayData(seat)
	if !playData.IsTing() {
		return tips
	}

	huTypes := c.play.paoHuTypes(seat)
	if len(huTypes) == 0 {
		return tips
	} else if c.play.PlayConf.OnlyZimo {
		tips = append(tips, mahjong.TipsOnlyZiMo)
		return tips
	}
	result := &mahjong.HuResult{
		HuTypes:   huTypes,
		TotalMuti: totalMuti(huTypes),
	}
	c.play.AddHuOperate(opt, seat, result, true)
	return tips
}
