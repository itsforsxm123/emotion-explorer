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
/*
type coloredButtonRenderer struct {
	fyne.WidgetRenderer
	background *canvas.Rectangle
	button     *widget.Button
}

func (r *coloredButtonRenderer) Refresh() {
	// Update background color if needed
	r.WidgetRenderer.Refresh()
}

type ColoredButton struct {
	widget.BaseWidget
	button *widget.Button
	color  color.Color
}

func NewColoredButton(text string, bgColor color.Color, tapped func()) *ColoredButton {
	cb := &ColoredButton{
		color: bgColor,
	}
	cb.ExtendBaseWidget(cb) // Important for custom widgets
	// Note: Standard button doesn't easily accept background color.
	// We'd likely embed or customize button logic here.
	// This is a placeholder concept. Fyne might require theme adjustments
	// or drawing primitives directly for full background color control.
	cb.button = widget.NewButton(text, tapped) // Use standard button internally for now
	return cb
}

func (cb *ColoredButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(cb.color)
	// Problem: How to layer button visuals ON TOP of bg?
	// May need fyne.Container with bg and button inside.
	// Or, draw text/icon manually on a colored rectangle.
	baseRenderer := cb.button.CreateRenderer() // Get standard button renderer

	// This is complex - standard renderer doesn't expose background easily.
	// A container approach is more typical:
	// container := container.NewMax(bg, cb.button) // Button overlays rectangle
	// return widget.NewSimpleRenderer(container)

	// For now, returning standard button renderer. Color won't apply yet.
	return baseRenderer
}
*/
// --- End custom widget example ---

// CreateSecondaryEmotionView generates the UI container displaying buttons for secondary emotions
// under a specific parent, along with a header and a back button.
func CreateSecondaryEmotionView(parentEmotion data.Emotion, secondaryEmotions []data.Emotion, goBack func()) fyne.CanvasObject {

	// --- Header ---
	// Display the name of the parent emotion
	headerLabel := widget.NewLabel(fmt.Sprintf("Exploring: %s", parentEmotion.Name))
	headerLabel.TextStyle = fyne.TextStyle{Bold: true} // Make it bold
	headerLabel.Alignment = fyne.TextAlignCenter       // Center align

	// --- Back Button ---
	backButton := widget.NewButton("<- Back to Primary", func() {
		log.Println("Back button clicked.")
		if goBack != nil {
			goBack() // Call the provided callback function
		} else {
			log.Println("Warning: goBack callback is nil in CreateSecondaryEmotionView.")
		}
	})

	// --- Secondary Emotion Buttons ---
	var secondaryItems []fyne.CanvasObject
	if len(secondaryEmotions) == 0 {
		// Display a message if no secondary emotions exist for this parent
		secondaryItems = append(secondaryItems, widget.NewLabel(fmt.Sprintf("No specific sub-emotions listed under %s.", parentEmotion.Name)))
	} else {
		// Create buttons for each secondary emotion
		for _, emotion := range secondaryEmotions {
			currentEmotion := emotion // Capture loop variable for closure

			secondaryButton := widget.NewButton(currentEmotion.Name, func() {
				// Action for secondary emotion button click
				log.Printf("Secondary Button '%s' (ID: %s, Parent: %s) clicked.\n",
					currentEmotion.Name, currentEmotion.ID, parentEmotion.Name)

				// --- TODO: Implement navigation to tertiary or detail view ---
			})
			secondaryItems = append(secondaryItems, secondaryButton)
		}
	}

	// Use GridWrap for secondary buttons as well, maybe with a slightly different size
	secondaryGrid := container.NewGridWrap(fyne.NewSize(140, 35), secondaryItems...)

	// --- Assemble the View ---
	// Use a VBox (Vertical Box) to stack the header, secondary grid, and back button
	viewLayout := container.NewVBox(
		headerLabel,
		widget.NewSeparator(), // Add a visual separator line
		secondaryGrid,
		widget.NewSeparator(), // Another separator
		backButton,
	)

	// Optional: Use Border layout for more control (e.g., back button fixed at bottom)
	// viewLayout := container.NewBorder(
	// 	headerLabel, // Top
	// 	backButton,  // Bottom
	// 	nil,         // Left
	// 	nil,         // Right
	// 	container.NewScroll(secondaryGrid), // Center (scrollable)
	// )

	return viewLayout
}
