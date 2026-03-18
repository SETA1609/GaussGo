package mods

import (
	"encoding/json"
	"fmt"
	"os"

	"gaussgo/internal/contracts"
)

type manifestJSON struct {
	ID           string                    `json:"id"`
	Name         string                    `json:"name"`
	Version      string                    `json:"version"`
	Entry        string                    `json:"entry"`
	Localization contracts.ModLocalization `json:"localization"`
	Dependencies []contracts.ModDependency `json:"dependencies"`
	Provides     []string                  `json:"provides"`
}

func ParseManifestFile(path string) (contracts.ModManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return contracts.ModManifest{}, fmt.Errorf("read manifest: %w", err)
	}

	var raw manifestJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return contracts.ModManifest{}, fmt.Errorf("parse manifest: %w", err)
	}

	manifest := contracts.ModManifest{
		ID:           raw.ID,
		Name:         raw.Name,
		Version:      raw.Version,
		Entry:        raw.Entry,
		Localization: raw.Localization,
		Dependencies: raw.Dependencies,
		Provides:     raw.Provides,
	}

	if err := ValidateManifest(manifest); err != nil {
		return contracts.ModManifest{}, err
	}

	return manifest, nil
}
