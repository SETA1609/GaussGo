package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/bootstrap"
	"gaussgo/internal/state"
	"gaussgo/internal/tui"
)

func main() {
	if !isQuietMode() {
		fmt.Println("GaussGo")
		fmt.Println("=======")
	}

	services := bootstrap.DefaultServices()
	runtime, err := bootstrap.BootstrapRuntime("mods", "states", "main_menu", services)
	if err != nil {
		var appErr apperrors.Error
		if errors.As(err, &appErr) && appErr.Code == apperrors.CodeNotFound && strings.Contains(appErr.Message, "read active pointer") {
			store := state.NewStore("states")
			st, createErr := store.Create("default")
			if createErr != nil {
				fmt.Println("startup failed creating default state:", createErr)
				return
			}
			if setErr := store.SetActive(st.StateID); setErr != nil {
				fmt.Println("startup failed setting active state:", setErr)
				return
			}

			runtime, err = bootstrap.BootstrapRuntime("mods", "states", "main_menu", services)
		}
	}
	if err != nil {
		fmt.Println("startup failed:", err)
		return
	}

	model, modelErr := tui.NewAppModel(runtime, "states")
	if modelErr != nil {
		fmt.Println("startup failed building TUI model:", modelErr)
		return
	}

	p := tea.NewProgram(model)
	if _, runErr := p.Run(); runErr != nil {
		fmt.Println("tui failed:", runErr)
		return
	}
}

func isQuietMode() bool {
	level := strings.ToLower(strings.TrimSpace(os.Getenv("GAUSSGO_LOG_LEVEL")))
	return envBool("GAUSSGO_QUIET") || envBool("GAUSSGO_SILENT") || level == "" || level == "silent" || level == "off" || level == "none"
}

func envBool(name string) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(name)))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}
