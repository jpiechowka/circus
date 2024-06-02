package manager

import (
	"log"

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
	log.Println("Implement me: select an appropriate worker")
}

func (m *Manager) UpdateTasks() {
	log.Println("Implement me: update tasks")
}

func (m *Manager) SendWork() {
	log.Println("Implement me: send work to workers")
}
