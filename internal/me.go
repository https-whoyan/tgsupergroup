package internal

import (
	"context"
	"github.com/https-whoyan/tgsupergroup/errors"
	"net/http"
)

type getMeResponse struct {
	basicResponse
	Result *struct {
		FirstName string `json:"first_name"`
	} `json:"result"`
}

func (r *requester) GetMe(ctx context.Context) (botName string, err error) {
	req, err := r.newRequest(ctx, http.MethodGet, getMeEndpoint, nil)
	if err != nil {
		return
	}
	var dst getMeResponse
	err = r.send(req, &dst)
	if err != nil {
		return
	}
	err = r.checkBasicResponse(dst.basicResponse)
	if err != nil {
		if err.Error() == notFoundDescription {
			return "", errors.ErrInvalidToken
		}
		return
	}
	return dst.Result.FirstName, nil
}
