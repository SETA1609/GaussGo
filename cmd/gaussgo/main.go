package main

import (
	"errors"
	"fmt"
	"strings"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/bootstrap"
	"gaussgo/internal/state"
)

func main() {
	fmt.Println("GaussGo")
	fmt.Println("=======")

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

	fmt.Printf("Discovered mods: %d\n", len(runtime.Diagnostics.DiscoveredMods))
	fmt.Printf("Load order: %v\n", runtime.Diagnostics.LoadOrder)
	fmt.Printf("Active state: %s\n", runtime.Context.ActiveStateID)
	fmt.Printf("Locale: %s\n", runtime.Context.CurrentLocale)
	fmt.Printf("Scene: %s\n", runtime.Controllers.Scene.Current())
	fmt.Println("Phase 3 runtime orchestration is active. TUI scenes are next.")
}
