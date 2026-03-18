package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gaussgo/internal/apperrors"
)

type unitFile struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Concepts []conceptFile `json:"concepts"`
}

type conceptFile struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
}

func loadUnitIDs(modsDir string, modID string) ([]string, error) {
	unitsDir := filepath.Join(modsDir, modID, "data", "units")
	entries, err := os.ReadDir(unitsDir)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "read units directory", err)
	}

	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := filepath.Ext(name)
		if ext != ".json" {
			continue
		}
		out = append(out, name[:len(name)-len(ext)])
	}
	return out, nil
}

func loadConceptsForUnit(modsDir string, modID string, unitID string) ([]string, map[string]conceptFile, error) {
	path := filepath.Join(modsDir, modID, "data", "units", unitID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "read unit file", err)
	}

	var unit unitFile
	if err := json.Unmarshal(data, &unit); err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeValidation, apperrors.ErrorTypeInput, "parse unit file", err)
	}
	if unit.ID == "" {
		unit.ID = unitID
	}

	ids := make([]string, 0, len(unit.Concepts))
	byID := make(map[string]conceptFile, len(unit.Concepts))
	for idx, concept := range unit.Concepts {
		cid := concept.ID
		if cid == "" {
			cid = fmt.Sprintf("concept-%d", idx+1)
		}
		ids = append(ids, cid)
		byID[cid] = concept
	}

	return ids, byID, nil
}
