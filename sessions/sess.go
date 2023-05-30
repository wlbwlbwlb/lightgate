package sessions

import (
	"fmt"
	"sync"

	"github.com/DarthPestilane/easytcp"
)

var mutex sync.RWMutex

var uidSessMapper map[int64]easytcp.Session

var idSessMapper map[string]easytcp.Session

var sidUidMapper map[string]int64

func OnLogin(userId int64, sess easytcp.Session) {
	fmt.Printf("user %d login\n", userId)

	mutex.Lock()
	defer mutex.Unlock()

	sessId := sess.ID().(string)
	idSessMapper[sessId] = sess
	sidUidMapper[sessId] = userId

	uidSessMapper[userId] = sess
}

func OnLogout(sess easytcp.Session) {
	uid, ok := sidUidMapper[sess.ID().(string)]
	if !ok {
		return
	}
	fmt.Printf("user %d logout\n", uid)

	mutex.Lock()
	defer mutex.Unlock()

	sessId := sess.ID().(string)
	delete(idSessMapper, sessId)
	delete(sidUidMapper, sessId)

	sess2, _ := uidSessMapper[uid]
	sessId2 := sess2.ID().(string)

	if sessId == sessId2 {
		delete(uidSessMapper, uid)
	}
}

//fatal error: user 10002175 logout
//concurrent map iteration and map write
