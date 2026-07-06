# cmldataresolve

`cmldataresolve` resolves values from a Camellia form contract plus form data.

The first target is `formtool.v1`: fields describe `jsonKey` and `dataType`,
while business data can store primitive values as strings. A caller can then
ask for a typed value without hand-written conversion code.

```go
resolver, err := cmldataresolve.New(contractJSON, dataJSON)
if err != nil {
    return err
}

days, err := resolver.Int32("leaveDays")
if err != nil {
    return err
}
```

## Supported Lookups

A field can be found by:

- `jsonKey`
- field code
- display name

If a `jsonKey` is nested, such as `employee.age`, the resolver first checks for
an exact key and then falls back to dotted object traversal.

## Supported Types

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float`, `float32`, `float64`
- `decimal`
- `bool`
- `date`, `datetime`
- `bytes`
- `array`

`object` is intentionally not a standard formtool field type. Contracts are JSON
objects, but individual data definitions should use scalar values, `bytes`, or
`array`.

`decimal` is returned as a validated string by `Any` and can be read as
`*big.Rat` through `TypedValue.DecimalRat`.

Static options are not tied to the `array` data type. Whether options are used
depends on the field `widgetType`.

For scalar `select` or `radio` fields, form data stores one selected
`optionCode`, while typed reads resolve the matched option `value`:

```json
{
  "jsonKey": "department",
  "dataType": "string",
  "widgetType": "select",
  "options": [
    { "code": "D003", "label": "R&D", "value": "003" },
    { "code": "D004", "label": "Manufacturing", "value": "004" }
  ]
}
```

Given data `{ "department": "D003" }`, use:

```go
codes, _ := resolver.SelectionCodes("department")   // []string{"D003"}
labels, _ := resolver.SelectionLabels("department") // []string{"R&D"}
value, _ := resolver.String("department")           // "003"
```

For `array` option widgets such as `multiSelect` or `checkbox`, form data
stores selected option-code arrays. The contract field should include
`arrayItemType`, `widgetType`, and an `options` snapshot:

```json
{
  "jsonKey": "levels",
  "dataType": "array",
  "arrayItemType": "int32",
  "widgetType": "multiSelect",
  "options": [
    { "code": "L1", "label": "Level 1", "value": "1" },
    { "code": "L2", "label": "Level 2", "value": "2" }
  ]
}
```

Given data `{ "levels": ["L2", "L1"] }`, use:

```go
codes, _ := resolver.SelectionCodes("levels")   // []string{"L2", "L1"}
labels, _ := resolver.SelectionLabels("levels") // []string{"Level 2", "Level 1"}
values, _ := resolver.Int32Slice("levels")      // []int32{2, 1}
```

For free array widgets such as `input` or `textarea`, form data stores the raw
typed array and values may repeat. Given an `arrayItemType` of `int32` and data
`{ "samples": [1, 2, 2, 3] }`, `resolver.Int32Slice("samples")` returns
`[]int32{1, 2, 2, 3}`.
## Contract Sources

The resolver reads fields from:

- `fields[]`
- `dataPoints[]`
- `schema.properties`
- nested `contract`, `templateContract`, or `formContract`

An instance contract may include `data`, `values`, `formData`, or string fields
such as `dataJson`. Use `FromInstanceContract` when the contract and data are in
the same JSON object.

## Module Path

The current module path is:

```text
github.com/wangWenCn/cmldataresolve
```

If this package is published under another GitHub organization, change the
module path before creating the first public tag.

