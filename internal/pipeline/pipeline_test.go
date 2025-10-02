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
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestParingApi1(t *testing.T) {
	words := [4]string{"one", "two", "three", "four"}
	pos := r.Intn(4)
	word := fmt.Sprintf("[\"%s\"]", words[pos])
	word, ok := extractWord(word)
	if ok != nil {
		t.Fatal(ok)
	}
	if word != words[pos] {
		t.Error(word)
	}
}

func TestParsingApi2(t *testing.T) {
	data := `[{"word": "hefte","length": 5,"category": "wordle","language": "en"}]`
	word, ok := extractWord(data)
	if ok != nil {
		t.Fatal(ok)
	}
	if word != "hefte" {
		t.Error(word)
	}
}

func TestGeneratingSentences(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	words := []string{"[\"one\"]", "[\"two\"]", "[\"three\"]", "[\"four\"]"}
	parsedWords := []string{"one", "two", "three", "four"}

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
			if !compareArrayToSentence(parsedWords, word) {
				fmt.Println(splitSentence)
				t.Fatal("Some extra word present in the sentence")
			}
		}
	}

}
