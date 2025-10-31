package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
	"github.com/kevin-chtw/tw_proto/game/pbmj"
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

	huTypes := c.play.paoHuTypes(seat)
	if len(huTypes) == 0 {
		return
	} else if c.play.PlayConf.OnlyZimo {
		opt.Tips = append(opt.Tips, mahjong.TipsOnlyZiMo)
		return
	}
	result := &pbmj.MJHuData{
		HuTypes: huTypes,
		Multi:   totalMuti(huTypes),
	}
	c.play.AddHuOperate(opt, seat, result, true)
}
