package writer

import (
	"context"

	"github.com/as-tool/as-etl-engine/core/plugin"
)

// Task: a task related to writing operations
type Task interface {
	plugin.Task

	// Start reading records from the receiver and write them
	StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
	// SupportFailOver: whether fault tolerance is supported, i.e., whether to retry after a failed write
	SupportFailOver() bool
}

// BaseTask: a fundamental task class that assists and simplifies the implementation of writing task interfaces
type BaseTask struct {
	*plugin.BaseTask
}

// NewBaseTask: a function or method to create a new instance of BaseTask
func NewBaseTask() *BaseTask {
	return &BaseTask{
		BaseTask: plugin.NewBaseTask(),
	}
}

// SupportFailOver: whether fault tolerance is supported, i.e., whether to retry after a failed write
func (b *BaseTask) SupportFailOver() bool {
	return false
}
