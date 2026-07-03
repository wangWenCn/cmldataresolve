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

func TestMissingField(t *testing.T) {
	resolver, err := New(map[string]any{"fields": []any{map[string]any{"jsonKey": "a", "dataType": "string"}}}, map[string]any{"a": "x"})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := resolver.String("b"); err == nil {
		t.Fatal("String(b) error = nil, want error")
	}
}
