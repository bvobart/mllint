package template

import "github.com/bvobart/mllint/api"

var RuleSomething = api.Rule{
	Name:    "Project should do Something",
	Slug:    "do-something",
	Details: "TODO: fill this in with details about what the rule checks, what the reasoning behind it is, where people can get more information on how to implement what this rule checks.",
	Weight:  1,
}
