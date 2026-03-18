package mods

import (
	"gaussgo/internal/contracts"
)

func BuildStatuses(manifests []contracts.ModManifest, enabled map[string]bool) ([]contracts.ModStatus, error) {
	order, err := ResolveLoadOrder(manifests)
	if err != nil {
		return nil, err
	}

	byID := make(map[string]contracts.ModManifest, len(manifests))
	for _, m := range manifests {
		byID[m.ID] = m
	}

	statuses := make([]contracts.ModStatus, 0, len(order))
	for _, id := range order {
		m := byID[id]
		available := true
		reason := ""
		for _, dep := range m.Dependencies {
			depEnabled, ok := enabled[dep.ID]
			if !ok || !depEnabled {
				available = false
				reason = "missing enabled dependency: " + dep.ID
				break
			}
		}

		isEnabled := enabled[id]
		if id == "core" {
			isEnabled = true
		}

		if !available {
			isEnabled = false
		}

		statuses = append(statuses, contracts.ModStatus{
			ModID:          id,
			Enabled:        isEnabled,
			Available:      available,
			ReasonDisabled: reason,
		})
	}

	return statuses, nil
}
