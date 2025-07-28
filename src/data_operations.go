package main

import (
	"io/ioutil"
	"strings"
)

// KeyData represents the structure of information stored for each key
type KeyData struct {
	InfoShort string `yaml:"infoShort,omitempty"`
	InfoLong  string `yaml:"infoLong,omitempty"`
	Example   string `yaml:"example,omitempty"`
}

// ProjectData represents the entire project data structure
type ProjectData map[string]interface{}

func getKeyData(projectData ProjectData, keyPath string) KeyData {
	// keyPath is a string representing the path to the key, e.g. "key1.keys.key2"
	// If keyPath is empty, use "info" for root-level project information
	if keyPath == "" {
		keyPath = "info"
	}
	
	// Split the keyPath into components
	keyParts := strings.Split(keyPath, ".")
	
	// Navigate to nested key
	current := projectData
	for i, key := range keyParts {
		if value, exists := current[key]; exists {
			if i == len(keyParts)-1 {
				// This is the final key
				if keyMap, ok := value.(map[interface{}]interface{}); ok {
					return KeyData{
						InfoShort: getStringValue(keyMap, "infoShort"),
						InfoLong:  getStringValue(keyMap, "infoLong"),
						Example:   getStringValue(keyMap, "example"),
					}
				}
			} else {
				// Navigate deeper
				if nextMap, ok := value.(map[interface{}]interface{}); ok {
					// Convert to ProjectData format
					converted := make(ProjectData)
					for k, v := range nextMap {
						if strKey, ok := k.(string); ok {
							converted[strKey] = v
						}
					}
					current = converted
				} else {
					break
				}
			}
		} else {
			break
		}
	}
	
	return KeyData{} // Return empty if key doesn't exist
}

func getStringValue(m map[interface{}]interface{}, key string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func setKeyData(projectData ProjectData, keyPath string, data KeyData) {
	// keyPath is a string representing the path to the key, e.g. "key1.keys.key2"
	// If keyPath is empty, use "info" for root-level project information
	if keyPath == "" {
		keyPath = "info"
	}
	
	// Split the keyPath into components
	keyParts := strings.Split(keyPath, ".")
	
	// Navigate/create nested structure
	current := projectData
	for i, key := range keyParts {
		if i == len(keyParts)-1 {
			// This is the final key, preserve existing structure and update only the info fields
			if existing, exists := current[key]; exists {
				// Key already exists, preserve existing structure
				// Convert existing data to map[string]interface{} regardless of type
				existingMap := make(map[string]interface{})
				
				// Handle both map[string]interface{} and map[interface{}]interface{}
				switch v := existing.(type) {
				case map[string]interface{}:
					// Copy all existing values
					for k, val := range v {
						existingMap[k] = val
					}
				case map[interface{}]interface{}:
					// Convert and copy all existing values
					for k, val := range v {
						if strKey, ok := k.(string); ok {
							existingMap[strKey] = val
						}
					}
				default:
					// Not a map, start fresh but this shouldn't happen
				}
				
				// Now update only the info fields
				existingMap["infoShort"] = data.InfoShort
				existingMap["infoLong"] = data.InfoLong
				existingMap["example"] = data.Example
				
				current[key] = existingMap
			} else {
				// Key doesn't exist, create new
				current[key] = map[string]interface{}{
					"infoShort": data.InfoShort,
					"infoLong":  data.InfoLong,
					"example":   data.Example,
				}
			}
		} else {
			// Navigate deeper, create if doesn't exist
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]interface{})
			}
			
			// Convert current[key] to map[string]interface{} for navigation
			switch v := current[key].(type) {
			case map[string]interface{}:
				current = v
			case map[interface{}]interface{}:
				// Convert to map[string]interface{}
				converted := make(map[string]interface{})
				for k, val := range v {
					if strKey, ok := k.(string); ok {
						converted[strKey] = val
					}
				}
				current[key] = converted
				current = converted
			default:
				// Not a map, replace with empty map
				newMap := make(map[string]interface{})
				current[key] = newMap
				current = newMap
			}
		}
	}
}

func parseEditedFile(filename string) (KeyData, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return KeyData{}, err
	}
	
	content := string(data)
	lines := strings.Split(content, "\n")
	
	var keyData KeyData
	var currentSection string
	var currentContent []string
	
	for _, line := range lines {
		// Skip comments and empty lines at the beginning
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		
		// Check if this is a section header
		trimmed := strings.TrimSpace(line)
		if trimmed == "infoShort:" {
			// Save previous section if any
			if currentSection != "" {
				setSection(&keyData, currentSection, strings.TrimSpace(strings.Join(currentContent, "\n")))
			}
			currentSection = "infoShort"
			currentContent = []string{}
		} else if trimmed == "infoLong:" {
			if currentSection != "" {
				setSection(&keyData, currentSection, strings.TrimSpace(strings.Join(currentContent, "\n")))
			}
			currentSection = "infoLong"
			currentContent = []string{}
		} else if trimmed == "example:" {
			if currentSection != "" {
				setSection(&keyData, currentSection, strings.TrimSpace(strings.Join(currentContent, "\n")))
			}
			currentSection = "example"
			currentContent = []string{}
		} else if currentSection != "" {
			// Add line to current section content
			currentContent = append(currentContent, line)
		}
	}
	
	// Don't forget the last section
	if currentSection != "" {
		setSection(&keyData, currentSection, strings.TrimSpace(strings.Join(currentContent, "\n")))
	}
	
	return keyData, nil
}

func setSection(keyData *KeyData, section, content string) {
	switch section {
	case "infoShort":
		keyData.InfoShort = content
	case "infoLong":
		keyData.InfoLong = content
	case "example":
		keyData.Example = content
	}
}
