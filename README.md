# confsnap

> Snapshot and diff configuration files across servers over SSH — useful for auditing config drift.

---

## Installation

```bash
go install github.com/youruser/confsnap@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/confsnap.git && cd confsnap && go build -o confsnap .
```

---

## Usage

Define your targets in a `confsnap.yaml` file:

```yaml
hosts:
  - address: web01.example.com
    user: admin
  - address: web02.example.com
    user: admin
files:
  - /etc/nginx/nginx.conf
  - /etc/ssh/sshd_config
```

Take a snapshot:

```bash
confsnap snapshot --config confsnap.yaml --out ./snapshots
```

Diff two snapshots to detect configuration drift:

```bash
confsnap diff ./snapshots/2024-01-10 ./snapshots/2024-01-17
```

Output highlights files that have changed across hosts, making it easy to spot unintended differences in your infrastructure.

---

## Flags

| Flag | Description |
|------|-------------|
| `--config` | Path to config file (default: `confsnap.yaml`) |
| `--out` | Output directory for snapshots |
| `--format` | Output format: `text`, `json` (default: `text`) |

---

## License

MIT © 2024 youruser