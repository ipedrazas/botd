package api

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func getHooks(redisClient *Redis) []string {
	keys, err := redisClient.Keys("*").Result()
	var hooks []string
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		hook, err := redisClient.Get(key).Result()
		if err != nil {
			panic(err)
		}
		hooks = append(hooks, hook)
	}

	return hooks
}

func setHook(redisClient *Redis, hooks string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	keyName := "hook-" + RandString(6)

	if err := redisClient.Set(keyName, hooks, 0).Err(); err != nil {
		return "", err
	}
	return keyName, nil
}
