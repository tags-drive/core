package web

import "sync"

func runPool(n int, data <-chan interface{}, worker func(data <-chan interface{})) {
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			worker(data)
			wg.Done()
		}()
	}

	wg.Wait()
}
