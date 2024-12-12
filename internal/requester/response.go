package requester

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type queryArgs map[string]string

type errResponse struct {
	ErrCode *int    `json:"error_code"`
	Desc    *string `json:"description"`
}

type basicResponse struct {
	errResponse
	Ok bool `json:"ok"`
}

func (r *requester) format(endpoints string) string {
	return basicURL + r.botToken + "/" + strings.TrimPrefix(endpoints, "/")
}

func (r *requester) newRequest(
	ctx context.Context, method string, endpoints string, queryArgs queryArgs,
) (*http.Request, error) {
	baseURL := r.format(endpoints)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q := parsedURL.Query()
	for key, value := range queryArgs {
		q.Set(key, value)
	}
	parsedURL.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *requester) send(req *http.Request, dst interface{}) error {
	r.mu.Lock()
	resp, err := r.httpClient.Do(req)
	time.Sleep(requestTiming)
	r.mu.Unlock()
	defer func(response *http.Response) {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}(resp)
	if err != nil {
		return err
	}
	if dst == nil {
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return errors.New(resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}
