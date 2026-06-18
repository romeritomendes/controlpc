package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/romeritomendes/controlpc/internal/client"
	"github.com/romeritomendes/controlpc/internal/protocol"
	"github.com/romeritomendes/controlpc/internal/server"
	"github.com/romeritomendes/controlpc/internal/system"
)

func main() {
	mode := flag.String("mode", "server", "Executar como 'server' ou 'client'")

	serverIP := flag.String("server", "localhost", "Endereço IP do servidor Mac (usado apenas no modo client)")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if *mode == "server" {
		logger.Info("Iniciando modo SERVIDOR (Mac)")

		// 1. Inicia o Hub
		hub := server.NewHub()
		go hub.Run()

		// 2. Inicia a captura global do sistema em uma Goroutine.
		// Passamos uma função de callback que será chamada sempre que o mouse mover.
		go system.StartInputCapture(logger, func(event protocol.Message) {
			// Envia o GOB para o Hub repassar aos clientes (Windows)
			hub.BroadcastMessage(event)
		})

		address := getLocalIP()
		if err := server.StartTCPServe(logger, hub, address); err != nil {
			logger.Error("Erro no servidor", slog.String("Error", err.Error()))
			os.Exit(1)
		}

	} else if *mode == "client" {
		wsURL := fmt.Sprintf("%s:3000", *serverIP)

		logger.Info("Iniciando modo CLIENTE (Windows)", slog.String("conectando_em", wsURL))
		client.Connect(logger, wsURL, "windows")

	} else {
		logger.Error("Modo inválido! Use -mode=server ou -mode=client", slog.String("voce_digitou", *mode))
		os.Exit(1)
	}
}

func getLocalIP() string {
	// Cria uma conexão UDP falsa (não envia dados, só descobre a rota)
	conn, err := net.Dial("udp", "8.8.8.8:3000")
	if err != nil {
		return "127.0.0.1:3000" // Fallback de segurança se estiver sem rede
	}
	defer conn.Close()

	// Extrai apenas a parte do IP (ignorando a porta)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return fmt.Sprintf("%s:3000", localAddr.IP.String())
}
