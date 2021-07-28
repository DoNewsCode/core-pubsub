
<div align="center">
  <h1>core-pubsub</h1>
  <p>
    <strong>A pubsub module for package <a href="github.com/DoNewsCode/core">Core</a> built upon
<a href="https://github.com/ThreeDotsLabs/watermill">Watermill</a>.</strong>
  </p>
  <p>

[![Build](https://github.com/DoNewsCode/core-queue/actions/workflows/go.yml/badge.svg)](https://github.com/DoNewsCode/core-queue/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/DoNewsCode/core-queue.svg)](https://pkg.go.dev/github.com/DoNewsCode/core-queue)
[![codecov](https://codecov.io/gh/DoNewsCode/core-queue/branch/master/graph/badge.svg)](https://codecov.io/gh/DoNewsCode/core-queue)
[![Go Report Card](https://goreportcard.com/badge/DoNewsCode/core-queue)](https://goreportcard.com/report/DoNewsCode/core-queue)
[![Sourcegraph](https://sourcegraph.com/github.com/DoNewsCode/core-queue/-/badge.svg)](https://sourcegraph.com/github.com/DoNewsCode/core-queue?badge)

 </p>
</div>

## Example
```go
package pubsub_test

import (
	"context"
	"fmt"
	"time"

	"github.com/DoNewsCode/core"
	pubsub "github.com/DoNewsCode/core-pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type MyModule struct {
	goch   *gochannel.GoChannel
	cancel func()
}

func (m *MyModule) ProvidePubSub(router *message.Router) {
	router.AddNoPublisherHandler(
		"example",
		"example-in",
		m.goch,
		func(msg *message.Message) error {
			fmt.Println(string(msg.Payload))
			m.cancel()
			return nil
		},
	)
}

func Example() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	goch := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	module := MyModule{goch, cancel}

	c := core.Default(core.WithInline("log.level", "none"))
	c.AddModuleFunc(pubsub.New)
	c.AddModule(&module)
	go func() {
		time.Sleep(time.Second)
		goch.Publish("example-in", message.NewMessage(watermill.NewUUID(), message.Payload("foo")))
	}()
	c.Serve(ctx)

	// Output:
	// foo
}
```
