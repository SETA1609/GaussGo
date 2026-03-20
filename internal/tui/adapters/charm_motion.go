package adapters

import (
	"github.com/charmbracelet/harmonica"

	"gaussgo/internal/tui/ports"
)

type CharmSpring struct {
	spring harmonica.Spring
}

func NewCharmSpring() ports.Spring {
	return &CharmSpring{spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.7)}
}

func (s *CharmSpring) Update(pos float64, vel float64, target float64) (float64, float64) {
	return s.spring.Update(pos, vel, target)
}
