package ports

type Program interface {
	Run() error
}

type ProgramFactory interface {
	New(model any) Program
}
