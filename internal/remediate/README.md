# remediate

The `remediate` package evaluates diff results against a set of rules and
applies automatic remediation actions when configuration drift is detected.

## Actions

| Action     | Description                                              |
|------------|----------------------------------------------------------|
| `rollback` | Restores the file to the most recent rollback snapshot.  |
| `notify`   | Records a drift event for downstream alerting.           |

## Usage

```go
store, _ := rollback.New("/var/lib/confsnap/rollback")

rules := []remediate.Rule{
    {Pattern: "/etc/nginx/*", Action: remediate.ActionRollback},
    {Pattern: "/etc/*",       Action: remediate.ActionNotify},
}

rm := remediate.New(rules, store)
results := rm.Evaluate(diffResults)

for _, r := range results {
    fmt.Printf("[%s] %s:%s — %s\n", r.Action, r.Host, r.Path, r.Message)
}
```

## Rule matching

Patterns follow standard shell glob syntax via `path.Match`.
Rules are evaluated in order; the first matching rule wins.

## Notes

- Only diff results where `Changed == true` are evaluated.
- Rollback requires a prior entry saved via `internal/rollback`.
- If no rollback entry exists, `Applied` is set to `false` and the
  reason is recorded in `Message`.
