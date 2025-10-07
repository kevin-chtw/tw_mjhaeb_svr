package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/mahjong"
)

const (
	HuTypePingHu      = iota //平胡 0
	HuTypeKaDang             //卡当 1
	HuTypeZiMo               //自摸 2
	HuTypeMoBao              //摸宝 3
	HuTypeLouBao             //搂宝 4
	HuTypeGuaDaFeng          //刮大风 5
	HuTypeBaoZhongBao        //宝中宝 6
)

var multiples = map[int32]int64{HuTypePingHu: 1, HuTypeKaDang: 2, HuTypeZiMo: 2, HuTypeMoBao: 3, HuTypeLouBao: 3, HuTypeGuaDaFeng: 6, HuTypeBaoZhongBao: 12}

func totalMuti(huTypes []int32) int64 {
	totalMuti := int64(1)
	for _, huType := range huTypes {
		if multiple, ok := multiples[huType]; ok {
			totalMuti *= multiple
		}
	}
	return totalMuti
}

func getHuTypes(huData *mahjong.HuData) []int32 {
	types := []int32{HuTypePingHu}
	if isKaDang(huData) {
		types[0] = HuTypeKaDang
	}
	return types
}
func isKaDang(huData *mahjong.HuData) bool {
	waitTile := huData.GetCurTile()
	if waitTile.IsHonor() {
		return false // 字牌不考虑卡当
	}

	tileMap := make(map[mahjong.Tile]int)
	for _, tile := range huData.Tiles {
		tileMap[tile]++
	}

	if count := tileMap[waitTile]; count != 1 && count != 4 {
		return false
	}

	color, point := waitTile.Info()
	// 检查卡张情况 (如3和5，听4)
	if point > 0 && point < 8 {
		if tileMap[mahjong.MakeTile(color, point-1)] > 0 && tileMap[mahjong.MakeTile(color, point+1)] > 0 {
			return true
		}
	}

	// 检查边张情况 (如01听2，78听6)
	if point == 2 && tileMap[mahjong.MakeTile(color, point-1)] > 0 && tileMap[mahjong.MakeTile(color, point-2)] > 0 {
		return true
	}
	if point == 6 && tileMap[mahjong.MakeTile(color, point+1)] > 0 && tileMap[mahjong.MakeTile(color, point+2)] > 0 {
		return true
	}
	return false
}
