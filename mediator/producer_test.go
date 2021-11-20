package mediator

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSendMessageToTopic(t *testing.T) {
	Convey("Sends a message to Kafka ", t, func() {
		var topic = "gitoperator-in"
		var msg = "Hello World!"
		err := SendMessageToTopic(msg, topic)
		So(err, ShouldBeNil)
	})
}

