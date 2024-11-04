package models

// Representa un espacio individual de estacionamiento
type ParkingSpot struct {
	ID       int  
	Occupied bool 
	CarID    int  // ID del auto que ocupa el espacio (0 si está vacío)
}

func NewParkingSpot(id int) *ParkingSpot {
	return &ParkingSpot{
		ID:       id,
		Occupied: false,
		CarID:    0,
	}
}

func (ps *ParkingSpot) Occupy(carID int) {
	ps.Occupied = true
	ps.CarID = carID
}

func (ps *ParkingSpot) Release() {
	ps.Occupied = false
	ps.CarID = 0
}

func (ps *ParkingSpot) IsOccupied() bool {
	return ps.Occupied
}

func (ps *ParkingSpot) GetCarID() int {
	return ps.CarID
}
