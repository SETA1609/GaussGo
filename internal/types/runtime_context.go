package types

type RuntimeContext struct {
	ActiveStateID    string
	EnabledMods      []string
	CurrentSceneID   string
	CurrentModID     string
	CurrentUnitID    string
	CurrentConceptID string
	CurrentLocale    string
}
