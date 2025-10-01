package netRoutines

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func FetchWord(ctx context.Context, url string) <-chan string {
	words := make(chan string, 2)
	go func() {
		defer close(words)

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
				log.Printf("Fetched word from \"%s\", written into chanel: %s", url, string(body))
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return words
}
