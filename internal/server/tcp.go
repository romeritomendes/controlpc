package server

import (
	"encoding/gob"
	"log/slog"
	"net"

	"github.com/romeritomendes/controlpc/internal/protocol"
)

func StartTCPServe(logger *slog.Logger, hub *Hub, port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer listener.Close()

	logger.Info("TCP Server started", slog.String("port", port))

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Connection rejected", slog.String("error", err.Error()))
			continue
		}

		go handleConnection(logger, hub, conn)
	}
}

func handleConnection(logger *slog.Logger, hub *Hub, conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	var msg protocol.Message
	if err := decoder.Decode(&msg); err != nil {
		logger.Error("Handshake failled", slog.String("error", err.Error()))
		conn.Close()
		return
	}

	if msg.Type != "handshake" {
		logger.Error("First message should be Handshake")
		conn.Close()
		return
	}

	client := &Client{
		Name: msg.ClientName,
		conn: conn,
		enc:  encoder,
		send: make(chan protocol.Message, 256),
	}

	hub.register <- client
	logger.Info("Client connected", slog.String("name", client.Name))

	go writePump(client)

	for {
		if err := decoder.Decode(&msg); err != nil {
			hub.unregister <- client
			conn.Close()
			break
		}
	}
}

func writePump(c *Client) {
	for msg := range c.send {
		c.enc.Encode(msg)
	}
}
