package mllint

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

type RunnerProgress struct {
	*uilive.Writer
}

func NewRunnerProgress() *RunnerProgress {
	writer := uilive.New()
	writer.RefreshInterval = time.Hour
	return &RunnerProgress{writer}
}

func (p *RunnerProgress) Done(task *RunnerTask) {
	color.New(color.Bold, color.FgGreen).Fprintf(p.Bypass(), "✔️ Done - %s\n", task.Linter.Name())
}

func (p *RunnerProgress) Update(running *sync.Map, nRunning int32) {
	if nRunning == 0 {
		p.updateDone()
	} else {
		p.updateRunning(running, nRunning)
	}
}

func (p *RunnerProgress) updateRunning(running *sync.Map, nRunning int32) {
	fmt.Fprintln(p)
	color.New(color.Bold).Fprintf(p.Newline(), "Linters remaining - %d\n", nRunning)

	running.Range(func(key, value interface{}) bool {
		task := value.(*RunnerTask)
		writer := p.Newline()
		color.New(color.FgYellow).Fprintf(writer, "⏳ Running - %s\n", task.Linter.Name())
		return true
	})

	p.Flush()
}

func (p *RunnerProgress) updateDone() {
	boldGreen := color.New(color.Bold, color.FgGreen)
	boldGreen.Fprintln(p, "✔️ All done!")
	fmt.Fprint(p, "\n---\n\n")
	p.Flush()
}
