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
- `object`
- `array`

`decimal` is returned as a validated string by `Any` and can be read as
`*big.Rat` through `TypedValue.DecimalRat`.

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
github.com/camellia-go/cmldataresolve
```

If this package is published under another GitHub organization, change the
module path before creating the first public tag.
