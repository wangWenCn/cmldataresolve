package cmldataresolve

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Field struct {
	Code          string   `json:"code,omitempty"`
	Name          string   `json:"name,omitempty"`
	JSONKey       string   `json:"jsonKey"`
	DataType      string   `json:"dataType"`
	ArrayItemType string   `json:"arrayItemType,omitempty"`
	WidgetType    string   `json:"widgetType,omitempty"`
	Options       []Option `json:"options,omitempty"`
}

type Option struct {
	Code     string `json:"code"`
	Label    string `json:"label"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled,omitempty"`
}

func normalizeDataType(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	v = strings.ReplaceAll(v, " ", "")
	v = strings.ReplaceAll(v, "_", "")
	v = strings.ReplaceAll(v, "-", "")
	switch v {
	case "", "string", "text", "varchar", "char":
		return "string"
	case "int", "integer", "signed":
		return "int"
	case "int8", "i8":
		return "int8"
	case "int16", "i16":
		return "int16"
	case "int32", "i32":
		return "int32"
	case "int64", "long", "i64":
		return "int64"
	case "uint", "unsigned":
		return "uint"
	case "uint8", "u8":
		return "uint8"
	case "uint16", "u16":
		return "uint16"
	case "uint32", "u32":
		return "uint32"
	case "uint64", "ulong", "u64":
		return "uint64"
	case "float", "number":
		return "float"
	case "float32", "single":
		return "float32"
	case "float64", "double":
		return "float64"
	case "decimal", "money", "numeric":
		return "decimal"
	case "bool", "boolean":
		return "bool"
	case "date":
		return "date"
	case "datetime", "timestamp", "time":
		return "datetime"
	case "bytes", "byte", "binary", "blob", "[]byte", "bytearray":
		return "bytes"
	case "array", "slice", "list", "jsonarray":
		return "array"
	default:
		return strings.TrimSpace(value)
	}
}

func fieldFromMap(raw map[string]any) (Field, bool) {
	jsonKey := firstString(raw, "jsonKey", "json_key", "key", "fieldKey")
	code := firstString(raw, "fieldCode", "field_code", "code", "dataPointCode", "data_point_code")
	name := firstString(raw, "displayName", "display_name", "name", "title", "label")
	dataType := firstString(raw, "dataType", "data_type", "type", "valueType", "value_type")
	arrayItemType := firstString(raw, "arrayItemType", "array_item_type", "itemType", "item_type")
	widgetType := firstString(raw, "widgetType", "widget_type", "controlType", "control_type")
	if jsonKey == "" {
		jsonKey = code
	}
	if jsonKey == "" {
		return Field{}, false
	}
	return Field{
		Code:          code,
		Name:          name,
		JSONKey:       jsonKey,
		DataType:      normalizeDataType(dataType),
		ArrayItemType: normalizeDataType(arrayItemType),
		WidgetType:    strings.TrimSpace(widgetType),
		Options:       optionsFromAny(raw["options"]),
	}, true
}

func fieldsFromContract(contract map[string]any) []Field {
	var out []Field
	out = append(out, fieldsFromArray(contract["fields"])...)
	out = append(out, fieldsFromArray(contract["dataPoints"])...)
	out = append(out, fieldsFromArray(contract["data_points"])...)
	out = append(out, fieldsFromSchema(contract["schema"])...)
	for _, key := range []string{"contract", "templateContract", "template_contract", "formContract", "form_contract"} {
		nested, ok := mapFromAny(contract[key])
		if ok {
			out = append(out, fieldsFromContract(nested)...)
			continue
		}
		if s, ok := contract[key].(string); ok && strings.TrimSpace(s) != "" {
			var decoded map[string]any
			if json.Unmarshal([]byte(s), &decoded) == nil {
				out = append(out, fieldsFromContract(decoded)...)
			}
		}
	}
	return dedupeFields(out)
}

func fieldsFromArray(value any) []Field {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	out := make([]Field, 0, len(items))
	for _, item := range items {
		raw, ok := mapFromAny(item)
		if !ok {
			continue
		}
		if field, ok := fieldFromMap(raw); ok {
			out = append(out, field)
		}
	}
	return out
}

func fieldsFromSchema(value any) []Field {
	schema, ok := mapFromAny(value)
	if !ok {
		return nil
	}
	props, ok := mapFromAny(schema["properties"])
	if !ok {
		return nil
	}
	out := make([]Field, 0, len(props))
	for jsonKey, propValue := range props {
		prop, _ := mapFromAny(propValue)
		field := Field{JSONKey: jsonKey, DataType: normalizeDataType(schemaPropDataType(prop))}
		if x, ok := mapFromAny(prop["x-camellia"]); ok {
			field.Code = firstString(x, "fieldCode", "code", "dataPointCode")
			field.Name = firstString(x, "displayName", "name", "title", "label")
			if dt := firstString(x, "dataType", "data_type"); dt != "" {
				field.DataType = normalizeDataType(dt)
			}
			if itemType := firstString(x, "arrayItemType", "array_item_type"); itemType != "" {
				field.ArrayItemType = normalizeDataType(itemType)
			}
			if widgetType := firstString(x, "widgetType", "widget_type"); widgetType != "" {
				field.WidgetType = strings.TrimSpace(widgetType)
			}
		}
		if itemType := firstString(prop, "x-camellia-arrayItemType", "x-camellia-arrayitemtype", "arrayItemType", "array_item_type"); itemType != "" {
			field.ArrayItemType = normalizeDataType(itemType)
		}
		field.Options = optionsFromSchemaProperty(prop)
		if field.Name == "" {
			field.Name = firstString(prop, "title", "name", "description")
		}
		if field.Code == "" {
			field.Code = firstString(prop, "fieldCode", "code")
		}
		if field.WidgetType == "" {
			field.WidgetType = firstString(prop, "x-camellia-widgetType", "x-camellia-widgettype", "widgetType", "widget_type")
		}
		out = append(out, field)
	}
	return out
}

func schemaPropDataType(prop map[string]any) string {
	if prop == nil {
		return "string"
	}
	if dt := firstString(prop, "x-camellia-dataType", "x-camellia-datatype", "dataType", "data_type"); dt != "" {
		return dt
	}
	t := strings.TrimSpace(fmt.Sprint(prop["type"]))
	format := strings.TrimSpace(fmt.Sprint(prop["format"]))
	switch t {
	case "integer":
		if format == "int32" {
			return "int32"
		}
		if format == "int64" {
			return "int64"
		}
		return "int"
	case "number":
		if format == "float" {
			return "float32"
		}
		if format == "double" {
			return "float64"
		}
		return "float"
	case "boolean":
		return "bool"
	case "array":
		return "array"
	case "string":
		if format == "date" {
			return "date"
		}
		if format == "date-time" || format == "datetime" {
			return "datetime"
		}
		if format == "binary" || format == "byte" {
			return "bytes"
		}
	}
	return t
}

func dedupeFields(fields []Field) []Field {
	seen := map[string]struct{}{}
	out := make([]Field, 0, len(fields))
	for _, field := range fields {
		field.JSONKey = strings.TrimSpace(field.JSONKey)
		field.Code = strings.TrimSpace(field.Code)
		field.Name = strings.TrimSpace(field.Name)
		field.DataType = normalizeDataType(field.DataType)
		if strings.EqualFold(field.DataType, "object") {
			continue
		}
		field.ArrayItemType = normalizeDataType(field.ArrayItemType)
		if field.DataType != "array" {
			field.ArrayItemType = ""
		}
		if !dataTypeAllowsOptions(field.DataType) {
			field.Options = nil
		}
		if field.JSONKey == "" {
			continue
		}
		key := strings.ToLower(field.JSONKey)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, field)
	}
	return out
}

func dataTypeAllowsOptions(dataType string) bool {
	switch normalizeDataType(dataType) {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float", "float32", "float64", "decimal", "array":
		return true
	default:
		return false
	}
}

func optionsFromAny(value any) []Option {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	out := make([]Option, 0, len(items))
	for index, item := range items {
		if raw, ok := mapFromAny(item); ok {
			code := firstString(raw, "code", "optionCode", "option_code", "key", "const")
			if code == "" {
				code = fmt.Sprint(index + 1)
			}
			label := firstString(raw, "label", "title", "name")
			if label == "" {
				label = code
			}
			value := firstString(raw, "value", "x-camellia-value", "x_camellia_value")
			if value == "" {
				value = code
			}
			out = append(out, Option{Code: code, Label: label, Value: value, Disabled: boolFromAny(raw["disabled"]) || boolFromAny(raw["x-camellia-disabled"]) || boolFromAny(raw["x_camellia_disabled"])})
			continue
		}
		code := strings.TrimSpace(fmt.Sprint(item))
		if code == "" {
			code = fmt.Sprint(index + 1)
		}
		out = append(out, Option{Code: code, Label: code, Value: code})
	}
	return out
}

func optionsFromSchemaProperty(prop map[string]any) []Option {
	if prop == nil {
		return nil
	}
	if items, ok := mapFromAny(prop["items"]); ok {
		if options := optionsFromAny(items["oneOf"]); len(options) > 0 {
			return options
		}
		if options := optionsFromAny(items["enum"]); len(options) > 0 {
			return options
		}
	}
	if options := optionsFromAny(prop["oneOf"]); len(options) > 0 {
		return options
	}
	return optionsFromAny(prop["enum"])
}

func boolFromAny(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(strings.TrimSpace(v), "true") || strings.TrimSpace(v) == "1"
	default:
		return false
	}
}
