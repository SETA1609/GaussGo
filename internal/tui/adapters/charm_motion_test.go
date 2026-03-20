package adapters

import "testing"

func TestCharmSpringUpdate(t *testing.T) {
	spring := NewCharmSpring()
	pos, vel := spring.Update(0, 0, 1)
	if pos == 0 && vel == 0 {
		t.Fatal("expected spring motion to update state")
	}
}
