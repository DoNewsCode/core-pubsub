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
