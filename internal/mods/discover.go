package mods

import (
	"fmt"
	"path/filepath"
	"sort"

	"gaussgo/internal/contracts"
)

func Discover(modsDir string) ([]contracts.ModManifest, error) {
	paths, err := filepath.Glob(filepath.Join(modsDir, "*", "manifest.json"))
	if err != nil {
		return nil, fmt.Errorf("discover manifests: %w", err)
	}

	sort.Strings(paths)

	manifests := make([]contracts.ModManifest, 0, len(paths))
	for _, manifestPath := range paths {
		m, err := ParseManifestFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", manifestPath, err)
		}

		dirID := filepath.Base(filepath.Dir(manifestPath))
		if m.ID != dirID {
			return nil, fmt.Errorf("%s: id must match directory name (%s)", manifestPath, dirID)
		}

		manifests = append(manifests, m)
	}

	if err := EnforceCorePresence(manifests); err != nil {
		return nil, err
	}

	return manifests, nil
}
