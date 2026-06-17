// Package protocol
// How to comunicate with the clients
package protocol

type InputEvent struct {
	EventType string `json:"type"` // Para garantir que apenas o alvo execute
	DeltaX    int    `json:"dx,omitempty"`
	DeltaY    int    `json:"dy,omitempty"`
	Button    string `json:"button,omitempty"`
	Key       string `json:"key,omitempty"`
}
