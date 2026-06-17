// internal/server/hub.go
package server

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Método exposto para o main.go conseguir enviar mensagens
func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			// A MÁGICA ESTÁ AQUI:
			// Itera sobre todos os clientes reais conectados no momento
			for client := range h.clients {
				select {
				case client.send <- message:
					// Mensagem enviada com sucesso para o buffer do cliente
				default:
					// Se o buffer do cliente estiver cheio (travou), fechamos a conexão dele de forma segura
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
