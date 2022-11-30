package threads

import "fmt"

var SafeMutex iSyncMutex = new(mutex)
var SafeGroup iSyncMutex = new(group)

type iSyncMutex interface {
	Run(func())
	SafeRun(func())
	Wait()
}

func GoSafe(call func()) {
	go RunSafe(call)
}


func RunSafe(call func()) {
	defer Recover()
	call()
}

func Recover(calls ...func()) {
	for _, call := range calls {
		call()
	}
	if err := recover(); err != nil {
		fmt.Println(`safe`, err)
	}
}
