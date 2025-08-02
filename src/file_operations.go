package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

func findProjectFile(project string) string {
	// Check if local .recall directory exists
	if _, err := os.Stat("./.recall"); err == nil {
		// Local .recall directory exists, prefer local storage
		localFile := fmt.Sprintf("./.recall/%s.yaml", project)
		return localFile
	}
	
	// No local .recall directory, check if project file exists locally first
	localFile := fmt.Sprintf("./.recall/%s.yaml", project)
	if _, err := os.Stat(localFile); err == nil {
		return localFile
	}
	
	// Fall back to global storage
	homeDir, _ := os.UserHomeDir()
	globalFile := fmt.Sprintf("%s/.recall/%s.yaml", homeDir, project)
	return globalFile
}

func loadProjectData(filename string) ProjectData {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, return empty project data
		return make(ProjectData)
	}
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("[ERROR] Could not read file: %v\n", err)
		return make(ProjectData)
	}
	
	var projectData ProjectData
	if err := yaml.Unmarshal(data, &projectData); err != nil {
		fmt.Printf("[ERROR] Could not parse YAML: %v\n", err)
		return make(ProjectData)
	}
	
	return projectData
}

func saveProjectData(filename string, projectData ProjectData) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	data, err := yaml.Marshal(projectData)
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(filename, data, 0644)
}

func createTempEditFile(data KeyData) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "recall_edit_*.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	content := fmt.Sprintf(`infoShort:
%s

infoLong:
%s

example:
%s
`, data.InfoShort, data.InfoLong, data.Example)
	
	if _, err := file.WriteString(content); err != nil {
		return "", err
	}
	
	return file.Name(), nil
}
