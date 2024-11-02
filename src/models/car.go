package models

import "time"

// Representa un automóvil en el sistema de estacionamiento
type Car struct {
	ID          int           // Identificador único del auto
	EntryTime   time.Time     // Hora de entrada al estacionamiento
	ParkingSpot int           // ID del espacio de estacionamiento asignado (-1 si no está estacionado)
	IsParked    bool          // Indica si el auto está estacionado
	IsWaiting   bool          // Indica si el auto está esperando para entrar
	ParkingTime time.Duration // Tiempo que el auto permanecerá estacionado
}

// Crea una nueva instancia de Car
func NewCar(id int) *Car {
	return &Car{
		ID:          id,
		ParkingSpot: -1,
		IsParked:    false,
		IsWaiting:   false,
	}
}

// Asigna un espacio de estacionamiento al auto
func (c *Car) Park(spotID int, parkingTime time.Duration) {
	c.ParkingSpot = spotID
	c.IsParked = true
	c.IsWaiting = false
	c.EntryTime = time.Now()
	c.ParkingTime = parkingTime
}

// Marca el auto como no estacionado
func (c *Car) Leave() {
	c.ParkingSpot = -1
	c.IsParked = false
	c.EntryTime = time.Time{}
	c.ParkingTime = 0
}

// Marca el auto como en espera
func (c *Car) SetWaiting(waiting bool) {
	c.IsWaiting = waiting
}

// Retorna el tiempo que el auto ha estado estacionado
func (c *Car) GetTimeParked() time.Duration {
	if !c.IsParked {
		return 0
	}
	return time.Since(c.EntryTime)
}

// Verifica si el auto debe salir basado en su tiempo de estacionamiento
func (c *Car) ShouldLeave() bool {
	return c.IsParked && c.GetTimeParked() >= c.ParkingTime
}
