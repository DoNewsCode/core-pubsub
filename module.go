// Package pubsub is a pubsub module for package Core build upon
// watermill(https://github.com/ThreeDotsLabs/watermill). It allows other modules
// to register pub sub handlers easily via an interface. See example for usage.
package pubsub

import (
	"context"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/log"
	"github.com/oklog/run"
)

// Provider is an interface that all modules supporting watermill pubsub should
// implement. The router is passed into the ProvidePubSub(), and downstream
// modules should register their handlers to the router.
type Provider interface {
	ProvidePubSub(*message.Router)
}

// Module is the base struct for pubsub module.
type Module struct {
	container contract.Container
	router    *message.Router
	logger    Logger
}

// ModuleIn is the parameters for New.
type ModuleIn struct {
	di.In

	Container contract.Container
	Router    *message.Router `optional:"true"`
	Logger    log.Logger
}

// New constructs a pubsub module.
func New(in ModuleIn) (module Module, err error) {
	logger := NewLogger(in.Logger)
	if in.Router == nil {
		if in.Router, err = message.NewRouter(message.RouterConfig{}, logger); err != nil {
			return Module{}, err
		}
	}
	return Module{container: in.Container, router: in.Router, logger: logger}, nil
}

// ProvideRunGroup scans modules from container and collect all pubsubs. It
// starts the pub/sub using the runner mechanism.
func (m Module) ProvideRunGroup(group *run.Group) {
	_ = m.container.Modules().Filter(func(p Provider) error {
		p.ProvidePubSub(m.router)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	group.Add(func() error {
		return m.router.Run(ctx)
	}, func(err error) {
		cancel()
		_ = m.router.Close()
	})
}
