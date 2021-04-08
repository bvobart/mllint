package categories

import "github.com/bvobart/mllint/api"

const (
	VersionControl      api.Category = "Version Control"
	FileFolderStructure api.Category = "File and Folder Structure"
	DependencyMgmt      api.Category = "Dependency Management"
	DataQuality         api.Category = "Data Quality"
	CodeQuality         api.Category = "Code Quality" // linting, CI usage, etc.
	Testing             api.Category = "Testing"
	Deployment          api.Category = "Deployment"
)

var All = []api.Category{
	VersionControl,
	FileFolderStructure,
	DependencyMgmt,
	DataQuality,
	CodeQuality,
	Testing,
	Deployment,
}

var BySlug = makeSlugMap()

func makeSlugMap() map[string]api.Category {
	res := map[string]api.Category{}
	for _, cat := range All {
		res[cat.Slug()] = cat
	}
	return res
}
