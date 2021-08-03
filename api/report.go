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

// OverallScore returns the weighted average of the scores of each rule, weighted with each rule's respective weight.
func (r Report) OverallScore() float64 {
	sumScores, sumWeights := 0.0, 0.0
	for rule, score := range r.Scores {
		sumScores += rule.Weight * score
		sumWeights += rule.Weight
	}

	if sumWeights == 0 {
		return 0
	}
	return sumScores / sumWeights
}

func NewReport() Report {
	return Report{
		Scores:  map[Rule]float64{},
		Details: map[Rule]string{},
	}
}

func MergeReports(finalReport Report, reports ...Report) Report {
	for _, report := range reports {
		for rule, score := range report.Scores {
			finalReport.Scores[rule] = score
		}
		for rule, details := range report.Details {
			finalReport.Details[rule] = details
		}
	}
	return finalReport
}
