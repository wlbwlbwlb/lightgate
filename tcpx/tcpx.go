package tcpx

import (
	"log"

	"github.com/wl955/lightgate/router"
	"github.com/wl955/lightgate/sess"

	"github.com/DarthPestilane/easytcp"
)

func Init(opts ...Option) (serve *easytcp.Server, e error) {
	custom := Options{}

	for _, opt := range opts {
		opt(&custom)
	}

	opt := easytcp.ServerOption{
		Codec: &easytcp.JsonCodec{},
	}

	serve = easytcp.NewServer(&opt)

	serve.OnSessionCreate = func(s easytcp.Session) {
		log.Printf("session created: %v\n", s.ID())
	}
	serve.OnSessionClose = func(s easytcp.Session) {
		log.Printf("session closed: %v\n", s.ID())
		if _, ok := s.ID().(int64); ok {
			sess.OnSessionClose(s)
		}
	}

	router.Init(serve)

	return
}
