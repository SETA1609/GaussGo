package mods

import (
	"fmt"
	"regexp"
	"strings"

	"gaussgo/internal/contracts"
)

var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func ValidateManifest(m contracts.ModManifest) error {
	if strings.TrimSpace(m.ID) == "" {
		return fmt.Errorf("invalid manifest: id is required")
	}
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("invalid manifest %s: name is required", m.ID)
	}
	if !semverRe.MatchString(strings.TrimSpace(m.Version)) {
		return fmt.Errorf("invalid manifest %s: version must be semver", m.ID)
	}
	if strings.TrimSpace(m.Entry) == "" {
		return fmt.Errorf("invalid manifest %s: entry is required", m.ID)
	}

	if strings.TrimSpace(m.Localization.Default) == "" {
		return fmt.Errorf("invalid manifest %s: localization.default is required", m.ID)
	}
	if len(m.Localization.Supported) == 0 {
		return fmt.Errorf("invalid manifest %s: localization.supported is required", m.ID)
	}
	if strings.TrimSpace(m.Localization.Path) == "" {
		return fmt.Errorf("invalid manifest %s: localization.path is required", m.ID)
	}

	foundDefault := false
	for _, locale := range m.Localization.Supported {
		if locale == m.Localization.Default {
			foundDefault = true
			break
		}
	}
	if !foundDefault {
		return fmt.Errorf("invalid manifest %s: localization.default must be in localization.supported", m.ID)
	}

	for _, dep := range m.Dependencies {
		if strings.TrimSpace(dep.ID) == "" {
			return fmt.Errorf("invalid manifest %s: dependency id is required", m.ID)
		}
		if dep.ID == m.ID {
			return fmt.Errorf("invalid manifest %s: self dependency is not allowed", m.ID)
		}
		if !isValidConstraint(dep.Version) {
			return fmt.Errorf("invalid manifest %s: dependency version constraint is invalid for %s", m.ID, dep.ID)
		}
	}

	if m.ID == "core" && len(m.Dependencies) > 0 {
		return fmt.Errorf("invalid manifest core: dependencies must be empty")
	}

	return nil
}

func isValidConstraint(in string) bool {
	in = strings.TrimSpace(in)
	if in == "" {
		return false
	}
	if semverRe.MatchString(in) {
		return true
	}
	for _, prefix := range []string{">=", "<=", ">", "<", "="} {
		if strings.HasPrefix(in, prefix) {
			return semverRe.MatchString(strings.TrimSpace(strings.TrimPrefix(in, prefix)))
		}
	}
	return false
}

func EnforceCorePresence(manifests []contracts.ModManifest) error {
	for _, m := range manifests {
		if m.ID == "core" {
			return nil
		}
	}

	return fmt.Errorf("core mod is required")
}
