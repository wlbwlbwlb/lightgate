package sess

import (
	"fmt"
	"sync"

	"github.com/DarthPestilane/easytcp"
)

var lock sync.Mutex

var storage map[int64]easytcp.Session

func OnLoginSuccess(uid int64, s easytcp.Session) {
	lock.Lock()
	defer lock.Unlock()
	s.SetID(uid)
	storage[uid] = s
}

func OnSessionClose(s easytcp.Session) {
	uid := s.ID().(int64)
	delete(storage, uid)
	fmt.Printf("user %d logout\n", uid)
}
