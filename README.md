# recall - A CLI Tool for Knowledge Management

A lightweight command-line tool that helps developers quickly store Project information in a nested data structure. Perfect for remembering function details, configuration snippets, project conventions, and implementation notes.

### Basic Usage
```bash
# Show help
recall

# Show general project information
recall myApp
# Output: "MyProject is a web application for task management..."

# Show specific function information
recall myApp database
recall myApp authentication
recall myApp deployment

# Edit information
recall --edit myApp database  # Edit info to database
recall --edit myApp           # Edit general project info
```

## Features

- **Local & Global Storage**: Store information locally per project (`./.recall/`) or globally (`~/.recall/`)
- **YAML-based**: Human-readable YAML files for easy editing and version control
- **Hierarchical Organization**: Structured data with nested keys and categories
- **Interactive Editing**: User-friendly editing interface with temporary documents
- **Quick Access**: Fast retrieval of project information without leaving the terminal

## Installation

```bash
# Install from source (placeholder - implement your preferred method)
git clone https://github.com/lbastigk/recall
cd recall
make install
```

## Quick Start

```bash
# Initialize recall in current project
recall --init

# Initialize global recall directory (optional)
recall --init-global
```

## Usage

### Basic Commands

```bash
recall                                      # Show usage help
recall <project>                            # Show all info from project.yaml
recall <project> <key>                      # Show specific key info
recall --edit <project> <key>               # Edit specific key
recall --init                               # Initialize local recall
recall --init-global                        # Initialize global recall
```

### Data Structure

Information is stored in YAML files with the following structure:

```yaml
# excerpt from .recall/recall.yaml:
info:
  example: |-
    recall --init

    recall --edit myProject

    recall myProject

    recall --edit myProject database connection

    recall myProject database connection
  infoLong: Recall is a command-line tool that helps manage project documentation
    and knowledge by storing information in YAML files with local/global storage priority.
    It allows you to create, edit, and retrieve project-specific information using
    a hierarchical key structure.
  infoShort: A CLI tool for project knowledge management
main:
  example: |
    func main() {
        settings := loadSettings()
        args := os.Args[1:]
        // Route to appropriate function based on args
    }
  infoLong: Handles command-line argument parsing and routes to appropriate functions
    based on the command pattern. Supports project info viewing, editing, and initialization
    commands.
  infoShort: Main entry point for the recall CLI application
  keys:
  ...
```

### Interactive Editing

When editing information, recall opens a user-friendly editor interface. Default editor is nano.

Example edit entry:
```
# Edit the information below and save the file
# Lines starting with # are comments and will be ignored
# Leave sections empty if not needed

infoShort:
Database connection utilities

infoLong:
Functions for connecting to the database, handling queries,
and managing connection pools.

example:
conn = Database.connect()
result = conn.query("SELECT * FROM users")
conn.close()
```

## Examples

### Basic Usage
```bash
# Show help
recall

# Show general information for myApp project
recall myApp

# Show specific function information
recall myApp database
recall myApp authentication
recall myApp deployment

# Edit information
recall --edit myApp database
```

### Common Use Cases
```bash
# Document database functions
recall --edit myApp database

# Remember deployment steps  
recall --edit myApp deployment

# Store API information
recall --edit myApp api

# Quick lookup
recall myApp database
# Shows: database connection info and examples
```

## File Locations

- **Local**: `./.recall/<project>.yaml` (project-specific)
- **Global**: `~/.recall/<project>.yaml` (accessible from anywhere)

The tool searches local storage first, then falls back to global storage.

## Configuration

Create `~/.recall/config.yaml` to customize behavior:

```yaml
editor: nano                    # Preferred editor for editing
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details