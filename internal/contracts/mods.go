package contracts

type ModDependency struct {
	ID      string
	Version string
}

type ModLocalization struct {
	Default   string
	Supported []string
	Path      string
}

type ModManifest struct {
	ID           string
	Name         string
	Version      string
	Entry        string
	Localization ModLocalization
	Dependencies []ModDependency
	Provides     []string
}

type ModStatus struct {
	ModID          string
	Enabled        bool
	Available      bool
	ReasonDisabled string
}

type ModsRepository interface {
	Discover() ([]ModManifest, error)
	ResolveLoadOrder([]ModManifest) ([]string, error)
	Refresh() ([]ModStatus, error)
}
