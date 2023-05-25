// The Scheduler: Selects candidate nodes, scores them, and picks the best one.
package scheduler

// Scheduler: The Scheduler type represents a scheduler that selects candidate nodes, scores them, and picks the best one.
type Scheduler interface {
	SelectCandidateNodes()
	Score()
	Pick()
}
