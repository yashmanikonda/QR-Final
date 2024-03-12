# QR-Final

QR-Final is a Go-based desktop application that generates QR codes for device identification. It uses the Go-QRCode library for QR code generation and the Go Qt library for creating the desktop user interface.

## Features

- Generates QR codes based on device ID stored in the Windows registry.
- Displays the QR code in a frameless Qt window with a countdown timer.
- Cleans up resources, including deleting the QR code file, on program exit or abrupt termination.
- Creates a Debug log to track program actions.

## Installation
**1. Clone Repo and Navigate to project directory**

**2. Initialize the Go module and download dependencies:**
go mod tidy

**3. Run the application:**
go run main.go
This command will generate a QR code based on the device ID and display it in a frameless window.

**4. Optionally, build the executable:**
go build main.go

**Configuration**
Modify the regKeyPath and regValueName constants in main.go to match your specific registry key and value for the device ID.
