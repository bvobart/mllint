package mllint

import "runtime"

// Watches the queue for new jobs, running or parking them as they come in / complete.
// Run in a new go-routine using `go r.queueWorker()`
func (r *MLLintRunner) queueWorker() {
	var next *RunnerTask
	parked := []*RunnerTask{}
	closed := false

	checkParked := func() {
		if r.nRunning >= int32(runtime.NumCPU()) && len(parked) > 0 {
			next, parked = parked[0], parked[1:]
			r.runTask(next)
		}
	}

	for {
		select {
		// when new task is scheduled...
		case task, open := <-r.queue:
			// if channel just closed and no tasks are running, signal that we're finished and exit
			if !open && r.nRunning == 0 {
				r.progress.AllTasksDone()
				close(r.closed)
				return
			}
			// if channel just closed, but there are still tasks running, signal to the next case that there will be no new tasks.
			if !open {
				closed = true
				checkParked()
				break
			}

			// if we're already running the maximum number of tasks, park it
			if r.nRunning >= int32(runtime.NumCPU()) {
				parked = append(parked, task)
				break
			}

			// otherwise just run the task
			r.runTask(task)

		// when a task starts awaiting results from tasks scheduled on a child linter
		case task := <-r.awaiting:
			r.nRunning--
			r.progress.TaskAwaiting(task)

			// if there are parked tasks, run one of them.
			checkParked()

		// when a task is done awaiting results and resumes its execution
		case task := <-r.resuming:
			r.nRunning++
			r.progress.TaskResuming(task)

		// when a task completes...
		case task := <-r.done:
			r.nRunning--
			r.progress.CompletedTask(task)

			// if the queue is closed and no other tasks are running, then signal that we're finished and exit
			if closed && r.nRunning == 0 {
				r.progress.AllTasksDone()
				close(r.closed)
				return
			}

			// else, if there are parked tasks, run one of them.
			checkParked()
		}
	}
}

// actually start running the task in a new go-routine
func (r *MLLintRunner) runTask(task *RunnerTask) {
	r.nRunning++
	r.progress.RunningTask(task)

	go func() {
		if l, ok := task.Linter.(WithRunner); ok {
			runner := childRunner{r, task}
			l.SetRunner(&runner)
		}

		report, err := task.Linter.LintProject(task.Project)
		task.Result <- LinterResult{Report: report, Err: err}
		r.done <- task
	}()
}
