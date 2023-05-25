// The Task: Is the smallest unit of work in an Orchestration system
// and typically runs in a container.
package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// State: The State type represents the states a task  goes through,
// from Pending to Scheduled to Running to Completed or Failed.
type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

// Task: The Task type represents a task that is to be run in a container.
type Task struct {
	Disk          int
	ExposedPorts  nat.PortSet
	FinishTime    time.Time
	ID            uuid.UUID
	Image         string
	Memory        int
	Name          string
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	State         State
}

// TaskEvent: The TaskEvent type represents an event that is emitted
// when a task changes state.
type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}
