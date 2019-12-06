package download

import (
	"sync"

	"github.com/schollz/progressbar/v2"
)

type progWaitGroup struct {
	pb *progressbar.ProgressBar
	wg sync.WaitGroup
}

func newProgWaitGroup(n int) *progWaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(n)
	return &progWaitGroup{
		pb: progressbar.New(n),
		wg: wg,
	}
}

func (pwg *progWaitGroup) Done() {
	pwg.wg.Done()
	pwg.pb.Add(1)
}

func (pwg *progWaitGroup) Wait() {
	pwg.wg.Wait()
	pwg.pb.Finish()
}
