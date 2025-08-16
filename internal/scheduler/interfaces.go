package scheduler

// Stoppable defines the interface for components that can be gracefully stopped
type Stoppable interface {
	Stop()
}
