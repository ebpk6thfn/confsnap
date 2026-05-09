# redact

The `redact` package masks sensitive values in configuration file content before
they are stored in snapshots, exported, or displayed in reports.

## Overview

When snapshotting files such as `/etc/app/config.ini` or `~/.ssh/config`, some
lines may contain secrets (passwords, API keys, tokens). The redactor scans
each line for `key=value` or `key: value` pairs whose key matches a configured
regular-expression pattern and replaces the value with a mask string.

## Usage

```go
import "github.com/yourorg/confsnap/internal/redact"

rules := []redact.Rule{
    {Pattern: `(?i)password`},
    {Pattern: `(?i)api[_-]?key`, Mask: "<api-key-hidden>"},
    {Pattern: `(?i)token`},
}

r, err := redact.New(rules)
if err != nil {
    log.Fatal(err)
}

redacted := r.ApplyString(fileContent)
```

## Rule fields

| Field     | Description                                          |
|-----------|------------------------------------------------------|
| `Pattern` | Go regular expression matched against the config key |
| `Mask`    | Replacement string (defaults to `***REDACTED***`)    |

## Supported formats

- `key=value` (INI, shell env files)
- `key: value` (YAML-style flat files)

Lines without a recognised separator (e.g. comments) are passed through
unchanged.
