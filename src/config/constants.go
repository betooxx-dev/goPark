package config

import (
	"fmt"
	"time"
)

const (
	TotalParkingSpots = 20
	TotalCarsToProcess = 100
	WindowCloseDelay = 5 * time.Second
)

// Configuración de tiempos
const (
	// Rango de tiempo que un auto permanece estacionado
	MinParkingTime = 3 * time.Second
	MaxParkingTime = 5 * time.Second

	// Rango de tiempo para la entrada de un auto
	MinEntryTime = 100 * time.Millisecond
	MaxEntryTime = 200 * time.Millisecond

	// Rango de tiempo para la salida de un auto
	MinExitTime = 100 * time.Millisecond
	MaxExitTime = 200 * time.Millisecond

	// Tiempo de espera entre intentos de estacionamiento
	RetryParkingDelay = 100 * time.Millisecond

	// Tiempo de espera para verificar cambios de dirección
	DirectionCheckDelay = 100 * time.Millisecond

	// Intervalo entre llegadas de nuevos autos
	MinCarArrivalInterval = 100 * time.Millisecond
	MaxCarArrivalInterval = 500 * time.Millisecond
)

// Parámetros de Poisson
const (
	// Lambda representa la tasa media de llegadas por segundo
	DefaultLambda = 0.167 // aproximadamente 5 autos por minuto
)

// Configuración de la interfaz gráfica
const (
	DefaultWindowWidth  = 600
	DefaultWindowHeight = 300

	ParkingGridColumns = 5
)

const (
	EmptySpotSymbol = "🅿️"
	CarSymbol       = "🚗"
)

// Estados y mensajes del sistema
const (
	WindowTitle = "goPark - Simulador de Estacionamiento"

	InitialStatusMessage      = "Estado: Iniciando..."
	SimulationCompleteMessage = "¡SIMULACIÓN COMPLETADA!\nTodos los autos (%d) han sido procesados.\nLa ventana se cerrará en %d segundos."
	StatusTemplate            = "Estado: %d/%d lugares ocupados | Procesados: %d/%d"

	StartButtonText  = "Iniciar Simulación"
	StopButtonText   = "Detener"
	ResumeButtonText = "Reanudar"
)

const (
	DirectionNoneText = "Libre"
	DirectionInText   = "Entrada"
	DirectionOutText  = "Salida"
)

var ErrorMessages = struct {
	ParkingFull   string
	InvalidSpotID string
	CarNotFound   string
	SystemStopped string
}{
	ParkingFull:   "El estacionamiento está lleno",
	InvalidSpotID: "ID de espacio de estacionamiento inválido",
	CarNotFound:   "Auto no encontrado",
	SystemStopped: "El sistema está detenido",
}

const (
	// Tamaño del buffer para el canal de semáforos
	SpotSemaphoreBuffer = TotalParkingSpots
)

func GetDirectionText(direction string) string {
	switch direction {
	case "in":
		return DirectionInText
	case "out":
		return DirectionOutText
	default:
		return DirectionNoneText
	}
}

func GetStatusMessage(occupied, total, waiting int, direction string, processed, target int) string {
	return fmt.Sprintf(
		StatusTemplate,
		occupied,
		total,
		processed, 
		target,
	)
}

func GetCompletionMessage() string {
	return fmt.Sprintf(
		SimulationCompleteMessage,
		TotalCarsToProcess,
		int(WindowCloseDelay.Seconds()),
	)
}
