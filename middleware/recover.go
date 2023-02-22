package middleware

import (
	"runtime/debug"

	"github.com/DarthPestilane/easytcp"
	"github.com/wl955/log"
)

func Recover() easytcp.MiddlewareFunc {
	return func(next easytcp.HandlerFunc) easytcp.HandlerFunc {
		return func(c easytcp.Context) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("PANIC, sessId=%v, e=%+v, stack=%s", c.Session().ID(), r, debug.Stack())
				}
			}()
			next(c)
		}
	}
}
