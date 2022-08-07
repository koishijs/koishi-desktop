package god

import (
	"github.com/samber/do"
	"sync"
)

type Task struct {
	// Id of the Task.
	//
	// Id range: 1-256
	Id uint8
}

// taskRegistry is the task registry
// of task manager ([god.Daemon]).
type taskRegistry struct {
	// The registry [sync.Mutex].
	mutex sync.Mutex

	// The internal [god.Task] registry.
	reg [256]*Task

	// Records next index to register Task.
	next uint8
}

func (registry *taskRegistry) Acquire(i *do.Injector) {
	// Lock taskRegistry.
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	// Acquire Task Id.
	id := registry.next
	for {
		// index = id - 1. Index range: 0-255
		if registry.reg[id-1] != nil {
			break
		}
		id++
	}
	do.ProvideValue(i, &Task{Id: id})
	registry.next++
}

func (registry *taskRegistry) Release(i *do.Injector) {
	// Lock taskRegistry.
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	task := do.MustInvoke[*Task](i)
	// index = id - 1. Index range: 0-255
	registry.reg[task.Id-1] = nil
}
