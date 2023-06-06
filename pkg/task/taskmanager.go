package task

import (
	"fmt"
	"sync"
)

type TaskManager struct {
	TaskStatus map[string]string
	mu         sync.Mutex
}

const (
	RUNNING = "running"
	SUCCESS = "success"
	FAILED  = "failed"
	PENDING = "pending"
)

var instance *TaskManager
var once sync.Once

func GetTaskManager() *TaskManager {
	once.Do(func() {
		instance = &TaskManager{
			TaskStatus: make(map[string]string),
		}
	})
	return instance
}

func (tm *TaskManager) AddTask(taskID string) {
	tm.TaskStatus[taskID] = PENDING
}

func (tm *TaskManager) ChangeTaskStatus(taskID string, status string) {
	tm.TaskStatus[taskID] = status
}

func (tm *TaskManager) GetTaskStatus(taskID string) string {
	fmt.Println(tm.TaskStatus)

	if status, ok := tm.TaskStatus[taskID]; ok {
		return status
	}

	return "" // 返回默认的空字符串，表示任务ID不存在
}
