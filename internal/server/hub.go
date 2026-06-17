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
	broadcast      chan []byte
	selectedClient *Client
	register       chan *Client
	unregister     chan *Client
	clients        map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		selectedClient: &Client{},
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
	}
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
			select {
			case h.selectedClient.send <- message:
			default:
				close(h.selectedClient.send)
				delete(h.clients, h.selectedClient)
			}
		}
	}
}

func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}
