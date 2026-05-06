# audit

The `audit` package provides a simple append-only audit logger for confsnap.
Each audit entry records when a config file was snapshotted, diffed, or
triggered an alert on a specific host.

## Usage

```go
logger, err := audit.New("/var/lib/confsnap/audit")
if err != nil {
    log.Fatal(err)
}

// Log a snapshot event
logger.Log("web01", "/etc/nginx/nginx.conf", "snapshot", "")

// Log a diff event with details
logger.Log("web01", "/etc/nginx/nginx.conf", "diff", "12 lines changed")

// Read today's entries
entries, err := logger.ReadDay(time.Now())
```

## File Format

Audit logs are stored as JSON-lines files, one per day:

```
audit-2024-01-15.jsonl
```

Each line is a JSON object:

```json
{"timestamp":"2024-01-15T10:30:00Z","host":"web01","file":"/etc/nginx/nginx.conf","event":"snapshot","details":""}
```

## Events

| Event      | Description                              |
|------------|------------------------------------------|
| `snapshot` | A config file was successfully captured  |
| `diff`     | A difference was detected between runs   |
| `alert`    | An alert threshold was exceeded          |
| `error`    | An error occurred during collection      |
