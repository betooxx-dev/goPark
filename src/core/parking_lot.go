package core

import (
	"sync"
	"time"

	"simulador/src/models"
)

// Representa la dirección del flujo de autos
type Direction string

const (
	DirectionNone Direction = "none"
	DirectionIn   Direction = "in"
	DirectionOut  Direction = "out"
)

// Contiene información sobre el estado actual del estacionamiento
type ParkingLotStatus struct {
	OccupiedSpots int
	CarsWaiting   int
	Direction     Direction
	TotalCars     int
	IsCompleted   bool
}

// Gestiona la lógica del estacionamiento
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

	// Callbacks para notificación de eventos
	onStatusChanged func(ParkingLotStatus)
	onSpotChanged   func(int, *models.ParkingSpot)
}

// Crea una nueva instancia de ParkingLot
func NewParkingLot(onStatusChanged func(ParkingLotStatus), onSpotChanged func(int, *models.ParkingSpot)) *ParkingLot {
	pl := &ParkingLot{
		spotsSem:        make(chan struct{}, 20),
		cars:            make(map[int]*models.Car),
		direction:       DirectionNone,
		isRunning:       false,
		onStatusChanged: onStatusChanged,
		onSpotChanged:   onSpotChanged,
	}

	// Inicializar espacios de estacionamiento
	for i := 0; i < 20; i++ {
		pl.spots[i] = models.NewParkingSpot(i)
		pl.spotsSem <- struct{}{}
	}

	return pl
}

// Retorna el estado actual del estacionamiento
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
		IsCompleted:   pl.totalCars >= 100 && pl.activeCount == 0,
	}
}

// Busca un espacio disponible
func (pl *ParkingLot) findAvailableSpot() int {
	for i, spot := range pl.spots {
		if !spot.IsOccupied() {
			return i
		}
	}
	return -1
}

// Intenta estacionar un auto en un espacio disponible
func (pl *ParkingLot) parkCar(carID int) (int, bool) {
	if !pl.isRunning {
		return -1, false
	}

	select {
	case <-pl.spotsSem:
		pl.mutex.Lock()
		spotID := pl.findAvailableSpot()
		if spotID != -1 {
			pl.spots[spotID].Occupy(carID)
			car := models.NewCar(carID)
			parkingTime := time.Duration(2+time.Now().UnixNano()%3) * time.Second
			car.Park(spotID, parkingTime)
			pl.cars[carID] = car
			pl.totalCars++
			pl.mutex.Unlock()

			if pl.onSpotChanged != nil {
				pl.onSpotChanged(spotID, pl.spots[spotID])
			}
			if pl.onStatusChanged != nil {
				pl.onStatusChanged(pl.GetStatus())
			}

			return spotID, true
		}
		pl.mutex.Unlock()
		pl.spotsSem <- struct{}{}
		return -1, false
	default:
		return -1, false
	}
}

// Libera un espacio de estacionamiento
func (pl *ParkingLot) leaveParkingSpot(spotID int, carID int) {
	pl.mutex.Lock()
	pl.spots[spotID].Release()
	delete(pl.cars, carID)
	pl.mutex.Unlock()
	pl.spotsSem <- struct{}{}

	if pl.onSpotChanged != nil {
		pl.onSpotChanged(spotID, pl.spots[spotID])
	}
	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}
}

// Gestiona la entrada de un auto al estacionamiento
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

	spotID, success := pl.parkCar(carID)
	if success {
		time.Sleep(time.Duration(500+time.Now().UnixNano()%1500) * time.Millisecond)

		go func() {
			car := pl.cars[carID]
			time.Sleep(car.ParkingTime)
			pl.ExitParking(carID, spotID)
		}()
	}

	pl.mutex.Lock()
	pl.carsWaiting--
	if pl.carsWaiting == 0 {
		pl.direction = DirectionNone
	}
	pl.mutex.Unlock()

	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	return success
}

// Gestiona la salida de un auto del estacionamiento
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

	pl.leaveParkingSpot(spotID, carID)
	time.Sleep(time.Duration(500+time.Now().UnixNano()%500) * time.Millisecond)

	pl.mutex.Lock()
	pl.direction = DirectionNone
	pl.activeCount--
	pl.mutex.Unlock()

	if pl.onStatusChanged != nil {
		pl.onStatusChanged(pl.GetStatus())
	}

	pl.wg.Done()
}

// Inicia la simulación
func (pl *ParkingLot) Start() {
	pl.mutex.Lock()
	pl.isRunning = true
	pl.mutex.Unlock()
}

// Detiene la simulación
func (pl *ParkingLot) Stop() {
	pl.mutex.Lock()
	pl.isRunning = false
	pl.mutex.Unlock()
}

// Retorna si la simulación está en ejecución
func (pl *ParkingLot) IsRunning() bool {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	return pl.isRunning
}

// Retorna el WaitGroup para sincronización
func (pl *ParkingLot) WaitGroup() *sync.WaitGroup {
	return &pl.wg
}
