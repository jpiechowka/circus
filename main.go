package main

import (
	"log"
	"time"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/jpiechowka/circus/manager"
	"github.com/jpiechowka/circus/node"
	"github.com/jpiechowka/circus/task"
	"github.com/jpiechowka/circus/worker"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	t := task.Task{
		ID:     uuid.New(),
		Name:   "task-test",
		State:  task.Pending,
		Image:  "image-test",
		Memory: 1024,
		Disk:   1,
	}

	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      t,
	}

	log.Printf("task: %v\n", t)
	log.Printf("task event: %v\n", te)

	w := worker.Worker{
		Name:  "worker-test",
		Queue: make(chan task.Task),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	log.Printf("worker: %v\n", w)

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

	log.Printf("manager: %v\n", m)

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

	log.Printf("node: %v\n", n)
	log.Printf("Create test contaner\n")
	dockerTask, result := createContainer()
	if result.Error != nil {
		log.Fatalf("%v\n", result.Error)
	}

	time.Sleep(time.Second * 13)
	log.Printf("Stopping container %s\n", result.ContainerID)
	result = stopContainer(dockerTask, result.ContainerID)
	if result.Error != nil {
		log.Fatalf("%v\n", result.Error)
	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	conf := task.Config{
		Name:  "test-container",
		Image: "postgres:latest",
		Env: []string{
			"POSTGRES_USER=circus",
			"POSTGRES_PASSWORD=p4ssw0rd",
		},
	}

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv)
	docker := task.Docker{
		Client: dockerClient,
		Config: conf,
	}

	result := docker.Run()
	if result.Error != nil {
		log.Printf("%v\n", result.Error)
		return nil, &task.DockerResult{Error: result.Error}
	}

	log.Printf("Container %s is running with config %v\n", result.ContainerID, conf)

	return &docker, &result
}

func stopContainer(docker *task.Docker, containerID string) *task.DockerResult {
	result := docker.Stop(containerID)
	if result.Error != nil {
		log.Printf("%v\n", result.Error)
		return &task.DockerResult{Error: result.Error}
	}

	log.Printf("Container %s has been stopped and removed\n", result.ContainerID)
	return &result
}
