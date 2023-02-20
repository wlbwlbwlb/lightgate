package router

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/wl955/lightgate/user/userRouter"
)

func Init(serve *easytcp.Server) {

	serve.AddRoute(1, func(c easytcp.Context) {
		c.SetResponseMessage(easytcp.NewMessage(2, []byte("pong")))
	})

	userRouter.Router(serve)

	//todo

}
