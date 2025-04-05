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

// --- Main Application State ---
var (
	emotionData     data.EmotionData
	primaryEmotions []data.Emotion
	mainWindow      fyne.Window
)

// --- Navigation/Callback Functions ---

// handlePrimaryEmotionSelected is called when a primary emotion button is clicked.
// It finds children and displays the secondary view.
func handlePrimaryEmotionSelected(selectedPrimary data.Emotion) {
	log.Printf("Navigating from Primary: Selected '%s' (ID: %s)\n", selectedPrimary.Name, selectedPrimary.ID)

	// 1. Find Children (Secondary Emotions)
	secondaryChildren := core.GetChildrenOf(selectedPrimary.ID, emotionData.Emotions)
	log.Printf("Found %d children for '%s'.\n", len(secondaryChildren), selectedPrimary.Name)

	// 2. Create Secondary View
	// MODIFIED: Pass a closure for onSecondaryEmotionSelected that captures selectedPrimary
	secondaryView := ui.CreateSecondaryEmotionView(
		selectedPrimary,   // The parent emotion
		secondaryChildren, // The children to display
		func(clickedSecondary data.Emotion) { // <<< THIS IS THE MISSING ARGUMENT from the error message
			handleSecondaryEmotionSelected(selectedPrimary, clickedSecondary) // Pass both primary parent and clicked secondary
		},
		navigateBackToPrimary, // Function to go back to the primary view
	)
	// 3. Update Window Content
	log.Println("Setting window content to Secondary View...")
	mainWindow.SetContent(secondaryView)
}

// handleSecondaryEmotionSelected is called when a secondary emotion button is clicked.
// It finds children and displays the tertiary view, or logs if it's a leaf node.
func handleSecondaryEmotionSelected(primaryParent data.Emotion, selectedSecondary data.Emotion) {
	log.Printf("Navigating from Secondary: Selected '%s' (ID: %s, Primary Parent: '%s')\n",
		selectedSecondary.Name, selectedSecondary.ID, primaryParent.Name)

	// 1. Find Children (Tertiary Emotions)
	tertiaryChildren := core.GetChildrenOf(selectedSecondary.ID, emotionData.Emotions)
	log.Printf("Found %d children for '%s'.\n", len(tertiaryChildren), selectedSecondary.Name)

	// 2. Check if children exist (is it a leaf node?)
	if len(tertiaryChildren) > 0 {
		// 3a. Create Tertiary View
		log.Printf("Navigating to Tertiary View for '%s'...", selectedSecondary.Name)
		tertiaryView := ui.CreateTertiaryEmotionView(
			selectedSecondary, // The parent (secondary emotion) for this view
			tertiaryChildren,  // The children (tertiary emotions) to display
			func() { // <<< NEW: Closure for the back button
				navigateBackToSecondary(primaryParent) // Pass the primary parent needed to rebuild the secondary view
			},
		)
		// 4a. Update Window Content
		mainWindow.SetContent(tertiaryView)
	} else {
		// 3b. Leaf Node: No further children
		log.Printf("Leaf Node: No further children found for '%s'. (Detail view TBD)\n", selectedSecondary.Name)
		// In a real app, you might show a detail panel here,
		// but for now, we just log and stay on the secondary screen.
		// Optionally, show a temporary notification/dialog? (Skipping for now)
	}
}

// navigateBackToPrimary is called by the "Back" button in the secondary view.
// It recreates and displays the primary view.
func navigateBackToPrimary() {
	log.Println("Navigating back to Primary View...")

	// 1. Re-create Primary View
	primaryView := ui.CreatePrimaryEmotionView(primaryEmotions, handlePrimaryEmotionSelected)

	// 2. Update Window Content
	log.Println("Setting window content back to Primary View...")
	mainWindow.SetContent(primaryView)
}

// navigateBackToSecondary is called by the "Back" button in the tertiary view.
// It recreates and displays the correct secondary view.
func navigateBackToSecondary(primaryParent data.Emotion) {
	log.Printf("Navigating back to Secondary View (Parent: '%s')...", primaryParent.Name)

	// 1. Find the children of the primary parent again (the secondary emotions)
	secondaryChildren := core.GetChildrenOf(primaryParent.ID, emotionData.Emotions)

	// 2. Re-create Secondary View
	// We need to pass the *exact same* arguments as when we created it the first time
	// in handlePrimaryEmotionSelected, including the closure for handling secondary clicks.
	secondaryView := ui.CreateSecondaryEmotionView(
		primaryParent,
		secondaryChildren,
		func(clickedSecondary data.Emotion) { // Recreate the callback closure
			handleSecondaryEmotionSelected(primaryParent, clickedSecondary)
		},
		navigateBackToPrimary, // Back from Secondary still goes to Primary
	)

	// 3. Update Window Content
	log.Println("Setting window content back to Secondary View...")
	mainWindow.SetContent(secondaryView)
}

// --- Main Function --- (No changes needed below this line)

func main() {
	// --- 1. Load Emotion Data ---
	log.Println("Loading emotion data...")
	var err error
	emotionData, err = data.LoadEmotions()
	if err != nil {
		log.Printf("FATAL: Failed to load emotion data: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Successfully loaded emotion data. Version: %s", emotionData.Metadata.Version)
	log.Printf("Found %d total emotions defined.", len(emotionData.Emotions))

	// --- 2. Initialize Fyne App ---
	myApp := app.New()
	mainWindow = myApp.NewWindow("Emotion Explorer")

	// --- 3. Get Primary Emotions ---
	log.Println("Extracting primary emotions...")
	primaryEmotions = core.GetPrimaryEmotions(emotionData.Emotions)
	log.Printf("Found %d primary emotions.", len(primaryEmotions))
	if len(primaryEmotions) == 0 {
		log.Println("Warning: No primary emotions found. Check emotions.json.")
	}

	// --- 4. Create the INITIAL Primary Emotion UI View ---
	log.Println("Creating initial primary emotion view...")
	initialPrimaryView := ui.CreatePrimaryEmotionView(primaryEmotions, handlePrimaryEmotionSelected)

	// --- 5. Set Initial Window Content and Show ---
	log.Println("Setting initial window content...")
	mainWindow.SetContent(initialPrimaryView)
	mainWindow.Resize(fyne.NewSize(600, 400))
	mainWindow.CenterOnScreen()
	log.Println("Showing window and running app...")
	mainWindow.ShowAndRun()

	log.Println("Application finished.")
}
