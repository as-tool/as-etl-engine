package reader

import (
	"context"

	"github.com/as-tool/as-etl/engine/common/config"
	"github.com/as-tool/as-etl/engine/core/plugin"
)

// Job: a unit of work
type Job interface {
	plugin.Job

	// Split the existing task into a number of tasks based on the Job, and pass them to each task mainly in the form of a configuration file.
	Split(ctx context.Context, number int) ([]*config.JSON, error)
}
