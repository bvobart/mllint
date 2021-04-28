package api

type CQLinterType string

type CQLinter interface {
	Type() CQLinterType
	Detect(project Project) bool
	Run(projectdir string) error // TODO: find a way to allow this to return its full output
}
