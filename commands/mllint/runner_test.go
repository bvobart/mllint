package mllint_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/api/mock_api"
	"github.com/bvobart/mllint/commands/mllint"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewRunner(t *testing.T) {
	r := mllint.NewRunner()
	require.NotNil(t, r)
}

// create a bunch of tests with linters that will simply return api.NewReport() and nil
func createTests(numTests int) map[string]testLinter {
	tests := make(map[string]testLinter, numTests)
	for i := 0; i < numTests; i++ {
		id := fmt.Sprint(i)
		tests[id] = testLinter{id, api.NewReport(), nil}
	}
	return tests
}

func TestMLLintRunner(t *testing.T) {
	// create amount of tests equal to `cpufactor` times the amount of available CPU threads
	cpufactor := 10
	numTests := cpufactor * runtime.NumCPU()
	maxScheduleTime := time.Duration(numTests) * (50 * time.Microsecond)     // scheduling a linter job should be really quick
	maxCompletionTime := time.Duration(cpufactor) * (102 * time.Millisecond) // 100 ms per task divided over runtime.NumCPU() threads, 2ms max scheduling overhead
	fmt.Println("Running ", numTests, " test linters with the mllint Runner")
	fmt.Println("- Max allowed schedule time:  ", maxScheduleTime)
	fmt.Println("- Max allowed completion time:", maxCompletionTime)
	fmt.Println()

	tests := createTests(numTests)

	ctrl := gomock.NewController(t)
	tester := testCtl{t, ctrl}
	startTime := time.Now()

	// create and start the runner
	runner := mllint.NewRunner()
	runner.Start()
	defer runner.Close()

	// schedule all the test linters
	ids := make([]string, 0, len(tests))
	tasks := make([]*mllint.RunnerTask, 0, len(tests))
	for id, test := range tests {
		ids = append(ids, id)
		tasks = append(tasks, tester.createTestLinterTask(runner, test))
	}

	require.WithinDuration(t, startTime, time.Now(), maxScheduleTime)

	completedIds := []string{}
	tasksCompleted := mllint.CollectTasks(tasks...)
	for {
		task, open := <-tasksCompleted
		if !open {
			break
		}

		tester.awaitTestLinterTask(tests[task.Id], task)
		completedIds = append(completedIds, task.Id)
	}

	require.ElementsMatch(t, ids, completedIds)

	fmt.Println("------", time.Since(startTime))
	require.WithinDuration(t, startTime, time.Now(), maxCompletionTime)
}

type testCtl struct {
	t    *testing.T
	ctrl *gomock.Controller
}

type testLinter struct {
	id     string
	report api.Report
	err    error
}

func (t testCtl) createTestLinterTask(runner *mllint.Runner, test testLinter) *mllint.RunnerTask {
	project := api.Project{Dir: fmt.Sprint("TestDir", test.id)}

	linter := mock_api.NewMockLinter(t.ctrl)
	linter.EXPECT().Name().AnyTimes().Return(fmt.Sprint("TestLinter", test.id))
	linter.EXPECT().LintProject(project).Times(1).DoAndReturn(func(project api.Project) (api.Report, error) {
		// sleep for 100 ms to simulate doing some actual linting
		time.Sleep(time.Millisecond * 100)
		return test.report, test.err
	})

	task := runner.RunLinter(fmt.Sprint(test.id), linter, project)
	require.NotNil(t.t, task)
	require.Equal(t.t, fmt.Sprint(test.id), task.Id)
	require.Equal(t.t, linter, task.Linter)
	require.Equal(t.t, project, task.Project)

	return task
}

func (t testCtl) awaitTestLinterTask(test testLinter, task *mllint.RunnerTask) {
	result := <-task.Result
	require.Equal(t.t, result.Report, test.report)
	require.Equal(t.t, result.Err, test.err)
}
