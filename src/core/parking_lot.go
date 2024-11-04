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
	TotalCars     int
	IsCompleted   bool
}

type ParkingLot struct {
	spots       [20]*models.ParkingSpot
	cars        map[int]*models.Car
	mutex       sync.Mutex
	entrance    sync.Mutex
	spotsSem    chan struct{}
	carsWaiting int
	direction   Direction
	totalCars   int
	isRunning   bool
	activeCount int
	wg          sync.WaitGroup

	onStatusChanged func(ParkingLotStatus)
	onSpotChanged   func(int, *models.ParkingSpot)
}

func NewParkingLot(onStatusChanged func(ParkingLotStatus), onSpotChanged func(int, *models.ParkingSpot)) *ParkingLot {
	pl := &ParkingLot{
		spotsSem:        make(chan struct{}, 20),
		cars:            make(map[int]*models.Car),
		direction:       DirectionNone,
		isRunning:       false,
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

	pl.mutex.Lock()
	if pl.direction == DirectionOut {
		pl.mutex.Unlock()
		return false
	}
	pl.direction = DirectionIn
	pl.carsWaiting++
	pl.activeCount++
	pl.mutex.Unlock()

	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	select {
	case <-pl.spotsSem:
		pl.mutex.Lock()
		spotID := pl.findAvailableSpot()
		if spotID != -1 {
			pl.spots[spotID].Occupy(carID)
			car := models.NewCar(carID)
			parkingTime := config.MinParkingTime +
				time.Duration(rand.Int63n(int64(config.MaxParkingTime-config.MinParkingTime)))
			car.Park(spotID, parkingTime)
			pl.cars[carID] = car
			pl.totalCars++
			pl.mutex.Unlock()

			if pl.onSpotChanged != nil {
				pl.onSpotChanged(spotID, pl.spots[spotID])
			}

			// Programar la salida del auto después del tiempo de estacionamiento
			go func() {
				time.Sleep(parkingTime)
				pl.ExitParking(carID, spotID)
			}()

			pl.mutex.Lock()
			pl.carsWaiting--
			if pl.carsWaiting == 0 {
				pl.direction = DirectionNone
			}
			pl.mutex.Unlock()

			if pl.onStatusChanged != nil {
				pl.onStatusChanged(pl.GetStatus())
			}

			return true
		}
		pl.mutex.Unlock()
		pl.spotsSem <- struct{}{}
	default:
	}

	pl.mutex.Lock()
	pl.carsWaiting--
	pl.activeCount--
	pl.mutex.Unlock()

	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	return false
}

func (pl *ParkingLot) ExitParking(carID int, spotID int) {
	pl.mutex.Lock()
	for pl.direction == DirectionIn {
		pl.mutex.Unlock()
		time.Sleep(100 * time.Millisecond)
		pl.mutex.Lock()
	}
	pl.direction = DirectionOut
	pl.mutex.Unlock()

	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	// Liberar el espacio
	pl.leaveParkingSpot(spotID, carID)
	time.Sleep(time.Duration(500+rand.Int63n(500)) * time.Millisecond)

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
	pl.spots[spotID].Release()
	delete(pl.cars, carID)
	pl.mutex.Unlock()
	pl.spotsSem <- struct{}{}

	if pl.onSpotChanged != nil {
		pl.onSpotChanged(spotID, pl.spots[spotID])
	}
}

func (pl *ParkingLot) GetStatus() ParkingLotStatus {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	occupiedSpots := 0
	for _, spot := range pl.spots {
		if spot.IsOccupied() {
			occupiedSpots++
		}
	}

	return ParkingLotStatus{
		OccupiedSpots: occupiedSpots,
		CarsWaiting:   pl.carsWaiting,
		Direction:     pl.direction,
		TotalCars:     pl.totalCars,
		IsCompleted:   pl.totalCars >= config.TotalCarsToProcess && pl.activeCount == 0,
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
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	return pl.isRunning
}

func (pl *ParkingLot) WaitGroup() *sync.WaitGroup {
	return &pl.wg
}
