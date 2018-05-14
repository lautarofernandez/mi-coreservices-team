package jq

import (
	"strings"
)

const (
	// RuleFieldSeparator is the separator used in our rule system to denote a nesting
	// level representing each step in a JSON object traversal.
	RuleFieldSeparator = "."
)

// Ruleset is the result object of compiling a list of rules for filtering JSON.
type Ruleset []ruleKey

// ruleKey represents within a Ruleset each of the rule nodes, it's used to
// traverse the rule chain and filter out a given map.
type ruleKey struct {
	key        string
	arrayChild bool
	child      *Ruleset
}

// ParseRules receives a list of rules as defined by this package, parses them and
// generates an optimized Ruleset. This Ruleset can later be used to filter out
// a map[string]interface{} resulting a new one with only the matched fields.
func ParseRules(rules []string) Ruleset {
	ruleset := Ruleset{}

	for _, rule := range rules {
		parseRule(rule, &ruleset)
	}

	return ruleset
}

// parseRule is the inner recursive function used in ParseRules for parsing a rule
// and appending it to the Resultset being generated.
func parseRule(rule string, out *Ruleset) {
	parts := strings.SplitN(rule, RuleFieldSeparator, 2)
	if len(parts) == 1 {
		*out = append(*out, ruleKey{key: parts[0]})
		return
	}

	root := parts[0]
	subRule := parts[1]
	arrayRoot := false

	if strings.HasSuffix(root, "[]") {
		arrayRoot = true
		root = strings.TrimSuffix(root, "[]")
	}

	for _, r := range *out {
		// The given key already exists in the out RuleSet, use that RuleSet instead
		if r.key == root {
			if r.child == nil {
				r.child = new(Ruleset)
			}

			parseRule(subRule, r.child)
			return
		}
	}

	// If the key does not exists in the ruleset, then we create it
	key := ruleKey{key: root, child: new(Ruleset), arrayChild: arrayRoot}
	*out = append(*out, key)
	parseRule(subRule, key.child)
}

// Filter receives a map, and using the precompiled Ruleset, iterates through it forming
// a new map with only the matched fields inside.
func (r *Ruleset) Filter(m map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}

	for _, i := range *r {
		v, resolved := resolveRule(m, i)
		if resolved {
			out[i.key] = v
		}
	}

	return out
}

func resolveRule(m map[string]interface{}, rule ruleKey) (interface{}, bool) {
	// Check to see if the given rule is contained within m. If the key is not contained
	// then we no longer need to go deeper into the rule web, because we can already
	// assert that the given rule will not be resolved.
	v, exists := m[rule.key]
	if !exists {
		return nil, false
	}

	// If the rule has no child, then we v directly, as it´s the result for the given key.
	if rule.child == nil {
		return v, true
	}

	// If the key exists (we have a v) and the given rule is not a leaf (has no children) then we
	// need to continue downwards into the rule system to find the rule's value. Given how
	// JSON works, there are 2 possible types that we can find on a given tree, the node
	// can be a map (JSON object), or a slice (JSON array). We must check to see
	// which path to take in order to process it accordingly.
	switch node := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{})

		// The rule we are executing has 1...N child rules, we need to execute each child rule
		// to `node` (subset of m that we want to filter) and accumulate it´s results.
		childRules := *rule.child
		for _, childRule := range childRules {
			v, resolved := resolveRule(node, childRule)
			if resolved {
				out[childRule.key] = v
			}
		}

		if len(out) == 0 {
			return nil, false
		}

		return out, true
	case []interface{}:
		// If the node is an array, we only go forward if the rule used to filtered it
		// is telling us that that node should be an array. If not we bail out.
		if !rule.arrayChild {
			return nil, false
		}

		// Given that the node is an array, the result of filtering it's inner objects is going
		// to be an array. We are going to instantiate an array which we'll uses to append the
		// result of filtering each of the inner objects with the current rule child's.
		arr := []interface{}{}
		for _, subNode := range node {
			out := make(map[string]interface{})

			// A subNode in JSON can be either an Object or an Array. In this filtering lib we're only
			// going to support filtering objects inside an array (not arrays inside arrays).
			if m, ok := subNode.(map[string]interface{}); ok {
				// The rule we are executing has 1...N child rules, we need to execute each child rule
				// to `m` (object inside node array we want to filter) and accumulate it´s results.
				childRules := *rule.child
				for _, childRule := range childRules {
					v, resolved := resolveRule(m, childRule)
					if resolved {
						out[childRule.key] = v
					}
				}
			}

			// After filtering the object inside the node array, we add it to the resulting arr
			// only if the new filtered object has elements in it. The filtering might result
			// in an empty object, and we don't want to have that as a result.
			if len(out) > 0 {
				arr = append(arr, out)
			}
		}

		// We don´t want to say we resolved a rule if the result is empty, so check for that.
		if len(arr) == 0 {
			return nil, false
		}

		return arr, true
	}

	return nil, false
}
