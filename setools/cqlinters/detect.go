package cqlinters

import (
	"github.com/bvobart/mllint/api"
)

const (
	TypePylint api.CQLinterType = "Pylint"
	TypeMypy   api.CQLinterType = "Mypy"
	TypeBlack  api.CQLinterType = "Black"
	TypeISort  api.CQLinterType = "isort"
	TypeBandit api.CQLinterType = "Bandit"
)

var (
	pylint api.CQLinter = Pylint{}
	mypy   api.CQLinter
	black  api.CQLinter
	isort  api.CQLinter
	bandit api.CQLinter
)

func Detect(project api.Project) []api.CQLinter {
	// TODO: implement
	return []api.CQLinter{}
}
