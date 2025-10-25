package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
)

type Play struct {
	*mahjong.Play
	game   *Game
	dealer *mahjong.Dealer
	bao    mahjong.Tile
}

func NewPlay(game *Game) *Play {
	p := &Play{
		game:   game,
		dealer: mahjong.NewDealer(game.Game),
		bao:    mahjong.TileNull,
	}
	p.Play = mahjong.NewPlay(p, game.Game, p.dealer)

	p.PlayConf = &mahjong.PlayConf{}
	p.RegisterSelfCheck(
		NewCheckerHu(p),
		mahjong.NewCheckerTing(p.Play),
		mahjong.NewCheckerKon(p.Play),
	)
	p.RegisterWaitCheck(
		mahjong.NewCheckerPao(p.Play),
		mahjong.NewCheckerChow(p.Play),
		mahjong.NewCheckerPon(p.Play),
		mahjong.NewCheckerZhiKon(p.Play),
		mahjong.NewCheckerChowTing(p.Play),
		mahjong.NewCheckerPonTing(p.Play),
	)
	return p
}

func (p *Play) CheckHu(data *mahjong.HuData) bool {
	if data.IsTing() && p.isBaoTile(p.GetCurTile()) {
		return true
	}
	return data.CanHu()
}

func (p *Play) GetExtraHuTypes(playData *mahjong.PlayData, self bool) []int32 {
	types := make([]int32, 0)
	if playData.IsTing() {
		if self || p.PlayConf.OnlyZimo {
			types = p.selfHuTypes()
		} else {
			types = p.paoHuTypes(playData.GetSeat())
		}
	}
	return types
}

func (p *Play) initBaoTile() {
	p.bao = p.dealer.LastTile()
}

func (p *Play) swapBaoTile() bool {
	if p.dealer.Count(p.bao) <= 1 {
		p.bao = p.dealer.SwapLastTile()
		return true
	}
	return false
}

func (p *Play) isBaoTile(tile mahjong.Tile) bool {
	if p.bao == tile {
		return true
	}
	if p.GetRule().GetValue(RuleHZMTF) != 0 {
		return tile == mahjong.TileZhong
	}
	return false
}

func (p *Play) selfHuTypes() []int32 {
	types := p.huTypes(p.GetCurSeat(), 4, true)
	if len(types) != 0 {
		types = append(types, HuTypeZiMo)
	}
	return types
}

func (p *Play) paoHuTypes(seat int32) []int32 {
	return p.huTypes(seat, 3, false)
}

func (p *Play) huTypes(seat int32, guaDaFengNum int, self bool) (types []int32) {
	types = make([]int32, 0)
	playData := p.GetPlayData(seat)
	callData := playData.GetCallData()

	if p.bao == p.GetCurTile() { //摸宝
		types = append(types, HuTypeMoBao)
		if len(callData) == 1 {
			types = []int32{HuTypeBaoZhongBao}
		}
		return
	}

	if _, ok := callData[p.GetCurTile()]; !ok {
		return
	}
	if len(callData) == 1 && isKaDang(mahjong.NewHuData(playData, self)) {
		types = append(types, HuTypeKaDang)
	} else {
		types = append(types, HuTypePingHu)
	}

	groups := playData.GetPonGroups()
	for _, g := range groups {
		if g.Tile == p.GetCurTile() {
			types = append(types, HuTypeGuaDaFeng)
		}
	}

	hands := playData.GetHandTiles()
	if mahjong.CountElement(hands, p.GetCurTile()) == guaDaFengNum {
		types = append(types, HuTypeGuaDaFeng)
	}
	return types
}
