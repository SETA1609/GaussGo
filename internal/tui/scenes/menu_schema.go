package scenes

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type MenuSchema struct {
	TitleKey    string           `json:"titleKey"`
	Title       string           `json:"title"`
	SubtitleKey string           `json:"subtitleKey"`
	Subtitle    string           `json:"subtitle"`
	Lines       []menuLineSchema `json:"lines"`
	Items       []menuItemSchema `json:"items"`
}

type menuLineSchema struct {
	Key      string `json:"key"`
	Fallback string `json:"fallback"`
	Prefix   string `json:"prefix"`
	Source   string `json:"source"`
}

type menuItemSchema struct {
	LabelKey   string `json:"labelKey"`
	Label      string `json:"label"`
	Action     string `json:"action"`
	Value      string `json:"value"`
	DisabledIf string `json:"disabledIf"`
}

func loadMenuSchemaFromFile(path string) (MenuSchema, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return MenuSchema{}, err
	}
	var schema MenuSchema
	if err := json.Unmarshal(content, &schema); err != nil {
		return MenuSchema{}, err
	}
	return schema, nil
}

func LoadMenuSchemas(modsDir string) (map[string]MenuSchema, error) {
	schemas := map[string]MenuSchema{}

	coreDir := filepath.Join(modsDir, "core", "ui", "scenes")
	entries, err := os.ReadDir(coreDir)
	if err != nil {
		return schemas, nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) != ".json" {
			continue
		}
		schema, readErr := loadMenuSchemaFromFile(filepath.Join(coreDir, name))
		if readErr != nil {
			return nil, readErr
		}
		sceneID := name[:len(name)-len(filepath.Ext(name))]
		schemas[sceneID] = schema
	}

	return schemas, nil
}

func mustMenuSchema(raw string) MenuSchema {
	var schema MenuSchema
	if err := json.Unmarshal([]byte(raw), &schema); err != nil {
		panic(err)
	}
	return schema
}

func buildFromSchema(s RuntimeSnapshot, modID string, schema MenuSchema) ViewModel {
	lines := make([]string, 0, len(schema.Lines))
	for _, line := range schema.Lines {
		base := tr(s, modID, line.Key, line.Fallback)
		suffix := valueFromSource(s, line.Source)
		if suffix != "" {
			base += suffix
		}
		lines = append(lines, base)
	}

	options := make([]Option, 0, len(schema.Items))
	for _, item := range schema.Items {
		options = append(options, Option{
			Label:    tr(s, modID, item.LabelKey, item.Label),
			Action:   Action{Type: item.Action, Value: item.Value},
			Disabled: isDisabledByCondition(s, item.DisabledIf),
		})
	}

	return ViewModel{
		Title:    tr(s, modID, schema.TitleKey, schema.Title),
		Subtitle: tr(s, modID, schema.SubtitleKey, schema.Subtitle),
		Lines:    lines,
		Options:  options,
	}
}

func valueFromSource(s RuntimeSnapshot, source string) string {
	switch source {
	case "activeStateID":
		return s.ActiveStateID
	case "currentLocale":
		return s.CurrentLocale
	case "currentModID":
		return s.CurrentModID
	default:
		return ""
	}
}

func isDisabledByCondition(s RuntimeSnapshot, condition string) bool {
	switch condition {
	case "locale_en":
		return s.CurrentLocale == "en"
	case "locale_es":
		return s.CurrentLocale == "es"
	default:
		return false
	}
}
