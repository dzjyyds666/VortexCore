package vortexUtil

import (
	"testing"

	"github.com/dzjyyds666/opensource/common"
	"github.com/smartystreets/goconvey/convey"
)

type Config struct {
	Redis struct {
		Url string `toml:"url" json:"url"`
	}
	Server struct {
		Port int    `toml:"port" json:"port"`
		Host string `toml:"host" json:"host"`
	}
}

func TestLoadConfigFromToml(t *testing.T) {
	convey.Convey("load config from toml", t, func() {
		t.Log("Loading configuration from config.toml")
		var config Config
		err := LoadConfigFromToml("./config.toml", &config)
		convey.So(err, convey.ShouldBeNil)
		t.Log("Config:", common.ToStringWithoutError(config))
	})
}
