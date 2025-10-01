package runner

import (
	"context"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/fanInOut"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/netRoutines"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/pipeline"
	"runtime"
)

func Run(ctx context.Context, urls []string, takes int) <-chan string {
	// Fetch + sentence generator
	channelsFetchedSentences := make([]<-chan string, len(urls))
	for i, url := range urls {
		channelsFetchedSentences[i] = pipeline.GenerateSentences(ctx, netRoutines.FetchWord(ctx, url))
	}
	channelFannedSentences := fanInOut.FanIn(ctx, channelsFetchedSentences...)

	// Split stage (fan-out)
	numCPU := runtime.NumCPU()
	channelsSplit := make([]<-chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		channelsSplit[i] = pipeline.SplitSentences(ctx, channelFannedSentences)
	}
	channelFannedWords := fanInOut.FanIn(ctx, channelsSplit...)

	// Filter stage (fan-out)
	channelsFilteringWords := make([]<-chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		channelsFilteringWords[i] = pipeline.FilterWords(ctx, channelFannedWords, 3)
	}
	channelWordsForTaken := fanInOut.FanIn(ctx, channelsFilteringWords...)

	// Take stage
	return pipeline.TakeWords(ctx, channelWordsForTaken, takes)
}
