package api

import "github.com/fatih/color"

type Issue struct {
	LinterName string
	Message    string
	Severity   Severity
}

type Severity string

const (
	SeverityError   Severity = "Error  "
	SeverityWarning Severity = "Warning"
	SeverityInfo    Severity = "Info   "
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return color.RedString(string(SeverityError))
	case SeverityWarning:
		return color.HiYellowString(string(SeverityWarning))
	case SeverityInfo:
		return color.BlueString(string(SeverityInfo))
	default:
		return string(s)
	}
}

func NewIssue(linterName string, severity Severity, msg string) Issue {
	return Issue{LinterName: linterName, Severity: severity, Message: msg}
}
