# cache

The `cache` package provides a lightweight, file-backed checksum cache for
**confsnap**. It prevents redundant SSH reads when a remote file has not
changed since the last snapshot run.

## How it works

1. Before fetching a file over SSH, the runner checks the cache using the
   `(host, path)` pair as a key.
2. If a **non-expired** entry exists and its checksum matches the last known
   value, the file fetch can be skipped.
3. After a successful fetch the runner calls `Set` to update the cache entry
   with the new checksum and a configurable TTL (in seconds).
4. Cache state is persisted to `<cache-dir>/cache.json` so it survives process
   restarts.

## TTL behaviour

| TTL value | Behaviour |
|-----------|------------------------------------------|
| `0`       | Entry never expires (use with caution)   |
| `> 0`     | Entry expires after the given seconds    |

## Usage

```go
c, err := cache.New("/var/lib/confsnap/cache")
if err != nil {
    log.Fatal(err)
}

if entry := c.Get("web01", "/etc/nginx/nginx.conf"); entry != nil {
    // checksum unchanged — skip SSH fetch
    fmt.Println("cache hit:", entry.Checksum)
} else {
    // fetch from remote, then store result
    _ = c.Set("web01", "/etc/nginx/nginx.conf", newChecksum, 3600)
}
```
