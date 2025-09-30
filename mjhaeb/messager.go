package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/game"
	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_proto/haebpb"
	"github.com/kevin-chtw/tw_proto/mjpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Messager struct {
	game *Game
	play *Play
}

func ToCallData(callData map[mahjong.Tile]map[mahjong.Tile]int64) map[int32]*mjpb.CallData {
	result := make(map[int32]*mjpb.CallData)
	for tile, callMap := range callData {
		callPb := &mjpb.CallData{
			CallTiles: make(map[int32]int64),
		}
		for tile, fan := range callMap {
			callPb.CallTiles[int32(tile)] = fan
		}
		result[int32(tile)] = callPb
	}
	return result
}

func NewMessager(game *Game) *Messager {
	return &Messager{
		game: game,
		play: game.Play,
	}
}

func (m *Messager) sendMsg(msg proto.Message, seat int32) error {
	data, err := anypb.New(msg)
	if err != nil {
		return err
	}
	ack := &haebpb.HAEBAck{Ack: data}
	m.game.Send2Player(ack, game.SeatAll)
	return nil
}

func (m *Messager) sendGameStartAck() {
	startAck := &mjpb.MJGameStartAck{
		Banker:    m.play.GetBanker(),
		TileCount: m.play.GetDealer().GetRestCount(),
		Scores:    m.play.GetCurScores(),
		Property:  m.game.GetRule().ToString(),
	}
	m.sendMsg(startAck, game.SeatAll)
}

func (m *Messager) sendOpenDoorAck() {
	count := m.game.GetPlayerCount()
	for i := range count {
		openDoor := &mjpb.MJOpenDoorAck{
			Seat:     i,
			Tiles:    m.play.GetPlayData(i).GetHandTilesInt32(),
			CallData: ToCallData(m.play.GetPlayData(i).GetCallDataMap()),
		}
		m.sendMsg(openDoor, i)
	}
}

func (m *Messager) sendAnimationAck() {
	animationAck := &mjpb.MJAnimationAck{
		Requestid: m.game.GetRequestID(game.SeatAll),
	}
	m.sendMsg(animationAck, game.SeatAll)
}

func (m *Messager) sendRequestAck(seat int32, operates *mahjong.Operates) {
	requestAck := &mjpb.MJRequestAck{
		Seat:        seat,
		RequestType: int32(operates.Value),
		Requestid:   m.game.GetRequestID(seat),
	}
	m.sendMsg(requestAck, seat)
}

func (m *Messager) sendDiscardAck() {
	discardAck := &mjpb.MJDiscardAck{
		Seat: m.play.GetCurSeat(),
		Tile: m.play.GetCurTile().ToInt32(),
	}
	m.sendMsg(discardAck, game.SeatAll)
}

func (m *Messager) sendTingAck(seat int32, tile mahjong.Tile) {
	tingAck := &mjpb.MJTingAck{
		Seat:     seat,
		Tile:     tile.ToInt32(),
		TianTing: m.play.GetPlayData(seat).IsTianTing(),
	}
	m.sendMsg(tingAck, game.SeatAll)
}

func (m *Messager) sendChowAck(seat int32, leftTile mahjong.Tile) {
	chowAck := &mjpb.MJChowAck{
		Seat:     seat,
		From:     m.play.GetCurSeat(),
		Tile:     m.play.GetCurTile().ToInt32(),
		LeftTile: leftTile.ToInt32(),
		CallData: ToCallData(m.play.GetPlayData(seat).GetCallDataMap()),
	}
	m.sendMsg(chowAck, game.SeatAll)
}

func (m *Messager) sendPonAck(seat int32) {
	ponAck := &mjpb.MJPonAck{
		Seat:     seat,
		From:     m.play.GetCurSeat(),
		Tile:     m.play.GetCurTile().ToInt32(),
		CallData: ToCallData(m.play.GetPlayData(seat).GetCallDataMap()),
	}
	m.sendMsg(ponAck, game.SeatAll)
}

func (m *Messager) sendKonAck(seat int32, tile mahjong.Tile, konType mahjong.KonType) {
	konAck := &mjpb.MJKonAck{
		Seat:    seat,
		From:    m.play.GetCurSeat(),
		Tile:    tile.ToInt32(),
		KonType: int32(konType),
	}
	m.sendMsg(konAck, game.SeatAll)
}

func (m *Messager) sendHuAck(huSeats []int32, paoSeat int32) {
	huAck := &mjpb.MJHuAck{
		PaoSeat: paoSeat,
		Tile:    m.play.GetCurTile().ToInt32(),
		HuData:  make([]*mjpb.MJHuData, len(huSeats)),
	}
	for i := range huSeats {
		huAck.HuData[i] = &mjpb.MJHuData{
			Seat:    huSeats[i],
			HuTypes: m.play.GetHuResult(huSeats[i]).HuTypes,
		}
	}
	m.sendMsg(huAck, game.SeatAll)
}

func (m *Messager) sendDrawAck(tile mahjong.Tile) {
	drawAck := &mjpb.MJDrawAck{
		Seat:     m.play.GetCurSeat(),
		Tile:     tile.ToInt32(),
		CallData: ToCallData(m.play.GetPlayData(m.play.GetCurSeat()).GetCallDataMap()),
	}
	m.sendMsg(drawAck, drawAck.Seat)
	drawAck.Tile = mahjong.TileNull.ToInt32()
	drawAck.CallData = nil
	for i := range m.game.GetPlayerCount() {
		if i != drawAck.Seat {
			m.sendMsg(drawAck, i)
		}
	}
}

func (m *Messager) sendResult(liuju bool) {
	resultAck := &mjpb.MJResultAck{
		Liuju:         liuju,
		PlayerResults: make([]*mjpb.MJPlayerResult, m.game.GetPlayerCount()),
	}
	for i := range resultAck.PlayerResults {
		resultAck.PlayerResults[i] = &mjpb.MJPlayerResult{
			Seat:     int32(i),
			CurScore: m.game.GetPlayer(int32(i)).GetCurScore(),
			WinScore: m.game.GetPlayer(int32(i)).GetScoreChangeWithTax(),
			Tiles:    m.play.GetPlayData(int32(i)).GetHandTilesInt32(),
		}
	}
	m.sendMsg(resultAck, game.SeatAll)
}
