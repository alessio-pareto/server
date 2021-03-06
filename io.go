package server

import (
	"io/ioutil"
	"sync"
)

type Mutex struct {
	waitGroup 	sync.WaitGroup
	c         	chan struct{}
}

func NewMutex() *Mutex {
	return &Mutex{
		waitGroup: 	sync.WaitGroup{},
		c:     		make(chan struct{}, 1),
	}
}

func (m *Mutex) CreateJobs(nJobs int) {
	m.waitGroup.Add(nJobs)
}

func (m *Mutex) Wait() {
	m.waitGroup.Wait()
}

func (m *Mutex) Done() {
	m.waitGroup.Done()
}

func (m *Mutex) SendSignal() {
	m.waitGroup.Add(1)
	m.c <- struct{}{}
	m.waitGroup.Wait()
}

func (m *Mutex) ListenForSignal() {
	<- m.c
}

func (srv *Server) ReadFileConcurrent(filePath string) ([]byte, error) {
	fm, ok := srv.fileMutexMap[filePath]
	if !ok {
		fm = new(sync.Mutex)
		srv.fileMutexMap[filePath] = fm
	}
	fm.Lock()

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fm.Unlock()
	return b, nil
}
