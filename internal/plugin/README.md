# plugin

The `plugin` package provides a lightweight, hook-based event system for **confsnap**.
It lets internal subsystems (and future external integrations) react to snapshot
lifecycle events without tight coupling.

## Event types

| Constant | When fired |
|---|---|
| `EventSnapshotTaken` | A single file snapshot has been captured from a host. |
| `EventDiffDetected` | A diff was found between the current and previous snapshot. |
| `EventRunComplete` | The entire runner sweep has finished. |

## Usage

```go
reg := plugin.New()

// Register a handler
err := reg.Register(plugin.EventDiffDetected, "slack-notifier", func(e plugin.Event) error {
    return slack.Notify(fmt.Sprintf("Drift on %s:%s", e.Host, e.Path))
})

// Emit from the runner after a diff is found
reg.Emit(plugin.Event{
    Type: plugin.EventDiffDetected,
    Host: result.Host,
    Path: result.Path,
    Meta: map[string]string{"lines_changed": strconv.Itoa(n)},
})
```

## Notes

- Handler registration is safe for concurrent use.
- All registered handlers for an event are called in order; errors are collected
  and returned as a single combined error — a failing handler does **not** stop
  subsequent handlers.
- Duplicate handler names for the same event type are rejected with an error.
