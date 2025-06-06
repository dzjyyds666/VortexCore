package vortexUtil

import (
	"github.com/dzjyyds666/opensource/common"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadConfigFromToml(t *testing.T) {
	convey.Convey("load config from toml", t, func() {
		toml, err := LoadConfigFromToml("./config.toml")
		convey.So(err, convey.ShouldBeNil)
		t.Log(common.ToStringWithoutError(toml))
	})
}
