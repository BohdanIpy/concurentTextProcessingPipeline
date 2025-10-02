package main

import (
	"context"
	"fmt"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/runner"
	"time"
)

func main() {
	// Comment if you want logs
	// log.SetOutput(io.Discard)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	urls := []string{
		// Uncoment to recv an API from local server
		// "https://random-word.ryanrk.com/api/en/word/random",
		// "https://random-word-api.vercel.app/api?words=1",
		// "https://random-words-api.kushcreates.com/api?language=en&words=1",
		//  Your custom API "http://localhost:8080/api/v1/array",
	}

	const numOfTakes = 4
	const minLen = 3
	counter := 0
	for word := range runner.Run(ctx, urls, numOfTakes, minLen) {
		fmt.Println(word)
		counter++
	}
	fmt.Println(counter)

	cancel()
	time.Sleep(time.Second * 30)
	/*
		Uncoment to check if all the created goroutines are destroyed
	*/
	/*
		fmt.Println("----", runtime.NumGoroutine())
		buf := make([]byte, 1<<16)
		n := runtime.Stack(buf, true)
		fmt.Println(string(buf[:n]))
	*/
}
