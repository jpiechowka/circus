package manager

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jpiechowka/circus/task"
)

type Manager struct {
	PendingQueue  chan task.Task
	TaskDb        map[string][]task.Task
	EventDb       map[string][]task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	fmt.Println("Implement me: select an appropriate worker")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("Implement me: update tasks")
}

func (m *Manager) SendWork() {
	fmt.Println("Implement me: send work to workers")
}