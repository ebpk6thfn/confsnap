# tag

The `tag` package provides a lightweight registry for associating string tags with hosts and file paths in confsnap.

## Overview

Tags enable grouping and filtering during snapshot runs, diff reports, and alerting. For example, you can tag all production hosts with `production` and query only those hosts when generating audit reports.

## Usage

```go
reg := tag.New()

// Tag hosts
reg.TagHost("web-01", "production", "web")
reg.TagHost("db-01", "production", "database")

// Tag files
reg.TagFile("/etc/nginx/nginx.conf", "nginx", "web")
reg.TagFile("/etc/ssh/sshd_config", "security")

// Query
hosts := reg.HostsWithTag("production") // ["db-01", "web-01"]
files := reg.FilesWithTag("web")         // ["/etc/nginx/nginx.conf"]
```

## Validation

Tags must be non-empty and must not contain whitespace:

```go
if err := tag.Validate([]string{"production", "web"}); err != nil {
    log.Fatal(err)
}
```

## Notes

- Tags are deduplicated automatically when added via `TagHost` or `TagFile`.
- All query results are returned in sorted order for deterministic output.
- The registry is not safe for concurrent writes; synchronise externally if needed.
