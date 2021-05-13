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

func CollectTasks(tasks ...*RunnerTask) chan *RunnerTask {
	c := collector{
		total:  len(tasks),
		done:   make(chan struct{}),
		funnel: make(chan *RunnerTask),
	}

	for _, task := range tasks {
		go c.awaitResult(task)
	}

	go c.awaitDone()
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

func (c *collector) awaitDone() {
	nDone := 0
	for {
		<-c.done
		nDone++

		if nDone == c.total {
			close(c.funnel)
			close(c.done)
			return
		}
	}
}
