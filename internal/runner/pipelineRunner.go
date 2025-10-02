package runner

import (
	"context"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/fanInOut"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/netRoutines"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/pipeline"
	"runtime"
)

func Run(ctx context.Context, urls []string, takes int, minLen int) <-chan string {
	channelsFetchedSentences := make([]<-chan string, len(urls))
	for i, url := range urls {
		channelsFetchedSentences[i] = pipeline.GenerateSentences(ctx, pipeline.ParseJsonBody(ctx, netRoutines.FetchWord(ctx, url)))
	}
	channelFannedSentences := fanInOut.FanIn(ctx, channelsFetchedSentences...)

	numCPU := runtime.NumCPU()
	channelsSplit := make([]<-chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		channelsSplit[i] = pipeline.SplitSentences(ctx, channelFannedSentences)
	}
	channelFannedWords := fanInOut.FanIn(ctx, channelsSplit...)

	channelsFilteringWords := make([]<-chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		channelsFilteringWords[i] = pipeline.FilterWords(ctx, channelFannedWords, minLen)
	}
	channelWordsForTaken := fanInOut.FanIn(ctx, channelsFilteringWords...)

	return pipeline.TakeWords(ctx, channelWordsForTaken, takes)
}
