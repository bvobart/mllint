package depmanagers

import (
	"io/ioutil"
	"path"
	"regexp"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
)

//---------------------------------------------------------------------------------------

type typeRequirementsTxt string

func (p typeRequirementsTxt) String() string {
	return string(p)
}

func (p typeRequirementsTxt) Detect(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "requirements.txt"))
}

func (p typeRequirementsTxt) For(project api.Project) api.DependencyManager {
	contents, err := ioutil.ReadFile(path.Join(project.Dir, "requirements.txt"))
	if err != nil {
		panic(err)
	}
	return RequirementsTxt{Project: project, reqFile: string(contents)}
}

//---------------------------------------------------------------------------------------

type RequirementsTxt struct {
	Project api.Project
	reqFile string
}

func (p RequirementsTxt) Type() api.DependencyManagerType {
	return TypeRequirementsTxt
}

func (p RequirementsTxt) HasDependency(dependency string) bool {
	// (?m) means multiline, i.e. ^ and $ will match on start and end of every line.
	matched, err := regexp.MatchString(`(?m)^\s*`+dependency, p.reqFile)
	return err == nil && matched
}

func (p RequirementsTxt) HasDevDependency(dependency string) bool {
	// a requirements.txt file has no concept of dev-dependencies unless users homebrew their own requirements-dev.txt or so,
	// so we just return false.
	return false
}

//---------------------------------------------------------------------------------------
//---------------------------------------------------------------------------------------

type typeSetupPy string

func (p typeSetupPy) String() string {
	return string(p)
}

func (p typeSetupPy) Detect(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "setup.py"))
}

func (p typeSetupPy) For(project api.Project) api.DependencyManager {
	return SetupPy{Project: project}
}

//---------------------------------------------------------------------------------------

type SetupPy struct {
	Project api.Project
}

func (p SetupPy) Type() api.DependencyManagerType {
	return TypeSetupPy
}

func (p SetupPy) HasDependency(dependency string) bool {
	// setup.py is a dynamic script, so this is too difficult to determine.
	return false
}

func (p SetupPy) HasDevDependency(dependency string) bool {
	// setup.py is a dynamic script, so this is too difficult to determine.
	return false
}

//---------------------------------------------------------------------------------------
