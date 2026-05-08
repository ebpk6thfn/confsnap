# rollback

The `rollback` package provides a simple mechanism for saving and restoring
configuration file snapshots as named rollback points.

## Overview

When confsnap detects drift, operators can save the last-known-good state of a
file as a rollback entry. That entry can later be retrieved to restore the
original content via SSH or manual intervention.

## Usage

```go
m, err := rollback.New("/var/lib/confsnap/rollback")
if err != nil {
    log.Fatal(err)
}

// Save a snapshot as a rollback point.
id, err := m.Save(snap)
if err != nil {
    log.Fatal(err)
}
fmt.Println("saved rollback:", id)

// Later, load it back.
entry, err := m.Load(id)
if err != nil {
    log.Fatal(err)
}
if entry != nil {
    fmt.Printf("restoring %s on %s\n", entry.Path, entry.Host)
    // push entry.Content back over SSH ...
}

// Clean up when no longer needed.
m.Delete(id)
```

## Entry ID format

Each rollback entry is identified by a deterministic string:

```
<host>__<path>__<timestamp>
```

For example:

```
web01___etc_nginx_nginx.conf__20240615T123000Z
```

Slashes and colons in the host and path components are replaced with
underscores so the ID is safe to use as a filename.

## Storage

Entries are stored as plain files under the configured directory, one file per
rollback point. The directory is created automatically if it does not exist.
