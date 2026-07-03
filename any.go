package cmldataresolve

import (
	"encoding/json"
	"fmt"
	"strings"
)

func normalizeAny(value any) (any, error) {
	switch v := value.(type) {
	case nil:
		return map[string]any{}, nil
	case []byte:
		return decodeJSONBytes(v)
	case string:
		return decodeJSONString(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return decodeJSONBytes(b)
	}
}

func decodeJSONString(value string) (any, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return map[string]any{}, nil
	}
	return decodeJSONBytes([]byte(trimmed))
}

func decodeJSONBytes(value []byte) (any, error) {
	var out any
	decoder := json.NewDecoder(strings.NewReader(string(value)))
	decoder.UseNumber()
	if err := decoder.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func mapFromAny(value any) (map[string]any, bool) {
	if value == nil {
		return nil, false
	}
	if m, ok := value.(map[string]any); ok {
		return m, true
	}
	normalized, err := normalizeAny(value)
	if err != nil {
		return nil, false
	}
	m, ok := normalized.(map[string]any)
	return m, ok
}

func dataMapFromAny(value any) (map[string]any, error) {
	normalized, err := normalizeAny(value)
	if err != nil {
		return nil, err
	}
	m, ok := normalized.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: data must be a JSON object", ErrInvalidContract)
	}
	return m, nil
}

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		raw, ok := m[key]
		if !ok || raw == nil {
			continue
		}
		value := strings.TrimSpace(fmt.Sprint(raw))
		if value != "" && value != "<nil>" {
			return value
		}
	}
	return ""
}

func extractDataFromContract(contract map[string]any) map[string]any {
	for _, key := range []string{"data", "values", "formData", "form_data"} {
		if data, ok := mapFromAny(contract[key]); ok {
			return data
		}
	}
	for _, key := range []string{"dataJson", "data_json", "valuesJson", "values_json"} {
		if s, ok := contract[key].(string); ok && strings.TrimSpace(s) != "" {
			if data, ok := mapFromAny(s); ok {
				return data
			}
		}
	}
	return nil
}
