package mods

import (
	"os"
	"path/filepath"
	"sort"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
)

func Discover(modsDir string) ([]contracts.ModManifest, error) {
	if _, err := os.Stat(modsDir); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "mods directory not found", err)
	}

	paths, err := filepath.Glob(filepath.Join(modsDir, "*", "manifest.json"))
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "discover manifests", err)
	}

	sort.Strings(paths)

	manifests := make([]contracts.ModManifest, 0, len(paths))
	for _, manifestPath := range paths {
		m, err := ParseManifestFile(manifestPath)
		if err != nil {
			return nil, apperrors.Wrap(apperrors.CodeValidation, apperrors.ErrorTypeInput, manifestPath+": invalid manifest", err)
		}

		dirID := filepath.Base(filepath.Dir(manifestPath))
		if m.ID != dirID {
			return nil, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, manifestPath+": id must match directory name ("+dirID+")")
		}

		manifests = append(manifests, m)
	}

	if err := EnforceCorePresence(manifests); err != nil {
		return nil, err
	}

	return manifests, nil
}
