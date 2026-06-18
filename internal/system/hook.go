package system

import (
	"log/slog"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/romeritomendes/controlpc/internal/protocol"
)

func StartInputCapture(logger *slog.Logger, sendEvent func(protocol.Message)) {
	isCapturing := false
	screenWidth, screenHeight := robotgo.GetScreenSize()
	centerX, centerY := screenWidth/2, screenHeight/2

	logger.Info("Motor INICIADO. Pressione 'F10' para alternar o controle.")

	s := hook.Start()
	defer hook.End()

	for ev := range s {
		// 1. O Atalho de Teclado (F10)
		if ev.Kind == hook.KeyDown && (ev.Keychar == 65535 || int(ev.Keycode) == 109) {
			isCapturing = !isCapturing
			if isCapturing {
				logger.Info("Modo Windows ATIVADO. Mouse preso no Mac.")
				robotgo.Move(centerX, centerY)
			} else {
				logger.Info("Modo Mac ATIVADO. Controle local restaurado.")
			}
			continue
		}

		if !isCapturing {
			continue
		}

		// 2. O Movimento do Mouse (Sem loop de feedback)
		if ev.Kind == hook.MouseMove {
			x, y := int(ev.X), int(ev.Y)

			// SEGREDO AQUI: Se o mouse já está no centro exato, ignora o evento.
			// Isso impede o loop infinito de feedback que fazia o mouse tremer.
			if x == centerX && y == centerY {
				continue
			}

			// Calcula o delta a partir do centro
			deltaX := x - centerX
			deltaY := y - centerY

			event := protocol.Message{
				Type:   "mouse_move",
				DeltaX: deltaX,
				DeltaY: deltaY,
			}
			sendEvent(event)

			robotgo.Move(centerX, centerY)
		}

		// 3. Os Cliques do Mouse (Botão Esquerdo)
		if ev.Kind == hook.MouseDown || ev.Kind == hook.MouseUp {
			action := "down"
			if ev.Kind == hook.MouseUp {
				action = "up"
			}

			// Botão 1 geralmente é o esquerdo no gohook
			button := "left"
			if ev.Button == 2 {
				button = "right"
			}

			event := protocol.Message{
				Type:   "mouse_click",
				Button: button, // "left" ou "right"
				Key:    action, // Usaremos o campo Key para avisar se apertou ("down") ou soltou ("up")
			}
			sendEvent(event)
		}
	}
}
