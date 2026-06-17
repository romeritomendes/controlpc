// Package client
// Client receive the MouseEvent
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

	// Tenta estabelecer a conexão WebSocket com o Mac
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Error("Erro ao conectar", slog.String("error", err.Error()))
		return
	}
	defer conn.Close()

	logger.Info("Conectado ao servidor com sucesso! Aguardando comandos...")

	// Inicia um loop infinito para ficar ouvindo as mensagens que chegam do Mac
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Conexão perdida ou erro ao ler", slog.String("error", err.Error()))
			break // Sai do loop se o servidor cair
		}

		// Transforma o JSON recebido na nossa struct InputEvent
		var event protocol.InputEvent
		if err := json.Unmarshal(message, &event); err != nil {
			logger.Error("Erro ao decodificar JSON do protocolo", slog.String("error", err.Error()))
			continue
		}

		// Chama a função que realmente interage com o Windows
		ExecuteCommand(logger, event)
	}
}

func ExecuteCommand(logger *slog.Logger, event protocol.InputEvent) {
	switch event.EventType {

	case "mouse_move":
		// Move o mouse a partir da posição ATUAL dele no Windows
		robotgo.MoveRelative(event.DeltaX, event.DeltaY)

	case "mouse_click":
		// Exemplo de como seria o clique
		// robotgo.Click(event.Button)

	default:
		logger.Warn("Evento desconhecido", slog.String("type", event.EventType))
	}
}
