// internal/ui/views.go
package ui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas" // Needed for Rectangle
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/itsforsxm123/emotion-explorer/internal/data"
)

// CreateEmotionListView generates a generic UI container displaying items (now tappable cards) for a list of emotions.
// It supports an optional title, an optional parent context,
// and callbacks for selection and back navigation.
func CreateEmotionListView(
	title string, // Title for the header (empty string for no header)
	parent *data.Emotion, // Optional parent context (can be nil)
	emotions []data.Emotion, // The list of emotions to display
	onSelected func(selectedEmotion data.Emotion), // Callback when an item is clicked
	goBack func(), // Callback for the back button (if nil, no back button is shown)
	backButtonLabel string, // Label for the back button
) fyne.CanvasObject {

	log.Printf("Creating generic list view: Title='%s', Parent='%v', #Emotions=%d, HasGoBack=%t",
		title, parent, len(emotions), goBack != nil)

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
			currentEmotion := emotion // Capture loop variable

			// --- Card Implementation Start ---

			// 1. Parse the color (handle errors)
			emotionColor, err := parseHexColor(currentEmotion.Color)
			if err != nil {
				log.Printf("Warning: Failed to parse color '%s' for emotion '%s': %v. Using default.",
					currentEmotion.Color, currentEmotion.Name, err)
				// Use a default color like gray if parsing fails
				emotionColor = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
			}

			// 2. Create a small rectangle for the color swatch
			colorSwatch := canvas.NewRectangle(emotionColor)
			// Define swatch size
			swatchSize := float32(20) // Adjust size as needed
			colorSwatch.SetMinSize(fyne.NewSize(swatchSize, swatchSize))

			// 3. Create the label for the emotion name
			nameLabel := widget.NewLabel(currentEmotion.Name)
			nameLabel.Alignment = fyne.TextAlignLeading      // Align text to the start (left)
			nameLabel.TextStyle = fyne.TextStyle{Bold: true} // Make name bold

			// 4. Arrange swatch and label horizontally using HBox
			// Add spacing for better layout
			cardContent := container.NewHBox(
				colorSwatch,
				layout.NewSpacer(), // Use spacer to push label right slightly or fixed width
				// layout.NewFixedWidthSpacer(theme.Padding()), // Alternative: use theme padding if theme imported
				nameLabel,
				layout.NewSpacer(), // Pushes content left
			)

			// 5. Create the Card visual element
			// Use container.NewPadded for automatic padding around the content
			cardVisual := widget.NewCard(
				"", // No title in the card header itself
				"", // No subtitle
				container.NewPadded(cardContent),
			)

			// 6. Define the tap action closure
			// This captures the currentEmotion for use when the card is tapped
			tapAction := func() {
				parentName := "N/A"
				if parent != nil {
					parentName = parent.Name
				}
				log.Printf("Card '%s' (ID: %s, Parent: %s) clicked via TappableCard. Triggering onSelected callback.\n",
					currentEmotion.Name, currentEmotion.ID, parentName)

				// Call the provided selection callback passed into CreateEmotionListView
				if onSelected != nil {
					onSelected(currentEmotion)
				} else {
					log.Println("Warning: onSelected callback is nil in CreateEmotionListView.")
				}
			}

			// 7. Create the TappableCard using our custom widget (from internal/ui/widgets.go)
			// This wraps the visual card and makes it interactive
			tappableWrapper := NewTappableCard(cardVisual, tapAction)

			// 8. Add the tappable wrapper (which contains the card) to the list of items
			contentItems = append(contentItems, tappableWrapper)

			// --- Card Implementation End ---
		}
	}

	// Use GridWrap layout for the emotion cards
	// Adjusted size for cards - might need tweaking based on content/font
	contentGrid := container.NewGridWrap(fyne.NewSize(200, 60), contentItems...) // Adjust size as needed

	// --- Assemble the View (Header/Footer/Content) ---
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
			goBack() // Call the provided callback
		})
		// Add some space above the back button
		bottomItems = append(bottomItems, layout.NewSpacer(), widget.NewSeparator(), backButton)
	}

	// Use a Border layout for structure: Header (Top), Content (Center), Back Button (Bottom)
	// Wrap the content grid in a Scroll container for cases where items overflow
	viewLayout := container.NewBorder(
		container.NewVBox(topItems...),    // Top: Header and separator (if any)
		container.NewVBox(bottomItems...), // Bottom: Separator and Back button (if any)
		nil,                               // Left
		nil,                               // Right
		container.NewScroll(contentGrid),  // Center: Scrollable grid of emotion cards
	)

	return viewLayout
}

// parseHexColor converts a hex color string (e.g., "#FF0000") to a color.Color.
// Returns an error if the format is invalid.
func parseHexColor(s string) (color.Color, error) {
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
