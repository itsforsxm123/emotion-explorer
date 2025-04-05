// internal/ui/widgets.go
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TappableCard is a simple custom widget that wraps any canvas object
// and makes it respond to tap events.
type TappableCard struct {
	widget.BaseWidget                   // Embed BaseWidget
	content           fyne.CanvasObject // The content to display (e.g., our card)
	onTapped          func()            // The function to call when tapped
}

// NewTappableCard creates a new TappableCard instance.
func NewTappableCard(content fyne.CanvasObject, onTapped func()) *TappableCard {
	tc := &TappableCard{
		content:  content,
		onTapped: onTapped,
	}
	tc.ExtendBaseWidget(tc) // Important: Set up BaseWidget
	return tc
}

// CreateRenderer returns the renderer for this widget.
// For TappableCard, the renderer simply displays the contained content.
func (tc *TappableCard) CreateRenderer() fyne.WidgetRenderer {
	// Use a SimpleRenderer that just renders the content object.
	return widget.NewSimpleRenderer(tc.content)
}

// Tapped is called when the TappableCard receives a tap event.
func (tc *TappableCard) Tapped(_ *fyne.PointEvent) {
	if tc.onTapped != nil {
		tc.onTapped() // Execute the stored callback function
	}
}

// TappedSecondary is called for right-click or alternative tap events (optional).
// We don't need it here, but it's part of the Tappable interface implicitly
// via the functions checked by the event system. We don't need to implement it
// if we only care about primary taps.

// Ensure TappableCard implements the fyne.Tappable interface implicitly
// by having the Tapped method. (No explicit 'implements' keyword in Go).
