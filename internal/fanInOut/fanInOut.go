package fanInOut

import (
	"context"
	"sync"
)

func FanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	fanIn := make(chan T, len(channels))
	var wg sync.WaitGroup
	transfer := func(c <-chan T) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// skip intentionally
			}
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case fanIn <- v:
				}
			}
		}
	}

	wg.Add(len(channels))

	for _, channel := range channels {
		go transfer(channel)
	}

	go func() {
		wg.Wait()
		close(fanIn)
	}()

	return fanIn
}
