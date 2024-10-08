package broadcastMq

import (
	"github.com/nsqio/go-nsq"
	"github.com/wlbwlbwlb/log"
	"github.com/wlbwlbwlb/mq"
)

func init() {
	mq.Sub("topic", "channel", nsq.HandlerFunc(handleMessage))
}

// HandleMessage implements the Handler interface.
func handleMessage(m *nsq.Message) (e error) {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return
	}

	// do whatever actual message processing is desired
	//err := processMessage(m.Body)
	log.Info(string(m.Body))

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return
}
