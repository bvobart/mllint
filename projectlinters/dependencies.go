package projectlinters

import (
	"fmt"

	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/utils/depsmgmt"
)

// TODO: include links to documentation sites where users can go for more info, e.g. about poetry / pipenv.

const (
	MsgUseDependencyManager = `Your project does not seem to be keeping track of its dependencies correctly,
		as no Pipfile, pyproject.toml or requirements.txt was found.
		The Python Packaging User Guide recommends using either Pipenv or Poetry as dependency managers.
		>  The recommendation is to use Pipenv if your project is an application
		>  and to use Poetry if it is a library or otherwise needs to be built into a Python package.`

	// TODO: detect when `pip freeze` was used for generating requirements.txt
	MsgNoRequirementsTxt = `Your project uses a  requirements.txt  file to manage dependencies.
		Using Pip in such a raw fashion can be hard to maintain, especially if you used 'pip freeze' to generate it.
		The Python Packaging User Guide recommends using either Pipenv or Poetry as dependency managers.
		>  The recommendation is to use Pipenv if your project is an application
		>  and to use Poetry if it is a library or otherwise needs to be built into a Python package.`

	MsgNoSetupPy = `Your project uses a  setup.py  file to manage dependencies.
		While using a setup.py is more maintainable than a requirements.txt, it may still be difficult to maintain.
		>  The Python Packaging User Guide recommends using either Pipenv or Poetry as dependency managers.
		>  Since you are already using a setup.py, you'll likely want to use Poetry as it is able to build and publish Python packages.`

	MsgUseSingleDependencyManager = `Your project appears to be using multiple dependency managers:  %s
		Using multiple dependency managers creates confusion regarding which to install and in what order,
		as well as which to place a new dependency in, or update an existing dependency in.
		>  It is recommended to use only one package manager, either Pipenv or Poetry.`

	MsgDontCombinePipenvSetupPy = `Your project appears to be using Pipenv to manage dependencies, but also has a  setup.py  file.
		Pipenv does support building and installing packages from a  setup.py  with 'pipenv install -e .',
		but it might be worth switching to Poetry. Poetry is very similar to Pipenv, but also
		supports building and publishing Python packages, which I presume is what you are using  setup.py  for now.`

	MsgDontCombineRequirementsTxtSetupPy = `Your project appears to be using both a  requirements.txt  and a  setup.py  file to manage dependencies.
		The Python Packaging User Guide recommends using either Pipenv or Poetry as dependency managers.
		>  The recommendation is to use Pipenv if your project is an application
		>  and to use Poetry if it is a library or otherwise needs to be built into a Python package.`

	MsgDontCombineRequirementsTxtPoetryPipenv = `Your project appears to be using both a  requirements.txt  as well as Poetry or Pipenv to manage dependencies.
		This is redundant and creates confusion regarding which to install and in what order if both need to be installed.
		It will also be confusing in which of the two a new dependency should be tracked, or an existing one updated.
		>  Migrate the rest of your dependencies to Pipenv / Poetry and remove the requirements.txt file.`

	MsgDontCombinePoetrySetupPy = `Your project appears to be using both Poetry and a  setup.py  file to manage dependencies and build into a Python package.
		This is redundant and creates confusion regarding how to build your package, edit its details, or where to place the code for building the package.
		>  Migrate any remaining logic from the setup.py file into Poetry and remove the setup.py`
)

const (
	RuleSingle                                 = "single"
	RuleNoRequirementsTxt                      = "no-requirements-txt"
	RuleNoSetupPy                              = "no-setup-py"
	RuleDontCombinePipenvSetupPy               = "no-combine-pipenv-setup-py"
	RuleDontCombinePoetrySetupPy               = "no-combine-poetry-setup-py"
	RuleDontCombineRequirementsTxtPoetryPipenv = "no-combine-poetry-pipenv"
	RuleDontCombineRequirementsTxtSetupPy      = "no-combine-requirements-txt-setup-py"
)

// UseDependencyManager is a linter to check whether projects use proper dependency management,
// There are two possible situations that this linter accepts without errors:
// - Project only uses Poetry
// - Project only uses Pipenv
//
// There are also situations where warnings are emitted:
// - Project uses only requirements.txt or only setup.py (recommend using Pipenv or Poetry respectively)
// - Project uses Pipenv, but also has a setup.py (since Pipenv doesn't support building and publishing packages directly from a Pipfile)
// - Project uses requirements.txt and setup.py (recommend using Poetry for dependency management)
// - Project uses both Pipenv or Poetry and a requirements.txt (creates confusion, remove the requirements.txt)
// - Project uses both Poetry and a setup.py (creates confusion, remove the setup.py)
//
// Finally, an error is emitted in the following situations:
// - Project is not using any dependency management.
type UseDependencyManager struct{}

func (l UseDependencyManager) Name() string {
	return "use-dependency-manager"
}

func (l UseDependencyManager) Rules() []string {
	return []string{RuleSingle, RuleNoRequirementsTxt, RuleNoSetupPy, RuleDontCombinePipenvSetupPy, RuleDontCombinePoetrySetupPy, RuleDontCombineRequirementsTxtPoetryPipenv, RuleDontCombineRequirementsTxtSetupPy}
}

func (l UseDependencyManager) LintProject(projectdir string) ([]api.Issue, error) {
	// detect dependency managers
	depmanagers := depsmgmt.Detect(projectdir)
	if len(depmanagers) == 0 {
		return []api.Issue{api.NewIssue(l.Name(), "", api.SeverityError, MsgUseDependencyManager)}, nil
	}

	// if using only one dependency manager
	if len(depmanagers) == 1 {
		switch depmanagers[0].Type() {
		// ... and it's Poetry or Pipenv, accept
		case depsmgmt.TypePipenv:
			fallthrough
		case depsmgmt.TypePoetry:
			return nil, nil
		// if using Pip requirements.txt, give warning
		case depsmgmt.TypeRequirementsTxt:
			return []api.Issue{api.NewIssue(l.Name(), RuleNoRequirementsTxt, api.SeverityWarning, MsgNoRequirementsTxt)}, nil
		// if using Pip setup.py, give warning
		case depsmgmt.TypeSetupPy:
			return []api.Issue{api.NewIssue(l.Name(), RuleNoSetupPy, api.SeverityWarning, MsgNoSetupPy)}, nil
		default:
			return nil, fmt.Errorf("new dependency manager %s was added, but %s linter was not updated", depmanagers[0].Type(), l.Name())
		}
	}

	// don't use multiple package managers
	types := types(depmanagers)
	issues := []api.Issue{}

	// don't combine Pipenv and setup.py, consider using Poetry instead, info
	if containsAll(types, depsmgmt.TypePipenv, depsmgmt.TypeSetupPy) {
		issues = append(issues, api.NewIssue(l.Name(), RuleDontCombinePipenvSetupPy, api.SeverityInfo, MsgDontCombinePipenvSetupPy))
	}

	// don't combine requirements.txt and setup.py, use Poetry, warning
	if containsAll(types, depsmgmt.TypeRequirementsTxt, depsmgmt.TypeSetupPy) {
		issues = append(issues, api.NewIssue(l.Name(), RuleDontCombineRequirementsTxtSetupPy, api.SeverityWarning, MsgDontCombineRequirementsTxtSetupPy))
	}

	// don't combine requirements.txt with Pipenv or Poetry, simply use Pipenv or Poetry, warning
	if contains(types, depsmgmt.TypeRequirementsTxt) && (contains(types, depsmgmt.TypePipenv) || contains(types, depsmgmt.TypePoetry)) {
		issues = append(issues, api.NewIssue(l.Name(), RuleDontCombineRequirementsTxtPoetryPipenv, api.SeverityWarning, MsgDontCombineRequirementsTxtPoetryPipenv))
	}

	// don't combine Poetry and setup.py, it's redundant, just use setup.py, info
	if containsAll(types, depsmgmt.TypePoetry, depsmgmt.TypeSetupPy) {
		issues = append(issues, api.NewIssue(l.Name(), RuleDontCombinePoetrySetupPy, api.SeverityInfo, MsgDontCombinePoetrySetupPy))
	}

	// add a generic warning if no more specific warning was emitted
	if len(issues) == 0 {
		issues = append(issues, api.NewIssue(l.Name(), RuleSingle, api.SeverityWarning, fmt.Sprintf(MsgUseSingleDependencyManager, types)))
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

func containsAll(types []depsmgmt.DependencyManagerType, targets ...depsmgmt.DependencyManagerType) bool {
	for _, target := range targets {
		if !contains(types, target) {
			return false
		}
	}
	return true
}
