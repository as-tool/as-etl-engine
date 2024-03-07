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
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/as-tool/as-etl-engine/common/config"

	"github.com/gorilla/handlers"
)

type enveronment struct {
	config *config.JSON
	engine *Engine
	err    error
	ctx    context.Context
	cancel context.CancelFunc
	server *http.Server
	addr   string
}

func NewEnveronment(filename string, addr string) (e *enveronment) {
	e = &enveronment{}
	var buf []byte
	buf, e.err = os.ReadFile(filename)
	if e.err != nil {
		return e
	}
	e.config, e.err = config.NewJSONFromBytes(buf)
	if e.err != nil {
		return e
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.addr = addr
	return e
}

func (e *enveronment) Build() error {
	return e.initEngine().initServer().startEngine().err
}

func (e *enveronment) initEngine() *enveronment {
	if e.err != nil {
		return e
	}
	e.engine = NewEngine(e.ctx, e.config)

	return e
}

func (e *enveronment) initServer() *enveronment {
	if e.err != nil {
		return e
	}
	if e.addr != "" {
		r := http.NewServeMux()
		recoverHandler := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))
		r.Handle("/metrics", NewMetricHandler(e.engine))
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		e.server = &http.Server{
			Addr:    e.addr,
			Handler: recoverHandler(r),
		}
		go func() {
			slog.Debug(fmt.Sprintf("listen begin: %v", e.addr))
			defer slog.Debug(fmt.Sprintf("listen end: %v", e.addr))
			if err := e.server.ListenAndServe(); err != nil {
				slog.Error("ListenAndServe fail. addr: %v err: %v", e.addr, err)
			}
			slog.Info(fmt.Sprintf("ListenAndServe success. addr: %v", e.addr))
		}()
	}

	return e
}

func (e *enveronment) startEngine() *enveronment {
	if e.err != nil {
		return e
	}
	go func() {
		statsTimer := time.NewTicker(5 * time.Second)
		defer statsTimer.Stop()
		exit := false
		for {
			select {
			case <-statsTimer.C:
			case <-e.ctx.Done():
				exit = true
			default:
			}
			if e.engine.Container != nil {
				fmt.Printf("%v\r", e.engine.Metrics().JSON())
			}

			if exit {
				return
			}
		}
	}()
	e.err = e.engine.Start()

	return e
}

func (e *enveronment) Close() {
	if e.server != nil {
		e.server.Shutdown(e.ctx)
	}

	if e.cancel != nil {
		e.cancel()
	}
}
