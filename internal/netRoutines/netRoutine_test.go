package netRoutines

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchWord(t *testing.T) {
	urls := []string{
		//"https://random-word.ryanrk.com/api/en/word/random",
		"https://random-word-api.vercel.app/api?words=1",
		"https://random-words-api.kushcreates.com/api?language=en&words=1",
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, url := range urls {
		go func() {
			words := FetchWord(ctx, url)
			for word := range words {
				fmt.Println(word)
			}
		}()
	}
	time.Sleep(time.Second * 3)
}

func arrayComparator[T comparable](arr1 []T, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

func TestOutputFetching(t *testing.T) {
	wordsToSend := []string{"alpha", "beta"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		words := wordsToSend
		_ = json.NewEncoder(w).Encode(words)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := FetchWord(ctx, ts.URL)

	select {
	case w, ok := <-ch:
		if !ok {
			t.Fatal("channel closed unexpectedly")
		}

		var words []string
		if err := json.Unmarshal([]byte(w), &words); err != nil {
			t.Fatal(err)
		}

		if !arrayComparator(words, wordsToSend) {
			t.Errorf("expected 'alpha', got %q", w)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for word")
	}
}
