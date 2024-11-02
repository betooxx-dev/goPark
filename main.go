package main

import (
	"log"
	"math/rand"
	"time"

	"simulador/src/ui"
)

func init() {
	// Inicializar el generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Configurar el log para desarrollo
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	// Crear los manejadores de eventos
	handlers := ui.NewHandlers()

	// Crear la interfaz gráfica
	gui := ui.NewGUI(handlers)

	// Establecer la referencia cruzada entre GUI y handlers
	handlers.SetGUI(gui)

	// Iniciar la aplicación
	log.Println("Iniciando simulador de estacionamiento...")
	gui.Run()
}
