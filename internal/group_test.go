package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type mockRequester struct {
	id    string
	calls atomic.Int32
}

func sleep() { time.Sleep(5 * time.Millisecond) }

func (m *mockRequester) SendMessageToChat(ctx context.Context, chatID ChatID, messageText string, args ...interface{}) error {
	m.calls.Add(1)
	sleep()
	return nil
}

func (m *mockRequester) SendMessageToTopic(ctx context.Context, chatID ChatID, topicID TopicThreadID, messageText string, args ...interface{}) error {
	m.calls.Add(1)
	sleep()
	return nil
}

func (m *mockRequester) GetMe(ctx context.Context) (string, error) {
	return m.id, nil
}

func (m *mockRequester) CheckTopic(ctx context.Context, topic *Topic) (bool, error) {
	return true, nil
}

func (m *mockRequester) CreateTopic(ctx context.Context, topic *Topic) (*TopicThreadID, error) {
	return nil, nil
}

func (m *mockRequester) GetChat(ctx context.Context, chatID ChatID) (*Chat, error) {
	return nil, nil
}

func (m *mockRequester) Close() error {
	return nil
}

func TestBotsGroup_SendMessageToChat(t *testing.T) {
	bot1 := &mockRequester{id: "bot1"}
	bot2 := &mockRequester{id: "bot2"}
	bot3 := &mockRequester{id: "bot3"}

	group := &botsGroup{
		heap: []*groupMember{
			{req: bot1, c: atomic.Int32{}, headIndex: 0},
			{req: bot2, c: atomic.Int32{}, headIndex: 1},
			{req: bot3, c: atomic.Int32{}, headIndex: 2},
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			_ = group.SendMessageToChat(ctx, 1, "Hello!")
		}()
	}
	wg.Wait()

	if bot1.calls.Load() != 1 || bot2.calls.Load() != 1 || bot3.calls.Load() != 1 {
		t.Errorf("unexpected call counts: bot1=%d, bot2=%d, bot3=%d",
			bot1.calls.Load(), bot2.calls.Load(), bot3.calls.Load())
	}
}

func TestBotsGroup_ErrorHandling(t *testing.T) {
	bot1 := &mockRequester{id: "bot1"}
	bot2 := &mockRequester{id: "bot2"}
	bot3 := &mockRequester{id: "bot3"}

	errorBot := &mockRequester{id: "errorBot"}

	group := &botsGroup{
		heap: []*groupMember{
			{req: bot1, c: atomic.Int32{}},
			{req: bot2, c: atomic.Int32{}},
			{req: bot3, c: atomic.Int32{}},
			{req: errorBot, c: atomic.Int32{}},
		},
	}
	err := group.SendMessageToChat(ctx, 1, "Hello!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bot1.calls.Load() != 1 && bot2.calls.Load() != 0 && bot3.calls.Load() != 0 {
		t.Errorf("unexpected call counts: bot1=%d, bot2=%d, bot3=%d",
			bot1.calls.Load(), bot2.calls.Load(), bot3.calls.Load())
	}
}
