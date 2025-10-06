package mjhaeb

import "github.com/kevin-chtw/tw_common/mahjong"

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
	p.Play = mahjong.NewPlay(game.Game, p.dealer)

	p.ExtraHuTypes = p
	p.PlayConf = &mahjong.PlayConf{}
	p.RegisterSelfCheck(&mahjong.CheckerHu{},
		&mahjong.CheckerTing{},
		&mahjong.CheckerKon{},
	)
	p.RegisterWaitCheck(
		&mahjong.CheckerPao{},
		&mahjong.CheckerChow{},
		&mahjong.CheckerPon{},
		&mahjong.CheckerZhiKon{},
		&mahjong.CheckerChowTing{},
		&mahjong.CheckerPonTing{},
	)
	return p
}

func (p *Play) SelfExtraFans() []int32 {
	hyTypes := make([]int32, 0)
	hyTypes = append(hyTypes, HuTypeZiMo)
	return hyTypes
}

func (p *Play) PaoExtraFans() []int32 {
	return []int32{}
}

func (p *Play) InitBaoTile() {
	p.bao = p.dealer.LastTile()
}
