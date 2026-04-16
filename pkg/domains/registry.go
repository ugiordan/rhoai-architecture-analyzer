package domains

import (
	"fmt"
	"sort"
)

var registry = map[string]DomainAnalyzer{}

// Register adds a domain analyzer to the global registry.
func Register(d DomainAnalyzer) {
	registry[d.Name()] = d
}

// Get returns the requested domain analyzers by name.
// Returns an error if any requested domain is not registered.
func Get(names []string) ([]DomainAnalyzer, error) {
	var result []DomainAnalyzer
	for _, name := range names {
		d, ok := registry[name]
		if !ok {
			return nil, fmt.Errorf("unknown domain: %q", name)
		}
		result = append(result, d)
	}
	return result, nil
}

// All returns all registered domain analyzers sorted by name.
func All() []DomainAnalyzer {
	result := make([]DomainAnalyzer, 0, len(registry))
	for _, d := range registry {
		result = append(result, d)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name() < result[j].Name()
	})
	return result
}

// ResolveDependencies expands a list of domain names to include all transitive
// dependencies. Returns an error if any dependency is not registered.
func ResolveDependencies(names []string) ([]string, error) {
	seen := make(map[string]bool)
	var resolve func(name string) error
	resolve = func(name string) error {
		if seen[name] {
			return nil
		}
		d, ok := registry[name]
		if !ok {
			return fmt.Errorf("unknown domain: %q", name)
		}
		seen[name] = true
		for _, dep := range d.Dependencies() {
			if err := resolve(dep); err != nil {
				return fmt.Errorf("domain %q requires %w", name, err)
			}
		}
		return nil
	}
	for _, name := range names {
		if err := resolve(name); err != nil {
			return nil, err
		}
	}
	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result, nil
}

// Names returns the names of all registered domains sorted alphabetically.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
