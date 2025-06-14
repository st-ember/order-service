package worker

import "sync"

func StartWorkerPool(workerFunc func(), workerCount int) {
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerFunc()
		}()
	}
	wg.Wait()
}
