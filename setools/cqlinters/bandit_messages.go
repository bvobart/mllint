package cqlinters

import "fmt"

type BanditMessage struct {
	TestID      string `yaml:"test_id"`
	TestName    string `yaml:"test_name"`
	Confidence  string `yaml:"issue_confidence"`
	Severity    string `yaml:"issue_severity"`
	Text        string `yaml:"issue_text"`
	MoreInfo    string `yaml:"more_info"`
	CodeSnippet string `yaml:"code"`
	Filename    string `yaml:"filename"`
	Line        int32  `yaml:"line_number"`
}

func (msg BanditMessage) String() string {
	return fmt.Sprint("`", msg.Filename, ":", msg.Line, "`", " - _(", msg.TestID, ", severity: ", msg.Severity, ", confidence: ", msg.Confidence, ")_ ", msg.Text, " See [here]("+msg.MoreInfo+") for more info")
}
