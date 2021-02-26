package projectlinters

import (
	"fmt"

	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/utils/depsmgmt"
)

const (
	MsgUseDependencyManager = `Your project does not seem to be keeping track of its dependencies correctly,
	as no Pipfile, pyproject.toml or requirements.txt was found.
	The Python Packaging User Guide recommends using either Pipenv or Poetry as dependency managers.`

	MsgDontUsePip = `Your project uses a requirements.txt file to manage dependencies.
	Using Pip in such a raw fashion can be hard to maintain, especially if you used 'pip freeze' to generate it.
	The Python Packaging User Guide recommends using either Pipenv or Poetry as dependency managers.`
	// TODO: detect when `pip freeze` was used for generating requirements.txt

	MsgUseSingleDependencyManager = `Your project appears to be using multiple dependency managers: %s.
	Using multiple dependency managers creates confusion regarding which to install and in what order,
	as well as which to place a new dependency in, or update an existing dependency in.
	It is recommended to use only one package manager, either Pipenv or Poetry.`
)

type UseDependencyManager struct{}

func (l UseDependencyManager) Name() string {
	return "use-dependency-manager"
}

func (l UseDependencyManager) LintProject(projectdir string) ([]api.Issue, error) {
	// detect dependency managers
	depmanagers := depsmgmt.Detect(projectdir)
	if len(depmanagers) == 0 {
		return []api.Issue{api.NewIssue(l.Name(), api.SeverityError, MsgUseDependencyManager)}, nil
	}

	// if one and it's Poetry or Pipenv, accept
	if len(depmanagers) == 1 {
		if depmanagers[0].Type() == depsmgmt.TypePipenv || depmanagers[0].Type() == depsmgmt.TypePoetry {
			return nil, nil
		}

		// if just using Pip give warning
		if depmanagers[0].Type() == depsmgmt.TypePip {
			return []api.Issue{api.NewIssue(l.Name()+"/no-pip", api.SeverityWarning, MsgDontUsePip)}, nil
		}
	}

	// don't use multiple package managers, add Pip warning if necessary
	types := types(depmanagers)
	issues := []api.Issue{api.NewIssue(l.Name()+"/single", api.SeverityError, fmt.Sprintf(MsgUseSingleDependencyManager, types))}
	if contains(types, depsmgmt.TypePip) {
		issues = append(issues, api.NewIssue(l.Name()+"/no-pip", api.SeverityWarning, MsgDontUsePip))
	}
	return issues, nil
}

func types(managers []depsmgmt.DependencyManager) []depsmgmt.DependencyManagerType {
	types := []depsmgmt.DependencyManagerType{}
	for _, manager := range managers {
		types = append(types, manager.Type())
	}
	return types
}

func contains(types []depsmgmt.DependencyManagerType, target depsmgmt.DependencyManagerType) bool {
	for _, typ := range types {
		if typ == target {
			return true
		}
	}
	return false
}
