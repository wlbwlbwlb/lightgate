package tcpx

import (
	"fmt"
	"sync"

	"github.com/DarthPestilane/easytcp"
)

var sessions = &SessionManager{
	storage: map[int64]easytcp.Session{},
}

type SessionManager struct {
	lock    sync.Mutex
	storage map[int64]easytcp.Session
}

func (p *SessionManager) onLoginSuccess(uid int64, sess easytcp.Session) {
	p.lock.Lock()
	defer p.lock.Unlock()

	sess.SetID(uid)
	p.storage[uid] = sess
}

func (p *SessionManager) onSessionClose(sess easytcp.Session) {
	uid := sess.ID().(int64)
	delete(sessions.storage, uid)
	fmt.Printf("user %d logout\n", uid)
}
