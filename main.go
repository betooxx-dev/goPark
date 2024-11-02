package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

type ParkingSpot struct {
    id       int
    occupied bool
    carID    int
}

type ParkingLot struct {
    spots        [20]ParkingSpot
    mutex        sync.Mutex
    entrance     sync.Mutex
    spotsSem     chan struct{}
    updateUI     func()
    carsWaiting  int
    direction    string
    totalCars    int
    statusLabel  *widget.Label
    isRunning    bool
    activeCount  int // Contador de coches actualmente en proceso
    wg           sync.WaitGroup
}

func NewParkingLot(updateUI func(), statusLabel *widget.Label) *ParkingLot {
    pl := &ParkingLot{
        spotsSem:    make(chan struct{}, 20),
        updateUI:    updateUI,
        direction:   "none",
        statusLabel: statusLabel,
        isRunning:   false,
    }
    
    for i := 0; i < 20; i++ {
        pl.spots[i] = ParkingSpot{id: i, occupied: false}
        pl.spotsSem <- struct{}{}
    }
    
    return pl
}

func (pl *ParkingLot) updateStatus() {
    occupiedSpots := 0
    for _, spot := range pl.spots {
        if spot.occupied {
            occupiedSpots++
        }
    }
    
    status := fmt.Sprintf("Estado: %d/20 lugares ocupados | %d autos esperando | Dirección: %s | Procesados: %d/100",
        occupiedSpots,
        pl.carsWaiting,
        pl.getDirectionText(),
        pl.totalCars)

    if pl.totalCars >= 100 && pl.activeCount == 0 {
        status += "\n¡Simulación Completada! Todos los autos han sido procesados."
    }
    
    pl.statusLabel.SetText(status)
}

func (pl *ParkingLot) getDirectionText() string {
    switch pl.direction {
    case "in":
        return "Entrada"
    case "out":
        return "Salida"
    default:
        return "Libre"
    }
}

func (pl *ParkingLot) findAvailableSpot() int {
    for i, spot := range pl.spots {
        if !spot.occupied {
            return i
        }
    }
    return -1
}

func (pl *ParkingLot) parkCar(carID int) (int, bool) {
    if !pl.isRunning {
        return -1, false
    }
    
    select {
    case <-pl.spotsSem:
        pl.mutex.Lock()
        spotID := pl.findAvailableSpot()
        if spotID != -1 {
            pl.spots[spotID].occupied = true
            pl.spots[spotID].carID = carID
            pl.totalCars++
            pl.mutex.Unlock()
            pl.updateStatus()
            return spotID, true
        }
        pl.mutex.Unlock()
        pl.spotsSem <- struct{}{}
        return -1, false
    default:
        return -1, false
    }
}

func (pl *ParkingLot) leaveParkingSpot(spotID int) {
    pl.mutex.Lock()
    pl.spots[spotID].occupied = false
    pl.spots[spotID].carID = 0
    pl.mutex.Unlock()
    pl.spotsSem <- struct{}{}
    pl.updateStatus()
}

func (pl *ParkingLot) enterParking(carID int) bool {
    if !pl.isRunning {
        return false
    }

    pl.mutex.Lock()
    if pl.direction == "out" {
        pl.mutex.Unlock()
        return false
    }
    pl.direction = "in"
    pl.carsWaiting++
    pl.activeCount++
    pl.updateStatus()
    pl.mutex.Unlock()
    
    spotID, success := pl.parkCar(carID)
    if success {
        time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
        pl.updateUI()
        
        go func() {
            parkingTime := time.Duration(rand.Intn(2000)+3000) * time.Millisecond
            time.Sleep(parkingTime)
            pl.exitParking(carID, spotID)
        }()
    }
    
    pl.mutex.Lock()
    pl.carsWaiting--
    if pl.carsWaiting == 0 {
        pl.direction = "none"
    }
    pl.updateStatus()
    pl.mutex.Unlock()
    
    return success
}

func (pl *ParkingLot) exitParking(carID int, spotID int) {
    pl.mutex.Lock()
    for pl.direction == "in" {
        pl.mutex.Unlock()
        time.Sleep(100 * time.Millisecond)
        pl.mutex.Lock()
    }
    pl.direction = "out"
    pl.updateStatus()
    pl.mutex.Unlock()
    
    pl.leaveParkingSpot(spotID)
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
    
    pl.mutex.Lock()
    pl.direction = "none"
    pl.activeCount--
    pl.updateStatus()
    pl.mutex.Unlock()
    
    pl.updateUI()
    pl.wg.Done()
}

func main() {
    myApp := app.New()
    window := myApp.NewWindow("Simulador de Estacionamiento")
    
    parkingContainer := container.NewGridWithColumns(5)
    statusLabel := widget.NewLabel("Estado: Iniciando...")
    
    spotLabels := make([]*widget.Label, 20)
    for i := 0; i < 20; i++ {
        spotLabels[i] = widget.NewLabel("🅿️")
        parkingContainer.Add(spotLabels[i])
    }
    
    var parkingLot *ParkingLot
    
    updateUI := func() {
        for i, spot := range parkingLot.spots {
            if spot.occupied {
                spotLabels[i].SetText(fmt.Sprintf("🚗\n%d", spot.carID))
            } else {
                spotLabels[i].SetText("🅿️")
            }
        }
    }
    
    parkingLot = NewParkingLot(updateUI, statusLabel)

    var stopButton *widget.Button
    
    stopButton = widget.NewButton("Detener", func() {
        parkingLot.isRunning = false
        stopButton.Disable()
    })
    stopButton.Disable()

    var startButton *widget.Button

    startButton = widget.NewButton("Iniciar Simulación", func() {
        startButton.Disable()
        stopButton.Enable()
        parkingLot.isRunning = true
        
        go func() {
            for i := 1; i <= 100 && parkingLot.isRunning; i++ {
                parkingLot.wg.Add(1)
                carID := i
                time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
                
                go func(id int) {
                    success := false
                    for !success && parkingLot.isRunning {
                        success = parkingLot.enterParking(id)
                        if !success {
                            time.Sleep(500 * time.Millisecond)
                        }
                    }
                    if !success {
                        parkingLot.wg.Done()
                    }
                }(carID)
            }
            
            // Esperar a que todos los autos terminen
            go func() {
                parkingLot.wg.Wait()
                stopButton.Disable()
                parkingLot.updateStatus()
            }()
        }()
    })
    
    buttonContainer := container.NewHBox(startButton, stopButton)
    
    mainContainer := container.NewVBox(
        statusLabel,
        parkingContainer,
        buttonContainer,
    )
    
    window.SetContent(mainContainer)
    window.Resize(fyne.NewSize(600, 300))
    window.ShowAndRun()
}