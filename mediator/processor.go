package mediator


import (
	"sync"
	utils "me/gitpoc/utils"
)

// https://golangbot.com/mutex/
func ProcessIncomingMessage(event utils.RecordEvent) {
	var w sync.WaitGroup
	var m sync.Mutex

	// TODO Logic: Gather requests in a microbatch by account and proceed

	go proceed(&w, &m)
    w.Wait()
}

func proceed(wg *sync.WaitGroup, m *sync.Mutex) {  
	m.Lock()

	// TODO Logic
	
    m.Unlock()
    wg.Done()   
}