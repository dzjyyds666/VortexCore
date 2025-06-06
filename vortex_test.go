package vortex

import (
	"context"
	vortexUtil "gothub.com/dzjyyds666/VortexCore/utils"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_Vortex(t *testing.T) {
	convey.Convey("Test_Vortex", t, func() {
		vortex := NewVortexCore(context.Background(),
			WithListenPort("18080"),
			WithDefaultLogger(),
			WithTransport(Transport.TCP),
			//WithHideBanner(true),
			WithProtocol(vortexUtil.Http1, vortexUtil.WebSocket),
		)
		vortex.Start()
	})
}
