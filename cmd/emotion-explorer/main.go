// cmd/emotion-explorer/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"time" // Make sure time is imported

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog" // <-- Import dialog
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"

	// Use your actual module path here
	"github.com/itsforsxm123/emotion-explorer/internal/core"
	"github.com/itsforsxm123/emotion-explorer/internal/data"
	"github.com/itsforsxm123/emotion-explorer/internal/journal" // Ensure journal is imported
	"github.com/itsforsxm123/emotion-explorer/internal/ui"
)

// --- Application State ---
type AppMode string

const (
	ModeBrowsing AppMode = "browsing"
	ModeLogging  AppMode = "logging"
)

var (
	emotionData     data.EmotionData
	primaryEmotions []data.Emotion
	mainWindow      fyne.Window
	myApp           fyne.App
	currentMode     AppMode = ModeBrowsing // Initialize in browsing mode
)

// --- Navigation/Callback Functions ---
// We will need to adapt these or create new ones for logging mode later
// ... (Keep existing functions for now) ...
func handlePrimaryEmotionSelected(selectedPrimary data.Emotion) {
	// This is the BROWSING mode handler
	if currentMode != ModeBrowsing {
		log.Println("Warning: handlePrimaryEmotionSelected called while not in browsing mode.")
		return // Or handle differently if needed
	}
	log.Printf("[Browse] Navigating from Primary: Selected '%s' (ID: %s)\n", selectedPrimary.Name, selectedPrimary.ID)
	secondaryChildren := core.GetChildrenOf(selectedPrimary.ID, emotionData.Emotions)
	log.Printf("[Browse] Found %d children for '%s'.\n", len(secondaryChildren), selectedPrimary.Name)
	secondaryView := ui.CreateEmotionListView(
		fmt.Sprintf("Exploring: %s", selectedPrimary.Name),
		&selectedPrimary,
		secondaryChildren,
		func(clickedSecondary data.Emotion) {
			handleSecondaryEmotionSelected(selectedPrimary, clickedSecondary) // Still uses browsing handler
		},
		navigateBackToPrimary, // Browsing back handler
		"<- Back to Primary",
	)
	log.Println("[Browse] Setting window content to Secondary View...")
	mainWindow.SetContent(secondaryView)
}
func handleSecondaryEmotionSelected(primaryParent data.Emotion, selectedSecondary data.Emotion) {
	// This is the BROWSING mode handler
	if currentMode != ModeBrowsing {
		log.Println("Warning: handleSecondaryEmotionSelected called while not in browsing mode.")
		return
	}
	log.Printf("[Browse] Navigating from Secondary: Selected '%s' (ID: %s, Primary Parent: '%s')\n",
		selectedSecondary.Name, selectedSecondary.ID, primaryParent.Name)
	tertiaryChildren := core.GetChildrenOf(selectedSecondary.ID, emotionData.Emotions)
	log.Printf("[Browse] Found %d children for '%s'.\n", len(tertiaryChildren), selectedSecondary.Name)
	if len(tertiaryChildren) > 0 {
		log.Printf("[Browse] Navigating to Tertiary View for '%s'...", selectedSecondary.Name)
		tertiaryView := ui.CreateEmotionListView(
			fmt.Sprintf("Exploring under: %s", selectedSecondary.Name),
			&selectedSecondary,
			tertiaryChildren,
			func(clickedTertiary data.Emotion) {
				handleTertiaryEmotionSelected(selectedSecondary, clickedTertiary) // Still uses browsing handler
			},
			func() {
				navigateBackToSecondary(primaryParent) // Browsing back handler
			},
			"<- Back to Secondary",
		)
		mainWindow.SetContent(tertiaryView)
	} else {
		log.Printf("[Browse] Leaf Node: No further children found for '%s'. (Detail view TBD)\n", selectedSecondary.Name)
	}
}
func handleTertiaryEmotionSelected(secondaryParent data.Emotion, selectedTertiary data.Emotion) {
	// This is the BROWSING mode handler
	if currentMode != ModeBrowsing {
		log.Println("Warning: handleTertiaryEmotionSelected called while not in browsing mode.")
		return
	}
	log.Printf("[Browse] Tertiary Button Clicked: '%s' (ID: %s, Secondary Parent: '%s'). No further navigation.\n",
		selectedTertiary.Name, selectedTertiary.ID, secondaryParent.Name)
}
func navigateBackToPrimary() {
	// This is the BROWSING mode handler
	if currentMode != ModeBrowsing {
		log.Println("Warning: navigateBackToPrimary called while not in browsing mode.")
		// If called during logging, maybe treat as cancel? For now, just log.
		return
	}
	log.Println("[Browse] Navigating back to Primary View...")
	// Re-create Primary View using Generic Function for BROWSING
	primaryView := createBrowsingPrimaryView() // Use helper
	log.Println("[Browse] Setting window content back to Primary View...")
	mainWindow.SetContent(primaryView)
}
func navigateBackToSecondary(primaryParent data.Emotion) {
	// This is the BROWSING mode handler
	if currentMode != ModeBrowsing {
		log.Println("Warning: navigateBackToSecondary called while not in browsing mode.")
		return
	}
	log.Printf("[Browse] Navigating back to Secondary View (Parent: '%s')...", primaryParent.Name)
	secondaryChildren := core.GetChildrenOf(primaryParent.ID, emotionData.Emotions)
	secondaryView := ui.CreateEmotionListView(
		fmt.Sprintf("Exploring: %s", primaryParent.Name),
		&primaryParent,
		secondaryChildren,
		func(clickedSecondary data.Emotion) {
			handleSecondaryEmotionSelected(primaryParent, clickedSecondary) // Browsing handler
		},
		navigateBackToPrimary, // Browsing back handler
		"<- Back to Primary",
	)
	log.Println("[Browse] Setting window content back to Secondary View...")
	mainWindow.SetContent(secondaryView)
}

// --- View Creation Helpers ---

// createBrowsingPrimaryView creates the main view for browsing emotions
func createBrowsingPrimaryView() fyne.CanvasObject {
	log.Println("Creating browsing primary view...")
	return ui.CreateEmotionListView(
		"Primary Emotions",           // Title
		nil,                          // Parent context
		primaryEmotions,              // Emotions
		handlePrimaryEmotionSelected, // Use BROWSING handler
		nil,                          // No back button from primary browsing
		"",                           // Back button label (not used)
	)
}

// showLoggingSelectionView will display the UI for selecting an emotion to log
// (We will implement the actual view creation logic in the next step)
func showLoggingSelectionView() {
	log.Println("Switching to Logging Mode - Showing emotion selection view (TBD)...")
	currentMode = ModeLogging
	// TODO: Replace this with the actual logging selection view
	// For now, just show primary emotions again, but with a different title/context
	// We need new handlers for logging mode clicks.
	placeholderView := ui.CreateEmotionListView(
		"Select Emotion to Log",  // New Title
		nil,                      // Parent context
		primaryEmotions,          // Start with primary emotions
		handleLogEmotionSelected, // Use NEW LOGGING handler
		cancelLogging,            // Use NEW CANCEL handler for back button
		"Cancel Logging",         // Back button label
	)

	mainWindow.SetTitle("Emotion Explorer - Logging") // Update window title
	mainWindow.SetContent(placeholderView)
}

// cancelLogging is called when the user cancels the logging process
func cancelLogging() {
	log.Println("Logging cancelled by user.")
	currentMode = ModeBrowsing
	mainWindow.SetTitle("Emotion Explorer")            // Reset window title
	mainWindow.SetContent(createBrowsingPrimaryView()) // Go back to browsing view
}

// handleLogEmotionSelected is the callback for when an emotion is selected in LOGGING mode
// (This is a placeholder - needs full implementation)
func handleLogEmotionSelected(selectedEmotion data.Emotion) {
	log.Printf("[Log Mode] Emotion selected: %s", selectedEmotion.Name)

	children := core.GetChildrenOf(selectedEmotion.ID, emotionData.Emotions)

	if len(children) == 0 {
		// *** This is a LEAF node - LOG IT! ***
		log.Printf("[Log Mode] Leaf emotion '%s' selected. Saving...", selectedEmotion.Name)

		entryToSave := journal.LogEntry{
			Timestamp:   time.Now(),
			EmotionID:   selectedEmotion.ID,
			EmotionName: selectedEmotion.Name,
			Notes:       "", // TODO: Add notes field later
		}

		err := journal.SaveLogEntry(entryToSave)
		if err != nil {
			log.Printf("ERROR: Failed to save log entry: %v", err)
			dialog.ShowError(fmt.Errorf("failed to save journal entry: %w", err), mainWindow)
		} else {
			log.Println("[Log Mode] Log entry saved successfully.")
			dialog.ShowInformation("Success", fmt.Sprintf("Logged '%s'", selectedEmotion.Name), mainWindow)
		}

		// --- Return to browsing mode ---
		currentMode = ModeBrowsing
		mainWindow.SetTitle("Emotion Explorer")
		mainWindow.SetContent(createBrowsingPrimaryView())
		// ---

	} else {
		// *** Not a leaf node - navigate deeper within logging mode ***
		log.Printf("[Log Mode] Navigating deeper from '%s'...", selectedEmotion.Name)
		// Need to show the children using the *logging* handlers
		loggingChildView := ui.CreateEmotionListView(
			fmt.Sprintf("Select under: %s", selectedEmotion.Name), // Title
			&selectedEmotion,         // Parent context
			children,                 // Emotions to display
			handleLogEmotionSelected, // Recursive call to LOGGING handler
			func() { // Go back within logging mode
				// If we go back from secondary log selection, show primary log selection
				// If we go back from tertiary log selection, show secondary log selection
				// For now, simplify: going back always goes to primary log selection
				// TODO: Implement proper back navigation within logging mode
				showLoggingSelectionView() // Go back to top level logging selection
			},
			"<- Back / Cancel", // Back button label
		)
		mainWindow.SetContent(loggingChildView)
	}
}

// --- Main Function ---

func main() {
	// --- 1. Load Emotion Data --- (No changes)
	// ...
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
	myApp = app.New()
	mainWindow = myApp.NewWindow("Emotion Explorer") // Initial title

	// --- 3. Get Primary Emotions --- (No changes)
	// ...
	log.Println("Extracting primary emotions...")
	primaryEmotions = core.GetPrimaryEmotions(emotionData.Emotions)
	log.Printf("Found %d primary emotions.", len(primaryEmotions))
	if len(primaryEmotions) == 0 {
		log.Println("Warning: No primary emotions found. Check emotions.json.")
	}

	// --- 4. Create the INITIAL UI View (Browsing Mode) ---
	initialBrowsingView := createBrowsingPrimaryView() // Use helper

	// --- 5. Set Initial Window Content and Size --- (No changes)
	// ...
	log.Println("Setting initial window content (Browsing Mode)...")
	mainWindow.SetContent(initialBrowsingView)
	mainWindow.Resize(fyne.NewSize(600, 400))
	mainWindow.CenterOnScreen()

	// --- 6. System Tray Setup ---
	if desk, ok := myApp.(desktop.App); ok {
		log.Println("System tray supported. Setting up...")

		// --- Create Menu Items ---
		showItem := fyne.NewMenuItem("Show Window", func() {
			log.Println("Showing main window via tray menu.")
			mainWindow.Show()
			mainWindow.RequestFocus()
		})

		// *** UPDATE Log Feeling Item Action ***
		logFeelingItem := fyne.NewMenuItem("Log Current Feeling...", func() {
			log.Println("Log Current Feeling... menu item clicked. Switching to Log Mode.")
			// Don't save dummy data anymore. Instead, trigger the logging UI flow.
			showLoggingSelectionView() // Call the function to display the selection UI

			// Ensure the main window is visible
			mainWindow.Show()
			mainWindow.RequestFocus()
		})
		// *** END UPDATE ***

		quitItem := fyne.NewMenuItem("Quit", func() {
			log.Println("Quitting application via tray menu.")
			myApp.Quit()
		})

		// --- Create the Menu --- (No changes)
		trayMenu := fyne.NewMenu("Emotion Explorer",
			showItem,
			logFeelingItem,
			fyne.NewMenuItemSeparator(),
			quitItem,
		)

		// Set the Tray Icon and Menu (No changes)
		trayIcon := theme.FyneLogo()
		desk.SetSystemTrayIcon(trayIcon)
		desk.SetSystemTrayMenu(trayMenu)

		// Intercept Window Close Requests (No changes)
		mainWindow.SetCloseIntercept(func() {
			// If in logging mode, maybe ask confirmation before hiding? For now, just hide.
			if currentMode == ModeLogging {
				log.Println("Main window close intercepted during logging. Hiding window (logging cancelled implicitly).")
				// Implicitly cancel logging if window is closed during the process
				cancelLogging() // Switch back to browsing mode
			} else {
				log.Println("Main window close intercepted during browsing. Hiding window.")
			}
			mainWindow.Hide()
		})

	} else {
		log.Println("System tray not supported on this system (app does not implement desktop.App).")
	}
	// --- End System Tray Setup ---

	// --- 7. Show Window and Run App ---
	log.Println("Showing window and running app...")
	mainWindow.ShowAndRun()

	log.Println("Application finished.")
}
