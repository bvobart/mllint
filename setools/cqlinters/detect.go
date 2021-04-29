package cqlinters

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
)

const (
	TypePylint api.CQLinterType = "pylint"
	TypeMypy   api.CQLinterType = "mypy"
	TypeBlack  api.CQLinterType = "black"
	TypeISort  api.CQLinterType = "isort"
	TypeBandit api.CQLinterType = "bandit"
)

var ByType = map[api.CQLinterType]api.CQLinter{
	TypePylint: Pylint{},
	// TypeMypy:   Mypy{},
	// TypeBlack:  Black{},
	// TypeISort:  ISort{},
	// TypeBandit: Bandit{},
}

func Detect(project api.Project) []api.CQLinter {
	res := []api.CQLinter{}
	for _, linter := range ByType {
		if linter.Detect(project) {
			res = append(res, linter)
		}
	}
	return res
}

func FromConfig(conf config.CodeQualityConfig) ([]api.CQLinter, error) {
	linters := []api.CQLinter{}
	notFound := []string{}

	for _, ltype := range conf.Linters {
		linter, ok := ByType[api.CQLinterType(strings.ToLower(ltype))]

		if ok {
			linters = append(linters, linter)
		} else {
			notFound = append(notFound, ltype)
		}
	}

	var err error
	if len(notFound) > 0 {
		err = fmt.Errorf("could not parse these code quality linters from mllint's config: %+v", notFound)
	}
	return linters, err
}
