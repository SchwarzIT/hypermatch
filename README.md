[![SIT](https://img.shields.io/badge/SIT-awesome-blueviolet.svg)](https://jobs.schwarz)
[![CI](https://github.com/SchwarzIT/hypermatch/actions/workflows/go-test.yml/badge.svg)](https://github.com/SchwarzIT/hypermatch/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/hypermatch)](https://goreportcard.com/report/github.com/SchwarzIT/hypermatch)
[![Go Reference](https://pkg.go.dev/badge/github.com/schwarzit/hypermatch.svg)](https://pkg.go.dev/github.com/schwarzit/hypermatch)
![License](https://img.shields.io/github/license/SchwarzIT/hypermatch)
![GitHub last commit](https://img.shields.io/github/last-commit/SchwarzIT/hypermatch)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/SchwarzIT/hypermatch)

![hypermatch logo](./logo/logo-small.png)

# Introduction
Hypermatch is a Go library that allows blazing fast matching of a large number of rules to events.
- It matches many thousands of events per second to huge numbers of rules in-memory with low latency
- Rules can be serialized to easy readable JSON objects
- Features an expressive rule syntax with support for: equals, prefix, suffix, wildcard, anything-but, all-of, any-of 

An event is a list of fields, which may be given as name/value pairs. A rule associates event field names with pattern to match the event values.

... [it's fast!](_benchmark/benchmark.md)

![example](./example.png)

# Quick Start

```go
import (
    hypermatch "github.com/SchwarzIT/hypermatch"
)

func main() {
    //Initialize hypermatch
    hm := hypermatch.NewHyperMatch()
    
    //Add a rule
    if err := hm.AddRule("markus_rule", hypermatch.ConditionSet{
        hypermatch.Condition{Path: "firstname", Pattern: hypermatch.Pattern{Type: hypermatch.PatternEquals, Value: "markus"}},
        hypermatch.Condition{Path: "lastname", Pattern: hypermatch.Pattern{Type: hypermatch.PatternEquals, Value: "troßbach"}},
        }); err != nil {
            panic(err)
    }
    
    //Test with match
    matchedRules := hm.Match([]hypermatch.Property{
        {Path: "firstname", Values: []string{"markus"}},
        {Path: "lastname", Values: []string{"troßbach"}},
    })
    log.Printf("Following rules matches: %v", matchedRules)
    
    //Test without match
    matchedRules = hm.Match([]hypermatch.Property{
        {Path: "firstname", Values: []string{"john"}},
        {Path: "lastname", Values: []string{"doe"}},
    })
    log.Printf("Following rules matches: %v", matchedRules)
}
```

# Documentation
## Example Event

An event is a list of field, expressible as a JSON object. Here is an example:

```javascript
{
        "name": "Too many parallel requests on system xy",
        "severity": "critical",
        "status": "firing",
        "message": "Lorem ipsum dolor sit amet, consetetur sadipscing elitr.",
        "team": "awesome-team",
        "application": "webshop",
        "component": "backend-service",
        "tags": [
            "shop",
            "backend"
        ]   
}
```

**This sample event is used as a reference throughout the documentation.**

## Matching Basics

Rules are a set of Conditions and are expressed by the **ConditionSet** type.
- **The matching is always done case-insensitive for all values!**
- Currently only string or string-arrays are supported as values

A condition exists of
- **Path**: Path in the event to match to
- **Pattern**: A patter to match the value at **Path**

The following example matches the example event from above:
```go
ConditionSet{
    {
        Path: "status",
        Pattern: Pattern{Type: PatternEquals, Value: "firing"},
    },
    {
        Path: "name",
        Pattern: Pattern{Type: PatternAnythingBut, Sub: []Pattern{
                {Type: PatternWildcard, Value: "TEST*"},
            },
        },
    },
    {
        Path: "severity",
        Pattern: Pattern{ Type: PatternAnyOf,
            Sub: []Pattern{
                {Type: PatternEquals, Value: "critical"},
                {Type: PatternEquals, Value: "warning"},
            },
        },
    },
    {
        Path: "tags",
        Pattern: Pattern{ Type: PatternAllOf,
            Sub: []Pattern{
                {Type: PatternEquals, Value: "shop"},
                {Type: PatternEquals, Value: "backend"},
            },
        },
    },
}
```

The rules and conditions are also expressible as JSON objects. The following JSON is the equivalent of the above Go notation for a **ConditionSet**:

```javascript
{
    "status": {
        "equals": "firing"
    },
    "name": {
        "anythingBut": [
            {"wildcard": "TEST*"}
        ]
    },
    "severity": {
        "anyOf": [
            {"equals": "critical"},
            {"equals": "warning"}
        ]
    },
    "tags": {
        "allOf": [
            {"equals": "shop"},
            {"equals": "backend"}
        ]
    }
}
```

The rest of the documentation uses JSON notation for easier readability.

## Matching syntax
### "equals" matching
**equals** checks the value of an attribute of the event for equality. The checks are always case-insensitive.

```javascript
{
    "status": {
        "equals": "firing"
    }
}
```

- If the value in the event is a **string**: the value is checked for equality to "firing".
- If the value in the event is a **string array**: it checks whether the array contains the value "firing".

### "prefix" matching
**prefix** checks the value of an attribute of the event for a given prefix. The checks are always case-insensitive.

```javascript
{
    "status": {
        "prefix": "fir"
    }
}
```

- If the value in the event is a **string**: the value is checked if it begins with "fir".
- If the value in the event is a **string array**: it checks whether the array contains a value which begins with "fir".

### "suffix" matching
**suffix** checks the value of an attribute of the event for a given suffix. The checks are always case-insensitive.

```javascript
{
    "status": {
        "suffix": "ing"
    }
}
```

- If the value in the event is a **string**: the value is checked if it ends with "ing".
- If the value in the event is a **string array**: it checks whether the array contains a value which end with "ing".

### "wildcard" matching
**wildcard** allows to check the value of an attribute of the event using wildcards. The checks are always case-insensitive.
- A wildcard is expressed using the character '*' and matches any number of characters (including zero) of the value.
- Wildcards may not be used directly one after the other.
```javascript
{
    "name": {
        "wildcard": "*parallel requests*"
    }
}
```
- If the value in the event is a **string**: the value is checked if it matches with "\*parallel requests\*".
- If the value in the event is a **string array**: it checks whether the array contains a value which matches with "\*parallel requests\*".

### "anythingBut" matching
**anythingBut** does correspond to a boolean "not". It negates the included matching condition and triggers in the opposite case.

```javascript
{
    "status": {
        "anythingBut": [
            {"equals": "firing"}
        ]
    }
}
```
- If the value in the event is a **string**: the value is checked if it is anything else but "firing".
- If the value in the event is a **string array**: it checks whether the array does **not** contain a value which is equal to "firing".


### "anyOf" matching
**anyOf** does correspond to a boolean "inclusive-or". It checks multiple conditions and matches if **any** of the conditions are true.

```javascript
{
    "status": {
        "anyOf": [
            {"equals": "firing"},
            {"equals": "resolved"}
        ]
    }
}
```
- If the value in the event is a **string**: the value is checked if it is equal to "firing" **or** "resolved"
- If the value in the event is a **string array**: it checks whether the array does contain a value which is equal to "firing" or "resolved" or both.


### "allOf" matching
**allOf** does correspond to a boolean "and". It checks multiple conditions and matches if **all** the conditions are true.

```javascript
{
    "tags": {
        "allOf": [
            {"equals": "shop"},
            {"equals": "backend"}
        ]
    }
}
```
- If the value in the event is a **string**: the value is checked if it is equal to "shop" **and** "backend" (which will never be the case ;-) )
- If the value in the event is a **string array**: it checks whether the array does contain both "shop" **and** "backend".

# Performance

**hypermatch** is designed to be blazing fast with very large numbers of rules.
Nevertheless, there are a few things to consider to get maximum performance:
- Shorten the number of fields inside the rules, the fewer conditions, the shorter is the path to find them out.
- Try to make the **paths** as diverse as possible in events and rules. The more heterogeneous fields, the higher the performance.
- Reduce the number of **anyOf** conditions wherever possible
