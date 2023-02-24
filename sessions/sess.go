package sessions

import (
	"fmt"
	"sync"

	"github.com/DarthPestilane/easytcp"
)

var lock sync.Mutex

var storage map[int64]easytcp.Session

func OnLogin(userId int64, sess easytcp.Session) {
	fmt.Printf("user %d login\n", userId)

	lock.Lock()
	defer lock.Unlock()

	sess.SetID(userId)
	storage[userId] = sess
}

func OnLogout(sess easytcp.Session) {
	uid := sess.ID().(int64)
	fmt.Printf("user %d logout\n", uid)

	lock.Lock()
	defer lock.Unlock()

	delete(storage, uid)
}
