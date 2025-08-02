package main



import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var version = "1.0.0"

func main() {
	// Load settings from ~/.recall/settings.yaml if it exists
	settings := loadSettings()

	args := os.Args[1:] // Skip the program name

	// Check if --edit flag is at the end
	editMode := false
	if len(args) > 1 && args[len(args)-1] == "--edit" {
		editMode = true
		args = args[:len(args)-1] // Remove --edit from the end
	}

	// Handle different argument patterns
	switch len(args) {
	case 0:
		if editMode {
			fmt.Println("[ERROR] --edit requires at least the project name")
			showUsage(settings)
			os.Exit(1)
		}
		showUsage(settings)
	case 1:
		if args[0] == "--init" {
			if editMode {
				fmt.Println("[ERROR] Cannot use --edit with --init")
				showUsage(settings)
				os.Exit(1)
			}
			initLocal(settings)
		} else if args[0] == "--init-global" {
			if editMode {
				fmt.Println("[ERROR] Cannot use --edit with --init-global")
				showUsage(settings)
				os.Exit(1)
			}
			initGlobal(settings)
		} else if args[0] == "--version" {
			fmt.Println(version)
		} else {
			project := args[0]
			keyPath := []string{}
			if editMode {
				editKey(settings, project, keyPath)
			} else {
				showKey(settings, project, keyPath)	// Empty keyPath means general project info
			}
		}
	default:
		// Handle --edit command at beginning (legacy support)
		if args[0] == "--edit" {
			if editMode {
				fmt.Println("[ERROR] Cannot use --edit both at beginning and end")
				showUsage(settings)
				os.Exit(1)
			}
			if len(args) < 2 {
				fmt.Println("[ERROR] --edit requires at least the project name")
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
			if editMode {
				editKey(settings, project, keyPath)
			} else {
				showKey(settings, project, keyPath)
			}
		}
	}
}

func showUsage(settings *Settings) {
	fmt.Println("recall CLI Tool - Version " + version)
	fmt.Println("Copyright (c) 2023 Your Name")
	fmt.Println("This tool helps you manage project knowledge efficiently.")
	fmt.Println("For more information, visit: https://github.com/yourusername/recall")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  recall                                Show this help")
	fmt.Println("  recall <project>                      Show general project info")
	fmt.Println("  recall <project> <key>                Show specific key info")
	fmt.Println("  recall <project> <key> <subkey>...    Show nested key info")
	fmt.Println("  recall --edit <project> <key>...      Edit specific key")
	fmt.Println("  recall <project> <key>... --edit      Edit specific key (alternative)")
	fmt.Println("  recall --init                         Initialize local recall")
	fmt.Println("  recall --init-global                  Initialize global recall")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  recall myApp")
	fmt.Println("  recall myApp database")
	fmt.Println("  recall myApp myClass myFunction")
	fmt.Println("  recall myApp myClass myFunction myVariable")
	fmt.Println("  recall --edit myApp deployment")
	fmt.Println("  recall myApp deployment --edit")
	fmt.Println()
	fmt.Printf("Settings: Editor=%s\n", settings.Editor)
}

func initLocal(settings *Settings) {
	fmt.Println("[INFO] Initializing local recall directory at ./.recall/...")
	if err := os.MkdirAll("./.recall", 0755); err != nil {
		fmt.Printf("[ERROR] Error creating ./.recall/ directory: %v\n", err)
		return
	}
	fmt.Println("[INFO] Created ./.recall/ directory")
}

func initGlobal(settings *Settings) {
	// 1.) Check if ~/.recall/ directory exists
	homeDir, _ := os.UserHomeDir()
	globalPath := homeDir + "/.recall"
	if _, err := os.Stat(globalPath); err == nil {
		// Directory exists, no need to create it
		fmt.Println("[INFO] Global recall directory already exists at " + globalPath)
	} else {
		// 2.) If not, create it
		if err := os.MkdirAll(globalPath, 0755); err != nil {
			fmt.Printf("[ERROR] Error creating "+globalPath+" directory: %v\n", err)
			return
		}
		fmt.Printf("[INFO] Created %s directory\n", globalPath)
	}

	// 3.) Check if settings.yaml exists in ~/.recall/
	settingsFile := globalPath + "/settings.yaml"
	if _, err := os.Stat(settingsFile); err == nil {
		// File exists, no need to create it
		fmt.Println("[INFO] Global settings file already exists at " + settingsFile)
	} else {
		// 4.) If not, create default settings.yaml
		defaultSettings := defaultSettings()
		data, err := yaml.Marshal(defaultSettings)
		if err != nil {
			fmt.Printf("[ERROR] Error marshaling settings: %v\n", err)
			return
		}
		
		if err := ioutil.WriteFile(settingsFile, data, 0644); err != nil {
			fmt.Printf("[ERROR] Error creating settings file: %v\n", err)
			return
		}
		fmt.Printf("[INFO] Created default settings file at %s\n", settingsFile)
	}
}

func showKey(settings *Settings, project string, keyPath []string) {
	var path string
	if len(keyPath) == 0 {
		fmt.Printf("Project: %s (general info)\n", project)
		path = "info" // Use "info" key for root-level project information
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
		fmt.Printf("Project: %s, Key: %s\n", project, strings.Join(keyPath, " → "))
	}

	// 1.) Find project file and load existing data
	projectFile := findProjectFile(project)
	projectData := loadProjectData(projectFile)
	
	// 2.) Check if project file exists
	if len(projectData) == 0 {
		fmt.Printf("[ERROR] Project '%s' not found. Use --edit to create it.\n", project)
		return
	}
	
	// 3.) Get key data
	keyData := getKeyData(projectData, path)
	
	// 4.) Display the information
	if keyData.InfoShort == "" && keyData.InfoLong == "" && keyData.Example == "" {
		if len(keyPath) == 0 {
			fmt.Printf("[INFO] No general info found for project '%s'. Use --edit to add it.\n", project)
		} else {
			fmt.Printf("[INFO] Key '%s' not found or is empty. Use --edit to create it.\n", strings.Join(keyPath, " → "))
		}
		return
	}

	if keyData.InfoShort != "" {
		fmt.Println()
		fmt.Printf("\033[1;32m#############################\033[0m\n")
		fmt.Printf("\033[1;32mShort:\033[0m\n%s\n", keyData.InfoShort)
	}
	if keyData.InfoLong != "" {
		fmt.Println()
		fmt.Printf("\033[1;32m#############################\033[0m\n")
		fmt.Printf("\033[1;32mDescription:\033[0m\n%s\n", keyData.InfoLong)
	}
	if keyData.Example != "" {
		fmt.Println()
		fmt.Printf("\033[1;32m#############################\033[0m\n")
		fmt.Printf("\033[1;32mExample:\033[0m\n%s\n", keyData.Example)
	}
	
	// 5.) Show available sub-keys if they exist
	showSubKeys(projectData, path)
}

func showSubKeys(projectData ProjectData, keyPath string) {
	var current map[string]interface{}
	
	if keyPath == "info" {
		// For "info" path, we want to show the root-level keys (excluding "info" itself)
		current = make(map[string]interface{})
		for k, v := range projectData {
			if k != "info" { // Exclude the info section itself
				current[k] = v
			}
		}
	} else if keyPath == "" {
		// Empty path means we're at the root, convert projectData to the right type
		current = make(map[string]interface{})
		for k, v := range projectData {
			current[k] = v
		}
	} else {
		// Navigate to the current key to check for sub-keys
		keyParts := strings.Split(keyPath, ".keys.")
		currentData := projectData
		
		// Navigate to the target key
		for _, key := range keyParts {
			if key == "" {
				continue
			}
			if value, exists := currentData[key]; exists {
				switch v := value.(type) {
				case map[string]interface{}:
					currentData = v
				case map[interface{}]interface{}:
					// Convert to map[string]interface{}
					converted := make(map[string]interface{})
					for k, val := range v {
						if strKey, ok := k.(string); ok {
							converted[strKey] = val
						}
					}
					currentData = converted
				default:
					return // Not a map, no sub-keys
				}
			} else {
				return // Key doesn't exist
			}
		}
		current = currentData
	}
	
	// Look for "keys" sub-section or direct keys at root level
	var keysToShow map[string]interface{}
	
	if keysSection, exists := current["keys"]; exists {
		// There's a "keys" section, use that
		switch keys := keysSection.(type) {
		case map[string]interface{}:
			keysToShow = keys
		case map[interface{}]interface{}:
			keysToShow = make(map[string]interface{})
			for k, v := range keys {
				if strKey, ok := k.(string); ok {
					keysToShow[strKey] = v
				}
			}
		}
	} else if keyPath == "" || keyPath == "info" {
		// At root level and no "keys" section, show all top-level keys
		keysToShow = make(map[string]interface{})
		for k, v := range current {
			// Skip info fields at root level
			if k != "infoShort" && k != "infoLong" && k != "example" {
				keysToShow[k] = v
			}
		}
	}
	//*
	if len(keysToShow) > 0 {
		fmt.Println()
		fmt.Printf("\033[1;32m#############################\033[0m\n")
		fmt.Printf("\033[1;32mAvailable sub-keys:\033[0m\n")
		for subKey := range keysToShow {
			fmt.Printf("  • %s\n", subKey)
		}
	}
	//*/
}

func editKey(settings *Settings, project string, keyPath []string) {
	var path string
	if len(keyPath) == 0 {
		fmt.Printf("[INFO] Editing general info for project: %s\n", project)
		fmt.Printf("[INFO] Edit the content between infoShort:, infoLong:, and example: sections\n")
		path = "" // Empty path means edit the root of the document
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
		fmt.Printf("[INFO] Editing project: %s, key: %s (using %s)\n", project, path, settings.Editor)
		fmt.Printf("[INFO] Edit the content between infoShort:, infoLong:, and example: sections\n")
	}

	// 1.) Find project file and load existing data
	projectFile := findProjectFile(project)
	projectData := loadProjectData(projectFile)
	
	// 2.) Get current key data or create new entry
	currentData := getKeyData(projectData, path)
	
	// 3.) Create temporary YAML file with current key info
	tempFile, err := createTempEditFile(currentData)
	if err != nil {
		fmt.Printf("[ERROR] Error creating temp file: %v\n", err)
		return
	}
	defer os.Remove(tempFile) // Clean up temp file when done
	
	// 4.) Use settings.Editor to open the file
	cmd := exec.Command(settings.Editor, tempFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("[ERROR] Error running editor: %v\n", err)
		return
	}
	
	// 5.) Read the edited file back
	editedData, err := parseEditedFile(tempFile)
	if err != nil {
		fmt.Printf("[ERROR] Error parsing edited file: %v\n", err)
		return
	}
	
	// 6.) Update the project data and save
	setKeyData(projectData, path, editedData)
	saveProjectData(projectFile, projectData)
	
	fmt.Printf("[INFO] Saved changes to %s\n", projectFile)
}

