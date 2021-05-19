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
	r := mllint.NewRunner(nil)
	require.NotNil(t, r)
}

func TestCollectTasksEmptyList(t *testing.T) {
	tasks := []*mllint.RunnerTask{}
	funnel := mllint.CollectTasks(tasks...)
	select {
	case task, open := <-funnel:
		require.False(t, open)
		require.Nil(t, task)
	default:
		t.Fatalf("expected a closed channel after running CollectTasks with empty list of tasks, but apparently it was open")
	}
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
	maxCompletionTime := time.Duration(cpufactor) * (104 * time.Millisecond) // 100 ms per task divided over runtime.NumCPU() threads, 4ms max scheduling overhead
	fmt.Println("Running ", numTests, " test linters with the mllint Runner")
	fmt.Println("- Max allowed completion time:", maxCompletionTime)
	fmt.Println()

	tests := createTests(numTests)

	ctrl := gomock.NewController(t)
	tester := testCtl{t, ctrl}

	// create and start the runner
	progress := mllint.NewBasicRunnerProgress()
	runner := mllint.NewRunner(progress)
	runner.Start()

	startTime := time.Now()

	// schedule all the test linters
	ids := make([]string, 0, len(tests))
	tasks := make([]*mllint.RunnerTask, 0, len(tests))
	for id, test := range tests {
		ids = append(ids, id)
		tasks = append(tasks, tester.createTestLinterTask(runner, test))
	}

	// await all tasks as they complete and collect their IDs
	completedIds := []string{}
	mllint.ForEachTask(mllint.CollectTasks(tasks...), func(task *mllint.RunnerTask, result mllint.LinterResult) {
		tester.checkTestLinterTaskResult(tests[task.Id], task, result)
		completedIds = append(completedIds, task.Id)
	})

	endTime := time.Now()

	runner.Close() // close here because we want to print normally afterwards.
	closeTime := time.Now()

	fmt.Println("Time to Completion:   ", endTime.Sub(startTime), " max:", maxCompletionTime)
	fmt.Println("Time to runner Close: ", closeTime.Sub(endTime))
	fmt.Println("-----------------------------------")
	fmt.Println("Total:                ", time.Since(startTime))
	fmt.Println()

	require.WithinDuration(t, endTime, startTime, maxCompletionTime)
	require.ElementsMatch(t, ids, completedIds)
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
	linter.EXPECT().Name().Times(1).Return(test.id)
	linter.EXPECT().LintProject(project).Times(1).DoAndReturn(func(project api.Project) (api.Report, error) {
		// sleep for 100 ms to simulate doing some actual linting
		time.Sleep(time.Millisecond * 100)
		return test.report, test.err
	})

	task := runner.RunLinter(fmt.Sprint(test.id), linter, project, mllint.DisplayName(fmt.Sprint("TestLinter", test.id)))
	require.NotNil(t.t, task)
	require.Equal(t.t, fmt.Sprint(test.id), task.Id)
	require.Equal(t.t, linter, task.Linter)
	require.Equal(t.t, project, task.Project)

	return task
}

func (t testCtl) checkTestLinterTaskResult(test testLinter, task *mllint.RunnerTask, result mllint.LinterResult) {
	require.Equal(t.t, result.Report, test.report)
	require.Equal(t.t, result.Err, test.err)
}
