// The Worker: Represents a physical or virtual machine where the worker and tasks will run.
package worker

import (
	"yorozuya/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

// Worker: The Worker type represents a worker that is to be run in a container.
type Worker struct {
	Db        map[uuid.UUID]task.Task
	Name      string
	Queue     queue.Queue
	TaskCount int
}
