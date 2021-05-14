package mllint_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bvobart/mllint/commands/mllint"
)

func TestProgressLive(t *testing.T) {
	progress := mllint.NewLiveRunnerProgress()
	progress.Start()

	numTasks := 10
	tasks := []*mllint.RunnerTask{}
	for i := 0; i < numTasks; i++ {
		task := &mllint.RunnerTask{Id: fmt.Sprint(i)}
		mllint.DisplayName(fmt.Sprint("Task ", i))(task)
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		progress.RunningTask(task)
		time.Sleep(time.Millisecond * 10)
	}
	for _, task := range tasks {
		progress.CompletedTask(task)
		time.Sleep(time.Millisecond * 10)
	}

	progress.AllTasksDone()
}
