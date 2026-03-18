package mods

import (
	"regexp"
	"strings"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
)

var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func ValidateManifest(m contracts.ModManifest) error {
	if strings.TrimSpace(m.ID) == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest: id is required")
	}
	if strings.TrimSpace(m.Name) == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": name is required")
	}
	if !semverRe.MatchString(strings.TrimSpace(m.Version)) {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": version must be semver")
	}
	if strings.TrimSpace(m.Entry) == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": entry is required")
	}

	if strings.TrimSpace(m.Localization.Default) == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": localization.default is required")
	}
	if len(m.Localization.Supported) == 0 {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": localization.supported is required")
	}
	if strings.TrimSpace(m.Localization.Path) == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": localization.path is required")
	}

	foundDefault := false
	for _, locale := range m.Localization.Supported {
		if locale == m.Localization.Default {
			foundDefault = true
			break
		}
	}
	if !foundDefault {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": localization.default must be in localization.supported")
	}

	seenDeps := make(map[string]struct{}, len(m.Dependencies))
	for _, dep := range m.Dependencies {
		if strings.TrimSpace(dep.ID) == "" {
			return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": dependency id is required")
		}
		if _, exists := seenDeps[dep.ID]; exists {
			return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": duplicate dependency id "+dep.ID)
		}
		seenDeps[dep.ID] = struct{}{}
		if dep.ID == m.ID {
			return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": self dependency is not allowed")
		}
		if !isValidConstraint(dep.Version) {
			return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest "+m.ID+": dependency version constraint is invalid for "+dep.ID)
		}
	}

	if m.ID == "core" && len(m.Dependencies) > 0 {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "invalid manifest core: dependencies must be empty")
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

	return apperrors.New(apperrors.CodeCoreRequired, apperrors.ErrorTypeDomain, "core mod is required")
}
