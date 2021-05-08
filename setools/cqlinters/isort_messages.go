package cqlinters

type ISortProblem struct {
	Path    string
	Message string
}

func (msg ISortProblem) String() string {
	return "`" + msg.Path + "` - " + msg.Message
}
