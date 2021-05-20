package api

import "fmt"

type DependencyManagerType interface {
	fmt.Stringer

	// Detect whether a project is using this type of dependency manager and return the dependency manager instantiated for this project if it is detected,
	// or an error if it is not detected.
	//
	// Generally, this is done by detecting whether the manager's configuration file exists.
	Detect(project Project) (DependencyManager, error)
}

type DependencyManager interface {
	// Dependencies returns a list of the names of all Python dependencies that this manager is tracking.
	Dependencies() []string

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
	if target == nil {
		return false
	}

	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func (list DependencyManagerList) ContainsType(target DependencyManagerType) bool {
	if target == nil {
		return false
	}

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
