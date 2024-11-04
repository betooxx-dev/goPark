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

// Crea una nueva instancia de la interfaz gráfica
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

// Configura todos los elementos de la interfaz
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

// Configura los botones y sus callbacks
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

// Actualiza el estado de habilitación de los botones
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

// Actualiza la etiqueta de estado
func (g *GUI) UpdateStatus(status string) {
	g.statusLabel.SetText(status)
}

// Actualiza la visualización de un espacio de estacionamiento
func (g *GUI) UpdateSpot(spotID int, spot *models.ParkingSpot) {
	if spot.IsOccupied() {
		g.parkingSpots[spotID].SetText(config.CarSymbol + "\n" + fmt.Sprintf("%d", spot.GetCarID()))
	} else {
		g.parkingSpots[spotID].SetText(config.EmptySpotSymbol)
	}
}

// Deshabilita todos los botones
func (g *GUI) DisableAllButtons() {
	g.updateButtonStates(false, false, false)
}

// Cierra la ventana de la aplicación
func (g *GUI) Close() {
	g.window.Close()
}

// Inicia la ejecución de la interfaz gráfica
func (g *GUI) Run() {
	g.window.ShowAndRun()
}
