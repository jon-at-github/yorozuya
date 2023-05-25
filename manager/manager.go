// The Manager: Keeps track of the state of the workers in the cluster.
package manager

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/jon-at-github/yorozuya/task"
)

// Manager: The Manager type represents a manager that keeps track of the state of the workers in the cluster.
type Manager struct {
	EventDb       map[string][]task.TaskEvent
	Pending       queue.Queue
	TaskDb        map[string][]task.Task
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
}
