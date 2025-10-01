package main

import (
	"context"
	"fmt"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/fanInOut"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/netRoutines"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/pipeline"
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

	channelsFetchedSentences := make([]<-chan string, len(urls))
	for i, url := range urls {
		channelsFetchedSentences[i] = pipeline.GenerateSentences(ctx, netRoutines.FetchWord(ctx, url))
	}
	channelFannedSentences := fanInOut.FanIn(ctx, channelsFetchedSentences...)

	numCPU := runtime.NumCPU()
	channelsSplit := make([]<-chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		channelsSplit[i] = pipeline.SplitSentences(ctx, channelFannedSentences)
	}
	channelFannedWords := fanInOut.FanIn(ctx, channelsFetchedSentences...)

	channelsFilteringWords := make([]<-chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		channelsFilteringWords[i] = pipeline.FilterWords(ctx, channelFannedWords, 3)
	}

	channelWordsForTaken := fanInOut.FanIn(ctx, channelsFilteringWords...)
	for word := range pipeline.TakeWords(ctx, channelWordsForTaken, 5) {
		fmt.Println(word)
	}

	time.Sleep(time.Second * 10)
	cancel()
	time.Sleep(time.Second * 10)
	fmt.Println("----", runtime.NumGoroutine())
}
