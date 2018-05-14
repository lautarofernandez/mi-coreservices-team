# JSON Filtering

This library that provides functionality to filter a `map[string]interface{}` (representing a JSON object) through a given set of rules. It returns a new map with only the attributes of the original that matches the rules given.

## Usage

```go
rules := []string{"id", "data[].amount"}

// Compile rules
ruleset := ParseRules(rules)

in := `{"id": 132456, "data": [{"amount": 12.12, "quantity": 5}, {"amount": 9.72, "quantity": 1}], "type": "payment"}`
var m map[string]interface{}
json.Unmarshal([]byte(in), &m)

// Using the compiled rules, filter a given object
filtered := ruleset.Filter(m)

f, _ := json.Marshal(filtered)

// `{"id": 132456, "data": [{"amount": 12.12}, {"amount":9.72}]}`
fmt.Println(string(f))
```

## Rule format

Rules are a sequence of dot separated (`.`) attributes corresponding to the nesting levels of a JSON object.

Rules can also attest that a given key contains an array as it's value, and that we want to filter out some element inside of that array (creating a new one as output).

This library does not support filtering arrays of arrays. The next JSON would be impossible to filter: `{"key":[[1,2,3],[4,5,6]]}`.

### Rule Examples

```
id
internal_id
type
site_id
user_id
version
schema_version
schema_original_version
date_created
last_modified
extra.actions
extra.payment
row.title
resources.main_resource.header
resources.other_resources[].header.id
resources.other_resources[].object.address.state
resources.other_resources[].not.exists
```