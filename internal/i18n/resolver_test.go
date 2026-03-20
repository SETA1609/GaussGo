package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolverFallbackOrder(t *testing.T) {
	modsDir := setupMods(t)
	r, err := NewResolver(modsDir, "en")
	if err != nil {
		t.Fatalf("new resolver failed: %v", err)
	}

	got := r.Resolve("linearAlgebra", "es", "menu.start_learning")
	if got != "Empezar a Aprender" {
		t.Fatalf("expected core es fallback, got %q", got)
	}
}

func TestResolverLiteralFallback(t *testing.T) {
	modsDir := setupMods(t)
	r, err := NewResolver(modsDir, "en")
	if err != nil {
		t.Fatalf("new resolver failed: %v", err)
	}

	const key = "missing.key.for.test"
	if got := r.Resolve("linearAlgebra", "es", key); got != key {
		t.Fatalf("expected literal fallback %q, got %q", key, got)
	}
}

func TestResolverNestedJSONKeys(t *testing.T) {
	modsDir := setupMods(t)
	r, err := NewResolver(modsDir, "en")
	if err != nil {
		t.Fatalf("new resolver failed: %v", err)
	}

	got := r.Resolve("linearAlgebra", "es", "nav.home")
	if got != "Inicio" {
		t.Fatalf("expected nested key translation Inicio, got %q", got)
	}
}

func TestResolverInterpolationPluralizationContext(t *testing.T) {
	modsDir := setupMods(t)
	r, err := NewResolver(modsDir, "en")
	if err != nil {
		t.Fatalf("new resolver failed: %v", err)
	}

	greeting := r.ResolveWithOptions("core", "en", "greeting", map[string]any{"name": "Ada"})
	if greeting != "Hello, Ada!" {
		t.Fatalf("expected interpolated greeting, got %q", greeting)
	}

	apples := r.ResolveWithOptions("core", "en", "apple", map[string]any{"count": 3})
	if apples != "3 apples" {
		t.Fatalf("expected pluralized apples, got %q", apples)
	}

	friend := r.ResolveWithOptions("core", "en", "friend", map[string]any{"context": "male"})
	if friend != "A boyfriend" {
		t.Fatalf("expected contextual translation, got %q", friend)
	}
}

func TestResolverNamespaceOrganization(t *testing.T) {
	modsDir := setupMods(t)
	r, err := NewResolver(modsDir, "en")
	if err != nil {
		t.Fatalf("new resolver failed: %v", err)
	}

	if got := r.Resolve("core", "en", "nav.home"); got != "Home" {
		t.Fatalf("expected common namespace nav.home, got %q", got)
	}
	if got := r.Resolve("core", "en", "auth.login"); got != "Log in" {
		t.Fatalf("expected auth namespace auth.login, got %q", got)
	}
}

func TestResolverSupportedLocales(t *testing.T) {
	modsDir := setupMods(t)
	r, err := NewResolver(modsDir, "en")
	if err != nil {
		t.Fatalf("new resolver: %v", err)
	}

	locales := r.SupportedLocales()

	// setupMods registers "en" and "es" across core and linearAlgebra
	if len(locales) != 2 {
		t.Fatalf("expected 2 supported locales, got %d: %v", len(locales), locales)
	}
	if locales[0] != "en" || locales[1] != "es" {
		t.Fatalf("expected [en es] sorted, got %v", locales)
	}

	// returned slice is a copy — mutations must not affect the resolver
	locales[0] = "xx"
	if r.SupportedLocales()[0] != "en" {
		t.Fatal("SupportedLocales returned a reference to internal slice")
	}
}

func setupMods(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	writeManifest(t, filepath.Join(root, "core"), `{
		"id":"core","name":"Core Runtime","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)
	writeLocale(t, filepath.Join(root, "core", "i18n", "en.json"), `{
		"menu": {"start_learning":"Start Learning", "exit":"Exit"},
		"nav": {"home":"Home"},
		"greeting":"Hello, {{name}}!",
		"apple_one":"{{count}} apple",
		"apple_other":"{{count}} apples",
		"friend":"A friend",
		"friend_male":"A boyfriend",
		"friend_female":"A girlfriend"
	}`)
	writeLocale(t, filepath.Join(root, "core", "i18n", "es.json"), `{
		"menu": {"start_learning":"Empezar a Aprender", "exit":"Salir"},
		"nav": {"home":"Inicio"},
		"greeting":"Hola, {{name}}!",
		"apple_one":"{{count}} manzana",
		"apple_other":"{{count}} manzanas",
		"friend":"Un amigo",
		"friend_male":"Un novio",
		"friend_female":"Una novia"
	}`)
	writeLocale(t, filepath.Join(root, "core", "i18n", "common", "en.json"), `{
		"nav": {"home":"Home", "about":"About", "contact":"Contact"}
	}`)
	writeLocale(t, filepath.Join(root, "core", "i18n", "common", "es.json"), `{
		"nav": {"home":"Inicio", "about":"Acerca", "contact":"Contacto"}
	}`)
	writeLocale(t, filepath.Join(root, "core", "i18n", "auth", "en.json"), `{
		"auth": {"login":"Log in", "logout":"Log out", "signup":"Sign up"}
	}`)
	writeLocale(t, filepath.Join(root, "core", "i18n", "auth", "es.json"), `{
		"auth": {"login":"Entrar", "logout":"Salir", "signup":"Crear cuenta"}
	}`)

	writeManifest(t, filepath.Join(root, "linearAlgebra"), `{
		"id":"linearAlgebra","name":"Linear Algebra","version":"0.1.0","entry":"data/units",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[{"id":"core","version":">=0.1.0"}],"provides":[]}`)
	writeLocale(t, filepath.Join(root, "linearAlgebra", "i18n", "common", "en.json"), `{
		"mod.name":"Linear Algebra"
	}`)
	writeLocale(t, filepath.Join(root, "linearAlgebra", "i18n", "common", "es.json"), `{
		"mod.name":"Algebra Lineal"
	}`)

	return root
}

func writeManifest(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest %s: %v", dir, err)
	}
}

func writeLocale(t *testing.T, file string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatalf("mkdir locale dir %s: %v", filepath.Dir(file), err)
	}
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("write locale file %s: %v", file, err)
	}
}
