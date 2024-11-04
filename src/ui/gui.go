package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"simulador/src/config"
	"simulador/src/models"
)

// Encapsula todos los elementos de la interfaz gráfica
type GUI struct {
	window       fyne.Window
	app          fyne.App
	parkingSpots []*widget.Label
	statusLabel  *widget.Label
	startButton  *widget.Button
	stopButton   *widget.Button
	resumeButton *widget.Button
	handlers     *Handlers
}

func NewGUI(handlers *Handlers) *GUI {
	gui := &GUI{
		app:          app.New(),
		handlers:     handlers,
		parkingSpots: make([]*widget.Label, config.TotalParkingSpots),
	}

	gui.window = gui.app.NewWindow(config.WindowTitle)
	gui.setupUI()

	return gui
}

func (g *GUI) setupUI() {
	// Crear el contenedor de espacios de estacionamiento
	parkingContainer := container.NewGridWithColumns(config.ParkingGridColumns)

	// Inicializar etiquetas de espacios
	for i := 0; i < config.TotalParkingSpots; i++ {
		g.parkingSpots[i] = widget.NewLabel(config.EmptySpotSymbol)
		parkingContainer.Add(g.parkingSpots[i])
	}

	// Crear etiqueta de estado
	g.statusLabel = widget.NewLabel(config.InitialStatusMessage)

	// Configurar botones
	g.setupButtons()

	// Crear contenedor de botones
	buttonContainer := container.NewHBox(
		g.startButton,
	)

	// Crear contenedor principal
	mainContainer := container.NewVBox(
		g.statusLabel,
		parkingContainer,
		buttonContainer,
	)

	// Configurar la ventana
	g.window.SetContent(mainContainer)
	g.window.Resize(fyne.NewSize(
		config.DefaultWindowWidth,
		config.DefaultWindowHeight,
	))
}

func (g *GUI) setupButtons() {
	g.startButton = widget.NewButton(config.StartButtonText, func() {
		g.handlers.HandleStart()
		g.updateButtonStates(false, true, false)
	})

	g.stopButton = widget.NewButton(config.StopButtonText, func() {
		g.handlers.HandleStop()
		g.updateButtonStates(false, false, true)
	})

	g.resumeButton = widget.NewButton(config.ResumeButtonText, func() {
		g.handlers.HandleResume()
		g.updateButtonStates(false, true, false)
	})

	// Estado inicial de los botones
	g.updateButtonStates(true, false, false)
}

func (g *GUI) updateButtonStates(startEnabled, stopEnabled, resumeEnabled bool) {
	g.startButton.Enable()
	if !startEnabled {
		g.startButton.Disable()
	}

	g.stopButton.Enable()
	if !stopEnabled {
		g.stopButton.Disable()
	}

	g.resumeButton.Enable()
	if !resumeEnabled {
		g.resumeButton.Disable()
	}
}

func (g *GUI) UpdateStatus(status string) {
	g.statusLabel.SetText(status)
}

func (g *GUI) UpdateSpot(spotID int, spot *models.ParkingSpot) {
	if spot.IsOccupied() {
		g.parkingSpots[spotID].SetText(config.CarSymbol + "\n" + fmt.Sprintf("%d", spot.GetCarID()))
	} else {
		g.parkingSpots[spotID].SetText(config.EmptySpotSymbol)
	}
}

func (g *GUI) DisableAllButtons() {
	g.updateButtonStates(false, false, false)
}

func (g *GUI) Close() {
	g.window.Close()
}

func (g *GUI) Run() {
	g.window.ShowAndRun()
}
