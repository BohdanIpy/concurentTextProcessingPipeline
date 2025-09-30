package main

import (
	"context"
	"fmt"
	"github.com/BohdanIpy/concurentTextProcessingPipeline/internal/netRoutines"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	for word := range netRoutines.FetchWord(ctx) {
		fmt.Println(word)
	}
}
