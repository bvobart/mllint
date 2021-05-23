package mllint_test

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/api/mock_api"
	"github.com/bvobart/mllint/commands/mllint"
	"github.com/bvobart/mllint/commands/mllint/mock_mllint"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewRunner(t *testing.T) {
	r := mllint.NewMLLintRunner(nil)
	require.NotNil(t, r)
}

func TestCollectTasksEmptyList(t *testing.T) {
	r := mllint.NewMLLintRunner(nil)
	tasks := []*mllint.RunnerTask{}
	funnel := r.CollectTasks(tasks...)
	select {
	case task, open := <-funnel:
		require.False(t, open)
		require.Nil(t, task)
	default:
		t.Fatalf("expected a closed channel after running CollectTasks with empty list of tasks, but apparently it was open")
	}
}

//---------------------------------------------------------------------------------------

// create a bunch of tests with linters that will simply return api.NewReport() and nil
func createTests(numTests int) *map[string]testLinter {
	lintProject := func(_ api.Project) {
		// sleep for 100 ms to simulate doing some actual linting
		time.Sleep(time.Millisecond * 100)
	}

	tests := make(map[string]testLinter, numTests)
	for i := 0; i < numTests; i++ {
		id := fmt.Sprint(i)
		tests[id] = testLinter{id, nil, lintProject, api.NewReport(), nil}
	}
	return &tests
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
	runner := mllint.NewMLLintRunner(progress)
	runner.Start()

	startTime := time.Now()

	// schedule all the test linters
	ids := make([]string, 0, len(*tests))
	tasks := make([]*mllint.RunnerTask, 0, len(*tests))
	for id, test := range *tests {
		ids = append(ids, id)
		tasks = append(tasks, tester.createTestLinterTask(runner, test))
	}

	// await all tasks as they complete and collect their IDs
	completedIds := []string{}
	mllint.ForEachTask(runner.CollectTasks(tasks...), func(task *mllint.RunnerTask, result mllint.LinterResult) {
		tester.checkTestLinterTaskResult((*tests)[task.Id], task, result)
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

//---------------------------------------------------------------------------------------

func createNestedLinterTest(t testCtl, maxDepth int32) testLinter {
	var runner mllint.Runner
	test := testLinter{id: "root", report: api.NewReport(), err: nil}
	test.setRunner = func(r mllint.Runner) {
		runner = r
	}

	var depth int32 = 0
	test.lintProject = func(project api.Project) {
		if atomic.LoadInt32(&depth) >= maxDepth {
			time.Sleep(time.Millisecond * 10)
			return
		}

		time.Sleep(time.Millisecond * 10)
		atomic.AddInt32(&depth, 1)

		childTest := test
		childTest.id = fmt.Sprint("child-", depth)

		task := t.createRecursiveTestTask(runner, childTest)
		taskCompleted := false
		mllint.ForEachTask(runner.CollectTasks(task), func(completedTask *mllint.RunnerTask, result mllint.LinterResult) {
			require.Equal(t.t, task, completedTask)
			require.False(t.t, taskCompleted)
			taskCompleted = true
			t.checkTestLinterTaskResult(test, completedTask, result)
		})
		require.True(t.t, taskCompleted)
	}

	return test
}

func TestMLLintRunnerChildLinters(t *testing.T) {
	// create a recursive linter that runs itself 42 times before exiting.
	maxDepth := 42
	maxCompletionTime := time.Duration(maxDepth) * (12 * time.Millisecond) // 10+2 ms per task and they're essentially run sequentially in this test.

	ctrl := gomock.NewController(t)
	tester := testCtl{t, ctrl}
	test := createNestedLinterTest(tester, int32(maxDepth))

	// create and start the runner
	progress := mllint.NewBasicRunnerProgress()
	runner := mllint.NewMLLintRunner(progress)
	runner.Start()

	startTime := time.Now()
	task := tester.createRecursiveTestTask(runner, test)
	taskCompleted := false
	mllint.ForEachTask(runner.CollectTasks(task), func(completedTask *mllint.RunnerTask, result mllint.LinterResult) {
		require.Equal(t, task, completedTask)
		require.False(t, taskCompleted)
		taskCompleted = true
		tester.checkTestLinterTaskResult(test, completedTask, result)
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
	require.True(t, taskCompleted)
}

//---------------------------------------------------------------------------------------

type testCtl struct {
	t    *testing.T
	ctrl *gomock.Controller
}

type testLinter struct {
	id          string
	setRunner   func(r mllint.Runner)
	lintProject func(project api.Project)
	report      api.Report
	err         error
}

func (t testCtl) createTestLinterTask(runner mllint.Runner, test testLinter) *mllint.RunnerTask {
	project := api.Project{Dir: fmt.Sprint("TestDir", test.id)}

	linter := mock_api.NewMockLinter(t.ctrl)
	linter.EXPECT().Name().Times(1).Return(test.id)
	linter.EXPECT().LintProject(project).Times(1).DoAndReturn(func(project api.Project) (api.Report, error) {
		test.lintProject(project)
		return test.report, test.err
	})

	task := runner.RunLinter(test.id, linter, project, mllint.DisplayName(fmt.Sprint("TestLinter", test.id)))
	require.NotNil(t.t, task)
	require.Equal(t.t, test.id, task.Id)
	require.Equal(t.t, linter, task.Linter)
	require.Equal(t.t, project, task.Project)

	return task
}

func (t testCtl) createRecursiveTestTask(runner mllint.Runner, test testLinter) *mllint.RunnerTask {
	project := api.Project{Dir: fmt.Sprint("TestDir", test.id)}

	linter := mock_mllint.NewMockLinterWithRunner(t.ctrl)
	linter.EXPECT().Name().Times(1).Return(test.id)
	linter.EXPECT().LintProject(project).Times(1).DoAndReturn(func(project api.Project) (api.Report, error) {
		test.lintProject(project)
		return test.report, test.err
	})
	linter.EXPECT().SetRunner(gomock.Any()).Times(1).Do(func(r mllint.Runner) {
		test.setRunner(r)
	})

	task := runner.RunLinter(test.id, linter, project, mllint.DisplayName(fmt.Sprint("RecursiveLinter - ", test.id)))
	require.NotNil(t.t, task)
	require.Equal(t.t, test.id, task.Id)
	require.Equal(t.t, linter, task.Linter)
	require.Equal(t.t, project, task.Project)

	return task
}

func (t testCtl) checkTestLinterTaskResult(test testLinter, task *mllint.RunnerTask, result mllint.LinterResult) {
	require.Equal(t.t, result.Report, test.report)
	require.Equal(t.t, result.Err, test.err)
}
