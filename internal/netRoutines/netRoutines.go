package netRoutines

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func FetchWord(ctx context.Context) <-chan string {
	words := make(chan string, 2)
	go func() {
		defer close(words)

		url := "https://random-word.ryanrk.com/api/en/word/random"

		for {

			select {
			case <-ctx.Done():
				return
			default:
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				log.Println("Request build error: ", err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Request build error: ", err)
				return
			}

			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()

			select {
			case <-ctx.Done():
				return
			case words <- string(body):
				log.Printf("Fetched word, written into chanel: %s", string(body))
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return words
}
