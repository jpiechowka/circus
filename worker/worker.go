package worker

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jpiechowka/circus/task"
)

type Worker struct {
	Name      string
	Queue     chan task.Task
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("Implement me: collect stats")
}

func (w *Worker) RunTask() {
	fmt.Println("Implement me: start or stop a task")
}

func (w *Worker) StartTask() {
	fmt.Println("Implement me: start a task")
}

func (w *Worker) StopTask() {
	fmt.Println("Implement me: stop a task")
}
