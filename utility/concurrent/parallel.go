package concurrent

import "sync"

func Parallel(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

func IterateInParallel(items []interface{}, iterator func(order interface{})) {
	if iterator == nil {
		return
	}

	var failure interface{}

	concurrency := 10
	segmentSize := len(items) / concurrency
	if len(items)%concurrency > 0 {
		segmentSize++
	}

	synchronizer := &sync.WaitGroup{}
	synchronizer.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer synchronizer.Done()

			defer func() {
				if reason := recover(); reason != nil {
					failure = reason
				}
			}()

			start := i * segmentSize
			end := start + segmentSize
			if end > len(items) {
				end = len(items)
			}

			for _, order := range items[start:end] {
				iterator(order)
			}
		}(i)
	}

	synchronizer.Wait()

	if failure != nil {
		panic(failure)
	}
}
