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

//---------------------------------------------------------------------------------------
