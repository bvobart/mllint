package api

// Report is the type of object returned by a Linter after linting a project.
type Report struct {
	// Scores maps each evaluated rule to a score
	Scores map[Rule]float64

	// Details contains any additional details to accompany a Rule's evaluation.
	// Typically, when a Linter detects that a project does not conform to a Rule,
	// it will want to provide some form of reasoning about it, pointers to which
	// parts of the project repo the Rule violation occcurs in, and what the user can
	// do to fix the issue.
	//
	// The mapped string may be formatted using Markdown.
	Details map[Rule]string
}

func NewReport() Report {
	return Report{
		Scores:  map[Rule]float64{},
		Details: map[Rule]string{},
	}
}
