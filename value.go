package cmldataresolve

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type TypedValue struct {
	field Field
	raw   any
}

func (v TypedValue) Field() Field {
	return v.field
}

func (v TypedValue) Raw() any {
	return v.raw
}

func (v TypedValue) DataType() string {
	return normalizeDataType(v.field.DataType)
}

func (v TypedValue) Any() (any, error) {
	switch v.DataType() {
	case "string":
		return v.String()
	case "int8", "int16", "int32":
		return v.Int32()
	case "int", "int64":
		return v.Int64()
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return v.Uint64()
	case "float", "float32", "float64":
		return v.Float64()
	case "decimal":
		return v.DecimalString()
	case "bool":
		return v.Bool()
	case "date", "datetime":
		return v.Time()
	case "bytes":
		return v.Bytes()
	case "object":
		return v.Object()
	case "array":
		return v.Array()
	default:
		return v.raw, nil
	}
}

func (v TypedValue) String() (string, error) {
	return rawString(v.raw), nil
}

func (v TypedValue) Int() (int, error) {
	i, err := v.Int64()
	if err != nil {
		return 0, err
	}
	if i > int64(math.MaxInt) || i < int64(math.MinInt) {
		return 0, fmt.Errorf("%w: %s overflows int", ErrConvertFailed, v.field.JSONKey)
	}
	return int(i), nil
}

func (v TypedValue) Int32() (int32, error) {
	i, err := v.Int64()
	if err != nil {
		return 0, err
	}
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0, fmt.Errorf("%w: %s overflows int32", ErrConvertFailed, v.field.JSONKey)
	}
	return int32(i), nil
}

func (v TypedValue) Int64() (int64, error) {
	switch x := v.raw.(type) {
	case json.Number:
		return x.Int64()
	case float64:
		if math.Trunc(x) != x {
			return 0, fmt.Errorf("%w: %s is not integer", ErrConvertFailed, v.field.JSONKey)
		}
		return int64(x), nil
	case float32:
		if math.Trunc(float64(x)) != float64(x) {
			return 0, fmt.Errorf("%w: %s is not integer", ErrConvertFailed, v.field.JSONKey)
		}
		return int64(x), nil
	case int:
		return int64(x), nil
	case int8:
		return int64(x), nil
	case int16:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case int64:
		return x, nil
	case uint:
		return int64(x), nil
	case uint8:
		return int64(x), nil
	case uint16:
		return int64(x), nil
	case uint32:
		return int64(x), nil
	case uint64:
		if x > math.MaxInt64 {
			return 0, fmt.Errorf("%w: %s overflows int64", ErrConvertFailed, v.field.JSONKey)
		}
		return int64(x), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(x), 10, 64)
	default:
		return strconv.ParseInt(rawString(x), 10, 64)
	}
}

func (v TypedValue) Uint64() (uint64, error) {
	switch x := v.raw.(type) {
	case json.Number:
		return strconv.ParseUint(x.String(), 10, 64)
	case uint:
		return uint64(x), nil
	case uint8:
		return uint64(x), nil
	case uint16:
		return uint64(x), nil
	case uint32:
		return uint64(x), nil
	case uint64:
		return x, nil
	case int, int8, int16, int32, int64:
		i, err := (TypedValue{field: v.field, raw: x}).Int64()
		if err != nil {
			return 0, err
		}
		if i < 0 {
			return 0, fmt.Errorf("%w: %s is negative", ErrConvertFailed, v.field.JSONKey)
		}
		return uint64(i), nil
	case string:
		return strconv.ParseUint(strings.TrimSpace(x), 10, 64)
	default:
		return strconv.ParseUint(rawString(x), 10, 64)
	}
}

func (v TypedValue) Float64() (float64, error) {
	switch x := v.raw.(type) {
	case json.Number:
		return x.Float64()
	case float32:
		return float64(x), nil
	case float64:
		return x, nil
	case int, int8, int16, int32, int64:
		i, err := (TypedValue{field: v.field, raw: x}).Int64()
		return float64(i), err
	case uint, uint8, uint16, uint32, uint64:
		u, err := (TypedValue{field: v.field, raw: x}).Uint64()
		return float64(u), err
	case string:
		return strconv.ParseFloat(strings.TrimSpace(x), 64)
	default:
		return strconv.ParseFloat(rawString(x), 64)
	}
}

func (v TypedValue) DecimalString() (string, error) {
	s := strings.TrimSpace(rawString(v.raw))
	if s == "" {
		return "", fmt.Errorf("%w: %s is empty decimal", ErrConvertFailed, v.field.JSONKey)
	}
	if _, ok := new(big.Rat).SetString(s); !ok {
		return "", fmt.Errorf("%w: %s is invalid decimal", ErrConvertFailed, v.field.JSONKey)
	}
	return s, nil
}

func (v TypedValue) DecimalRat() (*big.Rat, error) {
	s, err := v.DecimalString()
	if err != nil {
		return nil, err
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return nil, fmt.Errorf("%w: %s is invalid decimal", ErrConvertFailed, v.field.JSONKey)
	}
	return r, nil
}

func (v TypedValue) Bool() (bool, error) {
	switch x := v.raw.(type) {
	case bool:
		return x, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(x)) {
		case "true", "1", "yes", "y", "on":
			return true, nil
		case "false", "0", "no", "n", "off":
			return false, nil
		}
		return false, fmt.Errorf("%w: %s is invalid bool", ErrConvertFailed, v.field.JSONKey)
	default:
		s := strings.TrimSpace(rawString(x))
		return strconv.ParseBool(s)
	}
}

func (v TypedValue) Time() (time.Time, error) {
	switch x := v.raw.(type) {
	case time.Time:
		return x, nil
	case string:
		return parseTime(strings.TrimSpace(x))
	default:
		return parseTime(rawString(x))
	}
}

func (v TypedValue) Bytes() ([]byte, error) {
	switch x := v.raw.(type) {
	case []byte:
		out := make([]byte, len(x))
		copy(out, x)
		return out, nil
	case string:
		s := strings.TrimSpace(x)
		if decoded, err := base64.StdEncoding.DecodeString(s); err == nil {
			return decoded, nil
		}
		return []byte(x), nil
	default:
		return []byte(rawString(x)), nil
	}
}

func (v TypedValue) Object() (map[string]any, error) {
	if m, ok := mapFromAny(v.raw); ok {
		return m, nil
	}
	return nil, fmt.Errorf("%w: %s is not object", ErrConvertFailed, v.field.JSONKey)
}

func (v TypedValue) Array() ([]any, error) {
	switch x := v.raw.(type) {
	case []any:
		return x, nil
	case string:
		var out []any
		if err := json.Unmarshal([]byte(strings.TrimSpace(x)), &out); err != nil {
			return nil, err
		}
		return out, nil
	default:
		normalized, err := normalizeAny(x)
		if err != nil {
			return nil, err
		}
		out, ok := normalized.([]any)
		if !ok {
			return nil, fmt.Errorf("%w: %s is not array", ErrConvertFailed, v.field.JSONKey)
		}
		return out, nil
	}
}

func rawString(value any) string {
	switch x := value.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return x.String()
	case []byte:
		return string(x)
	default:
		b, err := json.Marshal(x)
		if err == nil && (strings.HasPrefix(string(b), "{") || strings.HasPrefix(string(b), "[")) {
			return string(b)
		}
		return fmt.Sprint(x)
	}
}

func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%w: empty time", ErrConvertFailed)
	}
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
	}
	var lastErr error
	for _, layout := range formats {
		t, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, fmt.Errorf("%w: %s: %v", ErrConvertFailed, value, lastErr)
}
