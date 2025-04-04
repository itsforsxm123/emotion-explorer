// cmd/emotion-explorer/main.go
package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/itsforsxm123/emotion-explorer/internal/core"
	"github.com/itsforsxm123/emotion-explorer/internal/data"
	"github.com/itsforsxm123/emotion-explorer/internal/ui"
)

// --- Main Application State (kept simple within main for now) ---
var (
	// Keep loaded data accessible
	emotionData data.EmotionData
	// Keep primary emotions list accessible for navigation
	primaryEmotions []data.Emotion
	// Keep main window accessible for content updates
	mainWindow fyne.Window
)

// --- Navigation/Callback Functions ---

// handlePrimaryEmotionSelected is called when a primary emotion button is clicked.
// It finds children and displays the secondary view.
func handlePrimaryEmotionSelected(selectedEmotion data.Emotion) {
	log.Printf("Navigating from Primary: Selected '%s' (ID: %s)\n", selectedEmotion.Name, selectedEmotion.ID)

	// 1. Find Children
	children := core.GetChildrenOf(selectedEmotion.ID, emotionData.Emotions)
	log.Printf("Found %d children for '%s'.\n", len(children), selectedEmotion.Name)

	// 2. Create Secondary View
	// Pass the selected parent, its children, and the function to go back
	secondaryView := ui.CreateSecondaryEmotionView(selectedEmotion, children, navigateBackToPrimary)

	// 3. Update Window Content
	log.Println("Setting window content to Secondary View...")
	mainWindow.SetContent(secondaryView)
}

// navigateBackToPrimary is called by the "Back" button in the secondary view.
// It recreates and displays the primary view.
func navigateBackToPrimary() {
	log.Println("Navigating back to Primary View...")

	// 1. Re-create Primary View
	// We need the original list of primaryEmotions and the handlePrimaryEmotionSelected callback
	primaryView := ui.CreatePrimaryEmotionView(primaryEmotions, handlePrimaryEmotionSelected)

	// 2. Update Window Content
	log.Println("Setting window content back to Primary View...")
	mainWindow.SetContent(primaryView)
}

// --- Main Function ---

func main() {
	// --- 1. Load Emotion Data ---
	log.Println("Loading emotion data...")
	var err error                          // Declare err here to avoid shadowing in the next line
	emotionData, err = data.LoadEmotions() // Assign to package-level var
	if err != nil {
		log.Printf("FATAL: Failed to load emotion data: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Successfully loaded emotion data. Version: %s", emotionData.Metadata.Version)
	log.Printf("Found %d total emotions defined.", len(emotionData.Emotions))

	// --- 2. Initialize Fyne App ---
	myApp := app.New()
	mainWindow = myApp.NewWindow("Emotion Explorer") // Assign to package-level var

	// --- 3. Get Primary Emotions ---
	log.Println("Extracting primary emotions...")
	primaryEmotions = core.GetPrimaryEmotions(emotionData.Emotions) // Assign to package-level var
	log.Printf("Found %d primary emotions.", len(primaryEmotions))
	if len(primaryEmotions) == 0 {
		log.Println("Warning: No primary emotions found. Check emotions.json.")
		// App will show "No primary emotions found." label from CreatePrimaryEmotionView
	}

	// --- 4. Create the INITIAL Primary Emotion UI View ---
	log.Println("Creating initial primary emotion view...")
	// Pass the list of emotions AND the callback function for when one is selected
	initialPrimaryView := ui.CreatePrimaryEmotionView(primaryEmotions, handlePrimaryEmotionSelected)

	// --- 5. Set Initial Window Content and Show ---
	log.Println("Setting initial window content...")
	mainWindow.SetContent(initialPrimaryView) // Set the initial view
	mainWindow.Resize(fyne.NewSize(600, 400))
	mainWindow.CenterOnScreen()
	log.Println("Showing window and running app...")
	mainWindow.ShowAndRun() // Start the Fyne event loop

	log.Println("Application finished.")
}
