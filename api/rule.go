package api

// Rule is a struct for defining what a rule looks like that `mllint` will check.
type Rule struct {
	// Slug should be a lowercased, dashed reference code, e.g. 'git-no-big-files'
	Slug string

	// Name should be a short sentence (<10 words) that concisely describes what this rule expects,
	// e.g. '.dvc folder should be comitted to Git' or 'Project should not use Git to track large files'
	Name string

	// Details should contain a longer, descriptive, Markdown-formatted text that explains the reasoning behind this rule,
	// as well as provide background info on the subject and pointers on how to fix violations of the rule.
	Details string

	// Weight determines the weight of this rule's score within its respective category.
	Weight float64 // TODO: figure out what to with this...

	// Whether this rule was explicitly disabled by the user.
	Disabled bool
}

func (r *Rule) Disable() {
	r.Disabled = true
}

func (r *Rule) Enable() {
	r.Disabled = false
}
