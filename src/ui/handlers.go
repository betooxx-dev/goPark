package ui

import (
	"simulador/src/config"
	"simulador/src/core"
	"simulador/src/models"
	"sync"
	"time"
)

// Maneja la lógica de eventos de la UI
type Handlers struct {
	parkingLot     *core.ParkingLot
	gui            *GUI
	currentCarID   int
	simulationDone chan bool
	mutex          sync.Mutex
}

// Crea una nueva instancia de Handlers
func NewHandlers() *Handlers {
	h := &Handlers{
		currentCarID:   1,
		simulationDone: make(chan bool),
	}

	// Crear instancia de ParkingLot con callbacks
	h.parkingLot = core.NewParkingLot(
		h.handleStatusChanged,
		h.handleSpotChanged,
	)

	return h
}

// Establece la referencia a la interfaz gráfica
func (h *Handlers) SetGUI(gui *GUI) {
	h.gui = gui
}

// Maneja los cambios de estado del estacionamiento
func (h *Handlers) handleStatusChanged(status core.ParkingLotStatus) {
	if status.IsCompleted {
		h.gui.UpdateStatus(config.GetCompletionMessage())
		h.handleSimulationComplete()
		return
	}

	statusMsg := config.GetStatusMessage(
		status.OccupiedSpots,
		config.TotalParkingSpots,
		status.CarsWaiting,
		string(status.Direction),
		status.TotalCars,
		config.TotalCarsToProcess,
	)
	h.gui.UpdateStatus(statusMsg)
}

// Maneja los cambios en los espacios de estacionamiento
func (h *Handlers) handleSpotChanged(spotID int, spot *models.ParkingSpot) {
	h.gui.UpdateSpot(spotID, spot)
}

// Maneja el evento de inicio de simulación
func (h *Handlers) HandleStart() {
	h.parkingLot.Start()
	go h.runSimulation()
}

// Maneja el evento de detención de simulación
func (h *Handlers) HandleStop() {
	h.parkingLot.Stop()
}

// Maneja el evento de reanudación de simulación
func (h *Handlers) HandleResume() {
	h.parkingLot.Start()
}

// Ejecuta la simulación principal
func (h *Handlers) runSimulation() {
	for h.currentCarID <= config.TotalCarsToProcess {
		if h.parkingLot.IsRunning() {
			h.mutex.Lock()
			carID := h.currentCarID
			h.currentCarID++
			h.mutex.Unlock()

			h.parkingLot.WaitGroup().Add(1)

			// Generar intervalo aleatorio entre llegadas de autos
			arrivalDelay := time.Duration(
				config.MinCarArrivalInterval.Nanoseconds() +
					time.Now().UnixNano()%
						(config.MaxCarArrivalInterval.Nanoseconds()-
							config.MinCarArrivalInterval.Nanoseconds()),
			)
			time.Sleep(arrivalDelay)

			go h.processCar(carID)
		} else {
			time.Sleep(config.DirectionCheckDelay)
		}
	}

	// Esperar a que todos los autos terminen
	go func() {
		h.parkingLot.WaitGroup().Wait()
		h.simulationDone <- true
	}()
}

// Maneja el proceso de un auto individual
func (h *Handlers) processCar(carID int) {
	success := false
	for !success {
		if h.parkingLot.IsRunning() {
			success = h.parkingLot.EnterParking(carID)
			if !success {
				time.Sleep(config.RetryParkingDelay)
			}
		} else {
			time.Sleep(config.DirectionCheckDelay)
		}
	}
}

// Maneja la finalización de la simulación
func (h *Handlers) handleSimulationComplete() {
	h.gui.DisableAllButtons()
	go func() {
		time.Sleep(config.WindowCloseDelay)
		h.gui.Close()
	}()
}
