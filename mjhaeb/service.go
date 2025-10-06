package mjhaeb

import (
	"slices"

	"github.com/kevin-chtw/tw_common/mahjong"
)

func init() {
	mahjong.Service = NewService()
}

type service struct {
	tiles        map[mahjong.Tile]int
	defaultRules []int
	huCore       *mahjong.HuCore
}

func NewService() mahjong.IService {
	s := &service{
		tiles:        make(map[mahjong.Tile]int),
		defaultRules: []int{10, 8},
		huCore:       mahjong.NewHuCore(14),
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

func (s *service) CheckHu(data *mahjong.HuData, rule *mahjong.Rule) (*mahjong.HuResult, bool) {
	if !s.huCore.CheckBasicHu(data.Tiles, data.LaiCount) {
		return nil, false
	}
	extraHuTypes := make([]int32, 0)
	if data.ExtraHuTypes != nil {
		extraHuTypes = slices.Clone(data.ExtraHuTypes)
	}
	result := &mahjong.HuResult{
		HuTypes:   extraHuTypes,
		TotalMuti: totalMuti(extraHuTypes),
	}
	return result, true
}

func (s *service) CheckCall(data *mahjong.HuData, rule *mahjong.Rule) map[mahjong.Tile]map[mahjong.Tile]int64 {
	callData := make(map[mahjong.Tile]map[mahjong.Tile]int64)
	count := len(data.Tiles) % 3
	switch count {
	case 2:
		tileSet := make(map[mahjong.Tile]bool)
		for _, tile := range data.Tiles {
			tileSet[tile] = true
		}

		tempData := *data
		for tile := range tileSet {
			tempData.Tiles = mahjong.RemoveElements(data.Tiles, tile, 1)
			fans := s.checkCalls(&tempData, rule)
			if len(fans) > 0 {
				callData[tile] = fans
			}
		}
	case 1:
		// 直接检查叫牌
		fans := s.checkCalls(data, rule)
		if len(fans) > 0 {
			callData[0] = fans
		}
	}

	return callData
}

func (s *service) checkCalls(data *mahjong.HuData, rule *mahjong.Rule) map[mahjong.Tile]int64 {
	mutils := make(map[mahjong.Tile]int64)
	testTiles := s.GetAllTiles(rule)
	originalTiles := slices.Clone(data.Tiles)
	for tile := range testTiles {
		data.Tiles = append(data.Tiles, tile)
		if result, ok := s.CheckHu(data, rule); ok {
			mutils[tile] = result.TotalMuti
		}
		data.Tiles = originalTiles
	}
	return mutils
}
