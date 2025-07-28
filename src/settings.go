package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
)

// Settings holds configuration for the recall application
type Settings struct {
	Editor string // Preferred editor (nano, vim, etc.)
}

// Default settings
func defaultSettings() *Settings {
	return &Settings{
		Editor: "nano",
	}
}

func loadSettings() *Settings {
	// Load settings from ~/.recall/settings.yaml
	homeDir, _ := os.UserHomeDir()
	settingsFile := homeDir + "/.recall/settings.yaml"
	// Try loading settings from the file
	// If file doesn't exist, return default settings
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		// Warn user to run `recall --init-global` to create settings
		fmt.Println("Settings file not found. Please run `recall --init-global` to create default settings.")
		fmt.Println("Using default settings.")
		fmt.Println("")

		// Return default settings
		return defaultSettings()
	}
	// Load settings from the file
	file, err := os.Open(settingsFile)
	if err != nil {
		fmt.Println("Error loading settings:", err)
		return defaultSettings()
	}
	defer file.Close()

	var settings Settings
	if err := yaml.NewDecoder(file).Decode(&settings); err != nil {
		fmt.Println("Error decoding settings:", err)
		return defaultSettings()
	}

	// Return loaded settings
	return &settings
}
