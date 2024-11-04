package config

import (
	"fmt"
	"time"
)

// Configuración del estacionamiento
const (
	// Capacidad total del estacionamiento
	TotalParkingSpots = 20

	// Número total de autos a procesar en la simulación
	TotalCarsToProcess = 100

	// Tiempo de espera para cerrar la ventana al finalizar
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

// Configuración de la interfaz gráfica
const (
	// Tamaño inicial de la ventana
	DefaultWindowWidth  = 600
	DefaultWindowHeight = 300

	// Número de columnas en la cuadrícula de estacionamiento
	ParkingGridColumns = 5
)

// Símbolos para la interfaz gráfica
const (
	EmptySpotSymbol = "🅿️"
	CarSymbol       = "🚗"
)

// Estados y mensajes del sistema
const (
	// Título de la ventana
	WindowTitle = "goPark - Simulador de Estacionamiento"

	// Mensajes de estado
	InitialStatusMessage      = "Estado: Iniciando..."
	SimulationCompleteMessage = "¡SIMULACIÓN COMPLETADA!\nTodos los autos (%d) han sido procesados.\nLa ventana se cerrará en %d segundos."
	StatusTemplate            = "Estado: %d/%d lugares ocupados | %d autos esperando | Dirección: %s | Procesados: %d/%d"

	// Textos de los botones
	StartButtonText  = "Iniciar Simulación"
	StopButtonText   = "Detener"
	ResumeButtonText = "Reanudar"
)

// Direcciones de flujo
const (
	DirectionNoneText = "Libre"
	DirectionInText   = "Entrada"
	DirectionOutText  = "Salida"
)

// Contiene mensajes de error comunes
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

// Configuración de canales y buffers
const (
	// Tamaño del buffer para el canal de semáforos
	SpotSemaphoreBuffer = TotalParkingSpots
)

// Convierte una dirección en su texto correspondiente
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

// Formatea el mensaje de estado
func GetStatusMessage(occupied, total, waiting int, direction string, processed, target int) string {
	return fmt.Sprintf(
		StatusTemplate,
		occupied,
		total,
		waiting,
		GetDirectionText(direction),
		processed,
		target,
	)
}

// Formatea el mensaje de finalización
func GetCompletionMessage() string {
	return fmt.Sprintf(
		SimulationCompleteMessage,
		TotalCarsToProcess,
		int(WindowCloseDelay.Seconds()),
	)
}
