package client

import (
	"encoding/gob"
	"log/slog"
	"net"

	"github.com/go-vgo/robotgo"
	"github.com/romeritomendes/controlpc/internal/protocol"
)

func Connect(logger *slog.Logger, address string, myName string) {
	logger.Info("Tentando conectar ao servidor...", slog.String("url", address))

	conn, err := net.Dial("tcp", address)
	if err != nil {
		logger.Error("Erro ao conectar", slog.String("error", err.Error()))
		return
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	handshake := protocol.Message{
		Type:       "handshake",
		ClientName: myName,
	}
	if err := encoder.Encode(handshake); err != nil {
		logger.Error("Handshake error")
		return
	}

	logger.Info("Handshake started", slog.String("name", myName))

	for {
		var msg protocol.Message
		if err := decoder.Decode(&msg); err != nil {
			logger.Error("Conexão perdida", slog.String("error", err.Error()))
			break
		}

		ExecuteCommand(logger, msg)
	}
}

func ExecuteCommand(logger *slog.Logger, event protocol.Message) {
	switch event.Type {
	case "mouse_move":
		// Lê a posição exata atual do Windows
		atualX, atualY := robotgo.Location()

		// Soma o movimento recebido do Mac
		novoX := atualX + event.DeltaX
		novoY := atualY + event.DeltaY

		// Move o mouse para a coordenada exata
		robotgo.Move(novoX, novoY)

	case "mouse_click":
		// O event.Key terá "down" (apertou) ou "up" (soltou)
		// O event.Button terá "left" ou "right"
		robotgo.Toggle(event.Button, event.Key)

	default:
		// Ignora eventos desconhecidos
	}
}
