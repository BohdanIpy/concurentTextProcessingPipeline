package netRoutines

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func fetchWord(ctx context.Context) <-chan string {
	words := make(chan string, 2)
	go func() {
		defer close(words)
		for {

			resp, err := http.Get("https://random-word.ryanrk.com/api/en/word/random")
			if err != nil {
				log.Fatal(err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Println(string(body))

			select {
			case <-ctx.Done():
				return
			case words <- string(body):
			}
		}
	}()
	return words
}
