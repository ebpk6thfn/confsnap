# retention

The `retention` package manages pruning of old snapshot files stored on disk.

## Overview

Over time, confsnap accumulates snapshot JSON files in the store directory.
The retention manager applies a configurable policy to remove files that are
too old or exceed a per-host-path count limit.

## Policy

| Field      | Description                                                   |
|------------|---------------------------------------------------------------|
| `MaxAge`   | Remove snapshots older than this duration. `0` = unlimited.  |
| `MaxCount` | Keep at most this many snapshots per host+path. `0` = unlimited. |

## Configuration (YAML)

```yaml
retention:
  max_age_days: 30
  max_count: 10
```

## Usage

```go
cfg := retention.DefaultConfig()
policy, err := cfg.ToPolicy()
if err != nil { /* handle */ }

m, err := retention.New(storeDir, policy)
if err != nil { /* handle */ }

removed, err := m.Prune(time.Now())
fmt.Printf("pruned %d snapshot(s)\n", len(removed))
```

## Grouping

Files are grouped by the `<host>_<escaped-path>` prefix of their filename
(the timestamp suffix is excluded). Count and age limits are applied
independently within each group so that fast-changing files on one host do
not cause snapshots for other hosts or paths to be removed prematurely.
