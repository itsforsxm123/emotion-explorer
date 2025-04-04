// internal/core/hierarchy.go
package core

import (
	"sort" // Import the sort package

	"github.com/itsforsxm123/emotion-explorer/internal/data" // Adjust import path if needed
)

// GetPrimaryEmotions filters the provided map of emotions and returns a slice
// containing only the primary emotions, sorted alphabetically by name.
// It returns an empty slice if the input map is nil or empty, or if no
// primary emotions are found.
func GetPrimaryEmotions(emotions map[string]data.Emotion) []data.Emotion {
	// Handle nil or empty map gracefully
	if len(emotions) == 0 {
		return []data.Emotion{}
	}

	primaryEmotions := make([]data.Emotion, 0) // Initialize with 0 capacity

	// Iterate through the map of all emotions
	for _, emotion := range emotions {
		// Check if the emotion's type is "primary"
		if emotion.Type == "primary" {
			primaryEmotions = append(primaryEmotions, emotion)
		}
	}

	// Sort the primary emotions alphabetically by name for consistent UI display
	sort.Slice(primaryEmotions, func(i, j int) bool {
		return primaryEmotions[i].Name < primaryEmotions[j].Name
	})

	return primaryEmotions
}

// GetChildrenOf finds all direct children of a given parent emotion ID.
// It searches the provided map of all emotions and returns a slice containing
// the child emotions, sorted alphabetically by name.
// Returns an empty slice if the parentID is not found, if the parent has no
// children, or if the allEmotions map is nil or empty.
func GetChildrenOf(parentID string, allEmotions map[string]data.Emotion) []data.Emotion {
	// Handle nil or empty map gracefully
	if len(allEmotions) == 0 {
		return []data.Emotion{}
	}

	children := make([]data.Emotion, 0) // Initialize slice for children

	// Iterate through all emotions in the map
	for _, emotion := range allEmotions {
		// Check if the emotion's ParentID matches the requested parentID
		if emotion.ParentID == parentID {
			children = append(children, emotion)
		}
	}

	// Sort the children alphabetically by name for consistent UI display
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name < children[j].Name
	})

	return children
}
