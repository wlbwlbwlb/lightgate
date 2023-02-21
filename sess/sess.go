package sess

import (
	"fmt"
	"github.com/DarthPestilane/easytcp"
)

var storage map[int64]easytcp.Session

func OnLoginSuccess(uid int64, s easytcp.Session) {
	s.SetID(uid)
	storage[uid] = s
}

func OnSessionClose(s easytcp.Session) {
	uid := s.ID().(int64)
	delete(storage, uid)
	fmt.Printf("user %d logout\n", uid)
}
