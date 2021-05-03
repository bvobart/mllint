package api

import "fmt"

type DependencyManagerType interface {
	fmt.Stringer

	// Detect whether a project is using this type of dependency manager.
	// Generally, this is done by detecting whether the manager's configuration file exists.
	Detect(project Project) bool

	// For instantiates an DependencyManager instance for that project. This will parse the manager's configuration file.
	// Note that this may panic if Detect has not previously been called.
	For(project Project) DependencyManager
}

type DependencyManager interface {
	// HasDependency should return true if this dependency manager is tracking this dependency.
	// This means it can either be in the regular dependencies, or dev dependencies.
	HasDependency(dependency string) bool

	// HasDevDependency should only return true if this dependency manager is tracking this dependency in its dev dependencies.
	HasDevDependency(dependency string) bool

	// Type returns the type of this DependencyManager.
	Type() DependencyManagerType
}

type DependencyManagerList []DependencyManager

// Main returns the first dependency manager in the list, under the assumption that that is the main / primary dependency manager used in the project.
func (list DependencyManagerList) Main() DependencyManager {
	if len(list) > 0 {
		return list[0]
	}
	return nil
}

func (list DependencyManagerList) Contains(target DependencyManager) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func (list DependencyManagerList) ContainsType(target DependencyManagerType) bool {
	for _, item := range list {
		if item.Type() == target {
			return true
		}
	}
	return false
}

func (list DependencyManagerList) ContainsAllTypes(targets ...DependencyManagerType) bool {
	for _, target := range targets {
		if !list.ContainsType(target) {
			return false
		}
	}
	return true
}
