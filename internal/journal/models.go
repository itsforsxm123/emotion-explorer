package journal // <--- THIS MUST BE THE VERY FIRST LINE

import "time"

// LogEntry represents a single recorded emotion instance.
type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	EmotionID   string    `json:"emotion_id"`      // Reference to data.Emotion.ID
	EmotionName string    `json:"emotion_name"`    // Denormalized for easier display
	Notes       string    `json:"notes,omitempty"` // Optional user notes
	// Optional: Intensity int `json:"intensity,omitempty"`
}
