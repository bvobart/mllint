package api

import (
	"github.com/bvobart/mllint/config"
)

type Project struct {
	Dir        string
	ConfigType config.FileType
	Reports    map[Category]Report
}
