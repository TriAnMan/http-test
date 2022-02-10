package worker

import (
	"context"
	"crypto/md5"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{}

const timeout = 10 * time.Second

func Do(job Job) (result Job) {
	job.Url = PrepareUrl(job.Url)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return makeRequest(ctx, job)
}

func makeRequest(ctx context.Context, job Job) (result Job){
	result = job

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, result.Url, nil)
	if err != nil {
		result.Error = err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Error = err
		return
	}

	// ensure a clean termination
	defer func(){
		// support HTTP/1.x "keep-alive"
		// see net/http.Response.Body
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		result.Error = HttpError(resp.Status)
		return
	}

	h := md5.New()
	if _, err = io.Copy(h, resp.Body); err != nil {
		result.Error = err
		return
	}

	result.Hash = h.Sum(nil)

	return
}
