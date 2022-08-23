package god

import (
	"sync"

	"github.com/samber/do"
)

type Task struct {
	// ID of the Task.
	//
	// ID range: 1-256
	ID uint8
}

// taskRegistry is the task registry
// of task manager ([god.daemonService]).
type taskRegistry struct {
	// The registry [sync.Mutex].
	mutex sync.Mutex

	// The internal [god.Task] registry.
	reg [256]*Task

	// Records next index to register Task.
	//
	// Index range: 0-255
	next uint8
}

func (registry *taskRegistry) Acquire(i *do.Injector) {
	// Lock taskRegistry.
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	// Acquire Task ID.
	id := registry.next + 1
	for {
		// index = id - 1. Index range: 0-255
		if registry.reg[id-1] == nil {
			break
		}
		id++
	}
	do.ProvideValue(i, &Task{ID: id})
	registry.next++
}

func (registry *taskRegistry) Release(i *do.Injector) {
	// Lock taskRegistry.
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	task := do.MustInvoke[*Task](i)
	// index = id - 1. Index range: 0-255
	registry.reg[task.ID-1] = nil
}
