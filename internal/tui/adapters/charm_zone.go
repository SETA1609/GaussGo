package adapters

import (
	zone "github.com/lrstanley/bubblezone/v2"

	"gaussgo/internal/tui/ports"
)

type CharmZone struct{}

func NewCharmZone() ports.Zone {
	return CharmZone{}
}

func (CharmZone) Init() {
	zone.NewGlobal()
}

func (CharmZone) Mark(id string, content string) string {
	return zone.Mark(id, content)
}

func (CharmZone) Scan(content string) string {
	return zone.Scan(content)
}
