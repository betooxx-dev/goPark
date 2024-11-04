package ui

import (
	"simulador/src/config"
	"simulador/src/core"
	"simulador/src/models"
	"sync"
	"time"
)

type Handlers struct {
	parkingLot     *core.ParkingLot
	gui            *GUI
	currentCarID   int
	simulationDone chan bool
	mutex          sync.Mutex
	isProcessing   bool
}

func NewHandlers() *Handlers {
	h := &Handlers{
		currentCarID:   1,
		simulationDone: make(chan bool),
		isProcessing:   false,
	}

	h.parkingLot = core.NewParkingLot(
		h.handleStatusChanged,
		h.handleSpotChanged,
	)

	return h
}

func (h *Handlers) SetGUI(gui *GUI) {
	h.gui = gui
}

func (h *Handlers) handleStatusChanged(status core.ParkingLotStatus) {
	if h.gui == nil {
		return
	}

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

func (h *Handlers) handleSpotChanged(spotID int, spot *models.ParkingSpot) {
	if h.gui == nil {
		return
	}
	h.gui.UpdateSpot(spotID, spot)
}

func (h *Handlers) HandleStart() {
	h.mutex.Lock()
	if h.isProcessing {
		h.mutex.Unlock()
		return
	}
	h.isProcessing = true
	h.mutex.Unlock()

	h.parkingLot.Start()
	go h.runSimulation()
}

func (h *Handlers) HandleStop() {
	h.parkingLot.Stop()
}

func (h *Handlers) HandleResume() {
	h.parkingLot.Start()
}

func (h *Handlers) runSimulation() {
	for h.currentCarID <= config.TotalCarsToProcess && h.parkingLot.IsRunning() {
		h.mutex.Lock()
		carID := h.currentCarID
		h.currentCarID++
		h.mutex.Unlock()

		h.parkingLot.WaitGroup().Add(1)

		arrivalDelay := time.Duration(
			config.MinCarArrivalInterval.Nanoseconds() +
				time.Now().UnixNano()%
					(config.MaxCarArrivalInterval.Nanoseconds()-
						config.MinCarArrivalInterval.Nanoseconds()),
		)
		time.Sleep(arrivalDelay)

		go h.processCar(carID)
	}

	go func() {
		h.parkingLot.WaitGroup().Wait()
		h.simulationDone <- true

		h.mutex.Lock()
		h.isProcessing = false
		h.mutex.Unlock()
	}()
}

func (h *Handlers) processCar(carID int) {
	maxRetries := 50
	retryCount := 0

	for retryCount < maxRetries {
		if !h.parkingLot.IsRunning() {
			time.Sleep(config.DirectionCheckDelay)
			continue
		}

		if h.parkingLot.EnterParking(carID) {
			return
		}

		retryCount++
		time.Sleep(config.RetryParkingDelay)
	}

	h.parkingLot.WaitGroup().Done()
}

func (h *Handlers) handleSimulationComplete() {
	if h.gui == nil {
		return
	}

	h.gui.DisableAllButtons()
	go func() {
		time.Sleep(config.WindowCloseDelay)
		h.gui.Close()
	}()
}
