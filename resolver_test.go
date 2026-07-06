package cmldataresolve

import (
	"encoding/json"
	"testing"
)

func TestResolverInt32FromFormToolFields(t *testing.T) {
	contract := map[string]any{
		"contractVersion": "formtool.v1",
		"fields": []any{
			map[string]any{
				"fieldCode":   "aaa",
				"displayName": "AAA",
				"jsonKey":     "aaa",
				"dataType":    "int32",
			},
		},
	}
	resolver, err := New(contract, map[string]any{"aaa": "45"})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	got, err := resolver.Int32("aaa")
	if err != nil {
		t.Fatalf("Int32() error = %v", err)
	}
	if got != 45 {
		t.Fatalf("Int32() = %d, want 45", got)
	}
	anyValue, err := resolver.Any("AAA")
	if err != nil {
		t.Fatalf("Any() error = %v", err)
	}
	if anyValue != int32(45) {
		t.Fatalf("Any() = %#v, want int32(45)", anyValue)
	}
}

func TestResolverFromSchemaAndNestedPath(t *testing.T) {
	contract := map[string]any{
		"schema": map[string]any{
			"properties": map[string]any{
				"employee.age": map[string]any{
					"type":   "integer",
					"format": "int32",
					"title":  "Age",
				},
			},
		},
	}
	resolver, err := New(contract, map[string]any{"employee": map[string]any{"age": json.Number("36")}})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	got, err := resolver.Int("employee.age")
	if err != nil {
		t.Fatalf("Int() error = %v", err)
	}
	if got != 36 {
		t.Fatalf("Int() = %d, want 36", got)
	}
}

func TestResolverFromInstanceContractData(t *testing.T) {
	instance := map[string]any{
		"fields": []any{
			map[string]any{"jsonKey": "ok", "dataType": "bool", "name": "Approved"},
			map[string]any{"jsonKey": "when", "dataType": "datetime"},
		},
		"data": map[string]any{"ok": "1", "when": "2026-07-02 12:30:00"},
	}
	resolver, err := FromInstanceContract(instance)
	if err != nil {
		t.Fatalf("FromInstanceContract() error = %v", err)
	}
	ok, err := resolver.Bool("Approved")
	if err != nil {
		t.Fatalf("Bool() error = %v", err)
	}
	if !ok {
		t.Fatalf("Bool() = false, want true")
	}
	when, err := resolver.Value("when")
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	if _, err := when.Time(); err != nil {
		t.Fatalf("Time() error = %v", err)
	}
}

func TestResolverArrayOptionsUseCodeAndTypedValue(t *testing.T) {
	contract := map[string]any{
		"fields": []any{
			map[string]any{
				"fieldCode":     "levels",
				"jsonKey":       "levels",
				"name":          "Levels",
				"dataType":      "array",
				"arrayItemType": "int32",
				"options": []any{
					map[string]any{"code": "L1", "label": "一级", "value": "1"},
					map[string]any{"code": "L2", "label": "二级", "value": "2"},
				},
			},
		},
	}
	resolver, err := New(contract, map[string]any{"levels": []any{"L2", "L1"}})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	codes, err := resolver.SelectionCodes("levels")
	if err != nil {
		t.Fatalf("SelectionCodes() error = %v", err)
	}
	if got, want := codes, []string{"L2", "L1"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("SelectionCodes() = %#v, want %#v", got, want)
	}
	labels, err := resolver.SelectionLabels("levels")
	if err != nil {
		t.Fatalf("SelectionLabels() error = %v", err)
	}
	if got, want := labels, []string{"二级", "一级"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("SelectionLabels() = %#v, want %#v", got, want)
	}
	values, err := resolver.Int32Slice("levels")
	if err != nil {
		t.Fatalf("Int32Slice() error = %v", err)
	}
	if got, want := values, []int32{2, 1}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Int32Slice() = %#v, want %#v", got, want)
	}
	anyValues, err := resolver.Any("levels")
	if err != nil {
		t.Fatalf("Any() error = %v", err)
	}
	list, ok := anyValues.([]any)
	if !ok || len(list) != 2 || list[0] != int32(2) || list[1] != int32(1) {
		t.Fatalf("Any(levels) = %#v, want []any{int32(2), int32(1)}", anyValues)
	}
}

func TestResolverArrayOptionsFromSchemaExtensions(t *testing.T) {
	contract := map[string]any{
		"schema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"deptCodes": map[string]any{
					"type":                     "array",
					"title":                    "Departments",
					"x-camellia-dataType":      "array",
					"x-camellia-arrayItemType": "int32",
					"items": map[string]any{
						"type": "string",
						"oneOf": []any{
							map[string]any{"const": "dev", "title": "Development", "x-camellia-value": "3"},
							map[string]any{"const": "qa", "title": "Quality", "x-camellia-value": "7"},
						},
					},
				},
			},
		},
	}
	resolver, err := New(contract, map[string]any{"deptCodes": []any{"qa", "dev"}})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	labels, err := resolver.SelectionLabels("deptCodes")
	if err != nil {
		t.Fatalf("SelectionLabels() error = %v", err)
	}
	if got, want := labels, []string{"Quality", "Development"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("SelectionLabels() = %#v, want %#v", got, want)
	}
	values, err := resolver.Int32Slice("deptCodes")
	if err != nil {
		t.Fatalf("Int32Slice() error = %v", err)
	}
	if got, want := values, []int32{7, 3}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Int32Slice() = %#v, want %#v", got, want)
	}
}

func TestResolverScalarOptionUsesCodeAndTypedValue(t *testing.T) {
	contract := map[string]any{
		"fields": []any{
			map[string]any{
				"fieldCode":  "department",
				"jsonKey":    "department",
				"name":       "Department",
				"dataType":   "string",
				"widgetType": "select",
				"options": []any{
					map[string]any{"code": "D003", "label": "R&D", "value": "003"},
					map[string]any{"code": "D004", "label": "Manufacturing", "value": "004"},
				},
			},
			map[string]any{
				"fieldCode":  "ageBand",
				"jsonKey":    "ageBand",
				"name":       "Age band",
				"dataType":   "int32",
				"widgetType": "radio",
				"options": []any{
					map[string]any{"code": "A", "label": "Adult", "value": "18"},
				},
			},
			map[string]any{
				"fieldCode":  "rawCode",
				"jsonKey":    "rawCode",
				"name":       "Raw code",
				"dataType":   "string",
				"widgetType": "input",
				"options": []any{
					map[string]any{"code": "X", "label": "Mapped", "value": "mapped-value"},
				},
			},
		},
	}
	resolver, err := New(contract, map[string]any{"department": "D003", "ageBand": "A", "rawCode": "X"})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	dept, err := resolver.String("department")
	if err != nil {
		t.Fatalf("String(department) error = %v", err)
	}
	if dept != "003" {
		t.Fatalf("String(department) = %q, want 003", dept)
	}
	labels, err := resolver.SelectionLabels("department")
	if err != nil {
		t.Fatalf("SelectionLabels(department) error = %v", err)
	}
	if len(labels) != 1 || labels[0] != "R&D" {
		t.Fatalf("SelectionLabels(department) = %#v, want R&D", labels)
	}
	age, err := resolver.Int32("ageBand")
	if err != nil {
		t.Fatalf("Int32(ageBand) error = %v", err)
	}
	if age != 18 {
		t.Fatalf("Int32(ageBand) = %d, want 18", age)
	}
	raw, err := resolver.String("rawCode")
	if err != nil {
		t.Fatalf("String(rawCode) error = %v", err)
	}
	if raw != "X" {
		t.Fatalf("String(rawCode) = %q, want raw input X", raw)
	}
}

func TestResolverFreeArrayAllowsRepeatedTypedValues(t *testing.T) {
	contract := map[string]any{
		"fields": []any{
			map[string]any{
				"fieldCode":     "samples",
				"jsonKey":       "samples",
				"name":          "Samples",
				"dataType":      "array",
				"arrayItemType": "int32",
				"widgetType":    "textarea",
			},
		},
	}
	resolver, err := New(contract, map[string]any{"samples": []any{float64(1), float64(2), float64(2), float64(3)}})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	values, err := resolver.Int32Slice("samples")
	if err != nil {
		t.Fatalf("Int32Slice(samples) error = %v", err)
	}
	if got, want := values, []int32{1, 2, 2, 3}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] || got[3] != want[3] {
		t.Fatalf("Int32Slice(samples) = %#v, want %#v", got, want)
	}
}

func TestObjectFieldTypeIsNotStandard(t *testing.T) {
	_, err := New(
		map[string]any{"fields": []any{map[string]any{"jsonKey": "payload", "dataType": "object"}}},
		map[string]any{"payload": map[string]any{"a": 1}},
	)
	if err == nil {
		t.Fatal("New() error = nil, want invalid contract because object fields are not standard")
	}
}

func TestMissingField(t *testing.T) {
	resolver, err := New(map[string]any{"fields": []any{map[string]any{"jsonKey": "a", "dataType": "string"}}}, map[string]any{"a": "x"})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := resolver.String("b"); err == nil {
		t.Fatal("String(b) error = nil, want error")
	}
}
