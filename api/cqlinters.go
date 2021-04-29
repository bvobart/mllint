package api

import "fmt"

type CQLinterType string

type CQLinter interface {
	fmt.Stringer
	Type() CQLinterType
	IsInstalled() bool

	Detect(project Project) bool
	Run(projectdir string) error // TODO: find a way to allow this to return its full output
}
