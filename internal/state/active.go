package state

import (
	"errors"
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
	tmpPath := path + ".tmp"

	f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "open active pointer temp file", err)
	}
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := f.Write([]byte(stateID + "\n")); err != nil {
		_ = f.Close()
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "write active pointer temp file", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "sync active pointer temp file", err)
	}
	if err := f.Close(); err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "close active pointer temp file", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "rename active pointer temp file", err)
	}

	return nil
}

func GetActive(statesDir string) (string, error) {
	path := filepath.Join(statesDir, activePointerFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", apperrors.Wrap(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "read active pointer", err)
		}
		return "", apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "read active pointer", err)
	}

	stateID := strings.TrimSpace(string(data))
	if stateID == "" {
		return "", apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "active pointer is empty")
	}

	return stateID, nil
}
