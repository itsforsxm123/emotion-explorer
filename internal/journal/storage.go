package journal // <-- Make sure this line is present

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync" // To prevent race conditions if called rapidly
	"time"
	// Import your data models if needed here, e.g.:
	// "github.com/itsforsxm123/emotion-explorer/internal/data"
)

const journalFilename = "journal.json"

var journalFilePath string  // Full path to the journal file
var journalMutex sync.Mutex // Mutex to protect file access

// init function to determine journal file path
func init() {
	// For simplicity now, place it next to the executable or in CWD
	// TODO: Use os.UserConfigDir() for a better location in the future
	// Example:
	// configDir, err := os.UserConfigDir()
	// if err == nil {
	//     journalDir := filepath.Join(configDir, "EmotionExplorer")
	//     if err := os.MkdirAll(journalDir, 0750); err == nil { // Ensure dir exists
	//         journalFilePath = filepath.Join(journalDir, journalFilename)
	//     } else {
	//         log.Printf("Warning: Could not create config directory '%s': %v. Using CWD.", journalDir, err)
	//     }
	// } else {
	//     log.Printf("Warning: Could not get user config directory: %v. Using CWD.", err)
	// }

	// Fallback to CWD if path wasn't set above
	// if journalFilePath == "" {
	cwd, err := os.Getwd() // Get current working directory
	if err != nil {
		log.Printf("Warning: Could not get current working directory for journal file: %v. Using filename only.", err)
		journalFilePath = journalFilename // Fallback
	} else {
		journalFilePath = filepath.Join(cwd, journalFilename)
	}
	// }
	log.Printf("Journal file path set to: %s", journalFilePath)
}

// loadJournalEntries reads the journal file and returns the list of entries.
// Returns an empty slice if the file doesn't exist or is empty/invalid.
func loadJournalEntries() ([]LogEntry, error) {
	journalMutex.Lock()         // Lock before reading
	defer journalMutex.Unlock() // Ensure unlock

	data, err := os.ReadFile(journalFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Journal file '%s' not found, starting fresh.", journalFilePath)
			return []LogEntry{}, nil // No file is not an error, just means no entries yet
		}
		log.Printf("Error reading journal file '%s': %v", journalFilePath, err)
		return nil, fmt.Errorf("reading journal file: %w", err) // Wrap error
	}

	if len(data) == 0 {
		log.Println("Journal file is empty, starting fresh.")
		return []LogEntry{}, nil // Empty file is okay
	}

	var entries []LogEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		log.Printf("Error unmarshalling journal JSON from '%s': %v", journalFilePath, err)
		// Return error to signal corruption
		return nil, fmt.Errorf("unmarshalling journal json: %w", err) // Wrap error
	}
	log.Printf("Loaded %d entries from journal file '%s'", len(entries), journalFilePath)
	return entries, nil
}

// SaveLogEntry appends a new entry to the journal file.
// It loads existing entries, appends the new one, and writes back.
func SaveLogEntry(newEntry LogEntry) error {
	journalMutex.Lock()         // Lock for the entire load-append-save operation
	defer journalMutex.Unlock() // Ensure unlock happens even on error/panic

	log.Printf("Attempting to save log entry: Emotion='%s', Time='%s'", newEntry.EmotionName, newEntry.Timestamp.Format(time.RFC3339))

	// --- Load existing ---
	// Note: loadJournalEntries already handles locking internally for the read,
	// but we need the lock around the whole process to prevent race conditions
	// between reading and writing back. We could refactor load to not lock
	// if it's only called from SaveLogEntry which already holds the lock.
	// For now, this nested locking is functionally okay, though slightly less efficient.

	var entries []LogEntry // Declare entries here

	// Read the file content directly within the main lock
	rawData, readErr := os.ReadFile(journalFilePath)
	if readErr != nil && !os.IsNotExist(readErr) {
		log.Printf("Error reading journal file '%s' before save: %v", journalFilePath, readErr)
		return fmt.Errorf("reading journal file before save: %w", readErr)
	}

	// Unmarshal if data exists
	if readErr == nil && len(rawData) > 0 {
		unmarshalErr := json.Unmarshal(rawData, &entries)
		if unmarshalErr != nil {
			log.Printf("Error unmarshalling existing journal JSON from '%s': %v. Starting fresh for this save.", journalFilePath, unmarshalErr)
			// Decide recovery strategy: Here we overwrite corrupted data.
			// Alternatively, could return error: return fmt.Errorf("unmarshalling existing journal: %w", unmarshalErr)
			entries = []LogEntry{} // Reset to empty if corrupt
		} else {
			log.Printf("Loaded %d existing entries from journal file '%s' for saving.", len(entries), journalFilePath)
		}
	} else {
		log.Println("Journal file empty or not found, initializing new entry list.")
		entries = []LogEntry{} // Ensure entries is an empty slice if file didn't exist or was empty
	}

	// --- Append the new entry ---
	entries = append(entries, newEntry)

	// --- Marshal the updated list back to JSON ---
	updatedData, marshalErr := json.MarshalIndent(entries, "", "  ") // Indent with 2 spaces
	if marshalErr != nil {
		log.Printf("Error marshalling updated journal entries to JSON: %v", marshalErr)
		return fmt.Errorf("marshalling updated journal: %w", marshalErr)
	}

	// --- Ensure the directory exists (important if using os.UserConfigDir) ---
	// dir := filepath.Dir(journalFilePath)
	// if err := os.MkdirAll(dir, 0750); err != nil {
	//  log.Printf("Error creating directory '%s': %v", dir, err)
	//  return fmt.Errorf("creating journal directory: %w", err)
	// }

	// --- Write the updated data back to the file (overwrite) ---
	// Use 0644 permissions (owner read/write, group/other read)
	writeErr := os.WriteFile(journalFilePath, updatedData, 0644)
	if writeErr != nil {
		log.Printf("Error writing updated journal file '%s': %v", journalFilePath, writeErr)
		return fmt.Errorf("writing updated journal file: %w", writeErr)
	}

	log.Printf("Successfully saved log entry. Total entries now: %d", len(entries))
	return nil
}

// Add a function to load entries for potential display later
// GetJournalEntries provides safe access to the loaded entries.
func GetJournalEntries() ([]LogEntry, error) {
	// loadJournalEntries handles locking internally
	return loadJournalEntries()
}
