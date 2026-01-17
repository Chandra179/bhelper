package collision

type GeneratorRegistry struct {
	generators map[string]IDGenerator
	order      []string
}

func NewGeneratorRegistry() *GeneratorRegistry {
	return &GeneratorRegistry{
		generators: make(map[string]IDGenerator),
		order:      make([]string, 0),
	}
}

func (r *GeneratorRegistry) Register(gen IDGenerator) {
	name := gen.Name()
	r.generators[name] = gen
	r.order = append(r.order, name)
}

func (r *GeneratorRegistry) Get(name string) (IDGenerator, bool) {
	gen, ok := r.generators[name]
	return gen, ok
}

func (r *GeneratorRegistry) List() []IDGenerator {
	result := make([]IDGenerator, 0, len(r.order))
	for _, name := range r.order {
		result = append(result, r.generators[name])
	}
	return result
}

func (r *GeneratorRegistry) Names() []string {
	return r.order
}
