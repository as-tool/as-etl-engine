// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"context"

	"github.com/as-tool/as-etl-engine/engine/common/config"
	"github.com/as-tool/as-etl-engine/engine/core/container"

	coreconst "github.com/as-tool/as-etl-engine/engine/common"
	"github.com/as-tool/as-etl-engine/engine/module/job"
	"github.com/as-tool/as-etl-engine/engine/module/taskgroup"

	"github.com/pingcap/errors"
)

// Model Mode
type Model string

// Container Work Mode
var (
	ModelJob       Model = "job"       // Work by Job
	ModelTaskGroup Model = "taskGroup" // Work by Task Group
)

// IsJob: whether to work by job
func (m Model) IsJob() bool {
	return m == ModelJob
}

// IsTaskGroup: whether to work by task group
func (m Model) IsTaskGroup() bool {
	return m == ModelTaskGroup
}

// Engine: execution engine
type Engine struct {
	container.Container
	ctx  context.Context
	conf *config.JSON
}

// NewEngine: create a new execution engine based on context ctx and JSON configuration conf
func NewEngine(ctx context.Context, conf *config.JSON) *Engine {
	return &Engine{
		ctx:  ctx,
		conf: conf,
	}
}

// Start: start
func (e *Engine) Start() (err error) {
	model := Model(e.conf.GetStringOrDefaullt(coreconst.DataxCoreContainerModel, string(ModelJob)))
	switch {
	case model.IsJob():
		e.Container, err = job.NewContainer(e.ctx, e.conf)
		if err != nil {
			return
		}
	case model.IsTaskGroup():
		e.Container, err = taskgroup.NewContainer(e.ctx, e.conf)
		if err != nil {
			return
		}
	default:
		return errors.Errorf("model is %v", model)
	}

	return e.Container.Start()
}
