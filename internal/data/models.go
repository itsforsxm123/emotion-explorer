// internal/data/models.go
package data

// EmotionData represents the entire structure of the emotions.json file.
type EmotionData struct {
	Metadata     Metadata               `json:"metadata"`
	EmotionTypes map[string]EmotionType `json:"emotionTypes"` // Map key is the type ID (e.g., "primary")
	Emotions     map[string]Emotion     `json:"emotions"`     // Map key is the emotion ID (e.g., "happy")
}

// Metadata holds information about the dataset version and source.
type Metadata struct {
	Version     string `json:"version"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

// EmotionType defines the characteristics of an emotion level (primary, secondary, etc.).
type EmotionType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

// Emotion represents a single emotion with its properties and relationship.
type Emotion struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`               // Corresponds to an EmotionType ID (e.g., "primary")
	Color    string `json:"color"`              // Hex color code
	ParentID string `json:"parentId,omitempty"` // Use omitempty as primary emotions won't have this
	// We can add fields here later if needed, e.g., to hold child emotions after processing
	// Children []*Emotion `json:"-"` // Ignored by JSON marshalling/unmarshalling
}
