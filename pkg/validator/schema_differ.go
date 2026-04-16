package validator

import (
	"fmt"
	"sort"
)

// SchemaChange represents a single change detected between two schemas.
type SchemaChange struct {
	ChangeType  string // e.g. "removed_field", "type_changed", "required_added", "field_added", "enum_changed", etc.
	Field       string // JSON path of the affected field
	Description string
}

// SchemaDiffResult aggregates all changes between two schemas.
type SchemaDiffResult struct {
	BreakingChanges []SchemaChange
	AdditiveChanges []SchemaChange
}

// IsCompatible returns true when there are no breaking changes.
func (r *SchemaDiffResult) IsCompatible() bool {
	return len(r.BreakingChanges) == 0
}

// DiffSchemas compares oldSchema and newSchema, returning classified changes.
func DiffSchemas(oldSchema, newSchema map[string]interface{}) *SchemaDiffResult {
	result := &SchemaDiffResult{}
	visited := make(map[visitKey]bool)
	diffObject(oldSchema, newSchema, "", result, 0, visited)
	return result
}

const maxDiffDepth = 50

// visitKey tracks which (old, new) map pairs have been compared to detect cycles.
type visitKey struct {
	oldPath string
	newPath string
}

func diffObject(old, new map[string]interface{}, path string, result *SchemaDiffResult, depth int, visited map[visitKey]bool) {
	if old == nil || new == nil {
		return
	}
	if depth > maxDiffDepth {
		return
	}

	// Cycle detection: skip if we've already compared this path pair
	key := visitKey{oldPath: path, newPath: path}
	if visited[key] {
		return
	}
	visited[key] = true

	// Type change check
	oldType, oldHas := old["type"].(string)
	newType, newHas := new["type"].(string)
	if oldHas && newHas && oldType != newType {
		p := path
		if p == "" {
			p = "/"
		}
		result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
			ChangeType:  "type_changed",
			Field:       p,
			Description: fmt.Sprintf("Type changed from '%s' to '%s'", oldType, newType),
		})
		return
	}

	oldProps := toMap(old["properties"])
	newProps := toMap(new["properties"])

	// Removed fields
	for name, oldProp := range oldProps {
		propPath := path + ".properties." + name
		newProp, exists := newProps[name]
		if !exists {
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "removed_field",
				Field:       propPath,
				Description: fmt.Sprintf("Field '%s' was removed", name),
			})
		} else {
			oldMap := toMap(oldProp)
			newMap := toMap(newProp)
			if oldMap != nil && newMap != nil {
				diffObject(oldMap, newMap, propPath, result, depth+1, visited)
			}
		}
	}

	// Added fields
	oldRequired := toStringSet(old["required"])
	newRequired := toStringSet(new["required"])
	for name := range newProps {
		propPath := path + ".properties." + name
		if _, exists := oldProps[name]; exists {
			continue
		}
		if newRequired[name] && !oldRequired[name] {
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "required_added",
				Field:       propPath,
				Description: fmt.Sprintf("New required field '%s' added", name),
			})
		} else {
			result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
				ChangeType:  "field_added",
				Field:       propPath,
				Description: fmt.Sprintf("New optional field '%s' added", name),
			})
		}
	}

	// Existing fields that became required
	for req := range newRequired {
		if oldRequired[req] {
			continue
		}
		if _, inOld := oldProps[req]; inOld {
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "required_added",
				Field:       fmt.Sprintf("%s.required.%s", path, req),
				Description: fmt.Sprintf("Existing field '%s' became required", req),
			})
		}
	}

	// Enum comparison
	diffEnum(old, new, path, result)

	// additionalProperties comparison
	diffAdditionalProperties(old, new, path, result, depth, visited)

	// Composition keywords: oneOf, allOf, anyOf
	for _, keyword := range []string{"oneOf", "allOf", "anyOf"} {
		diffComposition(old, new, path, keyword, result)
	}

	// items (array sub-schema)
	oldItems := toMap(old["items"])
	newItems := toMap(new["items"])
	if oldItems != nil && newItems != nil {
		diffObject(oldItems, newItems, path+".items", result, depth+1, visited)
	}
}

func diffEnum(old, new map[string]interface{}, path string, result *SchemaDiffResult) {
	oldEnum, oldOk := old["enum"]
	newEnum, newOk := new["enum"]
	if !oldOk || !newOk {
		return
	}
	oldSet := toInterfaceStringSet(oldEnum)
	newSet := toInterfaceStringSet(newEnum)

	var removed, added []string
	for v := range oldSet {
		if !newSet[v] {
			removed = append(removed, v)
		}
	}
	for v := range newSet {
		if !oldSet[v] {
			added = append(added, v)
		}
	}
	sort.Strings(removed)
	sort.Strings(added)

	if len(removed) > 0 {
		result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
			ChangeType:  "enum_changed",
			Field:       path + ".enum",
			Description: fmt.Sprintf("Enum values removed: %v", removed),
		})
	}
	if len(added) > 0 {
		result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
			ChangeType:  "enum_changed",
			Field:       path + ".enum",
			Description: fmt.Sprintf("Enum values added: %v", added),
		})
	}
}

func diffAdditionalProperties(old, new map[string]interface{}, path string, result *SchemaDiffResult, depth int, visited map[visitKey]bool) {
	oldAddl, oldHas := old["additionalProperties"]
	newAddl, newHas := new["additionalProperties"]

	apPath := path + ".additionalProperties"

	if oldHas && newHas {
		oldMap := toMap(oldAddl)
		newMap := toMap(newAddl)
		oldBool, oldIsBool := oldAddl.(bool)
		newBool, newIsBool := newAddl.(bool)

		switch {
		case oldMap != nil && newMap != nil:
			diffObject(oldMap, newMap, apPath, result, depth+1, visited)

		case oldMap != nil && newIsBool && !newBool:
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "additional_properties_restricted",
				Field:       apPath,
				Description: "additionalProperties changed from schema to false",
			})
		case oldMap != nil && newIsBool && newBool:
			result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
				ChangeType:  "additional_properties_relaxed",
				Field:       apPath,
				Description: "additionalProperties changed from schema to true (relaxed)",
			})
		case oldIsBool && !oldBool && newMap != nil:
			result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
				ChangeType:  "additional_properties_relaxed",
				Field:       apPath,
				Description: "additionalProperties changed from false to schema (relaxed)",
			})
		case oldIsBool && oldBool && newMap != nil:
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "additional_properties_restricted",
				Field:       apPath,
				Description: "additionalProperties changed from true to schema (restricting)",
			})
		case oldIsBool && oldBool && newIsBool && !newBool:
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "additional_properties_restricted",
				Field:       apPath,
				Description: "additionalProperties changed from true to false",
			})
		case oldIsBool && !oldBool && newIsBool && newBool:
			result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
				ChangeType:  "additional_properties_relaxed",
				Field:       apPath,
				Description: "additionalProperties changed from false to true",
			})
		}
	} else if oldHas && !newHas {
		// Removing additionalProperties: absent defaults to true in JSON Schema,
		// so removing false is a relaxation, removing true is a no-op.
		oldBool, oldIsBool := oldAddl.(bool)
		if oldIsBool && !oldBool {
			result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
				ChangeType:  "additional_properties_relaxed",
				Field:       apPath,
				Description: "additionalProperties removed (was false, now defaults to allowed)",
			})
		} else if !oldIsBool {
			// Was a schema object restricting additional properties, now removed (relaxed)
			result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
				ChangeType:  "additional_properties_relaxed",
				Field:       apPath,
				Description: "additionalProperties schema removed (now defaults to allowed)",
			})
		}
	} else if !oldHas && newHas {
		newBool, newIsBool := newAddl.(bool)
		if newIsBool && !newBool {
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "additional_properties_restricted",
				Field:       apPath,
				Description: "additionalProperties set to false (was unset/allowed)",
			})
		}
	}
}

func diffComposition(old, new map[string]interface{}, path, keyword string, result *SchemaDiffResult) {
	oldList := toSlice(old[keyword])
	newList := toSlice(new[keyword])

	kwPath := path + "." + keyword

	if len(oldList) > 0 && len(newList) == 0 {
		result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
			ChangeType:  "composition_removed",
			Field:       kwPath,
			Description: fmt.Sprintf("%s removed (had %d options)", keyword, len(oldList)),
		})
		return
	}

	if len(oldList) > 0 && len(newList) > 0 {
		if len(newList) < len(oldList) {
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "composition_restricted",
				Field:       kwPath,
				Description: fmt.Sprintf("%s options reduced from %d to %d", keyword, len(oldList), len(newList)),
			})
		} else if len(newList) > len(oldList) {
			result.AdditiveChanges = append(result.AdditiveChanges, SchemaChange{
				ChangeType:  "composition_expanded",
				Field:       kwPath,
				Description: fmt.Sprintf("%s options expanded from %d to %d", keyword, len(oldList), len(newList)),
			})
		}

		// Check type changes across sub-schemas
		oldTypes := map[string]bool{}
		newTypes := map[string]bool{}
		for _, item := range oldList {
			m := toMap(item)
			if m != nil {
				if t, ok := m["type"].(string); ok {
					oldTypes[t] = true
				}
			}
		}
		for _, item := range newList {
			m := toMap(item)
			if m != nil {
				if t, ok := m["type"].(string); ok {
					newTypes[t] = true
				}
			}
		}
		var removedTypes []string
		for t := range oldTypes {
			if !newTypes[t] {
				removedTypes = append(removedTypes, t)
			}
		}
		sort.Strings(removedTypes)
		if len(removedTypes) > 0 {
			result.BreakingChanges = append(result.BreakingChanges, SchemaChange{
				ChangeType:  "composition_type_changed",
				Field:       kwPath,
				Description: fmt.Sprintf("%s type options removed: %v", keyword, removedTypes),
			})
		}
	}
}

// --- helpers ---

// toMap converts an interface{} to map[string]interface{} if possible.
func toMap(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// toSlice converts an interface{} to []interface{} if possible.
func toSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return nil
}

// toStringSet converts an interface{} (expected []interface{} of strings) to a set.
func toStringSet(v interface{}) map[string]bool {
	set := map[string]bool{}
	if v == nil {
		return set
	}
	if arr, ok := v.([]interface{}); ok {
		for _, item := range arr {
			if s, ok := item.(string); ok {
				set[s] = true
			}
		}
	}
	return set
}

// toInterfaceStringSet converts an interface{} (expected []interface{}) to a set of type-qualified string keys.
// Type information is preserved so that integer 1 and string "1" are treated as distinct values.
func toInterfaceStringSet(v interface{}) map[string]bool {
	set := map[string]bool{}
	if arr, ok := v.([]interface{}); ok {
		for _, item := range arr {
			set[fmt.Sprintf("%T:%v", item, item)] = true
		}
	}
	return set
}
