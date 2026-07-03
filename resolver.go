package cmldataresolve

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Resolver struct {
	fields  []Field
	byAlias map[string]Field
	data    map[string]any
}

func New(contract any, data any) (*Resolver, error) {
	contractMap, ok := mapFromAny(contract)
	if !ok {
		return nil, fmt.Errorf("%w: contract must be a JSON object", ErrInvalidContract)
	}
	fields := fieldsFromContract(contractMap)
	if len(fields) == 0 {
		return nil, fmt.Errorf("%w: no fields found in contract", ErrInvalidContract)
	}
	dataMap, err := dataMapFromAny(data)
	if data == nil {
		dataMap = extractDataFromContract(contractMap)
		if dataMap == nil {
			dataMap = map[string]any{}
		}
	} else if err != nil {
		return nil, err
	}
	r := &Resolver{
		fields:  fields,
		byAlias: map[string]Field{},
		data:    dataMap,
	}
	for _, field := range fields {
		r.addAlias(field.JSONKey, field)
		r.addAlias(field.Code, field)
		r.addAlias(field.Name, field)
	}
	return r, nil
}

func FromJSON(contractJSON []byte, dataJSON []byte) (*Resolver, error) {
	return New(contractJSON, dataJSON)
}

func FromInstanceContract(instanceContract any) (*Resolver, error) {
	contractMap, ok := mapFromAny(instanceContract)
	if !ok {
		return nil, fmt.Errorf("%w: instance contract must be a JSON object", ErrInvalidContract)
	}
	data := extractDataFromContract(contractMap)
	return New(contractMap, data)
}

func (r *Resolver) Fields() []Field {
	out := make([]Field, len(r.fields))
	copy(out, r.fields)
	return out
}

func (r *Resolver) Field(key string) (Field, bool) {
	field, ok := r.byAlias[aliasKey(key)]
	return field, ok
}

func (r *Resolver) Value(key string) (TypedValue, error) {
	field, ok := r.Field(key)
	if !ok {
		return TypedValue{}, fmt.Errorf("%w: %s", ErrFieldNotFound, key)
	}
	raw, ok := lookupPath(r.data, field.JSONKey)
	if !ok || raw == nil {
		return TypedValue{}, fmt.Errorf("%w: %s", ErrValueMissing, field.JSONKey)
	}
	return TypedValue{field: field, raw: raw}, nil
}

func (r *Resolver) Raw(key string) (any, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.Raw(), nil
}

func (r *Resolver) Any(key string) (any, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.Any()
}

func (r *Resolver) String(key string) (string, error) {
	value, err := r.Value(key)
	if err != nil {
		return "", err
	}
	return value.String()
}

func (r *Resolver) Int(key string) (int, error) {
	value, err := r.Value(key)
	if err != nil {
		return 0, err
	}
	return value.Int()
}

func (r *Resolver) Int32(key string) (int32, error) {
	value, err := r.Value(key)
	if err != nil {
		return 0, err
	}
	return value.Int32()
}

func (r *Resolver) Int64(key string) (int64, error) {
	value, err := r.Value(key)
	if err != nil {
		return 0, err
	}
	return value.Int64()
}

func (r *Resolver) Uint64(key string) (uint64, error) {
	value, err := r.Value(key)
	if err != nil {
		return 0, err
	}
	return value.Uint64()
}

func (r *Resolver) Float64(key string) (float64, error) {
	value, err := r.Value(key)
	if err != nil {
		return 0, err
	}
	return value.Float64()
}

func (r *Resolver) Bool(key string) (bool, error) {
	value, err := r.Value(key)
	if err != nil {
		return false, err
	}
	return value.Bool()
}

func (r *Resolver) Array(key string) ([]any, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.Array()
}

func (r *Resolver) ArrayValues(key string) ([]any, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.ArrayValues()
}

func (r *Resolver) SelectionCodes(key string) ([]string, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.SelectionCodes()
}

func (r *Resolver) SelectionLabels(key string) ([]string, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.SelectionLabels()
}

func (r *Resolver) StringSlice(key string) ([]string, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.StringSlice()
}

func (r *Resolver) Int32Slice(key string) ([]int32, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.Int32Slice()
}

func (r *Resolver) Int64Slice(key string) ([]int64, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.Int64Slice()
}

func (r *Resolver) Uint64Slice(key string) ([]uint64, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.Uint64Slice()
}

func (r *Resolver) Float64Slice(key string) ([]float64, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.Float64Slice()
}

func (r *Resolver) BoolSlice(key string) ([]bool, error) {
	value, err := r.Value(key)
	if err != nil {
		return nil, err
	}
	return value.BoolSlice()
}

func (r *Resolver) JSON() ([]byte, error) {
	return json.Marshal(r.data)
}

func (r *Resolver) addAlias(alias string, field Field) {
	key := aliasKey(alias)
	if key == "" {
		return
	}
	if _, exists := r.byAlias[key]; !exists {
		r.byAlias[key] = field
	}
}

func aliasKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
