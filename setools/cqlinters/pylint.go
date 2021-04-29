package cqlinters

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

type Pylint struct{}

func (p Pylint) Type() api.CQLinterType {
	return TypePylint
}

func (p Pylint) String() string {
	return "Pylint"
}

func (p Pylint) Detect(project api.Project) bool {
	// TODO: implement
	return false
}

func (p Pylint) Run(projectdir string) error {
	// TODO: copy from python-ml-analysis
	return fmt.Errorf("not implemented")
}
