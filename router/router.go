package router

import (
	"github.com/wlbwlbwlb/lightgate/user/userRouter"

	"github.com/DarthPestilane/easytcp"
)

func Init(serve *easytcp.Server) {

	serve.AddRoute(1, func(c easytcp.Context) {
		c.SetResponseMessage(easytcp.NewMessage(2, []byte("ok")))
	})

	userRouter.Router(serve)

	//todo

}
