package sessions

import (
	"fmt"
	"sync"

	"github.com/DarthPestilane/easytcp"
)

var mutex sync.RWMutex

var uidSessMapper map[int64]easytcp.Session

var sidSessMapper map[string]easytcp.Session

var sidUidMapper map[string]int64

func OnLogin(userId int64, sess easytcp.Session) {
	fmt.Printf("user %d login\n", userId)

	mutex.Lock()
	defer mutex.Unlock()

	sessId := sess.ID().(string)
	sidSessMapper[sessId] = sess
	sidUidMapper[sessId] = userId

	uidSessMapper[userId] = sess
	//用户上线 todo
}

func OnLogout(sess easytcp.Session) {
	userId, ok := sidUidMapper[sess.ID().(string)]
	if !ok {
		return
	}
	fmt.Printf("user %d logout\n", userId)

	mutex.Lock()
	defer mutex.Unlock()

	sessId := sess.ID().(string)
	delete(sidSessMapper, sessId)
	delete(sidUidMapper, sessId)

	//写登出日志 todo

	sessNow, _ := uidSessMapper[userId]
	sessIdNow := sessNow.ID().(string)
	if sessId == sessIdNow {
		delete(uidSessMapper, userId)
		//用户离线 todo
	}
}

func OnKickout(userId int64) {
	sess, _ := uidSessMapper[userId]
	//通知old登出 todo
	sess.Close()
}

func OnLogoutNotify(userId int64, loc string) {
	//通知old其他服登出 todo
}

//fatal error: user 10002175 logout
//concurrent map iteration and map write
