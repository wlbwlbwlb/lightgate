package middleware

import (
	"runtime/debug"

	"github.com/DarthPestilane/easytcp"
	"github.com/wlbwlbwlb/log"
)

func Recover() easytcp.MiddlewareFunc {
	return func(next easytcp.HandlerFunc) easytcp.HandlerFunc {
		return func(c easytcp.Context) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("PANIC, e=%+v, stack=%s", r, debug.Stack())
				}
			}()
			next(c)
		}
	}
}
