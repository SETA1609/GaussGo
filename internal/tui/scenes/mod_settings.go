package scenes

import (
	"fmt"

	"gaussgo/internal/contracts"
)

func buildModSettings(s RuntimeSnapshot) ViewModel {
	modID := "core"
	lines := []string{tr(s, modID, "scene.mod_settings.help", "Manage enabled mods.")}
	if s.Flash != "" {
		lines = append(lines, s.Flash)
	}

	options := []Option{{Label: tr(s, modID, "menu.mod_settings.refresh", "Refresh Mods"), Action: Action{Type: ActionRefreshMods}}}
	for _, status := range s.ModStatuses {
		label := fmt.Sprintf("%s [%s]", status.ModID, statusState(s, status))
		action := Action{Type: ActionToggleMod, Value: status.ModID}
		disabled := status.ModID == "core" || !status.Available
		if !status.Available && status.ReasonDisabled != "" {
			label = label + " - " + status.ReasonDisabled
		}
		options = append(options, Option{Label: label, Action: action, Disabled: disabled})
	}
	options = append(options, Option{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}})

	return ViewModel{
		Title:    tr(s, modID, "scene.mod_settings.title", "Mod Settings"),
		Subtitle: tr(s, modID, "scene.mod_settings.subtitle", "Toggle non-core modules"),
		Lines:    lines,
		Options:  options,
	}
}

func statusState(s RuntimeSnapshot, status contracts.ModStatus) string {
	if status.ModID == "core" {
		return tr(s, "core", "status.locked", "locked")
	}
	if status.Enabled {
		return tr(s, "core", "status.enabled", "enabled")
	}
	if !status.Available {
		return tr(s, "core", "status.unavailable", "unavailable")
	}
	return tr(s, "core", "status.disabled", "disabled")
}
