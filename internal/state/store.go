package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
)

type Store struct {
	statesDir string
}

// NewStore builds a state store bound to a states directory.
//
// The directory is created lazily on first write.
func NewStore(statesDir string) *Store {
	return &Store{statesDir: statesDir}
}

// Create initializes a new state with required defaults and persists it.
func (s *Store) Create(profileName string) (contracts.State, error) {
	profileName = strings.TrimSpace(profileName)
	if profileName == "" {
		return contracts.State{}, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "profileName is required")
	}

	now := time.Now().UTC()
	stateID := slugify(profileName)
	state := contracts.State{
		SchemaVersion: SchemaVersion,
		StateID:       stateID,
		ProfileName:   profileName,
		CreatedAt:     now,
		UpdatedAt:     now,
		EnabledMods:   []string{"core"},
		UI: contracts.UIState{
			Locale:      DefaultLocale,
			Preferences: map[string]string{},
		},
		Progress: map[string]contracts.ModProgress{},
	}

	normalized, _, err := NormalizeAndValidate(state)
	if err != nil {
		return contracts.State{}, err
	}

	if err := s.Save(normalized); err != nil {
		return contracts.State{}, err
	}

	return normalized, nil
}

// Load reads a state from disk and applies normalization rules.
//
// Load does not write back normalized values automatically; callers can Save
// if they want normalization persisted immediately.
func (s *Store) Load(stateID string) (contracts.State, error) {
	path := s.stateFilePath(stateID)
	data, err := os.ReadFile(path)
	if err != nil {
		return contracts.State{}, apperrors.Wrap(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "read state file", err)
	}

	var file fileState
	if err := json.Unmarshal(data, &file); err != nil {
		return contracts.State{}, apperrors.Wrap(apperrors.CodeValidation, apperrors.ErrorTypeInput, "parse state file", err)
	}

	loaded, err := fromFileState(file)
	if err != nil {
		return contracts.State{}, apperrors.Wrap(apperrors.CodeValidation, apperrors.ErrorTypeInput, "decode state timestamps", err)
	}

	normalized, _, err := NormalizeAndValidate(loaded)
	if err != nil {
		return contracts.State{}, err
	}

	return normalized, nil
}

// Save validates, normalizes, and atomically persists a state file.
//
// Atomic write strategy: write temp file -> fsync -> rename.
func (s *Store) Save(state contracts.State) error {
	normalized, _, err := NormalizeAndValidate(state)
	if err != nil {
		return err
	}
	normalized.UpdatedAt = time.Now().UTC()

	if err := os.MkdirAll(s.statesDir, 0o755); err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "create states directory", err)
	}

	payload, err := json.MarshalIndent(toFileState(normalized), "", "  ")
	if err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "encode state", err)
	}

	finalPath := s.stateFilePath(normalized.StateID)
	tmpPath := finalPath + ".tmp"

	f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "open temp state file", err)
	}

	if _, err := f.Write(payload); err != nil {
		_ = f.Close()
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "write temp state file", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "sync temp state file", err)
	}
	if err := f.Close(); err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "close temp state file", err)
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "rename temp state file", err)
	}

	return nil
}

// SetActive updates the active state pointer for the store directory.
func (s *Store) SetActive(stateID string) error {
	return SetActive(s.statesDir, stateID)
}

// GetActive reads the active state pointer for the store directory.
func (s *Store) GetActive() (string, error) {
	return GetActive(s.statesDir)
}

func (s *Store) stateFilePath(stateID string) string {
	return filepath.Join(s.statesDir, fmt.Sprintf("state-%s.json", stateID))
}

func slugify(in string) string {
	in = strings.ToLower(strings.TrimSpace(in))
	replacer := strings.NewReplacer(" ", "-", "_", "-", "/", "-", "\\", "-")
	in = replacer.Replace(in)
	in = strings.Trim(in, "-")
	if in == "" {
		return "default"
	}
	return in
}
