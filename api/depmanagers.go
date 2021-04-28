package api

type DependencyManagerType string

type DependencyManager interface {
	Detect(project Project) bool
	Type() DependencyManagerType
}
