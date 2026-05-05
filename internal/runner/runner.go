package runner

import (
	"fmt"
	"path/filepath"

	"github.com/user/confsnap/internal/config"
	"github.com/user/confsnap/internal/diff"
	"github.com/user/confsnap/internal/snapshot"
	"github.com/user/confsnap/internal/ssh"
)

// Runner orchestrates snapshot collection and diffing across hosts.
type Runner struct {
	cfg      *config.Config
	store    *snapshot.Store
	storeDir string
}

// New creates a Runner for the given config and store directory.
func New(cfg *config.Config, storeDir string) (*Runner, error) {
	store, err := snapshot.NewStore(storeDir)
	if err != nil {
		return nil, fmt.Errorf("runner: init store: %w", err)
	}
	return &Runner{cfg: cfg, store: store, storeDir: storeDir}, nil
}

// Run connects to each host, fetches files, stores snapshots, and returns diff results.
func (r *Runner) Run() []diff.Result {
	var results []diff.Result

	for _, host := range r.cfg.Hosts {
		clientCfg := ssh.Config{
			Host:    host.Address,
			User:    host.User,
			KeyFile: host.KeyFile,
		}
		client, err := ssh.NewClient(clientCfg)
		if err != nil {
			for _, f := range r.cfg.Files {
				results = append(results, diff.Result{
					Host:  host.Address,
					Path:  f,
					Error: fmt.Errorf("ssh connect: %w", err),
				})
			}
			continue
		}

		for _, filePath := range r.cfg.Files {
			content, fetchErr := client.ReadFile(filePath)
			if fetchErr != nil {
				results = append(results, diff.Result{
					Host:  host.Address,
					Path:  filePath,
					Error: fmt.Errorf("read file: %w", fetchErr),
				})
				continue
			}

			newSnap := snapshot.New(host.Address, filePath, content)
			key := filepath.Join(host.Address, filePath)
			prev, _ := r.store.Load(key)

			result := diff.Compare(prev, newSnap)
			results = append(results, result)

			_ = r.store.Save(key, newSnap)
		}
	}

	return results
}
