package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ContractCheck describes the result of comparing one resource's schema.
type ContractCheck struct {
	Resource        string   `json:"resource"`
	Status          string   `json:"status,omitempty"`
	Message         string   `json:"message,omitempty"`
	BreakingChanges int      `json:"breaking_changes,omitempty"`
	AdditiveChanges int      `json:"additive_changes,omitempty"`
	IsCompatible    bool     `json:"is_compatible"`
	Consumers       []string `json:"consumers,omitempty"`
	Details         []struct {
		Path        string `json:"path"`
		Type        string `json:"type"`
		Description string `json:"description"`
	} `json:"details,omitempty"`
}

// AffectedConsumer records a consumer impacted by breaking changes.
type AffectedConsumer struct {
	Repo            string   `json:"repo"`
	Usage           string   `json:"usage"`
	Resource        string   `json:"resource"`
	BreakingChanges []string `json:"breaking_changes"`
}

// ContractResult is the top-level result of a contract check for a provider.
type ContractResult struct {
	Provider          string             `json:"provider"`
	IsCompatible      bool               `json:"is_compatible"`
	Checks            []ContractCheck    `json:"checks"`
	AffectedConsumers []AffectedConsumer `json:"affected_consumers"`
}

// dependencyGraph mirrors the YAML structure of dependency-graph.yaml.
type dependencyGraph struct {
	Contracts []struct {
		Provider  string `yaml:"provider"`
		Resource  string `yaml:"resource"`
		Consumers []struct {
			Repo  string `yaml:"repo"`
			Usage string `yaml:"usage"`
		} `yaml:"consumers"`
	} `yaml:"contracts"`
}

// CheckContract validates new schemas from a provider against stored baseline schemas
// and reports compatibility with consumers defined in the dependency graph.
func CheckContract(provider string, newSchemas []SchemaInfo, contractsDir string) (*ContractResult, error) {
	result := &ContractResult{
		Provider:     provider,
		IsCompatible: true,
	}

	depGraph, err := loadDependencyGraph(contractsDir)
	if err != nil {
		return nil, fmt.Errorf("loading dependency graph: %w", err)
	}

	for _, si := range newSchemas {
		resourceKey := si.ResourceKey
		newSchema := si.Schema

		oldSchema, err := loadStoredSchema(contractsDir, provider, resourceKey)
		if err != nil {
			return nil, fmt.Errorf("loading stored schema %s/%s: %w", provider, resourceKey, err)
		}

		if oldSchema == nil {
			log.Printf("No stored schema for %s/%s, skipping comparison (first extraction)", provider, resourceKey)
			result.Checks = append(result.Checks, ContractCheck{
				Resource:     resourceKey,
				Status:       "new",
				Message:      "First extraction, no baseline to compare",
				IsCompatible: true,
			})
			continue
		}

		diff := DiffSchemas(oldSchema, newSchema)
		consumers := findConsumers(depGraph, provider, resourceKey)

		consumerRepos := make([]string, 0, len(consumers))
		for _, c := range consumers {
			consumerRepos = append(consumerRepos, c.repo)
		}

		check := ContractCheck{
			Resource:        resourceKey,
			BreakingChanges: len(diff.BreakingChanges),
			AdditiveChanges: len(diff.AdditiveChanges),
			IsCompatible:    diff.IsCompatible(),
			Consumers:       consumerRepos,
		}

		if !diff.IsCompatible() {
			result.IsCompatible = false
			for _, c := range diff.BreakingChanges {
				check.Details = append(check.Details, struct {
					Path        string `json:"path"`
					Type        string `json:"type"`
					Description string `json:"description"`
				}{
					Path:        c.Field,
					Type:        c.ChangeType,
					Description: c.Description,
				})
			}

			breakingDescs := make([]string, len(diff.BreakingChanges))
			for i, c := range diff.BreakingChanges {
				breakingDescs[i] = c.Description
			}

			for _, consumer := range consumers {
				result.AffectedConsumers = append(result.AffectedConsumers, AffectedConsumer{
					Repo:            consumer.repo,
					Usage:           consumer.usage,
					Resource:        resourceKey,
					BreakingChanges: breakingDescs,
				})
			}
		}

		result.Checks = append(result.Checks, check)
	}

	return result, nil
}

func loadDependencyGraph(contractsDir string) (*dependencyGraph, error) {
	graphPath := filepath.Join(contractsDir, "dependency-graph.yaml")
	data, err := os.ReadFile(graphPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &dependencyGraph{}, nil
		}
		return nil, err
	}
	var dg dependencyGraph
	if err := yaml.Unmarshal(data, &dg); err != nil {
		return nil, err
	}
	return &dg, nil
}

func loadStoredSchema(contractsDir, provider, resourceKey string) (map[string]interface{}, error) {
	// Validate provider and resourceKey to prevent path traversal via "../"
	for _, part := range []string{provider, resourceKey} {
		if strings.Contains(part, "..") || strings.ContainsAny(part, `/\`) {
			return nil, fmt.Errorf("invalid path component: %q", part)
		}
	}

	schemaPath := filepath.Join(contractsDir, "schemas", provider, resourceKey+".json")

	// Verify the resolved path stays within the expected directory
	schemasDir := filepath.Join(contractsDir, "schemas")
	absSchema, err := filepath.Abs(schemaPath)
	if err != nil {
		return nil, err
	}
	absSchemas, err := filepath.Abs(schemasDir)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(absSchema, absSchemas+string(os.PathSeparator)) {
		return nil, fmt.Errorf("schema path escapes contracts directory: %s", schemaPath)
	}

	// Open and read the file. Path traversal is prevented by the component
	// validation and abs-path containment check above. Symlink following is
	// acceptable here because the path is already validated to stay within
	// the schemas directory.
	f, err := os.Open(schemaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", schemaPath)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return schema, nil
}

type consumerInfo struct {
	repo  string
	usage string
}

type consumerAccum struct {
	repo     string
	usages   []string
	seenUsage map[string]bool
}

func findConsumers(dg *dependencyGraph, provider, resourceKey string) []consumerInfo {
	byRepo := map[string]*consumerAccum{}

	for _, contract := range dg.Contracts {
		if contract.Provider != provider || contract.Resource != resourceKey {
			continue
		}
		for _, c := range contract.Consumers {
			acc, ok := byRepo[c.Repo]
			if !ok {
				acc = &consumerAccum{
					repo:      c.Repo,
					seenUsage: make(map[string]bool),
				}
				byRepo[c.Repo] = acc
			}
			if c.Usage != "" && !acc.seenUsage[c.Usage] {
				acc.seenUsage[c.Usage] = true
				acc.usages = append(acc.usages, c.Usage)
			}
		}
	}

	result := make([]consumerInfo, 0, len(byRepo))
	for _, acc := range byRepo {
		result = append(result, consumerInfo{
			repo:  acc.repo,
			usage: strings.Join(acc.usages, "; "),
		})
	}
	return result
}
