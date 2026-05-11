# template

The `template` package provides Go `text/template`-based rendering for confsnap diff reports.

## Usage

```go
import "github.com/yourorg/confsnap/internal/template"

src := `confsnap report — {{ .GeneratedAt | formatTime }}
Hosts : {{ .TotalHosts }}
Files : {{ .TotalFiles }}
Changed : {{ .Changed }}
Unchanged: {{ .Unchanged }}
{{ range $i, $r := .Results }}{{ inc $i }}. {{ $r.Host }} {{ $r.Path }}{{ if $r.Changed }} [CHANGED]{{ end }}
{{ end }}`

renderer, err := template.New(src)
if err != nil {
    log.Fatal(err)
}

data := template.BuildData(results)
if err := renderer.Render(os.Stdout, data); err != nil {
    log.Fatal(err)
}
```

## Template Variables

| Variable | Type | Description |
|---|---|---|
| `.GeneratedAt` | `time.Time` | UTC time the data was built |
| `.TotalHosts` | `int` | Number of distinct hosts in results |
| `.TotalFiles` | `int` | Total file results |
| `.Changed` | `int` | Count of changed files |
| `.Unchanged` | `int` | Count of unchanged files |
| `.Results` | `[]diff.Result` | Full list of diff results |

## Built-in Functions

| Function | Signature | Description |
|---|---|---|
| `formatTime` | `time.Time → string` | Formats a time as RFC3339 |
| `inc` | `int → int` | Increments an integer (useful for 1-based indexes) |
