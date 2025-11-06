package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
)

type CheckerPao struct{ play *Play } // 点炮检查器
func NewCheckerPao(play *Play) mahjong.CheckerWait {
	return &CheckerPao{play: play}
}
func (c *CheckerPao) Check(seat int32, opt *mahjong.Operates) {
	playData := c.play.GetPlayData(seat)
	if !playData.IsTing() {
		return
	}

	data := mahjong.NewHuData(playData, false)
	result, hu := data.CheckHu()
	if !hu {
		return
	}

	if c.play.PlayConf.OnlyZimo {
		opt.Tips = append(opt.Tips, mahjong.TipsOnlyZiMo)
		return
	}
	c.play.AddHuOperate(opt, seat, result, true)
}
