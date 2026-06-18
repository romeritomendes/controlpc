// internal/server/hub.go
package server

import (
	"encoding/gob"
	"net"

	"github.com/romeritomendes/controlpc/internal/protocol"
)

type Client struct {
	Name string
	conn net.Conn
	enc  *gob.Encoder
	send chan protocol.Message
}

type Hub struct {
	clients    map[string]*Client
	broadcast  chan protocol.Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan protocol.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

// Método exposto para o main.go conseguir enviar mensagens
func (h *Hub) BroadcastMessage(message protocol.Message) {
	h.broadcast <- message
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Name] = client

		case client := <-h.unregister:
			if _, ok := h.clients[client.Name]; ok {
				delete(h.clients, client.Name)
				close(client.send)
			}

		case message := <-h.broadcast:
			// A MÁGICA ESTÁ AQUI:
			// Itera sobre todos os clientes reais conectados no momento
			for _, client := range h.clients {
				select {
				case client.send <- message:
					// Mensagem enviada com sucesso para o buffer do cliente
				default:
					// Se o buffer do cliente estiver cheio (travou), fechamos a conexão dele de forma segura
					close(client.send)
					delete(h.clients, client.Name)
				}
			}
		}
	}
}
