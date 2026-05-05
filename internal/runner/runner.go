package runner

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/user/confsnap/internal/config"
	"github.com/user/confsnap/internal/diff"
	"github.com/user/confsnap/internal/snapshot"
	"github.com/user/confsnap/internal/ssh"
)

// Result holds the outcome of a single file capture on a host.
type Result struct {
	Host string
	Path string
	Diff diff.Result
	Err  error
}

// Runner orchestrates snapshot collection and diffing across hosts.
type Runner struct {
	cfg   *config.Config
	store *snapshot.Store
}

// New creates a Runner with the given config and store directory.
func New(cfg *config.Config, storeDir string) (*Runner, error) {
	store, err := snapshot.NewStore(storeDir)
	if err != nil {
		return nil, fmt.Errorf("runner: init store: %w", err)
	}
	return &Runner{cfg: cfg, store: store}, nil
}

// Run connects to each host, fetches each configured file, snapshots it,
// and returns a diff result per (host, file) pair.
func (r *Runner) Run() []Result {
	var results []Result

	for _, host := range r.cfg.Hosts {
		client, err := ssh.NewClient(ssh.Config{
			Host:    host.Address,
			User:    host.User,
			KeyFile: host.KeyFile,
			Timeout: 15 * time.Second,
		})
		if err != nil {
			for _, fp := range r.cfg.Files {
				results = append(results, Result{Host: host.Address, Path: fp, Err: err})
			}
			continue
		}

		for _, fp := range r.cfg.Files {
			res := r.processFile(client, host.Address, fp)
			results = append(results, res)
		}
	}

	return results
}

func (r *Runner) processFile(client *ssh.Client, host, path string) Result {
	contents, err := client.ReadFile(path)
	if err != nil {
		return Result{Host: host, Path: path, Err: fmt.Errorf("read %s: %w", path, err)}
	}

	current := snapshot.New(host, path, contents)

	key := filepath.Join(host, path)
	previous, _ := r.store.Load(key)

	if err := r.store.Save(key, current); err != nil {
		return Result{Host: host, Path: path, Err: fmt.Errorf("save snapshot: %w", err)}
	}

	dr := diff.Compare(previous, current)
	return Result{Host: host, Path: path, Diff: dr}
}
