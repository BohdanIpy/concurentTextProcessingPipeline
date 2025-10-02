package fanInOut

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func channelCreator(ctx context.Context, numOfElements int) <-chan int {
	channel := make(chan int)
	go func() {
		defer close(channel)
		for i := 0; i < numOfElements; i++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
			select {
			case <-ctx.Done():
				return
			case channel <- i:
			}
		}
	}()
	return channel
}

func TestFanIn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const numberOfChannels = 5
	const numberOfElements = 30
	channels := make([]<-chan int, numberOfChannels)
	for i := 0; i < numberOfChannels; i++ {
		channels[i] = channelCreator(ctx, numberOfElements)
	}
	channelMerged := FanIn(ctx, channels...)
	var count int = 0
	for range channelMerged {
		count++
	}
	if count != numberOfElements*numberOfChannels {
		t.Errorf("got %d elements, want %d", count, numberOfElements*numberOfChannels)
	}
}
