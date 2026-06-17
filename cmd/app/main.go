package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/romeritomendes/controlpc/internal/client"
	"github.com/romeritomendes/controlpc/internal/protocol"
	"github.com/romeritomendes/controlpc/internal/server"
	"github.com/romeritomendes/controlpc/internal/system"
)

func main() {
	mode := flag.String("mode", "server", "Executar como 'server' ou 'client'")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if *mode == "server" {
		logger.Info("Iniciando modo SERVIDOR (Mac)")

		// 1. Inicia o Hub
		hub := server.NewHub()
		go hub.Run()

		// 2. Inicia a captura global do sistema em uma Goroutine.
		// Passamos uma função de callback que será chamada sempre que o mouse mover.
		go system.StartInputCapture(logger, func(event protocol.InputEvent) {
			// Converte a struct de movimento em JSON
			data, err := json.Marshal(event)
			if err != nil {
				logger.Error("Erro ao empacotar evento", slog.String("error", err.Error()))
				return
			}
			// Envia o JSON para o Hub repassar aos clientes (Windows)
			hub.BroadcastMessage(data)
		})

		// 3. Registra a rota do WebSocket
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			server.ServeWS(logger, hub, w, r)
		})

		// 4. Trava a thread principal mantendo o servidor no ar
		if err := http.ListenAndServe(":3000", nil); err != nil {
			logger.Error("Erro no servidor", slog.String("Error", err.Error()))
			os.Exit(1)
		}

	} else {
		logger.Info("Iniciando modo CLIENTE (Windows)")
		client.Connect(logger, "ws://localhost:3000/ws")
	}
}
