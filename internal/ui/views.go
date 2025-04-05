// internal/ui/views.go
package ui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout" // Import for layout spacers etc.
	"fyne.io/fyne/v2/widget"

	"github.com/itsforsxm123/emotion-explorer/internal/data"
)

// CreateEmotionListView generates a generic UI container displaying buttons for a list of emotions.
// It supports an optional title, an optional parent context (for logging/future use),
// and callbacks for selection and back navigation.
func CreateEmotionListView(
	title string, // Title for the header (empty string for no header)
	parent *data.Emotion, // Optional parent context (can be nil)
	emotions []data.Emotion, // The list of emotions to display buttons for
	onSelected func(selectedEmotion data.Emotion), // Callback when an emotion button is clicked
	goBack func(), // Callback for the back button (if nil, no back button is shown)
	backButtonLabel string, // Label for the back button (e.g., "<- Back to Primary")
) fyne.CanvasObject {

	log.Printf("Creating generic list view: Title='%s', Parent='%v', #Emotions=%d, HasGoBack=%t",
		title, parent, len(emotions), goBack != nil)

	// --- Content Items (Buttons or Message) ---
	var contentItems []fyne.CanvasObject
	if len(emotions) == 0 {
		message := "No emotions found."
		if parent != nil {
			message = fmt.Sprintf("No specific sub-emotions listed under %s.", parent.Name)
		}
		contentItems = append(contentItems, widget.NewLabel(message))
		log.Printf("Warning: CreateEmotionListView called with 0 emotions for parent '%v'.", parent)
	} else {
		// Create buttons for each emotion
		for _, emotion := range emotions {
			currentEmotion := emotion // Capture loop variable

			button := widget.NewButton(currentEmotion.Name, func() {
				parentName := "N/A"
				if parent != nil {
					parentName = parent.Name
				}
				log.Printf("Button '%s' (ID: %s, Parent: %s) clicked. Triggering onSelected callback.\n",
					currentEmotion.Name, currentEmotion.ID, parentName)

				// Call the provided selection callback
				if onSelected != nil {
					onSelected(currentEmotion)
				} else {
					log.Println("Warning: onSelected callback is nil in CreateEmotionListView.")
				}
			})
			contentItems = append(contentItems, button)
		}
	}

	// Use GridWrap layout for the emotion buttons
	// Consistent sizing for all levels now
	contentGrid := container.NewGridWrap(fyne.NewSize(140, 35), contentItems...)

	// --- Assemble the View ---
	topItems := []fyne.CanvasObject{}
	bottomItems := []fyne.CanvasObject{}

	// Add Header if title is provided
	if title != "" {
		headerLabel := widget.NewLabel(title)
		headerLabel.TextStyle = fyne.TextStyle{Bold: true}
		headerLabel.Alignment = fyne.TextAlignCenter
		topItems = append(topItems, headerLabel, widget.NewSeparator())
	}

	// Add Back Button if goBack callback is provided
	if goBack != nil {
		backButton := widget.NewButton(backButtonLabel, func() {
			log.Println("Generic View: Back button clicked.")
			// Call the provided goBack callback
			goBack()
			// No nil check needed here, as we only add the button if goBack is not nil
		})
		// Add some space above the back button
		bottomItems = append(bottomItems, layout.NewSpacer(), widget.NewSeparator(), backButton)
	}

	// Use a Border layout for structure: Header (Top), Content (Center), Back Button (Bottom)
	viewLayout := container.NewBorder(
		container.NewVBox(topItems...),    // Top: Header and separator (if any)
		container.NewVBox(bottomItems...), // Bottom: Separator and Back button (if any)
		nil,                               // Left
		nil,                               // Right
		container.NewScroll(contentGrid),  // Center: Scrollable grid of emotion buttons
	)

	return viewLayout
}

// parseHexColor converts a hex color string (e.g., "#FF0000") to a color.Color.
// Returns an error if the format is invalid.
// (Keep this function as it's independent)
func parseHexColor(s string) (color.Color, error) {
	var r, g, b uint8
	var format string

	if len(s) == 0 {
		return color.Black, fmt.Errorf("empty color string")
	}

	if s[0] == '#' {
		s = s[1:]
	}

	switch len(s) {
	case 6:
		format = "%02x%02x%02x"
	case 3:
		format = "%1x%1x%1x"
		_, err := fmt.Sscanf(s, format, &r, &g, &b)
		if err != nil {
			return color.Black, fmt.Errorf("invalid shorthand hex color format: %w", err)
		}
		r = r*16 + r
		g = g*16 + g
		b = b*16 + b
		format = "%02x%02x%02x"
		s = fmt.Sprintf("%02x%02x%02x", r, g, b)
	default:
		return color.Black, fmt.Errorf("invalid hex color string length: %d", len(s))
	}

	_, err := fmt.Sscanf(s, format, &r, &g, &b)
	if err != nil {
		return color.Black, fmt.Errorf("invalid hex color format: %w", err)
	}

	return color.NRGBA{R: r, G: g, B: b, A: 255}, nil
}

// --- REMOVED CreatePrimaryEmotionView ---
// --- REMOVED CreateSecondaryEmotionView ---
// --- REMOVED CreateTertiaryEmotionView ---
