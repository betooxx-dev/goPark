package models

import "time"

type Car struct {
	ID          int
	EntryTime   time.Time
	ParkingSpot int
	IsParked    bool
	IsWaiting   bool
	ParkingTime time.Duration
}

func NewCar(id int) *Car {
	return &Car{
		ID:          id,
		ParkingSpot: -1,
		IsParked:    false,
		IsWaiting:   false,
	}
}

func (c *Car) Park(spotID int, parkingTime time.Duration) {
	c.ParkingSpot = spotID
	c.IsParked = true
	c.IsWaiting = false
	c.EntryTime = time.Now()
	c.ParkingTime = parkingTime
}

func (c *Car) Leave() {
	c.ParkingSpot = -1
	c.IsParked = false
	c.EntryTime = time.Time{}
	c.ParkingTime = 0
}

func (c *Car) SetWaiting(waiting bool) {
	c.IsWaiting = waiting
}

func (c *Car) GetTimeParked() time.Duration {
	if !c.IsParked {
		return 0
	}
	return time.Since(c.EntryTime)
}

func (c *Car) ShouldLeave() bool {
	return c.IsParked && c.GetTimeParked() >= c.ParkingTime
}
