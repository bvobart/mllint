package api

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"gitlab.com/bvobart/mllint/config"
)

func NewIssue(linter string, rule string, severity Severity, msg string) Issue {
	return Issue{Linter: linter, Rule: rule, Severity: severity, Message: msg}
}

// Issue represents an issue with the project that a linter has reported.
type Issue struct {
	// The name of the linter that recognised this issue
	Linter string

	// The name of the linter rule that recognised this issue (not including the linter name).
	// Leave this field empty if the linter only enforces one rule, or if this is the main rule of this linter.
	Rule string

	// The (often multi-line) message to the user about the issue and what they can do to fix it.
	// Any '>' characters that this message contains, will be converted into a HiYellow coloured '>' character,
	// so use that for providing actionable recommendations.
	Message string

	// The severity of the message: [Error, Warning, Info]
	Severity Severity
}

func (issue Issue) FullRule() string {
	if issue.Rule == "" {
		return issue.Linter
	}
	return issue.Linter + "/" + issue.Rule
}

func (issue Issue) String() string {
	fullrule := color.Set(color.Faint).Sprint(issue.FullRule())
	color.Unset()
	msg := strings.ReplaceAll(issue.Message, ">", color.HiYellowString(">"))
	return fmt.Sprintf("%s  %s  %s", issue.Severity.String(), fullrule, msg)
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

type IssueList []Issue

func (list IssueList) FilterEnabled(conf config.RuleConfig) IssueList {
	enabled := IssueList{}

	disabledRules := map[string]interface{}{}
	for _, rule := range conf.Disabled {
		disabledRules[rule] = struct{}{}
	}

	for _, issue := range list {
		if _, isDisabled := disabledRules[issue.FullRule()]; !isDisabled {
			enabled = append(enabled, issue)
		}
	}

	return enabled
}
