package linters

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"

	"github.com/bvobart/mllint/linters/ci"
	"github.com/bvobart/mllint/linters/codequality"
	"github.com/bvobart/mllint/linters/dependencymgmt"
	"github.com/bvobart/mllint/linters/testing"
	"github.com/bvobart/mllint/linters/versioncontrol"
)

//---------------------------------------------------------------------------------------

// ByCategory contains a linter for each implemented category.
var ByCategory = map[api.Category]api.Linter{
	categories.VersionControl:        versioncontrol.NewLinter(),
	categories.DependencyMgmt:        dependencymgmt.NewLinter(),
	categories.ContinuousIntegration: ci.NewLinter(),
	categories.CodeQuality:           codequality.NewLinter(),
	categories.Testing:               testing.NewLinter(),
}

// Disabled contains all the linters for categories that have been disabled using Disable or DisableAll.
var Disabled = map[api.Category]api.Linter{}

//---------------------------------------------------------------------------------------

// Configures all the linters in linters.ByCategory with the given config.
func ConfigureAll(conf *config.Config) error {
	for cat, linter := range ByCategory {
		configurable, ok := linter.(api.Configurable)
		if ok {
			if err := configurable.Configure(conf); err != nil {
				return fmt.Errorf("failed to configure linter for category %s: %w", cat, err)
			}
		}
	}
	return nil
}

//---------------------------------------------------------------------------------------

// Disables the linters for all categories or rules referenced by the given slugs.
// Basically just a for-loop around linters.Disable()
// Returns the amount of linting rules disabled.
func DisableAll(slugs []string) int {
	n := 0
	for _, slug := range slugs {
		n += Disable(slug)
	}
	return n
}

// Disables the linter for a specific category or rule by the category's or rule's slug.
// In case of a Rule slug, this should include the category name, e.g. 'version-control/code/git-no-big-files',
// such that the accompanying linter can be found and the rule can be disabled on it using linters.DisableRule.
// Returns the amount of linting rules disabled.
func Disable(slug string) int {
	// if the slug exactly matches a category's slug, disable the entire category
	if cat, found := categories.BySlug[slug]; found {
		return DisableCategory(cat)
	}

	// else, find the linter that checks these rules and disable the rule on it.
	if linter, found := GetLinter(slug); found {
		return DisableRule(linter, slug)
	}

	return 0
}

// DisableCategory disables the linter for a specific category, if a linter for that category
// is known in linters.ByCategory. Given a category `cat`, this method will remove
// linters.ByCategory[cat] and add it to linters.Disabled[cat]
func DisableCategory(cat api.Category) int {
	// only disable linter for a category if it is actually implemented
	linter, found := ByCategory[cat]
	if !found {
		return 0
	}

	Disabled[cat] = linter
	delete(ByCategory, cat)
	return len(linter.Rules())
}

// DisableRule disables a rule on the linter by means of the rule's slug.
// The slug used here should be the same as one of the `linter.Rules()[i].Slug`
func DisableRule(linter api.Linter, slug string) int {
	nDisabled := 0
	for _, rule := range linter.Rules() {
		if strings.HasPrefix(rule.Slug, slug) {
			rule.Disable()
			nDisabled++
		}
	}
	return nDisabled
}

//---------------------------------------------------------------------------------------

// FindRules finds all rules that match (start with) the given slug.
// E.g. `version-control/data` will return all the rules corresponding to data version control.
func FindRules(slug string) []*api.Rule {
	cat, ok := GetCategory(slug)
	if !ok {
		return []*api.Rule{}
	}

	linter, ok := ByCategory[cat]
	if !ok {
		return []*api.Rule{}
	}

	if slug == cat.Slug || slug == cat.Slug+"/" {
		return linter.Rules()
	}

	foundRules := []*api.Rule{}
	for _, rule := range linter.Rules() {
		if strings.HasPrefix(rule.Slug, slug) {
			foundRules = append(foundRules, rule)
		}
	}
	return foundRules
}

//---------------------------------------------------------------------------------------

// GetCategory parses the category from the given slug and returns the matched category and whether it matched.
// The category is the part of the slug up to the first `/`
func GetCategory(slug string) (api.Category, bool) {
	if slashIndex := strings.Index(slug, "/"); slashIndex != -1 {
		slug = slug[:slashIndex]
	}

	cat, ok := categories.BySlug[slug]
	return cat, ok
}

// GetLinter gets linter for the category of the rule referenced by the given slug.
func GetLinter(slug string) (api.Linter, bool) {
	cat, ok := GetCategory(slug)
	if !ok {
		return nil, false
	}

	linter, ok := ByCategory[cat]
	return linter, ok
}

// GetRule returns the exact rule referenced by the given slug, or nil if it is not found.
func GetRule(slug string) *api.Rule {
	linter, ok := GetLinter(slug)
	if !ok {
		return nil
	}

	for _, rule := range linter.Rules() {
		if rule.Slug == slug {
			return rule
		}
	}

	return nil
}

//---------------------------------------------------------------------------------------
