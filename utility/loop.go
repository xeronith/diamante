package utility

import "sync"

func Repeat(times int, iterator func(int)) {
	if iterator == nil {
		return
	}

	for i := 0; i < times; i++ {
		iterator(i)
	}
}

func RepeatParallel(times int, iterator func(int)) {
	if iterator == nil {
		return
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(times)

	for i := 0; i < times; i++ {
		go func(i int) {
			defer waitGroup.Done()
			iterator(i)
		}(i)
	}

	waitGroup.Wait()
}
