package models

// Representa un espacio individual de estacionamiento
type ParkingSpot struct {
	ID       int  // Identificador único del espacio
	Occupied bool // Indica si el espacio está ocupado
	CarID    int  // ID del auto que ocupa el espacio (0 si está vacío)
}

// Crea una nueva instancia de ParkingSpot
func NewParkingSpot(id int) *ParkingSpot {
	return &ParkingSpot{
		ID:       id,
		Occupied: false,
		CarID:    0,
	}
}

// Marca el espacio como ocupado por un auto específico
func (ps *ParkingSpot) Occupy(carID int) {
	ps.Occupied = true
	ps.CarID = carID
}

// Libera el espacio de estacionamiento
func (ps *ParkingSpot) Release() {
	ps.Occupied = false
	ps.CarID = 0
}

// Retorna si el espacio está ocupado
func (ps *ParkingSpot) IsOccupied() bool {
	return ps.Occupied
}

// Retorna el ID del auto que ocupa el espacio
func (ps *ParkingSpot) GetCarID() int {
	return ps.CarID
}
