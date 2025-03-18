package parallelx

import (
	"sync"

	"github.com/qf0129/gox/pkg/logx"
)

type ParallelFunc func() error

func RunParallel(handlers []ParallelFunc, concurrency ...int) []error {
	errs := make([]error, len(handlers))
	waitGroup := &sync.WaitGroup{}
	if len(concurrency) == 0 {
		concurrency = append(concurrency, len(handlers))
	}
	limitChan := make(chan struct{}, concurrency[0])
	for i := 0; i < len(handlers); i++ {
		waitGroup.Add(1)
		go func(index int) {
			limitChan <- struct{}{}
			defer func() {
				waitGroup.Done()
				<-limitChan
			}()
			if err := handlers[index](); err != nil {
				logx.Error("RunParallelErr", "index", index, "err", err)
				errs[index] = err
			}

		}(i)
	}
	waitGroup.Wait()
	return errs
}
