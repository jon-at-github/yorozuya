// The Node: Represents a physical or virtual machine where the worker and tasks will run.
package node

// Node: The Node type represents a physical or virtual machine where the worker and tasks will run.
type Node struct {
	Cores           int
	Disk            int
	DiskAllocated   int
	Ip              string
	Memory          int
	MemoryAllocated int
	Name            string
	Role            string
	TaskCount       int
}
