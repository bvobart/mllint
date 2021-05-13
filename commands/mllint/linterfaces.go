package mllint

import "github.com/bvobart/mllint/api"

type WithRunner interface {
	SetRunner(runner *Runner)
}

type LinterWithRunner interface {
	api.Linter
	WithRunner
}

type ConfigurableLinterWithRunner interface {
	api.ConfigurableLinter
	WithRunner
}
