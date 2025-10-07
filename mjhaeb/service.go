package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/mahjong"
)

func init() {
	mahjong.Service = NewService()
}

type service struct {
	tiles        map[mahjong.Tile]int
	defaultRules []int
}

func NewService() mahjong.IService {
	s := &service{
		tiles:        make(map[mahjong.Tile]int),
		defaultRules: []int{10, 8},
	}
	s.init()
	return s
}

func (s *service) init() {
	for color := mahjong.ColorCharacter; color <= mahjong.ColorDot; color++ {
		pc := mahjong.PointCountByColor[color]
		for i := range pc {
			tile := mahjong.MakeTile(color, i)
			s.tiles[tile] = 4
		}
	}
	s.tiles[mahjong.TileZhong] = 4
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

func (s *service) GetHuTypes(data *mahjong.HuData) []int32 {
	return getHuTypes(data)
}

func (s *service) TotalMuti(types []int32) int64 {
	return totalMuti(types)
}
