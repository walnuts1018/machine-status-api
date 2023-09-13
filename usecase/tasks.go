package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/walnuts1018/machine-status-api/infra/timeAlternative"
)

type TaskStatus string

const (
	Pending TaskStatus = "pending"
	Running TaskStatus = "running"
	Success TaskStatus = "success"
	Failure TaskStatus = "failure"
	Unknown TaskStatus = "unknown"
)

type tasks struct {
	ID           string
	Func         func() error
	Status       TaskStatus
	RegisteredAt *time.Time
	StartedAt    *time.Time
	FinishedAt   *time.Time
}

var taskQueue []*tasks
var mutex = sync.Mutex{}

func NewTasks(f func() error) *tasks {
	t := timeAlternative.Now()
	newtask := &tasks{
		ID:           uuid.Must(uuid.NewV4()).String(),
		Func:         f,
		Status:       Pending,
		RegisteredAt: &t,
		StartedAt:    nil,
		FinishedAt:   nil,
	}
	mutex.Lock()
	taskQueue = append(taskQueue, newtask)
	mutex.Unlock()
	return newtask
}

func Run(ctx context.Context) {
	for {
		lastIndex := 0
		select {
		case <-ctx.Done():
			return
		default:
			for i := lastIndex; i < len(taskQueue); i++ {
				if taskQueue[i].Status == Pending {
					doTask(taskQueue[i])
					lastIndex = i + 1
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func doTask(task *tasks) {
	task.Status = Running
	ts := timeAlternative.Now()
	task.StartedAt = &ts
	err := task.Func()
	tf := timeAlternative.Now()
	if err != nil {
		task.Status = Failure
	} else {
		task.Status = Success
	}
	task.FinishedAt = &tf
}

func FindTaskByID(id string) *tasks {
	for _, task := range taskQueue {
		if task.ID == id {
			return task
		}
	}
	return nil
}
