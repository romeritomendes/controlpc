package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Println("Rastreando a posição do mouse... (Pressione Ctrl+C para sair)")

	for {
		// robotgo.Location() retorna as coordenadas X e Y globais do mouse
		x, y := robotgo.Location()

		// Imprime as coordenadas. O \r faz com que a linha seja sobrescrita no terminal
		fmt.Printf("\rPosição atual: X: %d, Y: %d", x, y)

		// Uma pequena pausa para não sobrecarregar o uso da CPU
		time.Sleep(100 * time.Millisecond)
	}
}
