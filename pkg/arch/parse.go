package arch

import "encoding/json"

// Parse converts a raw architecture JSON map into typed Data.
// Returns a zero-value Data (not nil) if raw is nil or empty.
func Parse(raw map[string]interface{}) (*Data, error) {
	if len(raw) == 0 {
		return &Data{}, nil
	}

	// Round-trip through JSON to leverage struct tags for field mapping.
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var data Data
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
