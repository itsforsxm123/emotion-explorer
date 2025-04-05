// cmd/emotion-explorer/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"time" // Make sure time is imported

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container" // Import container
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout" // Import layout
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget" // Import widget

	// Use your actual module path here
	"github.com/itsforsxm123/emotion-explorer/internal/core"
	"github.com/itsforsxm123/emotion-explorer/internal/data"
	"github.com/itsforsxm123/emotion-explorer/internal/journal"
	"github.com/itsforsxm123/emotion-explorer/internal/ui"
)

// --- Application State ---

// AppMode defines the current operational mode of the application.
type AppMode int // Use int for enums, it's more idiomatic Go

const (
	ModeBrowsing AppMode = iota // Default mode: exploring emotions.
	ModeLogging                 // Mode for selecting an emotion to log.
)

const (
	appName         = "Emotion Explorer"
	logModeTitle    = appName + " - Logging..."
	browseModeTitle = appName
)

var (
	// Core App Components
	myApp      fyne.App
	mainWindow fyne.Window

	// Data
	emotionData     data.EmotionData // Consider if this needs to be global or passed around
	primaryEmotions []data.Emotion   // Cache primary emotions

	// UI Elements
	backButton       *widget.Button
	mainContentArea  *fyne.Container // The container holding the current view (center of border)
	mainBorderLayout *fyne.Container

	// State Management
	currentMode            AppMode              = ModeBrowsing
	navigationStack        *[]fyne.CanvasObject // Stack for browsing views
	loggingNavigationStack *[]fyne.CanvasObject // Stack for logging views
)

// --- Initialization ---

func main() {
	// 1. Initialize App and Load Data
	myApp = app.New()
	mainWindow = myApp.NewWindow(browseModeTitle) // Initial title

	if err := loadData(); err != nil {
		// Consider showing a dialog even before the main window is fully set up
		log.Printf("FATAL: Failed to load emotion data: %v\n", err)
		// dialog.ShowError(err, mainWindow) // This might fail if mainWindow isn't ready
		fmt.Fprintf(os.Stderr, "Error loading emotion data: %v\n", err) // Fallback to stderr
		os.Exit(1)
	}

	// 2. Initialize Navigation Stacks
	navStack := make([]fyne.CanvasObject, 0, 5) // Pre-allocate some capacity
	navigationStack = &navStack
	logNavStack := make([]fyne.CanvasObject, 0, 5)
	loggingNavigationStack = &logNavStack

	// 3. Setup Core UI Layout
	setupMainLayout() // Creates the border layout with back button and content area

	// 4. Push Initial View (Browsing Primary Emotions)
	initialBrowsingView := createEmotionListView("Primary Emotions", nil, primaryEmotions, handleEmotionSelected)
	pushView(initialBrowsingView, navigationStack) // Push to browsing stack initially

	// 5. Setup System Tray & Window Behavior
	setupSystemTray()
	setupWindowIntercepts()

	// 6. Resize, Center, Show, and Run
	mainWindow.Resize(fyne.NewSize(400, 500)) // Adjusted size
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()

	log.Println("Application finished.")
}

// loadData encapsulates the emotion data loading logic.
func loadData() error {
	log.Println("Loading emotion data...")
	var err error
	emotionData, err = data.LoadEmotions()
	if err != nil {
		return fmt.Errorf("failed to load emotions: %w", err)
	}
	log.Printf("Successfully loaded emotion data. Version: %s", emotionData.Metadata.Version)
	log.Printf("Found %d total emotions defined.", len(emotionData.Emotions))

	log.Println("Extracting primary emotions...")
	primaryEmotions = core.GetPrimaryEmotions(emotionData.Emotions) // Use loaded data
	log.Printf("Found %d primary emotions.", len(primaryEmotions))
	if len(primaryEmotions) == 0 {
		log.Println("Warning: No primary emotions found. Check emotions.json.")
	}
	return nil
}

// setupMainLayout creates the main window structure (border layout).
func setupMainLayout() {
	backButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), handleBack) // Use icon
	backButton.Disable()                                                            // Start disabled

	// This container will hold the dynamic content (emotion lists)
	mainContentArea = container.NewMax() // Use Max layout to fill available space

	// Create the main border layout
	border := container.NewBorder(
		container.NewHBox(backButton, layout.NewSpacer()), // Top: Back button aligned left
		nil,             // Bottom
		nil,             // Left
		nil,             // Right
		mainContentArea, // Center: Dynamic content goes here
	)
	mainBorderLayout = border // Store reference if needed, though direct access via mainWindow.Content() works
	mainWindow.SetContent(border)
	log.Println("Main layout setup complete.")
}

// --- Navigation Stack Management ---

// pushView adds a new view to the specified navigation stack and updates the UI.
func pushView(view fyne.CanvasObject, stack *[]fyne.CanvasObject) {
	*stack = append(*stack, view)
	log.Printf("Pushed view. Stack size: %d. Mode: %v", len(*stack), currentMode)
	updateContentFromActiveStack() // Update content based on the active stack
	updateBackButtonState()        // Update button state after push
}

// popView removes the top view from the specified navigation stack and updates the UI.
// Returns true if a pop occurred, false if the stack was empty or had only one item.
func popView(stack *[]fyne.CanvasObject) bool {
	if len(*stack) <= 1 {
		log.Printf("Pop requested on stack with size %d. Cannot pop.", len(*stack))
		return false // Cannot pop the last view
	}
	*stack = (*stack)[:len(*stack)-1] // Pop the last element
	log.Printf("Popped view. Stack size: %d. Mode: %v", len(*stack), currentMode)
	updateContentFromActiveStack() // Update content based on the active stack
	updateBackButtonState()        // Update button state after pop
	return true
}

// --- UI Update Logic ---

// updateContentFromActiveStack sets the main content area based on the top of the active stack.
func updateContentFromActiveStack() {
	var activeStack *[]fyne.CanvasObject
	if currentMode == ModeLogging {
		activeStack = loggingNavigationStack
	} else {
		activeStack = navigationStack
	}

	if len(*activeStack) == 0 {
		log.Println("Error: Active stack is empty, cannot update content.")
		// Show an error message or a placeholder in the UI?
		mainContentArea.Objects = []fyne.CanvasObject{widget.NewLabel("Error: No view available.")}
		mainContentArea.Refresh()
		return
	}

	// Get the top view from the active stack
	topView := (*activeStack)[len(*activeStack)-1]

	// Update the main content area
	mainContentArea.Objects = []fyne.CanvasObject{topView} // Replace objects in Max container
	mainContentArea.Refresh()
	log.Println("Main content area updated.")
}

// updateBackButtonState enables/disables the back button based on the active stack size.
func updateBackButtonState() {
	var activeStack *[]fyne.CanvasObject
	if currentMode == ModeLogging {
		activeStack = loggingNavigationStack
	} else {
		activeStack = navigationStack
	}

	if len(*activeStack) <= 1 {
		backButton.Disable()
		log.Println("Back button disabled.")
	} else {
		backButton.Enable()
		log.Println("Back button enabled.")
	}
}

// --- Event Handlers ---

// handleBack manages the back navigation logic for both modes.
func handleBack() {
	log.Println("Back button clicked.")
	if currentMode == ModeLogging {
		if !popView(loggingNavigationStack) {
			// If pop failed (we are at the root of logging), treat as cancel
			log.Println("Back clicked at root of logging stack. Cancelling logging.")
			switchToBrowsingMode() // Or could just stay here, depends on desired UX
		}
	} else {
		popView(navigationStack) // Pop the browsing stack
	}
}

// handleEmotionSelected is the central callback for emotion selection in ANY mode.
// It delegates to mode-specific handlers.
func handleEmotionSelected(selectedEmotion data.Emotion) {
	log.Printf("Emotion selected: '%s' (ID: %s) in Mode: %v", selectedEmotion.Name, selectedEmotion.ID, currentMode)
	if currentMode == ModeLogging {
		handleLogEmotionSelection(selectedEmotion)
	} else {
		handleBrowseEmotionSelection(selectedEmotion)
	}
}

// handleBrowseEmotionSelection handles navigation when an emotion is selected in browsing mode.
func handleBrowseEmotionSelection(selectedEmotion data.Emotion) {
	children := core.GetChildrenOf(selectedEmotion.ID, emotionData.Emotions)
	log.Printf("[Browse] Found %d children for '%s'.", len(children), selectedEmotion.Name)

	if len(children) > 0 {
		title := fmt.Sprintf("Exploring: %s", selectedEmotion.Name)
		// Create and push the new view onto the browsing stack
		childView := createEmotionListView(title, &selectedEmotion, children, handleEmotionSelected) // Use central handler
		pushView(childView, navigationStack)
	} else {
		// Leaf node in browsing mode - maybe show details in the future
		log.Printf("[Browse] Leaf Node: '%s'. (Detail view TBD)", selectedEmotion.Name)
		dialog.ShowInformation("Emotion Details", fmt.Sprintf("Selected: %s\n(More details could be shown here)", selectedEmotion.Name), mainWindow)
	}
}

// handleLogEmotionSelection handles navigation or saving when an emotion is selected in logging mode.
func handleLogEmotionSelection(selectedEmotion data.Emotion) {
	children := core.GetChildrenOf(selectedEmotion.ID, emotionData.Emotions)
	log.Printf("[Log] Found %d children for '%s'.", len(children), selectedEmotion.Name)

	if len(children) > 0 {
		// Navigate deeper within logging mode
		title := fmt.Sprintf("Log > %s > ...", selectedEmotion.Name)                                 // Shorter title
		childView := createEmotionListView(title, &selectedEmotion, children, handleEmotionSelected) // Use central handler
		pushView(childView, loggingNavigationStack)
	} else {
		// Leaf node selected in logging mode - Log it!
		log.Printf("[Log] Leaf Node: '%s'. Attempting to save.", selectedEmotion.Name)
		saveLoggedEmotion(selectedEmotion) // Encapsulate saving logic
		switchToBrowsingMode()             // Return to browsing after attempting save
	}
}

// saveLoggedEmotion handles the process of saving a selected emotion to the journal.
func saveLoggedEmotion(emotionToLog data.Emotion) {
	entry := journal.LogEntry{
		Timestamp:   time.Now(),
		EmotionID:   emotionToLog.ID,
		EmotionName: emotionToLog.Name,
		Notes:       "", // Notes field exists but is empty for now
	}

	err := journal.SaveLogEntry(entry)
	if err != nil {
		log.Printf("ERROR: Failed to save log entry for '%s': %v", emotionToLog.Name, err)
		dialog.ShowError(fmt.Errorf("failed to save journal entry: %w", err), mainWindow)
	} else {
		log.Printf("[Log] Entry for '%s' saved successfully.", emotionToLog.Name)
		dialog.ShowInformation("Logged", fmt.Sprintf("Successfully logged: %s", emotionToLog.Name), mainWindow)
	}
}

// --- Mode Switching Logic ---

// switchToLoggingMode prepares the UI for emotion logging.
func switchToLoggingMode() {
	if currentMode == ModeLogging {
		log.Println("Already in logging mode.")
		return // Avoid redundant setup
	}
	log.Println("Switching to Logging Mode...")
	currentMode = ModeLogging

	// Clear the previous logging stack to start fresh
	logNavStack := make([]fyne.CanvasObject, 0, 5)
	loggingNavigationStack = &logNavStack

	// Create and push the initial logging view (primary emotions)
	initialLogView := createEmotionListView("Select Feeling to Log", nil, primaryEmotions, handleEmotionSelected)
	pushView(initialLogView, loggingNavigationStack) // Push to the now active logging stack

	mainWindow.SetTitle(logModeTitle) // Update window title
	// updateContentFromActiveStack() is called by pushView
	// updateBackButtonState() is called by pushView
	mainWindow.Show()         // Ensure window is visible
	mainWindow.RequestFocus() // Bring to front
}

// switchToBrowsingMode returns the UI to the standard emotion browsing state.
func switchToBrowsingMode() {
	if currentMode == ModeBrowsing {
		log.Println("Already in browsing mode.")
		return
	}
	log.Println("Switching to Browsing Mode...")
	currentMode = ModeBrowsing

	// Clear the logging stack (optional, good for memory if logging stack could get deep)
	// logNavStack := make([]fyne.CanvasObject, 0, 5)
	// loggingNavigationStack = &logNavStack

	mainWindow.SetTitle(browseModeTitle) // Reset window title
	updateContentFromActiveStack()       // Display the top of the browsing stack
	updateBackButtonState()              // Update button based on browsing stack
	log.Println("Switched back to Browsing Mode.")
}

// --- View Creation Helper ---

// createEmotionListView wraps the call to the UI package's function.
// It now only needs the selection callback, as back is handled globally.
// NOTE: This assumes ui.CreateEmotionListView can be called without back button params.
// If ui.CreateEmotionListView *requires* back params, we need to adjust it or this wrapper.
// For now, let's assume the old ui.CreateEmotionListView is still used, taking nil/"" for back.
func createEmotionListView(
	title string,
	parent *data.Emotion, // Optional parent context
	emotions []data.Emotion,
	onSelect func(data.Emotion),
) fyne.CanvasObject {
	log.Printf("Creating view wrapper: '%s' with %d emotions.", title, len(emotions))
	// --- UPDATED CALL: Removed the nil and "" arguments ---
	return ui.CreateEmotionListView(
		title,
		parent,
		emotions,
		onSelect, // Pass the central selection handler
	)
}

// --- System Tray & Window Intercepts ---

func setupSystemTray() {
	if desk, ok := myApp.(desktop.App); ok {
		log.Println("System tray supported. Setting up...")
		m := fyne.NewMenu(appName,
			fyne.NewMenuItem("Show Window", func() {
				log.Println("Tray: Show Window clicked.")
				mainWindow.Show()
				mainWindow.RequestFocus() // Good practice to focus
			}),
			fyne.NewMenuItem("Log Current Feeling...", func() {
				log.Println("Tray: Log Current Feeling... clicked.")
				switchToLoggingMode() // Use the mode switch function
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				log.Println("Tray: Quit clicked.")
				myApp.Quit()
			}),
		)
		// Consider using a specific icon resource later
		desk.SetSystemTrayIcon(theme.FyneLogo())
		desk.SetSystemTrayMenu(m)
		log.Println("System tray menu set.")
	} else {
		log.Println("System tray not supported on this platform.")
	}
}

func setupWindowIntercepts() {
	// Intercept close requests
	mainWindow.SetCloseIntercept(func() {
		log.Println("Main window close intercepted.")
		if currentMode == ModeLogging {
			// Optional: Ask for confirmation before cancelling logging?
			// dialog.ShowConfirm("Cancel Log?", "Closing the window will cancel the current log entry. Proceed?", func(confirm bool) {
			// 	if confirm {
			// 		log.Println("Logging cancelled by closing window (confirmed).")
			// 		switchToBrowsingMode() // Switch back first
			// 		mainWindow.Hide()      // Then hide
			// 	} else {
			// 		log.Println("Window close cancelled by user.")
			// 	}
			// }, mainWindow)
			// --- For now, just cancel and hide ---
			log.Println("Window closed during logging. Cancelling log and hiding window.")
			switchToBrowsingMode() // Ensure state is reset
			mainWindow.Hide()
			// ---
		} else {
			log.Println("Hiding window (Browsing Mode).")
			mainWindow.Hide() // Default behavior: hide if tray is supported
		}
	})

	// Fallback if tray isn't supported (already handled by Fyne implicitly, but explicit is okay)
	if _, ok := myApp.(desktop.App); !ok {
		mainWindow.SetCloseIntercept(func() {
			log.Println("Close intercepted (no tray support). Quitting.")
			myApp.Quit()
		})
	}
	log.Println("Window close intercept setup complete.")
}
