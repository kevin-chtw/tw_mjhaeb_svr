package mjhaeb

import (
	"github.com/kevin-chtw/tw_common/game"
	"github.com/kevin-chtw/tw_common/mahjong"
	"github.com/kevin-chtw/tw_proto/haebpb"
)

type Messager struct {
	game *Game
	play *Play
}

func NewMessager(game *Game) *Messager {
	return &Messager{
		game: game,
		play: game.GetPlay().(*Play),
	}
}

func (m *Messager) sendGameStartAck() {
	startAck := &haebpb.HAEBGameStartAck{
		Banker:    m.play.GetBanker(),
		TileCount: m.play.GetDealer().GetRestCount(),
		Scores:    m.play.GetCurScores(),
	}
	ack := &haebpb.HAEBAck{HaebGameStartAck: startAck}
	m.game.Send2Player(ack, game.SeatAll)
}

func (m *Messager) sendOpenDoorAck() {
	count := m.game.GetPlayerCount()
	for i := range count {
		openDoor := &haebpb.HAEBOpenDoorAck{
			Seat:  i,
			Tiles: m.play.GetPlayData(i).GetHandTiles(),
		}
		ack := &haebpb.HAEBAck{HaebOpenDoorAck: openDoor}
		m.game.Send2Player(ack, i)
	}
}

func (m *Messager) sendAnimationAck() {
	animationAck := &haebpb.HAEBAnimationAck{
		Requestid: m.game.GetRequestID(game.SeatAll),
	}
	ack := &haebpb.HAEBAck{HaebAnimationAck: animationAck}
	m.game.Send2Player(ack, game.SeatAll)
}

func (m *Messager) sendRequestAck(seat int32, operates *mahjong.Operates) {
	requestAck := &haebpb.HAEBRequestAck{
		Seat:        seat,
		RequestType: int32(operates.Value),
		Requestid:   m.game.GetRequestID(seat),
	}
	ack := &haebpb.HAEBAck{HaebRequestAck: requestAck}
	m.game.Send2Player(ack, seat)
}

func (m *Messager) sendDiscardAck() {
	discardAck := &haebpb.HAEBDiscardAck{
		Seat: m.play.GetCurSeat(),
		Tile: m.play.GetCurTile(),
	}
	ack := &haebpb.HAEBAck{HaebDiscardAck: discardAck}
	m.game.Send2Player(ack, game.SeatAll)
}

func (m *Messager) sendPonAck(seat int32) {
	ponAck := &haebpb.HAEBPonAck{
		Seat: seat,
		From: m.play.GetCurSeat(),
		Tile: m.play.GetCurTile(),
	}
	ack := &haebpb.HAEBAck{HaebPonAck: ponAck}
	m.game.Send2Player(ack, game.SeatAll)
}

func (m *Messager) sendKonAck(seat, tile int32, konType mahjong.KonType) {
	konAck := &haebpb.HAEBKonAck{
		Seat:    seat,
		From:    m.play.GetCurSeat(),
		Tile:    tile,
		KonType: int32(konType),
	}
	ack := &haebpb.HAEBAck{HaebKonAck: konAck}
	m.game.Send2Player(ack, game.SeatAll)
}

func (m *Messager) sendHuAck(huSeats []int32, paoSeat int32) {
	huAck := &haebpb.HAEBHuAck{
		PaoSeat: paoSeat,
		Tile:    m.play.GetCurTile(),
		HuData:  make([]*haebpb.HAEBHuData, len(huSeats)),
	}
	for i := range huSeats {
		huAck.HuData[i] = &haebpb.HAEBHuData{
			Seat:    huSeats[i],
			HuTypes: m.play.GetHuResult(huSeats[i]).HuTypes,
		}
	}
	ack := &haebpb.HAEBAck{HaebHuAck: huAck}
	m.game.Send2Player(ack, game.SeatAll)
}

func (m *Messager) sendDrawAck(tile int32) {
	drawAck := &haebpb.HAEBDrawAck{
		Seat: m.play.GetCurSeat(),
		Tile: tile,
	}
	ack := &haebpb.HAEBAck{HaebDrawAck: drawAck}
	m.game.Send2Player(ack, drawAck.Seat)
	drawAck.Seat = mahjong.SeatNull
	for i := range m.game.GetPlayerCount() {
		if i != drawAck.Seat {
			m.game.Send2Player(ack, i)
		}
	}
}

func (m *Messager) sendResult(isLiuJu bool, paoSeat, paoCiSeat int32) {
	// 实现发送麻将结果
}
