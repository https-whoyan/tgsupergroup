package internal

import (
	"container/heap"
	"context"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/https-whoyan/tgsupergroup/errors"
)

type groupMember struct {
	req       Requester
	c         atomic.Int32
	headIndex int
}

type botsHeap []*groupMember

func (h botsHeap) Len() int           { return len(h) }
func (h botsHeap) Less(i, j int) bool { return h[i].c.Load() < h[j].c.Load() }
func (h botsHeap) Swap(i, j int) {
	h[i].headIndex, h[j].headIndex = j, i
	h[i], h[j] = h[j], h[i]
}

func (h *botsHeap) Push(x interface{}) {
	*h = append(*h, x.(*groupMember))
}
func (h *botsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type botsGroup struct {
	sync.Mutex
	heap botsHeap
}

func NewBotsGroup(ctx context.Context, parseMode ParseMode, botTokens ...string) (Requester, error) {
	requesters := make([]Requester, 0, len(botTokens))
	for _, token := range botTokens {
		r := NewRequester(token, &http.Client{}, parseMode)
		if _, err := r.GetMe(ctx); err == nil {
			requesters = append(requesters, r)
		}
	}
	if len(requesters) == 0 {
		return nil, errors.ErrInvalidToken
	}
	members := make([]*groupMember, len(requesters))
	for i, r := range requesters {
		members[i] = &groupMember{req: r, headIndex: i, c: atomic.Int32{}}
	}
	return &botsGroup{heap: members}, nil
}

func (b *botsGroup) getFirst() *groupMember {
	b.Lock()
	defer b.Unlock()
	first := b.heap[0]
	first.c.Add(1)
	heap.Fix(&b.heap, 0)
	return first
}

func (b *botsGroup) setAsDone(member *groupMember) {
	b.Lock()
	defer b.Unlock()
	member.c.Add(-1)
	heap.Fix(&b.heap, member.headIndex)
}

func (b *botsGroup) GetMe(ctx context.Context) (botName string, err error) {
	mem := b.getFirst()
	defer b.setAsDone(mem)
	return mem.req.GetMe(ctx)
}

func (b *botsGroup) GetChat(ctx context.Context, chatID ChatID) (*Chat, error) {
	mem := b.getFirst()
	defer b.setAsDone(mem)
	return mem.req.GetChat(ctx, chatID)
}

func (b *botsGroup) SendMessageToChat(
	ctx context.Context, chatID ChatID, messageText string, args ...interface{},
) error {
	mem := b.getFirst()
	defer b.setAsDone(mem)
	return mem.req.SendMessageToChat(ctx, chatID, messageText, args...)
}

func (b *botsGroup) SendMessageToTopic(
	ctx context.Context, chatID ChatID, topicID TopicThreadID, messageText string, args ...interface{},
) error {
	mem := b.getFirst()
	defer b.setAsDone(mem)
	return mem.req.SendMessageToTopic(ctx, chatID, topicID, messageText, args...)
}

func (b *botsGroup) CheckTopic(ctx context.Context, topic *Topic) (bool, error) {
	mem := b.getFirst()
	defer b.setAsDone(mem)
	return mem.req.CheckTopic(ctx, topic)
}

func (b *botsGroup) CreateTopic(ctx context.Context, topic *Topic) (topicID *TopicThreadID, err error) {
	mem := b.getFirst()
	defer b.setAsDone(mem)
	return mem.req.CreateTopic(ctx, topic)
}

func (b *botsGroup) Close() error {
	for _, r := range b.heap {
		_ = r.req.Close()
	}
	return nil
}
