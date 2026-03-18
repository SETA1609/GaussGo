package mods

import (
	"fmt"
	"sort"

	"gaussgo/internal/contracts"
)

func ResolveLoadOrder(manifests []contracts.ModManifest) ([]string, error) {
	if err := EnforceCorePresence(manifests); err != nil {
		return nil, err
	}

	manifestByID := make(map[string]contracts.ModManifest, len(manifests))
	for _, m := range manifests {
		manifestByID[m.ID] = m
	}

	inDegree := make(map[string]int, len(manifests))
	adj := make(map[string][]string, len(manifests))
	for _, m := range manifests {
		if _, ok := inDegree[m.ID]; !ok {
			inDegree[m.ID] = 0
		}
		for _, dep := range m.Dependencies {
			if _, exists := manifestByID[dep.ID]; !exists {
				return nil, fmt.Errorf("mod %s has missing dependency %s", m.ID, dep.ID)
			}
			adj[dep.ID] = append(adj[dep.ID], m.ID)
			inDegree[m.ID]++
		}
	}

	queue := make([]string, 0, len(manifests))
	if inDegree["core"] != 0 {
		return nil, fmt.Errorf("core must not depend on other mods")
	}
	queue = append(queue, "core")
	for id, degree := range inDegree {
		if id == "core" {
			continue
		}
		if degree == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue[1:])

	order := make([]string, 0, len(manifests))
	seen := make(map[string]bool, len(manifests))
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if seen[current] {
			continue
		}
		seen[current] = true
		order = append(order, current)

		for _, next := range adj[current] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
		if len(queue) > 1 {
			sort.Strings(queue)
		}
	}

	if len(order) != len(manifests) {
		return nil, fmt.Errorf("dependency cycle detected")
	}
	if order[0] != "core" {
		return nil, fmt.Errorf("core must be first in load order")
	}

	return order, nil
}
