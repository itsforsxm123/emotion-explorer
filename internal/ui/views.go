// internal/ui/views.go
package ui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/itsforsxm123/emotion-explorer/internal/data" // Use your module path
)

// CreateEmotionListView generates a generic UI container displaying items (tappable cards) for a list of emotions.
// It supports an optional title and an optional parent context.
// Back navigation is now handled globally by the main application structure.
func CreateEmotionListView(
	title string, // Title for the header (empty string for no header)
	parent *data.Emotion, // Optional parent context (can be nil)
	emotions []data.Emotion, // The list of emotions to display
	onSelected func(selectedEmotion data.Emotion), // Callback when an item is clicked
	// --- REMOVED goBack func() ---
	// --- REMOVED backButtonLabel string ---
) fyne.CanvasObject {

	// --- Log statement updated ---
	log.Printf("Creating generic list view: Title='%s', Parent='%v', #Emotions=%d",
		title, parent, len(emotions))

	// --- Content Items (Tappable Cards or Message) ---
	var contentItems []fyne.CanvasObject
	if len(emotions) == 0 {
		message := "No emotions found."
		if parent != nil {
			message = fmt.Sprintf("No specific sub-emotions listed under %s.", parent.Name)
		}
		contentItems = append(contentItems, widget.NewLabel(message))
		log.Printf("Warning: CreateEmotionListView called with 0 emotions for parent '%v'.", parent)
	} else {
		// Create card items for each emotion
		for _, emotion := range emotions {
			// --- Card Implementation (No changes needed here) ---
			currentEmotion := emotion // Capture loop variable
			emotionColor, err := parseHexColor(currentEmotion.Color)
			if err != nil {
				log.Printf("Warning: Failed to parse color '%s' for emotion '%s': %v. Using default.",
					currentEmotion.Color, currentEmotion.Name, err)
				emotionColor = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
			}
			colorSwatch := canvas.NewRectangle(emotionColor)
			swatchSize := float32(20)
			colorSwatch.SetMinSize(fyne.NewSize(swatchSize, swatchSize))
			nameLabel := widget.NewLabel(currentEmotion.Name)
			nameLabel.Alignment = fyne.TextAlignLeading
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			cardContent := container.NewHBox(
				colorSwatch,
				layout.NewSpacer(),
				nameLabel,
				layout.NewSpacer(),
			)
			cardVisual := widget.NewCard("", "", container.NewPadded(cardContent))
			tapAction := func() {
				parentName := "N/A"
				if parent != nil {
					parentName = parent.Name
				}
				log.Printf("Card '%s' (ID: %s, Parent: %s) clicked via TappableCard. Triggering onSelected callback.\n",
					currentEmotion.Name, currentEmotion.ID, parentName)
				if onSelected != nil {
					onSelected(currentEmotion)
				} else {
					log.Println("Warning: onSelected callback is nil in CreateEmotionListView.")
				}
			}
			tappableWrapper := NewTappableCard(cardVisual, tapAction)
			contentItems = append(contentItems, tappableWrapper)
			// --- Card Implementation End ---
		}
	}

	// Use GridWrap layout for the emotion cards
	contentGrid := container.NewGridWrap(fyne.NewSize(200, 60), contentItems...) // Adjust size as needed

	// --- Assemble the View (Header/Content) ---
	topItems := []fyne.CanvasObject{}
	// --- REMOVED bottomItems declaration ---

	// Add Header if title is provided
	if title != "" {
		headerLabel := widget.NewLabel(title)
		headerLabel.TextStyle = fyne.TextStyle{Bold: true}
		headerLabel.Alignment = fyne.TextAlignCenter
		topItems = append(topItems, headerLabel, widget.NewSeparator())
	}

	// --- REMOVED Back Button Logic ---
	// if goBack != nil { ... } block is removed

	// Use a Border layout for structure: Header (Top), Content (Center)
	// Wrap the content grid in a Scroll container
	viewLayout := container.NewBorder(
		container.NewVBox(topItems...), // Top: Header and separator (if any)
		// --- REMOVED Bottom parameter (was container.NewVBox(bottomItems...)) ---
		nil,                              // Bottom: Nothing here now
		nil,                              // Left
		nil,                              // Right
		container.NewScroll(contentGrid), // Center: Scrollable grid of emotion cards
	)

	return viewLayout
}

// parseHexColor function remains unchanged
func parseHexColor(s string) (color.Color, error) {
	// ... (implementation is the same) ...
	var r, g, b uint8
	var format string

	if len(s) == 0 {
		return color.Black, fmt.Errorf("empty color string")
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
