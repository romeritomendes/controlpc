// Package protocol
// How to comunicate with the clients
package protocol

type Message struct {
	Type       string // "handshake", "mouse_move", "mouse_click"
	ClientName string // Usado apenas no "handshake" para identificar a máquina
	DeltaX     int    // Usado no "mouse_move"
	DeltaY     int
	Button     string // Usado no "mouse_click"
	Key        string // "up" ou "down"
}
