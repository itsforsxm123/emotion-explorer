package main

import (
	"log" // For logging errors

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	// Import our internal data package
	// Use the module name defined in go.mod (e.g., "emotion-explorer")
	// followed by the path to the package
	"emotion-explorer/internal/data"
)

func main() {
	// 1. Create a new Fyne application
	myApp := app.New()

	// 2. Create a new window
	myWindow := myApp.NewWindow("Emotion Hierarchy Explorer")

	// 3. Load the emotion data
	emotionData, err := data.LoadEmotions()
	if err != nil {
		// If data loading fails, log the error and exit.
		// In a real app, you might show an error dialog.
		log.Fatalf("FATAL: Failed to load emotion data: %v", err)
	}

	// Log success for now (optional)
	log.Printf("Successfully loaded %d emotions.", len(emotionData.Emotions))

	// 4. Create a simple placeholder widget
	// We'll replace this soon with the actual emotion selection UI
	placeholderLabel := widget.NewLabel("Welcome to the Emotion Explorer! Data loaded.")

	// 5. Set the window content
	myWindow.SetContent(placeholderLabel)

	// 6. Set an initial size for the window
	myWindow.Resize(fyne.NewSize(600, 400)) // Adjust size as needed

	// 7. Show the window and run the application loop
	// This blocks until the application is closed.
	myWindow.ShowAndRun()

	// Code here will execute after the window is closed
	log.Println("Application finished.")
}
