package adapters

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"gaussgo/internal/tui/ports"
)

type BubbleTeaProgramFactory struct{}

func NewBubbleTeaProgramFactory() ports.ProgramFactory {
	return BubbleTeaProgramFactory{}
}

func (BubbleTeaProgramFactory) New(model any) ports.Program {
	teaModel, ok := model.(tea.Model)
	if !ok {
		return staticErrorProgram{err: fmt.Errorf("model does not implement tea.Model")}
	}
	return bubbleTeaProgram{program: tea.NewProgram(teaModel)}
}

type bubbleTeaProgram struct {
	program *tea.Program
}

func (p bubbleTeaProgram) Run() error {
	_, err := p.program.Run()
	return err
}

type staticErrorProgram struct {
	err error
}

func (p staticErrorProgram) Run() error {
	return p.err
}
