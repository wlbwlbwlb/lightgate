package tcpx

import (
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

func (p *SessionManager) onUserAuth(uid int64, sess easytcp.Session) {
	p.lock.Lock()
	defer p.lock.Unlock()

	sess.SetID(uid)
	p.storage[uid] = sess
}

func (p *SessionManager) onSessionClose(sess easytcp.Session) {
	delete(sessions.storage, sess.ID().(int64))
}
