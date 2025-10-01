package pipeline

import (
	"context"
	"math/rand"
	"strings"
	"time"
)

func GenerateSentences(ctx context.Context, words <-chan string) <-chan string {
	sentences := make(chan string)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		defer close(sentences)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// skipped intentionally
			}

			var builder strings.Builder
			length := r.Intn(7) + 3

			for i := 0; i < length; i++ {
				select {
				case <-ctx.Done():
					return
				case word, ok := <-words:
					if !ok {
						// upstream closed, stop sentence generation
						return
					}
					if i > 0 {
						builder.WriteByte(' ')
					}
					builder.WriteString(word)
				}
			}
			select {
			case <-ctx.Done():
				return
			case sentences <- builder.String():
			}
		}
	}()
	return sentences
}

func SplitSentences(ctx context.Context, sentences <-chan string) <-chan string {
	words := make(chan string, 10)
	go func() {
		defer close(words)
		for {
			select {
			case <-ctx.Done():
				return
			case sentence, ok := <-sentences:
				if !ok {
					return
				}
				for _, word := range strings.Fields(sentence) {
					select {
					case <-ctx.Done():
						return
					case words <- word:
					}
				}
			}
		}
	}()
	return words
}
