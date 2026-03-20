package adapters

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type minimalModel struct{}

func (minimalModel) Init() tea.Cmd { return tea.Quit }

func (minimalModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return minimalModel{}, tea.Quit }

func (minimalModel) View() tea.View { return tea.NewView("ok") }

func TestBubbleTeaProgramFactoryRejectsNonModel(t *testing.T) {
	factory := NewBubbleTeaProgramFactory()
	program := factory.New("not-a-model")
	if err := program.Run(); err == nil {
		t.Fatal("expected error when model is invalid")
	}
}

func TestBubbleTeaProgramFactoryAcceptsTeaModel(t *testing.T) {
	factory := NewBubbleTeaProgramFactory()
	program := factory.New(minimalModel{})
	if program == nil {
		t.Fatal("expected non-nil program")
	}
}
