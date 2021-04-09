package categories

import "github.com/bvobart/mllint/api"

// TODO: add category descriptions

var VersionControl = api.Category{
	Name:        "Version Control",
	Slug:        "version-control",
	Description: "TODO",
}

var FileStructure = api.Category{
	Name:        "File Structure",
	Slug:        "file-structure",
	Description: "TODO",
}

var DependencyMgmt = api.Category{
	Name:        "Dependency Management",
	Slug:        "dependency-management",
	Description: "TODO",
}

// linting, CI usage, etc.
var CodeQuality = api.Category{
	Name:        "Code Quality",
	Slug:        "code-quality",
	Description: "TODO",
}

var DataQuality = api.Category{
	Name:        "Data Quality",
	Slug:        "data-quality",
	Description: "TODO",
}

var Testing = api.Category{
	Name:        "Testing",
	Slug:        "testing",
	Description: "TODO",
}

var ContinuousIntegration = api.Category{
	Name:        "Continuous Integration",
	Slug:        "ci",
	Description: "TODO",
}

var Deployment = api.Category{
	Name:        "Deployment",
	Slug:        "deployment",
	Description: "TODO",
}

var All = []api.Category{
	VersionControl,
	FileStructure,
	DependencyMgmt,
	CodeQuality,
	DataQuality,
	Testing,
	ContinuousIntegration,
	Deployment,
}

var BySlug = makeSlugMap()

func makeSlugMap() map[string]api.Category {
	res := map[string]api.Category{}
	for _, cat := range All {
		res[cat.Slug] = cat
	}
	return res
}
