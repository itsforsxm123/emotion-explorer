// internal/ui/views.go
package ui

import (
	"fmt"
	"image/color"
	"log" // For logging button clicks initially

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/itsforsxm123/emotion-explorer/internal/data" // Import our data models
)

// CreatePrimaryEmotionView generates the UI container displaying buttons for each primary emotion.
// It takes a slice of primary emotions and returns a Fyne CanvasObject (the view).
func CreatePrimaryEmotionView(primaryEmotions []data.Emotion, onEmotionSelected func(emotion data.Emotion)) fyne.CanvasObject {
	// Handle empty input gracefully
	if len(primaryEmotions) == 0 {
		log.Println("Warning: CreatePrimaryEmotionView called with no primary emotions.")
		return widget.NewLabel("No primary emotions found.") // Display a message
	}

	items := []fyne.CanvasObject{} // Slice to hold the buttons

	// Iterate through the primary emotions to create a button for each
	for _, emotion := range primaryEmotions {
		// Capture the loop variable for the closure (important!)
		currentEmotion := emotion

		// Create a new button for the emotion
		button := widget.NewButton(currentEmotion.Name, func() {
			// Action to perform when the button is tapped
			log.Printf("Primary Button '%s' (ID: %s) clicked. Triggering callback.\n",
				currentEmotion.Name, currentEmotion.ID)

			// --- CALL THE CALLBACK ---
			// Check if the callback is provided before calling it
			if onEmotionSelected != nil {
				onEmotionSelected(currentEmotion) // Pass the selected emotion
			} else {
				log.Println("Warning: onEmotionSelected callback is nil in CreatePrimaryEmotionView.")
			}
		})

		items = append(items, button)
	}

	// Use GridWrap layout for responsive button arrangement
	gridContainer := container.NewGridWrap(fyne.NewSize(150, 40), items...)

	return gridContainer
}

// parseHexColor converts a hex color string (e.g., "#FF0000") to a color.Color.
// Returns an error if the format is invalid.
func parseHexColor(s string) (color.Color, error) {
	var r, g, b uint8
	var format string

	if len(s) == 0 {
		return color.Black, fmt.Errorf("empty color string") // Default or error
	}

	if s[0] == '#' {
		s = s[1:] // Remove leading '#'
	}

	switch len(s) {
	case 6: // RRGGBB
		format = "%02x%02x%02x"
	case 3: // RGB (shorthand) - Expand to RRGGBB
		format = "%1x%1x%1x" // Read single hex digits
		_, err := fmt.Sscanf(s, format, &r, &g, &b)
		if err != nil {
			return color.Black, fmt.Errorf("invalid shorthand hex color format: %w", err)
		}
		// Expand: e.g., F -> FF, A -> AA
		r = r*16 + r
		g = g*16 + g
		b = b*16 + b
		// Now format as RRGGBB for consistency in return type
		format = "%02x%02x%02x"
		s = fmt.Sprintf("%02x%02x%02x", r, g, b) // Recreate the 6-digit string
	default:
		return color.Black, fmt.Errorf("invalid hex color string length: %d", len(s))
	}

	// Scan the 6-digit hex string
	_, err := fmt.Sscanf(s, format, &r, &g, &b)
	if err != nil {
		return color.Black, fmt.Errorf("invalid hex color format: %w", err)
	}

	return color.NRGBA{R: r, G: g, B: b, A: 255}, nil // Return NRGBA (non-alpha-premultiplied) or RGBA
}

// --- Example of a custom widget for colored buttons (for future reference) ---
/* ... (custom widget code remains unchanged) ... */
// --- End custom widget example ---

// CreateSecondaryEmotionView generates the UI container displaying buttons for secondary emotions
// under a specific parent, along with a header and a back button.
// It now accepts an onSecondaryEmotionSelected callback to handle clicks on secondary emotion buttons.
func CreateSecondaryEmotionView(
	parentEmotion data.Emotion,
	secondaryEmotions []data.Emotion,
	onSecondaryEmotionSelected func(emotion data.Emotion), // <<< MODIFIED: Added callback
	goBack func(),
) fyne.CanvasObject {

	// --- Header ---
	headerLabel := widget.NewLabel(fmt.Sprintf("Exploring: %s", parentEmotion.Name))
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	headerLabel.Alignment = fyne.TextAlignCenter

	// --- Back Button ---
	backButton := widget.NewButton("<- Back to Primary", func() {
		log.Println("Back button clicked.")
		if goBack != nil {
			goBack()
		} else {
			log.Println("Warning: goBack callback is nil in CreateSecondaryEmotionView.")
		}
	})

	// --- Secondary Emotion Buttons ---
	var secondaryItems []fyne.CanvasObject
	if len(secondaryEmotions) == 0 {
		secondaryItems = append(secondaryItems, widget.NewLabel(fmt.Sprintf("No specific sub-emotions listed under %s.", parentEmotion.Name)))
	} else {
		for _, emotion := range secondaryEmotions {
			currentEmotion := emotion // Capture loop variable

			secondaryButton := widget.NewButton(currentEmotion.Name, func() {
				// --- MODIFIED: Call the new callback ---
				log.Printf("Secondary Button '%s' (ID: %s, Parent: %s) clicked. Triggering callback.\n",
					currentEmotion.Name, currentEmotion.ID, parentEmotion.Name)

				// Check if the callback is provided before calling it
				if onSecondaryEmotionSelected != nil {
					onSecondaryEmotionSelected(currentEmotion) // Pass the selected secondary emotion
				} else {
					// Log a warning if the callback is missing (helps debugging)
					log.Println("Warning: onSecondaryEmotionSelected callback is nil in CreateSecondaryEmotionView.")
				}
				// --- End Modification ---
			})
			secondaryItems = append(secondaryItems, secondaryButton)
		}
	}

	secondaryGrid := container.NewGridWrap(fyne.NewSize(140, 35), secondaryItems...)

	// --- Assemble the View ---
	viewLayout := container.NewVBox(
		headerLabel,
		widget.NewSeparator(),
		secondaryGrid,
		widget.NewSeparator(),
		backButton,
	)

	// Optional Border layout remains the same
	// viewLayout := container.NewBorder(...)

	return viewLayout
}

// CreateTertiaryEmotionView generates the UI container displaying buttons for tertiary emotions
// under a specific secondary parent, along with a header and a back button.
func CreateTertiaryEmotionView(
	parentEmotion data.Emotion, // The secondary emotion that is the parent of these tertiary ones
	tertiaryEmotions []data.Emotion,
	goBack func(), // Callback to go back to the Secondary View
) fyne.CanvasObject {

	log.Printf("Creating Tertiary View for parent '%s' with %d children.", parentEmotion.Name, len(tertiaryEmotions))

	// --- Header ---
	headerLabel := widget.NewLabel(fmt.Sprintf("Exploring under: %s", parentEmotion.Name))
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	headerLabel.Alignment = fyne.TextAlignCenter

	// --- Back Button ---
	backButton := widget.NewButton("<- Back to Secondary", func() {
		log.Println("Tertiary View: Back button clicked.")
		if goBack != nil {
			goBack()
		} else {
			log.Println("Warning: goBack callback is nil in CreateTertiaryEmotionView.")
		}
	})

	// --- Tertiary Emotion Buttons --- (Replaces placeholder)
	var tertiaryItems []fyne.CanvasObject
	if len(tertiaryEmotions) == 0 {
		// This case should ideally not be reached due to checks in main.go,
		// but handle defensively.
		tertiaryItems = append(tertiaryItems, widget.NewLabel(fmt.Sprintf("No specific sub-emotions listed under %s.", parentEmotion.Name)))
		log.Printf("Warning: CreateTertiaryEmotionView called for '%s' but received 0 tertiary emotions.", parentEmotion.Name)
	} else {
		// Create buttons for each tertiary emotion
		for _, emotion := range tertiaryEmotions {
			currentEmotion := emotion // Capture loop variable for closure

			tertiaryButton := widget.NewButton(currentEmotion.Name, func() {
				// Action for tertiary emotion button click (just log for now)
				log.Printf("Tertiary Button '%s' (ID: %s, Parent: %s) clicked. No further navigation implemented yet.\n",
					currentEmotion.Name, currentEmotion.ID, parentEmotion.Name)

				// --- TODO: Future - Implement navigation to detail view or handle leaf nodes if hierarchy deepens ---
			})
			tertiaryItems = append(tertiaryItems, tertiaryButton)
		}
	}

	// Use GridWrap for tertiary buttons, similar to secondary
	tertiaryGrid := container.NewGridWrap(fyne.NewSize(130, 35), tertiaryItems...) // Slightly smaller buttons maybe?

	// --- Assemble the View ---
	// Use a VBox to stack the header, tertiary grid, and back button
	viewLayout := container.NewVBox(
		headerLabel,
		widget.NewSeparator(), // Add a visual separator line
		tertiaryGrid,          // Use the grid of buttons
		widget.NewSeparator(), // Another separator
		backButton,
	)

	// Optional: Use Border layout for more control (e.g., back button fixed at bottom)
	// viewLayout := container.NewBorder(
	// 	headerLabel, // Top
	// 	backButton,  // Bottom
	// 	nil,         // Left
	// 	nil,         // Right
	// 	container.NewScroll(tertiaryGrid), // Center (scrollable) - Good idea if many buttons
	// )

	return viewLayout
}
