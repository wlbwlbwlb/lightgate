package sessions

import (
	"fmt"
	"sync"

	"github.com/DarthPestilane/easytcp"
)

var mutex sync.RWMutex

var storage map[int64]easytcp.Session

func OnLogin(userId int64, sess easytcp.Session) {
	fmt.Printf("user %d login\n", userId)

	mutex.Lock()
	defer mutex.Unlock()

	sess.SetID(userId)
	storage[userId] = sess
}

func OnLogout(sess easytcp.Session) {
	uid := sess.ID().(int64)
	fmt.Printf("user %d logout\n", uid)

	mutex.Lock()
	defer mutex.Unlock()

	delete(storage, uid)
}
