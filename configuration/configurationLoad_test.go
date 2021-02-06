package configuration

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigurationDevProfile(t *testing.T) {

	Convey("Loading server configuration values ", t, func() {
		os.Setenv("PROFILE", "dev")
		Reload()
		env := os.Getenv("PROFILE")
		So(env, ShouldEqual, "dev")
		os.Remove("PROFILE")
		Reload()
	})
}

func TestConfiguration(t *testing.T) {
	Convey("Reading server configuration values ", t, func() {
		os.Setenv("PROFILE", "dev")
		Reload()
		conf := GlobalConfiguration

		So(conf, ShouldNotBeNil)
		So(conf.Kafka.Bootstrapserver, ShouldEqual, "kafka:9094")
	})

}
