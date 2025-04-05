// cmd/emotion-explorer/main.go
package main

import (
	"fmt" // Import fmt for Sprintf
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/itsforsxm123/emotion-explorer/internal/core"
	"github.com/itsforsxm123/emotion-explorer/internal/data"
	"github.com/itsforsxm123/emotion-explorer/internal/ui"
)

// --- Main Application State --- (Remains the same)
var (
	emotionData     data.EmotionData
	primaryEmotions []data.Emotion
	mainWindow      fyne.Window
)

// --- Navigation/Callback Functions ---

// handlePrimaryEmotionSelected is called when a primary emotion button is clicked.
// It finds children and displays the secondary view using the generic function.
func handlePrimaryEmotionSelected(selectedPrimary data.Emotion) {
	log.Printf("Navigating from Primary: Selected '%s' (ID: %s)\n", selectedPrimary.Name, selectedPrimary.ID)

	// 1. Find Children (Secondary Emotions)
	secondaryChildren := core.GetChildrenOf(selectedPrimary.ID, emotionData.Emotions)
	log.Printf("Found %d children for '%s'.\n", len(secondaryChildren), selectedPrimary.Name)

	// 2. Create Secondary View using Generic Function
	secondaryView := ui.CreateEmotionListView(
		fmt.Sprintf("Exploring: %s", selectedPrimary.Name), // Title
		&selectedPrimary,  // Parent context
		secondaryChildren, // Emotions to display
		func(clickedSecondary data.Emotion) { // onSelected callback (closure)
			handleSecondaryEmotionSelected(selectedPrimary, clickedSecondary) // Pass context
		},
		navigateBackToPrimary, // goBack callback
		"<- Back to Primary",  // Back button label
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
		// 3a. Create Tertiary View using Generic Function
		log.Printf("Navigating to Tertiary View for '%s'...", selectedSecondary.Name)
		tertiaryView := ui.CreateEmotionListView(
			fmt.Sprintf("Exploring under: %s", selectedSecondary.Name), // Title
			&selectedSecondary, // Parent context
			tertiaryChildren,   // Emotions to display
			func(clickedTertiary data.Emotion) { // onSelected callback (closure)
				handleTertiaryEmotionSelected(selectedSecondary, clickedTertiary) // Pass context
			},
			func() { // goBack callback (closure)
				navigateBackToSecondary(primaryParent) // Pass context needed for back nav
			},
			"<- Back to Secondary", // Back button label
		)
		// 4a. Update Window Content
		mainWindow.SetContent(tertiaryView)
	} else {
		// 3b. Leaf Node: No further children
		log.Printf("Leaf Node: No further children found for '%s'. (Detail view TBD)\n", selectedSecondary.Name)
		// TODO: Implement Detail View navigation here later
	}
}

// handleTertiaryEmotionSelected is called when a tertiary emotion button is clicked.
// Currently, it just logs the selection as we don't navigate further down.
func handleTertiaryEmotionSelected(secondaryParent data.Emotion, selectedTertiary data.Emotion) {
	log.Printf("Tertiary Button Clicked: '%s' (ID: %s, Secondary Parent: '%s'). No further navigation.\n",
		selectedTertiary.Name, selectedTertiary.ID, secondaryParent.Name)
	// TODO: Implement Detail View navigation here later (if tertiary can also be leaf)
}

// navigateBackToPrimary is called by the "Back" button in the secondary view.
// It recreates and displays the primary view using the generic function.
func navigateBackToPrimary() {
	log.Println("Navigating back to Primary View...")

	// 1. Re-create Primary View using Generic Function
	primaryView := ui.CreateEmotionListView(
		"Primary Emotions",           // Title
		nil,                          // Parent context (none for primary)
		primaryEmotions,              // Emotions to display
		handlePrimaryEmotionSelected, // onSelected callback
		nil,                          // goBack callback (none for primary)
		"",                           // Back button label (not used)
	)

	// 2. Update Window Content
	log.Println("Setting window content back to Primary View...")
	mainWindow.SetContent(primaryView)
}

// navigateBackToSecondary is called by the "Back" button in the tertiary view.
// It recreates and displays the correct secondary view using the generic function.
func navigateBackToSecondary(primaryParent data.Emotion) {
	log.Printf("Navigating back to Secondary View (Parent: '%s')...", primaryParent.Name)

	// 1. Find the children of the primary parent again (the secondary emotions)
	secondaryChildren := core.GetChildrenOf(primaryParent.ID, emotionData.Emotions)

	// 2. Re-create Secondary View using Generic Function
	// We need to pass the *exact same* arguments as when we created it the first time
	// in handlePrimaryEmotionSelected, including the closures for handling clicks and back nav.
	secondaryView := ui.CreateEmotionListView(
		fmt.Sprintf("Exploring: %s", primaryParent.Name), // Title
		&primaryParent,    // Parent context
		secondaryChildren, // Emotions to display
		func(clickedSecondary data.Emotion) { // onSelected callback (closure)
			handleSecondaryEmotionSelected(primaryParent, clickedSecondary) // Pass context
		},
		navigateBackToPrimary, // goBack callback
		"<- Back to Primary",  // Back button label
	)

	// 3. Update Window Content
	log.Println("Setting window content back to Secondary View...")
	mainWindow.SetContent(secondaryView)
}

// --- Main Function --- (Only initial view creation changes)

func main() {
	// --- 1. Load Emotion Data --- (No changes)
	log.Println("Loading emotion data...")
	var err error
	emotionData, err = data.LoadEmotions()
	if err != nil {
		log.Printf("FATAL: Failed to load emotion data: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Successfully loaded emotion data. Version: %s", emotionData.Metadata.Version)
	log.Printf("Found %d total emotions defined.", len(emotionData.Emotions))

	// --- 2. Initialize Fyne App --- (No changes)
	myApp := app.New()
	mainWindow = myApp.NewWindow("Emotion Explorer")

	// --- 3. Get Primary Emotions --- (No changes)
	log.Println("Extracting primary emotions...")
	primaryEmotions = core.GetPrimaryEmotions(emotionData.Emotions)
	log.Printf("Found %d primary emotions.", len(primaryEmotions))
	if len(primaryEmotions) == 0 {
		log.Println("Warning: No primary emotions found. Check emotions.json.")
	}

	// --- 4. Create the INITIAL Primary Emotion UI View using Generic Function ---
	log.Println("Creating initial primary emotion view...")
	initialPrimaryView := ui.CreateEmotionListView(
		"Primary Emotions",           // Title
		nil,                          // Parent context
		primaryEmotions,              // Emotions
		handlePrimaryEmotionSelected, // onSelected callback
		nil,                          // goBack callback (no back button)
		"",                           // Back button label (not used)
	)

	// --- 5. Set Initial Window Content and Show --- (No changes)
	log.Println("Setting initial window content...")
	mainWindow.SetContent(initialPrimaryView)
	mainWindow.Resize(fyne.NewSize(600, 400)) // Keep consistent size
	mainWindow.CenterOnScreen()
	log.Println("Showing window and running app...")
	mainWindow.ShowAndRun()

	log.Println("Application finished.")
}
