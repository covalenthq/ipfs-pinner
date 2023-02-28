package dag

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	IPFS_HTTP_GATEWAYS = []string{"https://w3s.link/ipfs/%s", "https://dweb.link/ipfs/%s", "https://ipfs.io/ipfs/%s"}
)

type httpContentFetcher struct {
	cursor int
}

func NewHttpContentFetcher() *httpContentFetcher {
	return &httpContentFetcher{cursor: 0}
}

func (fetcher *httpContentFetcher) FetchCidViaHttp(ctx context.Context, cid string) ([]byte, error) {
	previous := fetcher.cursor

	for {
		content, err := fetcher.tryFetch(ctx, cid, IPFS_HTTP_GATEWAYS[fetcher.cursor])
		if err != nil {
			log.Printf("%s", err)
		} else {
			return content, nil
		}

		fetcher.cursor = (fetcher.cursor + 1) % len(IPFS_HTTP_GATEWAYS)
		log.Printf("value of cursor: %d", fetcher.cursor)
		if fetcher.cursor == previous {
			return emptyBytes, fmt.Errorf("exhausted listed gateways, but content not found")
		}
	}
}

func (fetcher *httpContentFetcher) tryFetch(ctx context.Context, cid string, url string) ([]byte, error) {
	url = fmt.Sprintf(url, cid)
	log.Printf("trying out %s", url)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := fetcher.Get(timeoutCtx, url)
	if err != nil {
		return emptyBytes, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return ioutil.ReadAll(resp.Body)
	} else {
		return emptyBytes, fmt.Errorf("status from GET %s is %d", url, resp.StatusCode)
	}
}

func (fetcher *httpContentFetcher) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}
