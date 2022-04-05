package http

import (
	"net/http"
	"time"

	"github.com/ybbus/httpretry"
)

func NewHttpClient(client *http.Client) *http.Client {
	if client == nil {
		client = &http.Client{}
	}
	return httpretry.NewCustomClient(
		client,
		// retry x times
		httpretry.WithMaxRetryCount(4),
		// retry on status == 429, if status >= 500, if err != nil, or if response was nil (status == 0)
		httpretry.WithRetryPolicy(func(statusCode int, err error) bool {
			//fmt.Println("retrying now ")
			//fmt.Println("YYYY-MM-DD hh:mm:ss : ", time.Now().Format("2006-01-02 15:04:05"))
			return err != nil || statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError || statusCode == 0
		}),
		httpretry.WithBackoffPolicy(httpretry.ExponentialBackoff(5*time.Second, time.Minute, 2*time.Second)),
	)
}
