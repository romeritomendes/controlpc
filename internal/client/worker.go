package client

import (
	"encoding/json"
	"log/slog"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/romeritomendes/controlpc/internal/protocol"
)

func Connect(logger *slog.Logger, url string) {
	logger.Info("Tentando conectar ao servidor...", slog.String("url", url))

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Error("Erro ao conectar", slog.String("error", err.Error()))
		return
	}
	defer conn.Close()

	logger.Info("Conectado ao servidor com sucesso! Aguardando comandos...")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Conexão perdida", slog.String("error", err.Error()))
			break
		}

		var event protocol.InputEvent
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		ExecuteCommand(logger, event)
	}
}

func ExecuteCommand(logger *slog.Logger, event protocol.InputEvent) {
	switch event.EventType {
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
