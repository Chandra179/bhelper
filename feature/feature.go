package feature

// Feature represents a plugin that can be executed by the CLI
type Feature interface {
	// ID returns unique identifier for the feature
	ID() string

	// Name returns display name
	Name() string

	// Description returns short description
	Description() string

	// Help returns detailed help text
	Help() string

	// Execute runs the feature with given input
	Execute(input string) (string, error)

	// Examples returns usage examples
	Examples() []Example
}

// Example represents a usage example
type Example struct {
	Input       string
	Description string
}

// FeatureRegistry manages all available features
type FeatureRegistry struct {
	features map[string]Feature
	order    []string // Preserve registration order
}

// NewFeatureRegistry creates a new feature registry
func NewFeatureRegistry() *FeatureRegistry {
	return &FeatureRegistry{
		features: make(map[string]Feature),
		order:    make([]string, 0),
	}
}

// Register adds a feature to the registry
func (r *FeatureRegistry) Register(f Feature) {
	r.features[f.ID()] = f
	r.order = append(r.order, f.ID())
}

// Get retrieves a feature by ID
func (r *FeatureRegistry) Get(id string) (Feature, bool) {
	f, ok := r.features[id]
	return f, ok
}

// List returns all features in registration order
func (r *FeatureRegistry) List() []Feature {
	result := make([]Feature, 0, len(r.order))
	for _, id := range r.order {
		result = append(result, r.features[id])
	}
	return result
}
