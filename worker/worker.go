package worker

import (
	"log"

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
	log.Println("Implement me: collect stats")
}

func (w *Worker) RunTask() {
	log.Println("Implement me: start or stop a task")
}

func (w *Worker) StartTask() {
	log.Println("Implement me: start a task")
}

func (w *Worker) StopTask() {
	log.Println("Implement me: stop a task")
}
