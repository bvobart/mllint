package linters

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"

	"github.com/bvobart/mllint/linters/ci"
	"github.com/bvobart/mllint/linters/common"
	"github.com/bvobart/mllint/linters/dependencymgmt"
	"github.com/bvobart/mllint/linters/versioncontrol"
)

// ByCategory contains a linter for each implemented category.
var ByCategory = map[api.Category]api.Linter{
	categories.VersionControl:        versioncontrol.NewLinter(),
	categories.DependencyMgmt:        dependencymgmt.NewLinter(),
	categories.ContinuousIntegration: ci.NewLinter(),
}

// Disabled contains all the linters for categories that have been disabled using Disable or DisableAll.
var Disabled = map[api.Category]api.Linter{}

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

// Disables the linters for all categories or rules referenced by the given slugs.
// Basically just a for-loop around linters.Disable()
func DisableAll(slugs []string) {
	for _, slug := range slugs {
		Disable(slug)
	}
}

// Disables the linter for a specific category or rule by the category's or rule's slug.
// In case of a Rule slug, this should include the category name, e.g. 'version-control/code/git-no-big-files',
// such that the accompanying linter can be found and the rule can be disabled on it using linters.DisableRule.
func Disable(slug string) {
	// if the slug exactly matches a category's slug, disable the entire category
	if cat, found := categories.BySlug[slug]; found {
		DisableCategory(cat)
		return
	}

	// else, trim the referenced category, get the accompanying linter and disable the rule on that linter.
	if cat, found := GetCategory(slug); found {
		if linter, lfound := ByCategory[cat]; lfound {
			ruleSlug := strings.TrimPrefix(slug, cat.Slug+"/")
			DisableRule(linter, ruleSlug)
		}
		return
	}
}

// DisableCategory disables the linter for a specific category, if a linter for that category
// is known in linters.ByCategory. Given a category `cat`, this method will remove
// linters.ByCategory[cat] and add it to linters.Disabled[cat]
func DisableCategory(cat api.Category) {
	// only disable linter for a category if it is actually implemented
	linter, found := ByCategory[cat]
	if !found {
		return
	}

	Disabled[cat] = linter
	delete(ByCategory, cat)
}

// DisableRule disables a rule on the linter by means of the rule's slug.
// The slug used here should be the same as one of the `linter.Rules()[i].Slug`
func DisableRule(linter api.Linter, slug string) {
	for _, rule := range linter.Rules() {
		if rule.Slug == slug {
			if compLinter, ok := linter.(*common.CompositeLinter); ok {
				compLinter.DisableRule(rule)
				return
			}

			rule.Disable()
			return
		}
	}
}

func GetCategory(slug string) (api.Category, bool) {
	slashIndex := strings.Index(slug, "/")
	if slashIndex == -1 {
		return api.Category{}, false
	}

	cat, ok := categories.BySlug[slug[:slashIndex]]
	return cat, ok
}

func GetRule(slug string) (api.Rule, bool) {
	cat, ok := GetCategory(slug)
	if !ok {
		return api.Rule{}, false
	}

	linter, ok := ByCategory[cat]
	if !ok {
		return api.Rule{}, false
	}

	ruleSlug := strings.TrimPrefix(slug, cat.Slug+"/")
	for _, rule := range linter.Rules() {
		if rule.Slug == ruleSlug {
			return *rule, true
		}
	}

	return api.Rule{}, false
}
