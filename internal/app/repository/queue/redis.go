package queue

import "github.com/redis/go-redis/v9"

// redisQueue is a struct that implements the Repository interface using Redis as the underlying storage mechanism.
type redisQueue struct {
	c         *redis.Client
	queueName string
}

// NewRedisQueue creates a new instance of redisQueue.
// Parameters:
//   - c: A pointer to a redis.Client instance that will be used to interact with the Redis server.
//   - queueName: The name of the Redis list that will be used as the message queue.
//
// Returns a Repository interface that can be used to push messages to the Redis queue.
func NewRedisQueue(c *redis.Client, queueName string) Repository {
	return &redisQueue{
		c:         c,
		queueName: queueName,
	}
}
