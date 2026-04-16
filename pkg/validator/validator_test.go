package validator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const testCRDYAML = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: widgets.example.com
spec:
  group: example.com
  names:
    kind: Widget
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                replicas:
                  type: integer
                name:
                  type: string
    - name: v2
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                replicas:
                  type: integer
                name:
                  type: string
                extra:
                  type: boolean
`

func TestExtractSchemasFromCRD(t *testing.T) {
	dir := t.TempDir()
	crdPath := filepath.Join(dir, "widget_crd.yaml")
	if err := os.WriteFile(crdPath, []byte(testCRDYAML), 0644); err != nil {
		t.Fatal(err)
	}

	schemas, err := ExtractSchemasFromCRD(crdPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(schemas) != 2 {
		t.Fatalf("expected 2 schemas, got %d", len(schemas))
	}
	if schemas[0].ResourceKey != "widget.v1" {
		t.Errorf("expected resource_key widget.v1, got %s", schemas[0].ResourceKey)
	}
	if schemas[1].ResourceKey != "widget.v2" {
		t.Errorf("expected resource_key widget.v2, got %s", schemas[1].ResourceKey)
	}
	if schemas[0].Group != "example.com" {
		t.Errorf("expected group example.com, got %s", schemas[0].Group)
	}
	if schemas[0].Kind != "Widget" {
		t.Errorf("expected kind Widget, got %s", schemas[0].Kind)
	}
}

func TestExtractSchemasFromDir(t *testing.T) {
	dir := t.TempDir()
	crdDir := filepath.Join(dir, "config", "crd", "bases")
	if err := os.MkdirAll(crdDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(crdDir, "widget.yaml"), []byte(testCRDYAML), 0644); err != nil {
		t.Fatal(err)
	}

	schemas, err := ExtractSchemasFromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(schemas) != 2 {
		t.Fatalf("expected 2 schemas, got %d", len(schemas))
	}
}

func TestDiffSchemasIdentical(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
		},
	}
	result := DiffSchemas(schema, schema)
	if !result.IsCompatible() {
		t.Error("identical schemas should be compatible")
	}
	if len(result.BreakingChanges) != 0 {
		t.Errorf("expected 0 breaking changes, got %d", len(result.BreakingChanges))
	}
	if len(result.AdditiveChanges) != 0 {
		t.Errorf("expected 0 additive changes, got %d", len(result.AdditiveChanges))
	}
}

func TestDiffSchemasRemovedField(t *testing.T) {
	old := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
			"age":  map[string]interface{}{"type": "integer"},
		},
	}
	new := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
		},
	}
	result := DiffSchemas(old, new)
	if result.IsCompatible() {
		t.Error("removing a field should be a breaking change")
	}
	if len(result.BreakingChanges) != 1 {
		t.Fatalf("expected 1 breaking change, got %d", len(result.BreakingChanges))
	}
	if result.BreakingChanges[0].ChangeType != "removed_field" {
		t.Errorf("expected change type removed_field, got %s", result.BreakingChanges[0].ChangeType)
	}
}

func TestDiffSchemasAddedField(t *testing.T) {
	old := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
		},
	}
	new := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name":  map[string]interface{}{"type": "string"},
			"email": map[string]interface{}{"type": "string"},
		},
	}
	result := DiffSchemas(old, new)
	if !result.IsCompatible() {
		t.Error("adding an optional field should be compatible")
	}
	if len(result.AdditiveChanges) != 1 {
		t.Fatalf("expected 1 additive change, got %d", len(result.AdditiveChanges))
	}
	if result.AdditiveChanges[0].ChangeType != "field_added" {
		t.Errorf("expected change type field_added, got %s", result.AdditiveChanges[0].ChangeType)
	}
}

func TestDiffSchemasTypeChange(t *testing.T) {
	old := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"count": map[string]interface{}{"type": "integer"},
		},
	}
	new := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"count": map[string]interface{}{"type": "string"},
		},
	}
	result := DiffSchemas(old, new)
	if result.IsCompatible() {
		t.Error("type change should be breaking")
	}
	found := false
	for _, c := range result.BreakingChanges {
		if c.ChangeType == "type_changed" {
			found = true
		}
	}
	if !found {
		t.Error("expected a type_changed breaking change")
	}
}

func TestDiffSchemasEnumRemoval(t *testing.T) {
	old := map[string]interface{}{
		"type": "string",
		"enum": []interface{}{"a", "b", "c"},
	}
	new := map[string]interface{}{
		"type": "string",
		"enum": []interface{}{"a", "b"},
	}
	result := DiffSchemas(old, new)
	if result.IsCompatible() {
		t.Error("removing enum values should be breaking")
	}
	found := false
	for _, c := range result.BreakingChanges {
		if c.ChangeType == "enum_changed" {
			found = true
		}
	}
	if !found {
		t.Error("expected an enum_changed breaking change")
	}
}

func TestDiffSchemasAdditionalPropertiesRestriction(t *testing.T) {
	old := map[string]interface{}{
		"type":                 "object",
		"additionalProperties": true,
	}
	new := map[string]interface{}{
		"type":                 "object",
		"additionalProperties": false,
	}
	result := DiffSchemas(old, new)
	if result.IsCompatible() {
		t.Error("restricting additionalProperties should be breaking")
	}
	found := false
	for _, c := range result.BreakingChanges {
		if c.ChangeType == "additional_properties_restricted" {
			found = true
		}
	}
	if !found {
		t.Error("expected additional_properties_restricted breaking change")
	}
}

func TestDiffSchemasCompositionReduced(t *testing.T) {
	old := map[string]interface{}{
		"oneOf": []interface{}{
			map[string]interface{}{"type": "string"},
			map[string]interface{}{"type": "integer"},
			map[string]interface{}{"type": "boolean"},
		},
	}
	new := map[string]interface{}{
		"oneOf": []interface{}{
			map[string]interface{}{"type": "string"},
		},
	}
	result := DiffSchemas(old, new)
	if result.IsCompatible() {
		t.Error("reducing composition options should be breaking")
	}
	var foundRestricted, foundTypeChanged bool
	for _, c := range result.BreakingChanges {
		if c.ChangeType == "composition_restricted" {
			foundRestricted = true
		}
		if c.ChangeType == "composition_type_changed" {
			foundTypeChanged = true
		}
	}
	if !foundRestricted {
		t.Error("expected composition_restricted breaking change")
	}
	if !foundTypeChanged {
		t.Error("expected composition_type_changed breaking change")
	}
}

func TestCheckContract(t *testing.T) {
	dir := t.TempDir()

	// Create dependency graph
	depGraphContent := `contracts:
  - provider: my-operator
    resource: widget.v1
    consumers:
      - repo: consumer-a
        usage: reads spec.replicas
      - repo: consumer-b
        usage: watches widgets
`
	if err := os.WriteFile(filepath.Join(dir, "dependency-graph.yaml"), []byte(depGraphContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create stored schema (baseline)
	schemaDir := filepath.Join(dir, "schemas", "my-operator")
	if err := os.MkdirAll(schemaDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"spec": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"replicas": map[string]interface{}{"type": "integer"},
					"name":     map[string]interface{}{"type": "string"},
				},
			},
		},
	}
	oldData, _ := json.Marshal(oldSchema)
	if err := os.WriteFile(filepath.Join(schemaDir, "widget.v1.json"), oldData, 0644); err != nil {
		t.Fatal(err)
	}

	// Test compatible change (add optional field)
	newSchemas := []SchemaInfo{{
		ResourceKey: "widget.v1",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"spec": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"replicas": map[string]interface{}{"type": "integer"},
						"name":     map[string]interface{}{"type": "string"},
						"extra":    map[string]interface{}{"type": "boolean"},
					},
				},
			},
		},
	}}
	cr, err := CheckContract("my-operator", newSchemas, dir)
	if err != nil {
		t.Fatal(err)
	}
	if !cr.IsCompatible {
		t.Error("adding optional field should be compatible")
	}
	if len(cr.AffectedConsumers) != 0 {
		t.Errorf("expected 0 affected consumers, got %d", len(cr.AffectedConsumers))
	}

	// Test breaking change (remove field)
	breakingSchemas := []SchemaInfo{{
		ResourceKey: "widget.v1",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"spec": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{"type": "string"},
					},
				},
			},
		},
	}}
	cr, err = CheckContract("my-operator", breakingSchemas, dir)
	if err != nil {
		t.Fatal(err)
	}
	if cr.IsCompatible {
		t.Error("removing field should be incompatible")
	}
	if len(cr.AffectedConsumers) != 2 {
		t.Errorf("expected 2 affected consumers, got %d", len(cr.AffectedConsumers))
	}

	// Test new resource (no baseline)
	newResourceSchemas := []SchemaInfo{{
		ResourceKey: "gadget.v1",
		Schema: map[string]interface{}{
			"type": "object",
		},
	}}
	cr, err = CheckContract("my-operator", newResourceSchemas, dir)
	if err != nil {
		t.Fatal(err)
	}
	if !cr.IsCompatible {
		t.Error("new resource with no baseline should be compatible")
	}
	if len(cr.Checks) != 1 || cr.Checks[0].Status != "new" {
		t.Error("expected check with status 'new'")
	}
}
