# bhelper

A powerful terminal-based CLI utility built with Go and Bubble Tea, designed to be your daily helper for text analysis, time conversions, and ID collision analysis.

## Features

bhelper provides an interactive terminal interface with four main analytical tools:

### 🔤 Character Analyzer
- Analyzes text encoding and character properties
- Shows decimal, hexadecimal, and binary representations
- Displays UTF-8, UTF-16, UTF-32 byte counts
- Perfect for debugging encoding issues and understanding character representations

### 🕐 Timezone Analyzer  
- Converts dates (dd-mm-yyyy format) to Unix timestamps
- Displays comprehensive time information including timezone offset, UTC time, day of week, ISO week number, Julian day, season, and leap year status
- Essential for programming timestamp conversions

### ⏱️ Time Converter
- Converts time values between different units (nanoseconds to hours)
- Supports units: ns, µs/us, ms, s, min/m, h/hr
- Useful for understanding time relationships and debugging timing issues

### 🎯 Collision Analyzer
- Analyzes collision probability for ID generation systems
- Supports Base64, Base62, and Snowflake ID formats
- Uses mathematical calculations (birthday paradox) and actual simulations
- Provides time-to-collision estimates at different probability levels

## Installation

### Prerequisites
- Go 1.25.5 or higher
- Terminal environment

## Usage

### Running the Application

```bash
# Run directly with make
make r
# Equivalent to: go run .

# Or run the built binary
./bhelper
```

### Interactive Interface

The application provides an intuitive terminal interface:

- **Feature Selection**: Navigate with arrow keys (↑↓), select with Enter
- **Help Screens**: Press `H` or `?` for detailed feature help
- **Input Execution**: Type your input and press Enter to execute
- **History Navigation**: `Ctrl+Z` (undo), `Ctrl+Y` (redo)

### Usage Examples

#### Character Analysis
```
Input: Hello
Output: Character encoding analysis with decimal/hex/binary representations

Input: 🚀
Output: Unicode character analysis with byte counts for different UTF formats
```

#### Timezone Conversion
```
Input: 16-01-2026
Output: Unix timestamp, UTC time, timezone offset, day of week, and more
```

#### Time Conversion
```
Input: 100ms
Output: Conversion to nanoseconds, microseconds, seconds, minutes, and hours
```

#### Collision Analysis
```
Input: base64:10:1000/sec
Output: Collision probability analysis for Base64 IDs with 10-character length at 1000 IDs per second

Input: snowflake:64:10000/sec
Output: Collision analysis for Snowflake IDs with 64-bit format at 10000 IDs per second
```

## Project Structure

```
bhelper/
├── main.go                    # Entry point, feature registration
├── cli.go                     # Main TUI model and UI logic
├── styles.go                  # Lipgloss styling definitions
├── history.go                 # Input history management
├── Makefile                   # Build and run commands
├── go.mod/go.sum              # Go module dependencies
├── feature/                   # Core feature package
│   ├── feature.go            # Feature interface and registry
│   ├── character.go          # Text encoding analyzer
│   ├── timezone.go           # Unix timestamp converter
│   └── time/                 # Time conversion package
│       ├── converter.go      # Time unit converter
│       └── converter_test.go # Tests
└── feature/collision/         # ID collision analysis package
    ├── analyzer.go           # Main collision analysis logic
    ├── parser.go             # Input parsing
    ├── generator.go          # ID generation interfaces
    ├── formatter.go          # Output formatting
    └── registry.go           # Generator registry
```

## Architecture

bhelper uses a clean, modular architecture:

- **Plugin-based Design**: Features implement a common `Feature` interface
- **Registry Pattern**: Centralized feature management and discovery
- **TUI Framework**: Built with Bubble Tea for responsive terminal interface
- **Modular Structure**: Each feature is self-contained with comprehensive tests
