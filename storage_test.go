package tgsupergroup

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

var ctx = context.Background()

func newClient() *redis.Client {
	loadDotEnv()
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	user := os.Getenv("REDIS_USER")
	pass := os.Getenv("REDIS_PASS")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}

	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: user,
		Password: pass,
		DB:       db,
	})
	exceptedPong, err := conn.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	if exceptedPong != "PONG" {
		panic("expected PONG")
	}
	return conn
}

func processErr(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func TestRedisCacher(t *testing.T) {
	// Save 4 topic:
	// 2 chats:
	// One chatID has a one topic
	// Second chat ID have a three topic, two topics have a same name
	client := newClient()
	storage := NewRedisStorage(client)

	// Safe first one topic chat
	processErr(storage.Save(ctx, &Topic{
		ChatID:   1,
		Name:     "testing_name",
		ThreadID: 1,
	}), t)
	// Test GetAll and GetByName
	topics, err := storage.GetAll(ctx, 1)
	processErr(err, t)
	if topics == nil {
		t.Error("topics is nil")
	}
	if topics.Len() != 1 {
		t.Error("topics length is not 1")
	}
	if topics.GetID("testing_name") != 1 {
		t.Error("topics does not contain testing_name")
	}
	topicByName, err := storage.GetByName(ctx, 1, "testing_name")
	processErr(err, t)
	if !reflect.DeepEqual(topicByName, &Topic{
		ChatID:   1,
		Name:     "testing_name",
		ThreadID: 1,
	}) {
		t.Error("topic by name is not equal to testing_name")
	}

	// Safe secondChat
	processErr(storage.Save(ctx, &Topic{ChatID: 2, Name: "testing_name_1", ThreadID: 4}), t)
	processErr(storage.Save(ctx, &Topic{ChatID: 2, Name: "testing_name_1", ThreadID: 1}), t)
	processErr(storage.Save(ctx, &Topic{ChatID: 2, Name: "testing_name_2", ThreadID: 2}), t)
	topics, err = storage.GetAll(ctx, 2)
	processErr(err, t)
	if topics.Len() != 2 {
		t.Error("topics length is not 2")
	}
	if topics.GetID("testing_name_1") != 1 {
		t.Error("testing_name_1 threadID is not 1")
	}
	if topics.GetID("testing_name_2") != 2 {
		t.Error("testing_name_1 threadID is not 2")
	}
	topicByName, err = storage.GetByName(ctx, 2, "testing_name_1")
	processErr(err, t)
	if !reflect.DeepEqual(topicByName, &Topic{
		ChatID:   2,
		Name:     "testing_name_1",
		ThreadID: 1,
	}) {
		t.Error("topic by name is not equal to testing_name_1")
	}
}
