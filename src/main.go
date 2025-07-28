package main

import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

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
	if err := os.MkdirAll("./.recall", 0755); err != nil {
		fmt.Printf("Error creating ./.recall/ directory: %v\n", err)
		return
	}
	fmt.Println("Created ./.recall/ directory")
}

func initGlobal(settings *Settings) {
	// 1.) Check if ~/.recall/ directory exists
	homeDir, _ := os.UserHomeDir()
	globalPath := homeDir + "/.recall"
	if _, err := os.Stat(globalPath); err == nil {
		// Directory exists, no need to create it
		fmt.Println("Global recall directory already exists at " + globalPath)
	} else {
		// 2.) If not, create it
		if err := os.MkdirAll(globalPath, 0755); err != nil {
			fmt.Printf("Error creating "+globalPath+" directory: %v\n", err)
			return
		}
		fmt.Printf("Created %s directory\n", globalPath)
	}

	// 3.) Check if settings.yaml exists in ~/.recall/
	settingsFile := globalPath + "/settings.yaml"
	if _, err := os.Stat(settingsFile); err == nil {
		// File exists, no need to create it
		fmt.Println("Global settings file already exists at " + settingsFile)
	} else {
		// 4.) If not, create default settings.yaml
		defaultSettings := defaultSettings()
		data, err := yaml.Marshal(defaultSettings)
		if err != nil {
			fmt.Printf("Error marshaling settings: %v\n", err)
			return
		}
		
		if err := ioutil.WriteFile(settingsFile, data, 0644); err != nil {
			fmt.Printf("Error creating settings file: %v\n", err)
			return
		}
		fmt.Printf("Created default settings file at %s\n", settingsFile)
	}
}

func showKey(settings *Settings, project string, keyPath []string) {
	path := strings.Join(keyPath, ".keys.")
	fmt.Printf("Project: %s, Key: %s\n", project, path)

	// TODO: Load and display specific nested key from project.yaml
	// Look in ./.recall/ first, then ~/.recall/

	// 1.) Check if ./.recall/<project>.yaml exists

	// 2.) If not, check ~/.recall/<project>.yaml

	// 3.) If found, read the file and parse the YAML

	// 4.) Display the key info, or an error if not found

}

func editKey(settings *Settings, project string, keyPath []string) {
	var path string
	if len(keyPath) == 0 {
		fmt.Printf("Editing general info for project: %s\n", project)
		path = "info"
	} else {
		// For nested keys, we need to construct the path properly
		// e.g., keyPath ["foo", "bar"] becomes "foo.keys.bar"
		if len(keyPath) == 1 {
			path = keyPath[0]
		} else {
			// For multiple levels, insert ".keys." between them
			path = keyPath[0]
			for i := 1; i < len(keyPath); i++ {
				path += ".keys." + keyPath[i]
			}
		}
		fmt.Printf("Editing project: %s, key: %s (using %s)\n", project, path, settings.Editor)
	}

	// 1.) Find project file and load existing data
	projectFile := findProjectFile(project)
	projectData := loadProjectData(projectFile)
	
	// 2.) Get current key data or create new entry
	currentData := getKeyData(projectData, path)
	
	// 3.) Create temporary YAML file with current key info
	tempFile, err := createTempEditFile(currentData)
	if err != nil {
		fmt.Printf("Error creating temp file: %v\n", err)
		return
	}
	defer os.Remove(tempFile) // Clean up temp file when done
	
	// 4.) Use settings.Editor to open the file
	cmd := exec.Command(settings.Editor, tempFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running editor: %v\n", err)
		return
	}
	
	// 5.) Read the edited file back
	editedData, err := parseEditedFile(tempFile)
	if err != nil {
		fmt.Printf("Error parsing edited file: %v\n", err)
		return
	}
	
	// 6.) Update the project data and save
	setKeyData(projectData, path, editedData)
	saveProjectData(projectFile, projectData)
	
	fmt.Printf("✓ Saved changes to %s\n", projectFile)
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
