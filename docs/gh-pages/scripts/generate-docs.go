package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/utils"
)

// This script will use mllint's source code to build Markdown files that contain documentation for each of mllint's rules and categories.
// These Markdown files will be placed in `../content/docs/categories` and `../content/docs/rules`
func main() {
	docspath := findDocsPath()
	if !utils.FolderExists(docspath) {
		panic(fmt.Errorf("Wtf, folder '%s' does not exist...", docspath))
	}

	ccf := categoryContentFactory{}
	ccf.GenerateDocs(docspath)

	rcf := ruleContentFactory{}
	rcf.GenerateDocs(docspath)
}

// Returns the equivalent of '../content/docs' relative to this script.
func findDocsPath() string {
	_, scriptPath, _, _ := runtime.Caller(0) //nolint:dogsled
	return path.Join(path.Dir(path.Dir(scriptPath)), "content", "docs")
}

//---------------------------------------------------------------------------------------

type categoryContentFactory struct{}

func (ccf categoryContentFactory) GenerateDocs(docspath string) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	for _, cat := range categories.All {
		output, err := ccf.createOutputFile(docspath, cat)
		checkErr(err)
		defer output.Close()

		_, err = output.WriteString(ccf.buildHeader(cat))
		checkErr(err)
		_, err = output.WriteString(ccf.buildContent(cat))
		checkErr(err)
	}
}

func (ccf categoryContentFactory) createOutputFile(docspath string, cat api.Category) (*os.File, error) {
	outputPath := path.Join(docspath, "categories", cat.Slug) + ".md"
	output, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file for category %s at %s: %w", cat.Name, outputPath, err)
	}
	return output, nil
}

func (ccf categoryContentFactory) buildHeader(cat api.Category) string {
	return fmt.Sprintf(`---
title: "Category — %s"
description: "- `+"`%s`"+`"
weight: 9
showtoc: false
---

`, cat.Name, cat.Slug)
}

func (ccf categoryContentFactory) buildContent(cat api.Category) string {
	return cat.Description
}

//---------------------------------------------------------------------------------------

type ruleContentFactory struct{}

func (rcf ruleContentFactory) GenerateDocs(docspath string) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	for cat, linter := range linters.ByCategory {
		for _, rule := range linter.Rules() {
			output, err := rcf.createOutputFile(docspath, rule)
			checkErr(err)
			defer output.Close()

			_, err = output.WriteString(rcf.buildHeader(rule, cat))
			checkErr(err)
			_, err = output.WriteString(rcf.buildContent(rule))
			checkErr(err)
		}
	}
}

func (rcf ruleContentFactory) createOutputFile(docspath string, rule *api.Rule) (*os.File, error) {
	outputPath := path.Join(docspath, "rules", rule.Slug) + ".md"

	// create dir for the rule if it doesn't exist yet
	if outputDir := path.Dir(outputPath); !utils.FolderExists(outputDir) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file for rule %s at %s: %w", rule.Name, outputPath, err)
	}
	return output, nil
}

func (rcf ruleContentFactory) buildHeader(rule *api.Rule, cat api.Category) string {
	return fmt.Sprintf(`---
title: "Rule — %s — %s"
description: |
             - `+"`%s`"+`
             - `+"`weight: %.1f`"+`
weight: 9
showtoc: false
---

`, cat.Name, rule.Name, rule.Slug, rule.Weight)
}

func (rcf ruleContentFactory) buildContent(rule *api.Rule) string {
	return rule.Details
}
