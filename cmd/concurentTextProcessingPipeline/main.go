package main

import (
	"context"
	"fmt"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/netRoutines"
	"runtime"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	urls := []string{
		//"https://random-word.ryanrk.com/api/en/word/random",
		"https://random-word-api.vercel.app/api?words=1",
		"https://random-words-api.kushcreates.com/api?language=en&words=1",
	}

	for _, url := range urls {
		go func() {
			for word := range netRoutines.FetchWord(ctx, url) {
				fmt.Println(word)
			}
		}()
	}

	time.Sleep(time.Second * 10)
	cancel()
	time.Sleep(time.Second * 10)
	fmt.Println("----", runtime.NumGoroutine())
}
