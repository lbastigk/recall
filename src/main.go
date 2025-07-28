package main

import (
	"fmt"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
)

// Settings holds configuration for the recall application
type Settings struct {
	Editor       string 	// Preferred editor (nano, vim, etc.)
}

// Default settings
func defaultSettings() *Settings {
	return &Settings{
		Editor:     "nano",
	}
}

func main() {
	// Load settings from ~/.recall/settings.yaml if it exists
	settings := loadSettings()

	args := os.Args[1:] // Skip the program name

	// Handle different argument patterns
	switch len(args) {
	case 0:
		showUsage(settings)
	case 1:
		if args[0] == "--init" {
			initLocal(settings)
		} else if args[0] == "--init-global" {
			initGlobal(settings)
		} else {
			project := args[0]
			keyPath := []string{}
			showKey(settings, project, keyPath)	// Empty keyPath means general project info
		}
	default:
		// Handle --edit command
		if args[0] == "--edit" {
			if len(args) < 2 {
				fmt.Println("Error: --edit requires at least project the project name")
				showUsage(settings)
				os.Exit(1)
			}
			// recall --edit <project> <key> [nested keys...]
			project := args[1]
			keyPath := args[2:]
			editKey(settings, project, keyPath)
		} else {
			// recall <project> <key> [nested keys...]
			project := args[0]
			keyPath := args[1:]
			showKey(settings, project, keyPath)
		}
	}
}

func showUsage(settings *Settings) {
	fmt.Println("recall - A CLI Tool for Project Knowledge Management")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  recall                                Show this help")
	fmt.Println("  recall <project>                      Show general project info")
	fmt.Println("  recall <project> <key>                Show specific key info")
	fmt.Println("  recall <project> <key> <subkey>...    Show nested key info")
	fmt.Println("  recall --edit <project> <key>...      Edit specific key")
	fmt.Println("  recall --init                         Initialize local recall")
	fmt.Println("  recall --init-global                  Initialize global recall")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  recall myApp")
	fmt.Println("  recall myApp database")
	fmt.Println("  recall myApp myClass myFunction")
	fmt.Println("  recall myApp myClass myFunction myVariable")
	fmt.Println("  recall --edit myApp deployment")
	fmt.Println()
	fmt.Printf("Settings: Editor=%s\n", settings.Editor)
}

func initLocal(settings *Settings) {
	fmt.Println("Initializing local recall directory at ./.recall/...")
	// TODO: Create .recall/ directory
	fmt.Println("Created ./.recall/ directory")
}

func initGlobal(settings *Settings) {
	homeDir, _ := os.UserHomeDir()
	globalPath := homeDir + "/.recall"
	fmt.Printf("Initializing global recall directory at %s...\n", globalPath)
	// TODO: Create ~/.recall/ directory and settings.yaml
	fmt.Printf("Created %s directory\n", globalPath)
}

func showKey(settings *Settings, project string, keyPath []string) {
	if len(keyPath) == 0 {
		fmt.Printf("Project: %s (showing general info)\n", project)
	} else {
		path := strings.Join(keyPath, ".")
		fmt.Printf("Project: %s, Key: %s\n", project, path)
	}
	// TODO: Load and display specific nested key from project.yaml
	// Look in ./.recall/ first, then ~/.recall/
}

func editKey(settings *Settings, project string, keyPath []string) {
	if len(keyPath) == 0 {
		fmt.Printf("Editing general info for project: %s\n", project)
	} else {
		path := strings.Join(keyPath, ".")
		fmt.Printf("Editing project: %s, key: %s (using %s)\n", project, path, settings.Editor)
	}
	// TODO: Open editor for specific nested key
	// Use settings.Editor to open the editor
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

// Example structure of myProject.yaml:
// recall.yaml
// info:
//   name: recall
//   description: CLI tool for managing project knowledge
//   type: project
//   usage: ./recall <project> [key] [subkey...]
// keys:
//   main:
//     description: Main entry point for the project
//     type: function
//     keys:
//       args:
//         description: Command line arguments
//         type: list
//       exampleVariable:
//         description: Example variable for demonstration
//         type: string
//         value: "default value"
//   initLocal:
//     description: Initializes local recall directory
//     type: function
//   initGlobal:
//     description: Initializes global recall directory
//     type: function
//   showKey:
//     description: Displays specific key info
//     type: function
//   editKey:
//     description: Opens editor for specific key
//     type: function
