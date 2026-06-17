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

	logger.Info("Preparando motor de captura global...")

	// Tenta registrar o atalho
	hook.Register(hook.KeyDown, []string{"ctrl", "shift", "s"}, func(e hook.Event) {
		isCapturing = !isCapturing
		if isCapturing {
			logger.Info("Modo Windows ATIVADO. Mouse preso no Mac.")
			robotgo.Move(centerX, centerY)
		} else {
			logger.Info("Modo Mac ATIVADO. Controle local restaurado.")
		}
	})

	// Inicia a escuta
	s := hook.Start()
	defer hook.End()

	logger.Info("Motor de captura INICIADO. Aguardando eventos de hardware...")

	// Loop infinito lendo TUDO o que acontece no Mac
	for ev := range s {
		// DEBUG: Se você descomentar a linha abaixo, ele vai imprimir QUALQUER tecla/clique
		// logger.Debug("Evento de hardware detectado", slog.Any("kind", ev.Kind), slog.Int("keycode", int(ev.Keycode)))

		if !isCapturing {
			continue
		}

		if ev.Kind == hook.MouseMove {
			deltaX := int(ev.X) - centerX
			deltaY := int(ev.Y) - centerY

			if deltaX != 0 || deltaY != 0 {
				event := protocol.InputEvent{
					EventType: "mouse_move",
					DeltaX:    deltaX,
					DeltaY:    deltaY,
				}

				sendEvent(event)
				robotgo.Move(centerX, centerY)
			}
		}
	}
}
