package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jpiechowka/circus/manager"
	"github.com/jpiechowka/circus/node"
	"github.com/jpiechowka/circus/task"
	"github.com/jpiechowka/circus/worker"
)

func main() {
	t := task.Task{
		Id:     uuid.New(),
		Name:   "task-test",
		State:  task.Pending,
		Image:  "image-test",
		Memory: 1024,
		Disk:   1,
	}

	te := task.TaskEvent{
		Id:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      t,
	}

	fmt.Printf("task: %v\n", t)
	fmt.Printf("task event: %v\n", te)

	w := worker.Worker{
		Name:  "worker-test",
		Queue: make(chan task.Task),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	fmt.Printf("worker: %v\n", w)

	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		PendingQueue: make(chan task.Task),
		TaskDb:       make(map[string][]task.Task),
		EventDb:      make(map[string][]task.TaskEvent),
		Workers:      []string{w.Name},
	}

	fmt.Printf("manager: %v\n", m)

	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "node-test",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}

	fmt.Printf("node: %v\n", n)
}
