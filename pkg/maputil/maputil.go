// Package maputil provides type-safe accessors for map[string]interface{}
// data structures, used across renderer, sbom, report, and other packages
// that consume JSON-unmarshaled component-architecture.json data.
package maputil

import "fmt"

// GetStr returns a string value from a map, or fallback if missing or wrong type.
func GetStr(m map[string]interface{}, key, fallback string) string {
	if m == nil {
		return fallback
	}
	v, ok := m[key]
	if !ok {
		return fallback
	}
	s, ok := v.(string)
	if !ok {
		return fallback
	}
	return s
}

// GetMap returns a nested map from a map, or nil if missing or wrong type.
func GetMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	mm, ok := v.(map[string]interface{})
	if ok {
		return mm
	}
	return nil
}

// GetSlice returns a slice of maps from a map, handling both
// []map[string]interface{} and []interface{} (from JSON unmarshal).
func GetSlice(m map[string]interface{}, key string) []map[string]interface{} {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	switch typed := v.(type) {
	case []map[string]interface{}:
		return typed
	case []interface{}:
		out := make([]map[string]interface{}, 0, len(typed))
		for _, item := range typed {
			if mm, ok := item.(map[string]interface{}); ok {
				out = append(out, mm)
			}
		}
		return out
	}
	return nil
}

// GetStringSlice returns a string slice from a map, handling both
// []string and []interface{} (from JSON unmarshal).
func GetStringSlice(m map[string]interface{}, key string) []string {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	switch typed := v.(type) {
	case []string:
		return typed
	case []interface{}:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// GetInt returns an int value from a map, handling float64 (JSON numbers),
// int, and int64 types.
func GetInt(m map[string]interface{}, key string) int {
	if m == nil {
		return 0
	}
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	}
	return 0
}

// GetBool returns a bool value from a map, or fallback if missing or wrong type.
func GetBool(m map[string]interface{}, key string, fallback bool) bool {
	if m == nil {
		return fallback
	}
	v, ok := m[key]
	if !ok {
		return fallback
	}
	b, ok := v.(bool)
	if !ok {
		return fallback
	}
	return b
}

// MapVal returns a string representation of a map value, or "-" if missing.
func MapVal(m map[string]interface{}, key string) string {
	if m == nil {
		return "-"
	}
	v, ok := m[key]
	if !ok {
		return "-"
	}
	return fmt.Sprintf("%v", v)
}
