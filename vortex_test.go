package vortex

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_Vortex(t *testing.T) {
	convey.Convey("Test_Vortex", t, func() {
		vortex := NewVortexCore(context.Background(),
			WithListenPort("8080"),
			WithDefaultLogger(),
			WithTransport(Transport.TCP),
			WithProtocol(http1, webSocket),
		)

		vortex.Start()
	})
}
