// Package labelset provides key-value label management for hosts and files,
// enabling flexible grouping and querying beyond simple tags.
package labelset

import (
	"fmt"
	"sort"
	"strings"
)

// LabelSet holds a map of key-value labels.
type LabelSet map[string]string

// Manager manages labels attached to named entities (hosts or file paths).
type Manager struct {
	entities map[string]LabelSet
}

// New returns a new Manager.
func New() *Manager {
	return &Manager{entities: make(map[string]LabelSet)}
}

// Set attaches a key-value label to an entity. Overwrites existing value for key.
func (m *Manager) Set(entity, key, value string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("label key must not be empty")
	}
	if _, ok := m.entities[entity]; !ok {
		m.entities[entity] = make(LabelSet)
	}
	m.entities[entity][key] = value
	return nil
}

// Get returns the value for a label key on an entity, and whether it was found.
func (m *Manager) Get(entity, key string) (string, bool) {
	ls, ok := m.entities[entity]
	if !ok {
		return "", false
	}
	v, ok := ls[key]
	return v, ok
}

// Labels returns a copy of all labels for an entity.
func (m *Manager) Labels(entity string) LabelSet {
	ls, ok := m.entities[entity]
	if !ok {
		return LabelSet{}
	}
	copy := make(LabelSet, len(ls))
	for k, v := range ls {
		copy[k] = v
	}
	return copy
}

// Remove deletes a label key from an entity.
func (m *Manager) Remove(entity, key string) {
	if ls, ok := m.entities[entity]; ok {
		delete(ls, key)
	}
}

// Match returns all entity names whose labels satisfy all provided key-value pairs.
func (m *Manager) Match(selector LabelSet) []string {
	var results []string
	for entity, ls := range m.entities {
		if matchesAll(ls, selector) {
			results = append(results, entity)
		}
	}
	sort.Strings(results)
	return results
}

func matchesAll(ls, selector LabelSet) bool {
	for k, v := range selector {
		if ls[k] != v {
			return false
		}
	}
	return true
}
