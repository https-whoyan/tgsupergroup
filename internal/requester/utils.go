package requester

import (
	"errors"
	"fmt"
)

func (r *requester) escapeF(messageText string, args ...interface{}) string {
	return r.parseMode.EscapeText(fmt.Sprintf(messageText, args))
}

const (
	messageThreadNotFound   = "Bad Request: message thread not found"
	chatNotFoundDescription = "Bad Request: chat not found"
	notFound                = "Not Found"
)

var (
	threadNotFoundErr = errors.New("thread not found")
)

func (*requester) checkBasicResponse(resp basicResponse) (parsedErr error) {
	if resp.Ok {
		return nil
	}
	respErrDesc := *resp.Desc
	switch respErrDesc {
	case messageThreadNotFound:
		return threadNotFoundErr
	case chatNotFoundDescription:
		return errChatNotFound
	}
	return errors.New(respErrDesc)
}
