package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
)

func init() {
	mahjong.Service = NewService()
}

type service struct {
	tiles        map[mahjong.Tile]int
	defaultRules []int
	fdRules      map[string]int32
}

func NewService() mahjong.IService {
	s := &service{
		tiles:        make(map[mahjong.Tile]int),
		defaultRules: []int{10, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		fdRules:      make(map[string]int32),
	}
	s.initTiles()
	s.initFdRules()
	return s
}

func (s *service) initTiles() {
	for color := mahjong.ColorCharacter; color <= mahjong.ColorDot; color++ {
		pc := mahjong.PointCountByColor[color]
		for i := range pc {
			tile := mahjong.MakeTile(color, i)
			s.tiles[tile] = 4
		}
	}
	s.tiles[mahjong.TileZhong] = 4
}

func (s *service) initFdRules() {
	s.fdRules["huytpe"] = RuleHuType
	s.fdRules["loubao"] = RuleLouBao
	s.fdRules["baozhongbao"] = RuleBaoZhongBao
	s.fdRules["hzmtf"] = RuleHZMTF
	s.fdRules["guadafen"] = RuleGuaDaFeng
	s.fdRules["btcsp"] = RuleTingChuShouPao
	s.fdRules["37jia"] = Rule37Jia
	s.fdRules["dandiaojia"] = RuleDanDiaoJia
	s.fdRules["duipengjia"] = RuleDuiPengJia
	s.fdRules["dmtkj"] = RuleDuoMianTingJia
	s.fdRules["7dui"] = Rule7Dui
	s.fdRules["konsoufen"] = RuleKonSouFen
	s.fdRules["qinyise"] = RuleQinYiSe
	s.fdRules["mengqin3bei"] = RuleMengQin3Bei
	s.fdRules["xiandahoumo"] = RuleXianDaHouMo
}

func (s *service) GetAllTiles(conf *mahjong.Rule) map[mahjong.Tile]int {
	return s.tiles
}

func (s *service) GetHandCount() int {
	return 13
}

func (s *service) GetDefaultRules() []int {
	return s.defaultRules
}

func (s *service) GetFdRules() map[string]int32 {
	return s.fdRules
}

func (s *service) GetHuTypes(data *mahjong.HuData) []int32 {
	return getHuTypes(data)
}

func (s *service) TotalMuti(types []int32) int64 {
	return totalMuti(types)
}
