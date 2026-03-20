package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	i18next "codeberg.org/lil5/i18next_go"

	"gaussgo/internal/contracts"
	"gaussgo/internal/mods"
)

var _ contracts.I18nResolver = (*Resolver)(nil)

type Resolver struct {
	engine        *i18next.I18next
	defaultLocale string
	modDefaults   map[string]string
	supported     []string
	resources     map[string]map[string]map[string]string
	orderByMod    map[string][]string
}

// NewResolver discovers all mods in the given directory and loads their
// localization resources. It supports both flat locale files and nested
// namespace directories within each mod's localization path.
func NewResolver(modsDir string, defaultLocale string) (*Resolver, error) {
	if strings.TrimSpace(defaultLocale) == "" {
		defaultLocale = "en"
	}

	manifests, err := mods.Discover(modsDir)
	if err != nil {
		return nil, err
	}

	resources := map[string]map[string]map[string]string{}
	languages := map[string]struct{}{defaultLocale: {}}
	namespaceSet := map[string]struct{}{}
	modDefaults := make(map[string]string, len(manifests))
	orderByMod := map[string][]string{}

	for _, manifest := range manifests {
		modDefaults[manifest.ID] = manifest.Localization.Default
		orderByMod[manifest.ID] = appendIfMissing(orderByMod[manifest.ID], manifest.ID)

		baseLocaleDir := filepath.Join(modsDir, manifest.ID, manifest.Localization.Path)

		for _, locale := range manifest.Localization.Supported {
			locale = strings.TrimSpace(locale)
			if locale == "" {
				continue
			}
			languages[locale] = struct{}{}

			legacyPath := filepath.Join(baseLocaleDir, locale+".json")
			if entries, ok := readLocaleEntries(legacyPath); ok {
				addResource(resources, locale, manifest.ID, entries)
				namespaceSet[manifest.ID] = struct{}{}
			}
		}

		namespaceDirs, err := os.ReadDir(baseLocaleDir)
		if err != nil {
			continue
		}
		for _, nsDir := range namespaceDirs {
			if !nsDir.IsDir() {
				continue
			}
			nsName := strings.TrimSpace(nsDir.Name())
			if nsName == "" {
				continue
			}
			fullNamespace := manifest.ID + "." + nsName
			orderByMod[manifest.ID] = appendIfMissing(orderByMod[manifest.ID], fullNamespace)

			for _, locale := range manifest.Localization.Supported {
				locale = strings.TrimSpace(locale)
				if locale == "" {
					continue
				}
				languages[locale] = struct{}{}
				nsPath := filepath.Join(baseLocaleDir, nsName, locale+".json")
				if entries, ok := readLocaleEntries(nsPath); ok {
					addResource(resources, locale, fullNamespace, entries)
					namespaceSet[fullNamespace] = struct{}{}
				}
			}
		}
	}

	orderByMod["core"] = appendIfMissing(orderByMod["core"], "core.common")
	orderByMod["core"] = appendIfMissing(orderByMod["core"], "core")

	namespaces := make([]string, 0, len(namespaceSet))
	for ns := range namespaceSet {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)

	supported := make([]string, 0, len(languages))
	for locale := range languages {
		supported = append(supported, locale)
	}
	sort.Strings(supported)

	engine := i18next.Init(i18next.Options{
		Languages:   supported,
		FallbackLng: defaultLocale,
		Ns:          namespaces,
		DefaultNs:   "core",
		Resources:   resources,
	})

	return &Resolver{
		engine:        engine,
		defaultLocale: defaultLocale,
		modDefaults:   modDefaults,
		supported:     supported,
		resources:     resources,
		orderByMod:    orderByMod,
	}, nil
}

// Resolve provides a simple interface for translating a key within a mod's scope.
func (r *Resolver) Resolve(modID string, locale string, key string) string {
	return r.ResolveWithOptions(modID, locale, key, nil)
}

// ResolveWithOptions performs a full localization lookup with fallback logic.
// The search order is:
// 1. Target locale in the specific mod's namespaces.
// 2. Target locale in the "core" mod's namespaces.
// 3. Mod's default locale in the mod's namespaces.
// 4. "core" mod's default locale in the "core" namespaces.
// It also handles interpolation and pluralization based on the provided options.
func (r *Resolver) ResolveWithOptions(modID string, locale string, key string, options map[string]any) string {
	if r == nil || r.engine == nil || strings.TrimSpace(key) == "" {
		return key
	}

	locale = strings.TrimSpace(locale)
	if locale == "" {
		locale = r.defaultLocale
	}
	modID = strings.TrimSpace(modID)
	if modID == "" {
		modID = "core"
	}

	modDefault := r.modDefaults[modID]
	if modDefault == "" {
		modDefault = r.defaultLocale
	}
	coreDefault := r.modDefaults["core"]
	if coreDefault == "" {
		coreDefault = r.defaultLocale
	}

	sequence := [][2]string{
		{locale, modID},
		{locale, "core"},
		{modDefault, modID},
		{coreDefault, "core"},
	}
	candidates := keyCandidates(key, options)

	for _, step := range sequence {
		for _, namespace := range r.namespacesFor(step[1]) {
			if !r.hasNamespace(step[0], namespace) {
				continue
			}
			ctx := r.engine.GetContext(step[0])
			payload := map[string]any{"ns": namespace}
			for k, v := range options {
				payload[k] = v
			}
			for _, candidate := range candidates {
				resolved := ctx.T(candidate, payload)
				if resolved != candidate {
					return resolved
				}
			}
		}
	}

	return key
}

// SupportedLocales returns a sorted list of all locale codes discovered
// across all loaded mod manifests. Returns a copy of the internal slice.
func (r *Resolver) SupportedLocales() []string {
	if r == nil {
		return nil
	}
	return append([]string(nil), r.supported...)
}

func (r *Resolver) hasNamespace(locale string, namespace string) bool {
	if r == nil {
		return false
	}
	ns, ok := r.resources[locale]
	if !ok {
		return false
	}
	_, ok = ns[namespace]
	return ok
}

func (r *Resolver) Resource(locale string, namespace string, key string) (string, error) {
	if !r.hasNamespace(locale, namespace) {
		return "", fmt.Errorf("resource namespace not found: locale=%s namespace=%s", locale, namespace)
	}
	v := r.resources[locale][namespace][key]
	if strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("resource key not found: locale=%s namespace=%s key=%s", locale, namespace, key)
	}
	return v, nil
}

func (r *Resolver) namespacesFor(modID string) []string {
	list := append([]string(nil), r.orderByMod[modID]...)
	if len(list) == 0 {
		list = append(list, modID)
	}
	if modID != "core" {
		list = append(list, r.orderByMod["core"]...)
		list = appendIfMissing(list, "core")
	}
	return unique(list)
}

func addResource(resources map[string]map[string]map[string]string, locale string, namespace string, entries map[string]string) {
	if resources[locale] == nil {
		resources[locale] = map[string]map[string]string{}
	}
	resources[locale][namespace] = entries
}

func readLocaleEntries(path string) (map[string]string, bool) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var raw map[string]any
	if err := json.Unmarshal(content, &raw); err != nil {
		return nil, false
	}
	entries := map[string]string{}
	flattenEntries("", raw, entries)
	return entries, true
}

func appendIfMissing(in []string, v string) []string {
	for _, item := range in {
		if item == v {
			return in
		}
	}
	return append(in, v)
}

func unique(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, value := range in {
		if strings.TrimSpace(value) == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func keyCandidates(key string, options map[string]any) []string {
	out := make([]string, 0, 4)
	if strings.TrimSpace(key) == "" {
		return out
	}

	if options != nil {
		if context, ok := options["context"]; ok {
			if contextStr, ok := context.(string); ok && strings.TrimSpace(contextStr) != "" {
				out = append(out, key+"_"+contextStr)
			}
		}
		if countVal, ok := options["count"]; ok {
			if isSingularCount(countVal) {
				out = appendIfMissing(out, key+"_one")
			} else {
				out = appendIfMissing(out, key+"_other")
			}
		}
	}

	out = appendIfMissing(out, key)
	return out
}

func isSingularCount(v any) bool {
	switch n := v.(type) {
	case int:
		return n == 1
	case int8:
		return n == 1
	case int16:
		return n == 1
	case int32:
		return n == 1
	case int64:
		return n == 1
	case uint:
		return n == 1
	case uint8:
		return n == 1
	case uint16:
		return n == 1
	case uint32:
		return n == 1
	case uint64:
		return n == 1
	case float32:
		return n == 1
	case float64:
		return n == 1
	default:
		return false
	}
}

func flattenEntries(prefix string, in map[string]any, out map[string]string) {
	for key, value := range in {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch typed := value.(type) {
		case string:
			out[fullKey] = typed
		case map[string]any:
			flattenEntries(fullKey, typed, out)
		}
	}
}
