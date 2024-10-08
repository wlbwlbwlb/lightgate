package mytcp

import (
	"log"

	"github.com/wlbwlbwlb/lightgate/middleware"
	"github.com/wlbwlbwlb/lightgate/router"
	"github.com/wlbwlbwlb/lightgate/sessions"

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

	serve.OnSessionCreate = func(sess easytcp.Session) {
		log.Printf("session %v created\n", sess.ID())
	}
	serve.OnSessionClose = func(sess easytcp.Session) {
		log.Printf("session %v closed\n", sess.ID())
		sessions.OnLogout(sess)
	}

	serve.Use(middleware.Log(), middleware.Recover())

	router.Init(serve)

	return
}
