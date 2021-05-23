package mllint

func ForEachTask(tasks chan *RunnerTask, f func(task *RunnerTask, result LinterResult)) {
	for {
		task, open := <-tasks
		if !open {
			return
		}

		f(task, <-task.Result)
	}
}

func collectTasks(onComplete func(), tasks ...*RunnerTask) chan *RunnerTask {
	if len(tasks) == 0 {
		funnel := make(chan *RunnerTask)
		close(funnel)
		return funnel
	}

	c := collector{
		total:  len(tasks),
		done:   make(chan struct{}),
		funnel: make(chan *RunnerTask),
	}

	for _, task := range tasks {
		go c.awaitResult(task)
	}

	go c.awaitDone(onComplete)
	return c.funnel
}

type collector struct {
	total  int
	done   chan struct{}
	funnel chan *RunnerTask
}

func (c *collector) awaitResult(task *RunnerTask) {
	result := <-task.Result
	task.Result <- result
	c.funnel <- task
	c.done <- struct{}{}
}

func (c *collector) awaitDone(onComplete func()) {
	nDone := 0
	for {
		<-c.done
		nDone++

		if nDone == c.total {
			onComplete()
			close(c.funnel)
			close(c.done)
			return
		}
	}
}
