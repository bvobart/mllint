package cqlinters

import "github.com/bvobart/mllint/api"

const (
	TypePylint api.CQLinterType = "pylint"
	TypeMypy   api.CQLinterType = "mypy"
	TypeBlack  api.CQLinterType = "black"
	TypeISort  api.CQLinterType = "isort"
	TypeBandit api.CQLinterType = "bandit"
)

var ByType = map[api.CQLinterType]api.CQLinter{
	TypePylint: Pylint{},
	TypeMypy:   Mypy{},
	TypeBlack:  Black{},
	TypeISort:  ISort{},
	// TypeBandit: Bandit{},
}
