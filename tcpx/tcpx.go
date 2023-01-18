package tcpx

import (
	"log"

	"github.com/DarthPestilane/easytcp"
)

func Init(opts ...Option) (serve *easytcp.Server, e error) {
	custom := Options{}

	for _, opt := range opts {
		opt(&custom)
	}

	opt := easytcp.ServerOption{
		Packer: easytcp.NewDefaultPacker(),
		Codec:  &easytcp.JsonCodec{},
	}

	serve = easytcp.NewServer(&opt)

	serve.OnSessionCreate = func(sess easytcp.Session) {
		log.Printf("session created: %v\n", sess.ID())
	}
	serve.OnSessionClose = func(sess easytcp.Session) {
		log.Printf("session closed: %v\n", sess.ID())
		if _, ok := sess.ID().(int64); ok {
			sessions.onSessionClose(sess)
		}
	}

	serve.AddRoute(1, func(c easytcp.Context) {
		c.SetResponseMessage(easytcp.NewMessage(2, []byte("pong")))
	})

	addRoute(serve)

	return
}
