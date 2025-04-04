// internal/data/loader_test.go
package data

import (
	"testing" // Import the standard Go testing package
)

// TestLoadEmotions tests the LoadEmotions function.
func TestLoadEmotions(t *testing.T) {
	// Call the function we want to test
	data, err := LoadEmotions()

	// 1. Check for unexpected errors during loading/parsing
	if err != nil {
		// t.Fatalf fails the test immediately and prints the message
		t.Fatalf("LoadEmotions() returned an unexpected error: %v", err)
	}

	// 2. Perform basic sanity checks on the loaded data

	// Check metadata
	if data.Metadata.Version != "1.0" {
		t.Errorf("Expected Metadata.Version '1.0', but got '%s'", data.Metadata.Version)
	}
	if data.Metadata.Source != "Feelings Wheel" {
		t.Errorf("Expected Metadata.Source 'Feelings Wheel', but got '%s'", data.Metadata.Source)
	}

	// Check if primary emotion type exists
	primaryType, ok := data.EmotionTypes["primary"]
	if !ok {
		t.Fatalf("EmotionType 'primary' not found in EmotionTypes map")
	}
	if primaryType.Name != "Primary Emotions" {
		t.Errorf("Expected primary type name 'Primary Emotions', got '%s'", primaryType.Name)
	}

	// Check if a specific primary emotion exists
	happyEmotion, ok := data.Emotions["happy"]
	if !ok {
		t.Fatalf("Emotion 'happy' not found in Emotions map")
	}
	if happyEmotion.Name != "Happy" {
		t.Errorf("Expected happy emotion name 'Happy', got '%s'", happyEmotion.Name)
	}
	if happyEmotion.Type != "primary" {
		t.Errorf("Expected happy emotion type 'primary', got '%s'", happyEmotion.Type)
	}
	if happyEmotion.Color != "#F29727" {
		t.Errorf("Expected happy emotion color '#F29727', got '%s'", happyEmotion.Color)
	}
	if happyEmotion.ParentID != "" { // Primary emotions should have no parent ID
		t.Errorf("Expected happy emotion ParentID to be empty, got '%s'", happyEmotion.ParentID)
	}

	// Check if a specific secondary emotion exists and has the correct parent
	playfulEmotion, ok := data.Emotions["playful"]
	if !ok {
		t.Fatalf("Emotion 'playful' not found in Emotions map")
	}
	if playfulEmotion.Name != "Playful" {
		t.Errorf("Expected playful emotion name 'Playful', got '%s'", playfulEmotion.Name)
	}
	if playfulEmotion.Type != "secondary" {
		t.Errorf("Expected playful emotion type 'secondary', got '%s'", playfulEmotion.Type)
	}
	if playfulEmotion.ParentID != "happy" {
		t.Errorf("Expected playful emotion ParentID 'happy', got '%s'", playfulEmotion.ParentID)
	}

	// Check if a specific tertiary emotion exists and has the correct parent
	arousedEmotion, ok := data.Emotions["aroused"]
	if !ok {
		t.Fatalf("Emotion 'aroused' not found in Emotions map")
	}
	if arousedEmotion.Name != "Aroused" {
		t.Errorf("Expected aroused emotion name 'Aroused', got '%s'", arousedEmotion.Name)
	}
	if arousedEmotion.Type != "tertiary" {
		t.Errorf("Expected aroused emotion type 'tertiary', got '%s'", arousedEmotion.Type)
	}
	if arousedEmotion.ParentID != "playful" {
		t.Errorf("Expected aroused emotion ParentID 'playful', got '%s'", arousedEmotion.ParentID)
	}

	// Optional: Check the total number of emotions loaded (adjust number if your JSON changes)
	// expectedEmotionCount := 136 // Count items in your "emotions" object
	// if len(data.Emotions) != expectedEmotionCount {
	//  t.Errorf("Expected %d emotions, but loaded %d", expectedEmotionCount, len(data.Emotions))
	// }

	// If we reached here without t.Fatalf, the basic checks passed.
	// More specific checks can be added as needed.
	t.Log("LoadEmotions basic checks passed.") // t.Log only shows up when running tests with -v flag
}
