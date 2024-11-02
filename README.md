# goPark 🚗

Un simulador de estacionamiento concurrente desarrollado en Go utilizando Fyne para la interfaz gráfica. goPark simula la gestión de un estacionamiento con múltiples autos entrando y saliendo simultáneamente, demostrando conceptos de programación concurrente y sincronización.

![image](https://github.com/user-attachments/assets/65bb3f1e-86d9-4602-80fa-ce4aef3a0907)

## 🚀 Características

- Simulación en tiempo real de 100 autos
- Gestión concurrente de 20 espacios de estacionamiento
- Interfaz gráfica intuitiva con Fyne
- Control de dirección de flujo (entrada/salida)
- Sistema de semáforos para control de acceso
- Visualización en tiempo real del estado del estacionamiento

## 📋 Requisitos Previos

- Go 1.21 o superior
- Fyne y sus dependencias

Para sistemas basados en Linux, necesitarás:
```bash
sudo apt-get install gcc libgl1-mesa-dev xorg-dev
```

## ⚡ Instalación

```bash
# Clonar el repositorio
git clone https://github.com/betooxx-dev/goPark.git

# Navegar al directorio
cd goPark

# Instalar dependencias
go mod tidy

# Ejecutar la aplicación
go run main.go
```

## 🏗️ Arquitectura

El proyecto sigue una arquitectura en capas:

```
goPark/
├── src/
│   ├── config/                      
│   │   └── constants.go
│   ├── core/
│   │   └── parking_lot.go
│   ├── models/                
│   │   ├── car.go
│   │   └── parking_spot.go
│   └── ui/
│       ├── gui.go
│       └── handlers.go
├── go.mod
├── go.sum
├── main.go                     
└── README.md
```

## 🎮 Uso

1. Ejecuta la aplicación
2. Presiona "Iniciar Simulación"
3. Observa cómo los autos entran y salen del estacionamiento
4. Usa los botones "Detener" y "Reanudar" para controlar la simulación
5. La simulación termina después de procesar 100 autos

## 🛠️ Tecnologías Utilizadas

- [Go](https://golang.org/) - Lenguaje de programación
- [Fyne](https://fyne.io/) - Framework GUI
- Goroutines y Channels para concurrencia
- Mutex para sincronización
- Semáforos para control de acceso
