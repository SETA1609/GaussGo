package mods

import (
	"testing"

	"gaussgo/internal/contracts"
)

func validManifest() contracts.ModManifest {
	return contracts.ModManifest{
		ID:      "linearAlgebra",
		Name:    "Linear Algebra",
		Version: "0.1.0",
		Entry:   "data/units",
		Localization: contracts.ModLocalization{
			Default:   "en",
			Supported: []string{"en", "es"},
			Path:      "i18n",
		},
		Dependencies: []contracts.ModDependency{{ID: "core", Version: ">=0.1.0"}},
		Provides:     []string{"learning.units"},
	}
}

func TestValidateManifestSuccess(t *testing.T) {
	m := validManifest()
	if err := ValidateManifest(m); err != nil {
		t.Fatalf("expected valid manifest, got error: %v", err)
	}
}

func TestValidateManifestInvalidSemver(t *testing.T) {
	m := validManifest()
	m.Version = "v1"
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected semver error")
	}
}

func TestValidateManifestInvalidLocalization(t *testing.T) {
	m := validManifest()
	m.Localization.Default = "fr"
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected localization default mismatch error")
	}
}

func TestValidateManifestRejectsCoreDependencies(t *testing.T) {
	m := validManifest()
	m.ID = "core"
	m.Dependencies = []contracts.ModDependency{{ID: "linearAlgebra", Version: ">=0.1.0"}}
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected core dependency error")
	}
}

func TestValidateManifestRejectsDuplicateDependencyID(t *testing.T) {
	m := validManifest()
	m.Dependencies = []contracts.ModDependency{
		{ID: "core", Version: ">=0.1.0"},
		{ID: "core", Version: ">=0.1.0"},
	}
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected duplicate dependency id error")
	}
}

func TestValidateManifestRejectsSelfDependency(t *testing.T) {
	m := validManifest()
	m.Dependencies = []contracts.ModDependency{{ID: "linearAlgebra", Version: ">=0.1.0"}}
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected self dependency error")
	}
}

func TestValidateManifestRejectsMissingRequiredFields(t *testing.T) {
	m := validManifest()
	m.ID = ""
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected missing id error")
	}

	m = validManifest()
	m.Name = ""
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected missing name error")
	}

	m = validManifest()
	m.Entry = ""
	if err := ValidateManifest(m); err == nil {
		t.Fatal("expected missing entry error")
	}
}
