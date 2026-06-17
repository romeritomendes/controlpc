package system

import (
	"log/slog"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/romeritomendes/controlpc/internal/protocol"
)

func StartInputCapture(logger *slog.Logger, sendToClient func(protocol.InputEvent)) {
	isCapturing := false

	screenWidth, screenHeight := robotgo.GetScreenSize()
	centerX, centerY := screenWidth/2, screenHeight/2

	logger.Info("Global Hook. Press 'ctrl+shift+s' to change the control")

	hook.Register(hook.KeyDown, []string{"ctrl", "shift", "s"}, func(e hook.Event) {
		isCapturing = !isCapturing
		if isCapturing {
			logger.Info("Modo Windows ATIVADO. Mouse preso no Mac.")
			robotgo.Move(centerX, centerY) // Joga pro centro
		} else {
			logger.Info("Modo Mac ATIVADO. Controle local restaurado.")
		}
	})

	// Inicia o canal de escuta de todos os eventos do sistema
	s := hook.Start()
	defer hook.End()

	for ev := range s {
		if !isCapturing {
			continue // Se não está capturando, ignora e deixa o Mac usar o mouse normalmente
		}

		if ev.Kind == hook.MouseMove {
			// Calcula o Delta (o quanto o mouse tentou sair do centro)
			deltaX := int(ev.X) - centerX
			deltaY := int(ev.Y) - centerY

			// Só envia se houve movimento real
			if deltaX != 0 || deltaY != 0 {
				event := protocol.InputEvent{
					EventType: "mouse_move",
					DeltaX:    deltaX,
					DeltaY:    deltaY,
				}

				// Envia para o WebSocket
				sendToClient(event)

				// O SEGREDO: Joga o mouse de volta pro centro imediatamente
				robotgo.Move(centerX, centerY)
			}
		}

		// Você faria a mesma coisa para cliques de mouse e digitação de teclado aqui
		// if ev.Kind == hook.MouseDown { ... }
	}
}
