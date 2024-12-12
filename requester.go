package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type queryArgs map[string]string
type ChatID = int64
type MessageID = int64

type requester struct {
	mu         sync.Mutex
	basicURL   string
	botToken   string
	botName    string
	parseMode  ParseMode
	httpClient *http.Client
}

func newRequester(
	botToken string,
	httpCli *http.Client,
	parseMode ParseMode,
	botName string,
) *requester {
	if httpCli == nil {
		httpCli = http.DefaultClient
	}
	return &requester{
		basicURL:   basicURL,
		botToken:   botToken,
		httpClient: httpCli,
		parseMode:  parseMode,
		botName:    botName,
	}
}

type getMeResponse struct {
	basicResponse
	Result *struct {
		FirstName string `json:"first_name"`
	} `json:"result"`
}

func (r *requester) getMe(ctx context.Context) (botName string, err error) {
	req, err := r.newRequest(
		ctx,
		http.MethodGet,
		pingEndpoint,
		nil,
	)
	if err != nil {
		return
	}
	var dst getMeResponse
	err = r.send(req, &dst)
	if err != nil {
		return
	}
	if !dst.Ok {
		return "", errors.New(*dst.Desc)
	}
	return dst.Result.FirstName, nil
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
