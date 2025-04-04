// internal/data/loader.go
package data

import (
	"embed" // Required for embedding files
	"encoding/json"
	"fmt" // For formatting error messages
)

// Embed the JSON file directly from the current directory.
// Ensure exactly one space after 'embed'.
//
//go:embed emotions.json
var embeddedJSON embed.FS

// LoadEmotions reads and parses the embedded emotions.json file.
func LoadEmotions() (EmotionData, error) {
	var emotionData EmotionData

	// Read the file by its base name from the embed FS.
	bytes, err := embeddedJSON.ReadFile("emotions.json")
	if err != nil {
		return EmotionData{}, fmt.Errorf("failed to read embedded file 'emotions.json': %w", err)
	}

	err = json.Unmarshal(bytes, &emotionData)
	if err != nil {
		return EmotionData{}, fmt.Errorf("failed to unmarshal emotions.json: %w", err)
	}

	return emotionData, nil
}

// --- Optional Helper Functions (We can add these later as needed) ---

// // GetPrimaryEmotions filters and returns only the primary emotions from the loaded data.
// func GetPrimaryEmotions(data EmotionData) []Emotion {
// 	var primaries []Emotion
// 	for _, emotion := range data.Emotions {
// 		if emotion.Type == data.EmotionTypes["primary"].ID { // Assuming "primary" is the ID in emotionTypes
// 			primaries = append(primaries, emotion)
// 		}
// 	}
// 	// Optionally sort them alphabetically or by some other criteria
// 	// sort.Slice(primaries, func(i, j int) bool {
// 	// 	return primaries[i].Name < primaries[j].Name
// 	// })
// 	return primaries
// }

// // GetChildrenOf returns secondary/tertiary emotions that have the given emotionID as their parent.
// func GetChildrenOf(data EmotionData, emotionID string) []Emotion {
// 	var children []Emotion
// 	for _, emotion := range data.Emotions {
// 		if emotion.ParentID == emotionID {
// 			children = append(children, emotion)
// 		}
// 	}
// 	// Optionally sort children
// 	// sort.Slice(children, func(i, j int) bool {
// 	// 	return children[i].Name < children[j].Name
// 	// })
// 	return children
// }
