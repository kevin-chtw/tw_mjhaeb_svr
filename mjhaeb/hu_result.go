package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/gamebase/mahjong"
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
	if len(huData.PlayData.GetCallData()) > 1 { //卡当仅一个听口
		return false
	}

	waitTile := huData.CurTile
	point := waitTile.Point()

	if point <= 0 || point >= 8 { // 在0-8范围内，1-7需要检查相邻牌
		return false
	}

	if huData.CheckShun(waitTile, point-1, point+1) {
		return true
	}

	switch point {
	case 2:
		return huData.CheckShun(waitTile, point-1, point-2) && !huData.CheckShun(waitTile, point+1, point+2)
	case 6:
		return !huData.CheckShun(waitTile, point-1, point-2) && huData.CheckShun(waitTile, point+1, point+2)
	default:
		return false
	}
}
