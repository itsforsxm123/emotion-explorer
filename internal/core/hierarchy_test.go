// internal/core/hierarchy_test.go
package core_test // Use _test package for black-box testing

import (
	"testing"

	// Import the package we are testing
	core "github.com/itsforsxm123/emotion-explorer/internal/core"
	// Import the data package for the Emotion struct
	"github.com/itsforsxm123/emotion-explorer/internal/data"
	// Import testify/assert for readable assertions
	"github.com/stretchr/testify/assert"
)

// TestGetPrimaryEmotions tests the GetPrimaryEmotions function with various scenarios.
func TestGetPrimaryEmotions(t *testing.T) {

	// --- Test Data Setup ---

	// Define some sample emotions for testing
	emotionJoy := data.Emotion{ID: "joy", Name: "Joy", Type: "primary", Color: "#FFD700"}
	emotionSadness := data.Emotion{ID: "sadness", Name: "Sadness", Type: "primary", Color: "#ADD8E6"}
	emotionAnger := data.Emotion{ID: "anger", Name: "Anger", Type: "primary", Color: "#FF0000"}
	emotionFear := data.Emotion{ID: "fear", Name: "Fear", Type: "primary", Color: "#800080"}

	emotionContentment := data.Emotion{ID: "contentment", Name: "Contentment", Type: "secondary", Color: "#FFFFE0", ParentID: "joy"}
	emotionGrief := data.Emotion{ID: "grief", Name: "Grief", Type: "secondary", Color: "#A9A9A9", ParentID: "sadness"}
	emotionRage := data.Emotion{ID: "rage", Name: "Rage", Type: "tertiary", Color: "#DC143C", ParentID: "anger"} // Assuming tertiary exists or just different type

	// --- Test Cases ---

	testCases := []struct {
		name           string                  // Name of the test case
		inputEmotions  map[string]data.Emotion // Input map for GetPrimaryEmotions
		expectedOutput []data.Emotion          // Expected slice of primary emotions (must be sorted by Name)
	}{
		{
			name: "Happy Path - Mixed Emotions",
			inputEmotions: map[string]data.Emotion{
				"joy":         emotionJoy,
				"sadness":     emotionSadness,
				"anger":       emotionAnger,
				"contentment": emotionContentment, // Secondary
				"grief":       emotionGrief,       // Secondary
				"rage":        emotionRage,        // Tertiary/Other
			},
			// Expected output should only contain primary emotions, sorted alphabetically by Name
			expectedOutput: []data.Emotion{
				emotionAnger, // Anger comes before Joy
				emotionJoy,
				emotionSadness,
			},
		},
		{
			name: "Edge Case - Only Primary Emotions",
			inputEmotions: map[string]data.Emotion{
				"fear":    emotionFear,
				"joy":     emotionJoy,
				"sadness": emotionSadness,
				"anger":   emotionAnger,
			},
			// Expected output should be all input emotions, sorted alphabetically by Name
			expectedOutput: []data.Emotion{
				emotionAnger,
				emotionFear,
				emotionJoy,
				emotionSadness,
			},
		},
		{
			name: "Edge Case - No Primary Emotions",
			inputEmotions: map[string]data.Emotion{
				"contentment": emotionContentment, // Secondary
				"grief":       emotionGrief,       // Secondary
				"rage":        emotionRage,        // Tertiary/Other
			},
			// Expected output should be an empty slice
			expectedOutput: []data.Emotion{},
		},
		{
			name:           "Edge Case - Empty Input Map",
			inputEmotions:  map[string]data.Emotion{}, // Empty map
			expectedOutput: []data.Emotion{},          // Expect empty slice
		},
		{
			name:           "Edge Case - Nil Input Map",
			inputEmotions:  nil,              // Nil map
			expectedOutput: []data.Emotion{}, // Expect empty slice
		},
	}

	// --- Run Test Cases ---

	for _, tc := range testCases {
		// Run each test case as a sub-test
		t.Run(tc.name, func(t *testing.T) {
			// Call the function under test
			actualOutput := core.GetPrimaryEmotions(tc.inputEmotions)

			// --- Assertions ---
			// Check if the actual output matches the expected output.
			// assert.Equal checks for equality of type, length, capacity, and element values in order.
			// This works perfectly because GetPrimaryEmotions guarantees a sorted slice.
			assert.Equal(t, tc.expectedOutput, actualOutput)

			// Optional: Add a specific check for length if needed, though assert.Equal covers it.
			assert.Len(t, actualOutput, len(tc.expectedOutput), "Length of output slice should match expected")

			// Optional: Check if the output is actually sorted (as a sanity check, though covered by Equal)
			// This requires Go 1.21+ for slices.IsSortedFunc
			// isSorted := slices.IsSortedFunc(actualOutput, func(a, b data.Emotion) int {
			// 	return cmp.Compare(a.Name, b.Name)
			// })
			// assert.True(t, isSorted, "Output slice should be sorted by name")
			// If using older Go, implement manual sort check or rely on assert.Equal with pre-sorted expectedOutput.

		})
	}
}

// TestGetChildrenOf tests the GetChildrenOf function.
func TestGetChildrenOf(t *testing.T) {

	// --- Test Data Setup ---
	// Re-use some emotions and add more specific parent/child relationships

	// Primary
	emotionJoy := data.Emotion{ID: "joy", Name: "Joy", Type: "primary", Color: "#FFD700"}
	emotionSadness := data.Emotion{ID: "sadness", Name: "Sadness", Type: "primary", Color: "#ADD8E6"}
	emotionAnger := data.Emotion{ID: "anger", Name: "Anger", Type: "primary", Color: "#FF0000"}
	emotionFear := data.Emotion{ID: "fear", Name: "Fear", Type: "primary", Color: "#800080"} // No children in this test data

	// Secondary (Children of Joy)
	emotionContentment := data.Emotion{ID: "contentment", Name: "Contentment", Type: "secondary", Color: "#FFFFE0", ParentID: "joy"}
	emotionOptimism := data.Emotion{ID: "optimism", Name: "Optimism", Type: "secondary", Color: "#FFFF00", ParentID: "joy"}
	emotionZest := data.Emotion{ID: "zest", Name: "Zest", Type: "secondary", Color: "#FFA500", ParentID: "joy"}

	// Secondary (Children of Sadness)
	emotionGrief := data.Emotion{ID: "grief", Name: "Grief", Type: "secondary", Color: "#A9A9A9", ParentID: "sadness"}
	emotionDisappointment := data.Emotion{ID: "disappointment", Name: "Disappointment", Type: "secondary", Color: "#808080", ParentID: "sadness"}

	// Tertiary (Child of Grief - to ensure only direct children are found)
	emotionSorrow := data.Emotion{ID: "sorrow", Name: "Sorrow", Type: "tertiary", Color: "#696969", ParentID: "grief"}

	// Map representing all emotions in our test dataset
	allTestEmotions := map[string]data.Emotion{
		// Primary
		"joy":     emotionJoy,
		"sadness": emotionSadness,
		"anger":   emotionAnger, // Has no children defined here
		"fear":    emotionFear,  // Has no children defined here
		// Secondary (Joy)
		"contentment": emotionContentment,
		"optimism":    emotionOptimism,
		"zest":        emotionZest,
		// Secondary (Sadness)
		"grief":          emotionGrief,
		"disappointment": emotionDisappointment,
		// Tertiary (under Grief)
		"sorrow": emotionSorrow,
	}

	// --- Test Cases ---
	testCases := []struct {
		name             string                  // Name of the test case
		parentID         string                  // Input: Parent ID to find children for
		inputAllEmotions map[string]data.Emotion // Input: Map of all emotions
		expectedOutput   []data.Emotion          // Expected slice of direct children (sorted by Name)
	}{
		{
			name:             "Parent with multiple children (Joy)",
			parentID:         "joy",
			inputAllEmotions: allTestEmotions,
			// Expected: Children of 'joy', sorted alphabetically by name
			expectedOutput: []data.Emotion{
				emotionContentment, // C
				emotionOptimism,    // O
				emotionZest,        // Z
			},
		},
		{
			name:             "Parent with multiple children (Sadness)",
			parentID:         "sadness",
			inputAllEmotions: allTestEmotions,
			// Expected: Children of 'sadness', sorted alphabetically by name
			expectedOutput: []data.Emotion{
				emotionDisappointment, // D
				emotionGrief,          // G
			},
		},
		{
			name:             "Parent with no direct children (Anger)",
			parentID:         "anger", // Anger exists but has no children in the map
			inputAllEmotions: allTestEmotions,
			expectedOutput:   []data.Emotion{}, // Expect empty slice
		},
		{
			name:             "Parent with children who also have children (Grief)",
			parentID:         "grief", // Grief is a child of Sadness, and parent of Sorrow
			inputAllEmotions: allTestEmotions,
			// Expected: Only direct children of 'grief' (Sorrow)
			expectedOutput: []data.Emotion{
				emotionSorrow,
			},
		},
		{
			name:             "Parent ID does not exist",
			parentID:         "nonexistent_id",
			inputAllEmotions: allTestEmotions,
			expectedOutput:   []data.Emotion{}, // Expect empty slice
		},
		{
			name:             "Empty input map",
			parentID:         "joy", // Parent ID doesn't matter if map is empty
			inputAllEmotions: map[string]data.Emotion{},
			expectedOutput:   []data.Emotion{}, // Expect empty slice
		},
		{
			name:             "Nil input map",
			parentID:         "joy", // Parent ID doesn't matter if map is nil
			inputAllEmotions: nil,
			expectedOutput:   []data.Emotion{}, // Expect empty slice
		},
	}

	// --- Run Test Cases ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function under test
			actualOutput := core.GetChildrenOf(tc.parentID, tc.inputAllEmotions)

			// --- Assertions ---
			// Check if the actual output slice matches the expected output slice.
			// This checks length, order, and element values.
			assert.Equal(t, tc.expectedOutput, actualOutput)

			// Optional: Explicit length check (covered by assert.Equal)
			assert.Len(t, actualOutput, len(tc.expectedOutput), "Length of output slice should match expected")
		})
	}
}
