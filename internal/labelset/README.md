# labelset

The `labelset` package provides structured key-value label management for confsnap entities such as hosts and file paths.

## Overview

Unlike simple tags (single strings), labels are key-value pairs that allow richer metadata and precise selector-based queries.

## Usage

```go
m := labelset.New()

// Attach labels to a host
m.Set("web-01", "env",    "prod")
m.Set("web-01", "region", "us-east")
m.Set("db-01",  "env",    "prod")
m.Set("db-01",  "region", "eu-west")

// Retrieve a single label
val, ok := m.Get("web-01", "env") // "prod", true

// Get all labels for an entity
ls := m.Labels("web-01") // LabelSet{"env":"prod", "region":"us-east"}

// Remove a label
m.Remove("web-01", "region")

// Find all entities matching a selector
hosts := m.Match(labelset.LabelSet{"env": "prod"})
// ["db-01", "web-01"]
```

## Selector Matching

`Match` accepts a `LabelSet` selector and returns all entity names whose labels contain **all** of the specified key-value pairs. An empty selector matches every entity.

Results are returned in sorted order for deterministic output.
