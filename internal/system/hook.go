package system

import (
	"log/slog"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/romeritomendes/controlpc/internal/protocol"
)

func StartInputCapture(logger *slog.Logger, sendEvent func(protocol.InputEvent)) {
	isCapturing := false
	screenWidth, screenHeight := robotgo.GetScreenSize()
	centerX, centerY := screenWidth/2, screenHeight/2

	// Variáveis para calcular o Delta sem forçar o mouse pro centro toda hora
	lastX, lastY := centerX, centerY

	logger.Info("Motor INICIADO. Pressione 'F10' para alternar o controle.")

	s := hook.Start()
	defer hook.End()

	for ev := range s {
		// 1. O Atalho de Teclado (F10)
		if ev.Kind == hook.KeyDown && int(ev.Keycode) == 109 {
			isCapturing = !isCapturing
			if isCapturing {
				logger.Info("Modo Windows ATIVADO. Mouse preso no Mac.")
				robotgo.Move(centerX, centerY)
				lastX, lastY = centerX, centerY // Reseta a posição base
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
			deltaX := int(ev.X) - lastX
			deltaY := int(ev.Y) - lastY

			if deltaX != 0 || deltaY != 0 {
				event := protocol.InputEvent{
					EventType: "mouse_move",
					DeltaX:    deltaX,
					DeltaY:    deltaY,
				}
				sendEvent(event)

				lastX, lastY = int(ev.X), int(ev.Y)

				// Bounding Box: Só joga pro centro se chegar perto da borda (evita o loop)
				if lastX < 100 || lastX > screenWidth-100 || lastY < 100 || lastY > screenHeight-100 {
					robotgo.Move(centerX, centerY)
					lastX, lastY = centerX, centerY
				}
			}
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

			event := protocol.InputEvent{
				EventType: "mouse_click",
				Button:    button, // "left" ou "right"
				Key:       action, // Usaremos o campo Key para avisar se apertou ("down") ou soltou ("up")
			}
			sendEvent(event)
		}
	}
}
