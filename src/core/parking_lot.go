package core

import (
	"math/rand"
	"simulador/src/config"
	"simulador/src/models"
	"sync"
	"time"
)

type Direction string

const (
	DirectionNone Direction = "none"
	DirectionIn   Direction = "in"
	DirectionOut  Direction = "out"
)

type ParkingLotStatus struct {
	OccupiedSpots int
	CarsWaiting   int
	Direction     Direction
	TotalCars     int        // Autos que han entrado al estacionamiento
	ProcessedCars int        // Total de autos procesados (incluyendo rechazados)
	IsCompleted   bool
}

type ParkingLot struct {
	spots         [20]*models.ParkingSpot
	cars          map[int]*models.Car
	mutex         sync.RWMutex
	spotsSem      chan struct{}
	carsWaiting   int
	direction     Direction
	totalCars     int        // Autos que han entrado al estacionamiento
	processedCars int        // Total de autos procesados (incluyendo rechazados)
	isRunning     bool
	activeCount   int
	wg            sync.WaitGroup

	onStatusChanged func(ParkingLotStatus)
	onSpotChanged   func(int, *models.ParkingSpot)
}

func NewParkingLot(onStatusChanged func(ParkingLotStatus), onSpotChanged func(int, *models.ParkingSpot)) *ParkingLot {
	pl := &ParkingLot{
		spotsSem:      make(chan struct{}, 20),
		cars:          make(map[int]*models.Car),
		direction:     DirectionNone,
		isRunning:     false,
		onStatusChanged: onStatusChanged,
		onSpotChanged:   onSpotChanged,
	}

	for i := 0; i < 20; i++ {
		pl.spots[i] = models.NewParkingSpot(i)
		pl.spotsSem <- struct{}{}
	}

	return pl
}

func (pl *ParkingLot) findAvailableSpot() int {
	for i, spot := range pl.spots {
		if !spot.IsOccupied() {
			return i
		}
	}
	return -1
}

func (pl *ParkingLot) EnterParking(carID int) bool {
	if !pl.isRunning {
		return false
	}

	// Incrementar contadores bajo el mutex
	pl.mutex.Lock()
	pl.carsWaiting++
	pl.activeCount++
	currentDirection := pl.direction
	if currentDirection == DirectionNone {
		pl.direction = DirectionIn
	}
	pl.mutex.Unlock()

	// Notificar cambio de estado
	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	// Si la dirección es salida, esperar un tiempo y reintentar
	if currentDirection == DirectionOut {
		time.Sleep(config.DirectionCheckDelay)
		pl.mutex.Lock()
		pl.carsWaiting--
		pl.activeCount--
		pl.mutex.Unlock()
		return false
	}

	success := false
	// Intentar obtener un espacio
	select {
	case <-pl.spotsSem:
		pl.mutex.Lock()
		spotID := pl.findAvailableSpot()
		if spotID != -1 {
			// Ocupar el espacio
			pl.spots[spotID].Occupy(carID)
			car := models.NewCar(carID)
			parkingTime := config.MinParkingTime +
				time.Duration(rand.Int63n(int64(config.MaxParkingTime-config.MinParkingTime)))
			car.Park(spotID, parkingTime)
			pl.cars[carID] = car
			pl.totalCars++     // Incrementar cuando el auto entra
			pl.processedCars++ // Incrementar processedCars solo cuando realmente entra
			success = true
			
			// Actualizar contadores
			pl.carsWaiting--
			if pl.carsWaiting == 0 {
				pl.direction = DirectionNone
			}
			pl.mutex.Unlock()

			// Notificar cambios
			if pl.onSpotChanged != nil {
				pl.onSpotChanged(spotID, pl.spots[spotID])
			}
			if pl.onStatusChanged != nil {
				pl.onStatusChanged(pl.GetStatus())
			}

			// Programar la salida del auto
			go func() {
				time.Sleep(parkingTime)
				pl.ExitParking(carID, spotID)
			}()
		} else {
			pl.mutex.Unlock()
			pl.spotsSem <- struct{}{} // Devolver el token si no se pudo usar
		}
		
	case <-time.After(100 * time.Millisecond): // Timeout para evitar bloqueos
		// No se pudo obtener espacio, actualizar contadores
		pl.mutex.Lock()
		pl.carsWaiting--
		pl.activeCount--
		pl.mutex.Unlock()
	}

	// Si el auto no pudo entrar, incrementar processedCars solo si es la última vez que lo intentará
	if !success {
		pl.mutex.Lock()
		// Verificar si este es el último intento del auto
		if pl.processedCars < config.TotalCarsToProcess {
			pl.processedCars++
		}
		pl.mutex.Unlock()

		// Notificar cambio final de estado
		if pl.onStatusChanged != nil {
			pl.onStatusChanged(pl.GetStatus())
		}
	}

	return success
}

func (pl *ParkingLot) ExitParking(carID int, spotID int) {
	pl.mutex.Lock()
	if pl.direction == DirectionIn && pl.carsWaiting > 0 {
		pl.mutex.Unlock()
		time.Sleep(config.DirectionCheckDelay)
		pl.mutex.Lock()
	}
	pl.direction = DirectionOut
	pl.mutex.Unlock()

	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	exitTime := config.MinExitTime +
		time.Duration(rand.Int63n(int64(config.MaxExitTime-config.MinExitTime)))
	time.Sleep(exitTime)

	pl.leaveParkingSpot(spotID, carID)

	pl.mutex.Lock()
	pl.direction = DirectionNone
	pl.activeCount--
	pl.mutex.Unlock()

	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	pl.wg.Done()
}

func (pl *ParkingLot) leaveParkingSpot(spotID int, carID int) {
	pl.mutex.Lock()
	if spot := pl.spots[spotID]; spot != nil && spot.GetCarID() == carID {
		spot.Release()
		delete(pl.cars, carID)
		pl.mutex.Unlock()
		pl.spotsSem <- struct{}{}

		if pl.onSpotChanged != nil {
			pl.onSpotChanged(spotID, pl.spots[spotID])
		}
	} else {
		pl.mutex.Unlock()
	}
}

func (pl *ParkingLot) GetStatus() ParkingLotStatus {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()

	occupiedSpots := 0
	for _, spot := range pl.spots {
		if spot.IsOccupied() {
			occupiedSpots++
		}
	}

	return ParkingLotStatus{
		OccupiedSpots:  occupiedSpots,
		CarsWaiting:    pl.carsWaiting,
		Direction:      pl.direction,
		TotalCars:      pl.totalCars,
		ProcessedCars:  pl.processedCars,
		IsCompleted:    pl.processedCars >= config.TotalCarsToProcess && pl.activeCount == 0,
	}
}

func (pl *ParkingLot) Start() {
	pl.mutex.Lock()
	pl.isRunning = true
	pl.mutex.Unlock()
}

func (pl *ParkingLot) Stop() {
	pl.mutex.Lock()
	pl.isRunning = false
	pl.mutex.Unlock()
}

func (pl *ParkingLot) IsRunning() bool {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()
	return pl.isRunning
}

func (pl *ParkingLot) WaitGroup() *sync.WaitGroup {
	return &pl.wg
}