package middleware

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/wlbwlbwlb/log"
)

func Log() easytcp.MiddlewareFunc {
	return func(next easytcp.HandlerFunc) easytcp.HandlerFunc {
		return func(c easytcp.Context) {
			req := c.Request()
			if req != nil {
				log.Infof("Recv, id=%d", req.ID())
			}
			defer func() {
				resp := c.Response()
				if resp != nil {
					log.Infof("Send, id=%d", resp.ID())
				}
			}()
			next(c)
		}
	}
}
