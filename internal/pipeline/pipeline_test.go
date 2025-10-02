package pipeline

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(42))
}

func TestParseJsonBody(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	words := [4]string{"one", "two", "three", "four"}
	pos := r.Intn(4)
	// testing data
	word := fmt.Sprintf("[\"%s\"]", words[pos])
	data := `[{"word": "hefte","length": 5,"category": "wordle","language": "en"}]`
	//
	channelData := make(chan string)

	result := ParseJsonBody(ctx, channelData)

	go func() {
		defer close(channelData)
		channelData <- data
		channelData <- word
	}()

	for res := range result {
		if res != "hefte" && res != words[pos] {
			t.Fatal("Not parsed properly", res)
		} else {
			fmt.Println(res)
		}
	}
}

func TestGeneratingSentences(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	words := []string{"one", "two", "three", "four"}

	compareArrayToSentence := func(arr []string, elem string) bool {
		for a := range arr {
			if arr[a] == elem {
				return true
			}
		}
		return false
	}

	channelWords := make(chan string)

	go func() {
		defer close(channelWords)
		for i := 0; i < 50; i++ {
			channelWords <- words[r.Intn(len(words))]
		}
	}()

	sentences := GenerateSentences(ctx, channelWords)
	for i := 0; i < 5; i++ {
		sentence := <-sentences
		if sentence == "" {
			t.Fatal("Empty sentence")
		}
		splitSentence := strings.Split(sentence, " ")
		if len(splitSentence) < 3 || len(splitSentence) > 9 {
			t.Fatal("Sentence too long or too short")
		}
		for _, word := range splitSentence {
			if !compareArrayToSentence(words, word) {
				fmt.Println(splitSentence)
				t.Fatal("Some extra word present in the sentence")
			}
		}
	}
}

func generateSentences() (string, []string) {
	words := []string{"one", "two", "three", "four"}
	var builder strings.Builder
	arr := make([]string, r.Intn(7)+3)
	for i := 0; i < len(arr); i++ {
		word := words[r.Intn(len(words))]
		arr[i] = word
		if i != 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(word)
	}
	return builder.String(), arr
}

func removeOccurenceInArray[T comparable](arr []T, target T) ([]T, bool) {
	for i, v := range arr {
		if v == target {
			return append(arr[:i], arr[i+1:]...), true
		}
	}
	return arr, false
}

func TestSplitSentences(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < 10; i++ {
		toParse, arrRes := generateSentences()
		querChan := make(chan string, 1)
		querChan <- toParse
		resChan := SplitSentences(ctx, querChan)
		close(querChan)
		for res := range resChan {
			newArrRes, present := removeOccurenceInArray(arrRes, res)
			if !present {
				t.Fatalf("Unexpected word %q, remaining expected: %v", res, arrRes)
			}
			arrRes = newArrRes
		}
		if len(arrRes) != 0 {
			t.Fatalf("Not all words were emitted, leftover: %v", arrRes)
		}
	}
}

func wordsGenerator(ctx context.Context) <-chan string {
	words := []string{"55555", "1", "22", "333", "4444", "666666", "7777777", "999999999"}
	channelWords := make(chan string)
	go func() {
		defer close(channelWords)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				channelWords <- words[r.Intn(len(words))]
			}
		}
	}()
	return channelWords
}

func TestFilterWords(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const minLen = 4
	for w := range FilterWords(ctx, wordsGenerator(ctx), minLen) {
		if len(w) < minLen {
			t.Fatalf("The word is shorter than min len: %s", w)
		}
	}
}

func TestTakeWords(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	current := 0
	const sizeOfTakes = 50
	for range TakeWords(ctx, wordsGenerator(ctx), sizeOfTakes) {
		current++
	}
	if current != sizeOfTakes {
		t.Fatalf("Should have taken %d words, but got %d", sizeOfTakes, current)
	}
}
