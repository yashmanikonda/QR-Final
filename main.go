package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"golang.org/x/sys/windows/registry"
)

const (
	regKeyPath     = `Software\BioSecQRID`
	regValueName   = "PlainText"
	qrCodeFile     = "output_qr_code.png"
	logFileName    = "QR_debug.log"
	sleepDuration  = 30 * time.Second
	updateInterval = 1 * time.Second
)

var exitChan = make(chan os.Signal, 1)

func main() {
	// Initialize the log file
	logFile, err := os.Create(logFileName)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Set up signal handling for cleanup on program exit
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the application
	runApplication()
}

func runApplication() {
	// Check if the QR code file already exists, delete it if it does
	if fileExists(qrCodeFile) {
		err := deleteQRCodeFile(qrCodeFile)
		if err != nil {
			log.Fatal("Error deleting existing QR code file:", err)
		}
	}

	// Set the countdown duration
	duration := sleepDuration

	// Read the device ID from the registry
	deviceID, err := GetDeviceIDFromRegistry()
	if err != nil {
		log.Fatal("Error reading device ID from registry:", err)
	}

	// Generate the QR code
	err = GenerateQRCode(deviceID, qrCodeFile)
	if err != nil {
		log.Fatal("Error generating QR code:", err)
	}
	fmt.Printf("QR code for Device ID '%s' has been generated and saved as '%s'.\n", deviceID, qrCodeFile)
	log.Println("QR code generated.")

	// Create and run the Qt application
	widgetApp := widgets.NewQApplication(len(os.Args), os.Args)

	// Get the screen geometry to calculate the center (primary screen)
	screenGeometry := widgets.QApplication_Desktop().ScreenGeometry(-1)

	// Create a main window
	mainWindow := widgets.NewQMainWindow(nil, 0)
	mainWindow.SetWindowTitle("QR Code Viewer")

	// Remove the title bar
	mainWindow.SetWindowFlags(core.Qt__FramelessWindowHint | core.Qt__WindowStaysOnTopHint)

	// Calculate the position to center the window
	windowWidth := 400  // Set your desired window width
	windowHeight := 400 // Set your desired window height
	screenCenter := screenGeometry.Center()
	windowRect := core.NewQRect4(
		screenCenter.X()-windowWidth/2,
		screenCenter.Y()-windowHeight/2,
		windowWidth,
		windowHeight,
	)
	mainWindow.SetGeometry(windowRect)

	// Create a QLabel to display the image
	imageLabel := widgets.NewQLabel(nil, 0)
	pixmap := gui.NewQPixmap()
	pixmap.Load(qrCodeFile, "", core.Qt__AutoColor)
	pixmap = pixmap.Scaled2(windowWidth, windowHeight, core.Qt__KeepAspectRatio, core.Qt__SmoothTransformation)
	imageLabel.SetPixmap(pixmap)

	// Create a QLabel for the countdown timer
	timerLabel := widgets.NewQLabel(nil, 0)
	timerLabel.SetAlignment(core.Qt__AlignCenter)

	// Set up the main window layout
	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(imageLabel, 0, core.Qt__AlignCenter)
	layout.AddWidget(timerLabel, 0, core.Qt__AlignCenter) // Add the timer label below the image
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	mainWindow.SetCentralWidget(widget)

	// Show the main window
	mainWindow.Show()

	// Start the countdown timer
	go startCountdownTimer(timerLabel, duration, mainWindow)

	// Run the application event loop
	widgetApp.Exec()

	// Cleanup when the application exits
	cleanup()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func deleteQRCodeFile(filename string) error {
	// Delete the QR code file
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	log.Println("QR code file deleted.")
	return nil
}

func cleanup() {
	deleteQRCodeFile(qrCodeFile)
	log.Println("Application cleanup complete.")
}

func GetDeviceIDFromRegistry() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, regKeyPath, registry.READ)
	if err != nil {
		return "", err
	}
	defer k.Close()

	deviceID, _, err := k.GetStringValue(regValueName)
	if err != nil {
		return "", err
	}

	return deviceID, nil
}

func GenerateQRCode(deviceID, filename string) error {
	return qrcode.WriteFile(deviceID, qrcode.Low, 256, filename)
}

func startCountdownTimer(label *widgets.QLabel, duration time.Duration, mainWindow *widgets.QMainWindow) {
	for remainingTime := duration; remainingTime > 0; remainingTime -= updateInterval {
		time.Sleep(updateInterval)
		updateTimerLabel(label, remainingTime)
	}

	updateTimerLabel(label, 0)

	// Close the main window
	mainWindow.Close()
}

func updateTimerLabel(label *widgets.QLabel, remainingTime time.Duration) {
	seconds := int(remainingTime.Seconds())
	timerText := fmt.Sprintf("Time remaining: %02d:%02d", seconds/60, seconds%60)
	label.SetText(timerText)
}
