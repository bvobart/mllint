package mllint_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bvobart/mllint/commands/mllint"
	"github.com/gosuri/uilive"
)

func TestProgressTest(t *testing.T) {
	t.Skip("Skipping this uilive.Newline() demo")

	writer := uilive.New()
	writer2 := writer.Newline()
	// start listening for updates and render
	writer.Start()

	for i := 0; i <= 100; i++ {
		fmt.Fprintf(writer, "Downloading.. (%d/%d) GB\n", i, 100)
		fmt.Fprintf(writer2, "Downloading.. (%d/%d) MB\n", i, 100)
		time.Sleep(time.Millisecond * 5)
	}

	fmt.Fprintln(writer, "Finished: Downloaded 100GB")
	writer.Stop() // flush and stop rendering
}

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
