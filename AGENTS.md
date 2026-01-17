# AGENTS.md

This file provides guidance for agentic coding assistants working on the bhelper codebase.

## Build, Test, and Lint Commands

### Dependency Management
```bash
make i
# Equivalent to: go mod tidy && go mod vendor
```

### Running the Application
```bash
make r
# Equivalent to: go run .
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./feature/collision/

# Run a single test
go test -run TestGenerate ./feature/collision/

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Building
```bash
go build -o bhelper .
```

## Project Architecture

**bhelper** is a CLI application built with Go and Bubble Tea TUI framework. It uses a plugin-based architecture where features are registered and executed through a registry system.

- **Root package**: Main CLI logic and TUI models
- **feature/**: Core feature interface and registry
- **feature/collision/**: Collision probability analyzer (example complex feature)
- **feature/character.go**: Text encoding analyzer
- **feature/timezone.go**: Timezone and timestamp analyzer

### Key Interfaces

```go
// All features must implement this interface
type Feature interface {
    ID() string
    Name() string
    Description() string
    Help() string
    Execute(input string) (string, error)
    Examples() []Example
}
```

## Code Style Guidelines

### Package Organization

- **Main package**: Root directory (CLI logic, TUI models)
- **feature/**: Core feature abstraction
- **feature/[featurename]/**: Individual feature packages with tests

### Import Ordering

1. Standard library imports (alphabetical)
2. Third-party imports (alphabetical)
3. Project imports (alphabetical)

Use descriptive aliases when helpful:
```go
tea "github.com/charmbracelet/bubbletea"
```

### Naming Conventions

- **Types/Interfaces**: PascalCase (`CollisionAnalyzer`, `FeatureRegistry`)
- **Exported functions**: PascalCase (`NewFeatureRegistry`)
- **Private functions**: camelCase (`calculateUTF16Bytes`)
- **Variables**: camelCase (`selectedIndex`, `config`)
- **Constants**: PascalCase or ALL_CAPS (`ModeFeatureList`, `MAX_SIZE`)
- **Constructors**: `New[TypeName]()` pattern

### Struct Definitions

```go
// TypeName describes what this struct does
type TypeName struct {
    field1 Type    // Description of field1
    field2 *Type   // Description of field2 (pointer)
}
```

### Constructors

```go
// NewTypeName creates a new TypeName with validation
func NewTypeName(param Type) (*TypeName, error) {
    if param <= 0 {
        return nil, fmt.Errorf("param must be positive, got %d", param)
    }
    return &TypeName{
        field: param,
    }, nil
}
```

### Error Handling

- Functions return `(result, error)` pattern
- Wrap errors with context using `%w`:
```go
return "", fmt.Errorf("failed to create generator: %w", err)
```
- Validate early and return errors immediately
- Don't ignore errors

### Method Organization

Group methods logically within a file:
1. Interface implementation methods (`ID()`, `Name()`, etc.)
2. Public methods
3. Private helper functions

### String Building

Use `strings.Builder` for concatenation:
```go
var b strings.Builder
b.WriteString("prefix")
b.WriteString(str)
return b.String()
```

Use `fmt.Sprintf` for simple formatting only.

### Testing Guidelines

**Test file naming**: `[source]_test.go` (e.g., `analyzer_test.go`)

**Test function naming**: `Test[FunctionName]` or `Test[TypeName][Method]`

```go
func TestNewFeatureRegistry(t *testing.T) {
    registry := NewFeatureRegistry()
    if registry == nil {
        t.Fatal("Expected non-nil registry")
    }
}

func TestParseInputError(t *testing.T) {
    _, err := ParseInput("invalid")
    if err == nil {
        t.Error("Expected error for invalid input")
    }
}
```

- Use `t.Fatalf()` for setup failures (should stop test)
- Use `t.Errorf()` for assertion failures (continue test)
- Test both success and error paths
- Test edge cases (empty input, invalid types, boundary values)

### Feature Registration

Features are registered in `main.go`:
```go
registry := feature.NewFeatureRegistry()
registry.Register(feature.NewCharacterAnalyzer())
registry.Register(feature.NewTimezoneAnalyzer())
registry.Register(collision.NewCollisionAnalyzer())
```

### TUI Models (Bubble Tea)

Implement the tea.Model interface:
```go
func (m Model) Init() tea.Cmd { ... }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (m Model) View() string { ... }
```

### Comments

- Package comments describe purpose
- Exported items have comments
- Keep comments concise and focused on "why" not "what"
- No inline comments for obvious code

### Dependencies

- **Bubble Tea**: TUI framework (`github.com/charmbracelet/bubbletea`)
- **Lipgloss**: Styling (`github.com/charmbracelet/lipgloss`)
- **Bubbles**: UI components (`github.com/charmbracelet/bubbles`)
- **Snowflake**: ID generation (`github.com/bwmarrin/snowflake`)

## File Structure Notes

- `vendor/` directory exists but is gitignored (local development)
- `styles.go` defines reusable TUI styles
- `history.go` provides undo/redo functionality
- `cli.go` contains main TUI model with mode switching

## Common Patterns

### Registry Pattern
Used for managing features and ID generators:
```go
registry := NewRegistry()
registry.Register(item)
items := registry.List()
item, ok := registry.Get(id)
```

### Switch on Type
Handle different message types in Bubble Tea:
```go
switch msg := msg.(type) {
case tea.KeyMsg:
    return c.updateFeatureList(msg)
case tea.WindowSizeMsg:
    return c, nil
}
```

### String Formatting for Output
Use descriptive section headers and consistent formatting:
```go
result.WriteString(fmt.Sprintf("Section:  %s\n", value))
result.WriteString(fmt.Sprintf("Subsection:  %s\n", value2))
```
