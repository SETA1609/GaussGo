package state

import (
	"os"
	"path/filepath"
	"strings"

	"gaussgo/internal/apperrors"
)

const activePointerFile = "active.txt"

func SetActive(statesDir string, stateID string) error {
	stateID = strings.TrimSpace(stateID)
	if stateID == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "stateID is required")
	}

	if err := os.MkdirAll(statesDir, 0o755); err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "create states directory", err)
	}

	path := filepath.Join(statesDir, activePointerFile)
	if err := os.WriteFile(path, []byte(stateID+"\n"), 0o644); err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "write active pointer", err)
	}

	return nil
}

func GetActive(statesDir string) (string, error) {
	path := filepath.Join(statesDir, activePointerFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", apperrors.Wrap(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "read active pointer", err)
	}

	stateID := strings.TrimSpace(string(data))
	if stateID == "" {
		return "", apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "active pointer is empty")
	}

	return stateID, nil
}
