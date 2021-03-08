package api

import "github.com/fatih/color"

// Issue represents an issue with the project that a linter has reported.
type Issue struct {
	// The name of the linter that recognised this issue
	Linter string
	// The name of the linter rule that recognised this issue (not including the linter name).
	// Leave this field empty if the linter only enforces one rule, or if this is the main rule of this linter.
	Rule string
	// The message to the user about the issue and what they can do to fix it.
	Message string
	// The severity of the message: [Error, Warning, Info]
	Severity Severity
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

func NewIssue(linter string, rule string, severity Severity, msg string) Issue {
	return Issue{Linter: linter, Rule: rule, Severity: severity, Message: msg}
}
