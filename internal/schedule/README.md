# schedule

The `schedule` package provides simple interval-based scheduling for automated
confsnap snapshot runs.

## Intervals

| Constant           | Value      | Description                        |
|--------------------|------------|------------------------------------|
| `IntervalHourly`   | `hourly`   | Run every hour                     |
| `IntervalDaily`    | `daily`    | Run every 24 hours                 |
| `IntervalWeekly`   | `weekly`   | Run every 7 days                   |
| `IntervalManual`   | `manual`   | Never auto-run; trigger manually   |

## Usage

```go
s, err := schedule.New(schedule.IntervalDaily, time.Now())
if err != nil {
    log.Fatal(err)
}

if s.Due(time.Now()) {
    // run snapshots
    if err := s.Advance(time.Now()); err != nil {
        log.Fatal(err)
    }
}
```

## Integration with runner

The schedule can be persisted alongside the snapshot store and checked at
startup to decide whether a new snapshot cycle should be triggered.
