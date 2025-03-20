package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
	"github.com/go-redis/redis/v8"
)

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = errors.New("cache miss")

// RedisClientInterface defines the interface for Redis client operations
type RedisClientInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Ping() error
	Close() error
	FlushDB() error
}

// RedisClient represents a Redis client
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStub creates a new Redis stub for testing
func NewRedisStub() *RedisClient {
	log.Println("Creating Redis stub")
	return &RedisClient{
		client: nil,
		ctx:    context.Background(),
	}
}

// NewRedisClient creates a new Redis client with retry mechanism
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	// Log connection details (without password)
	if password != "" {
		log.Printf("Connecting to Redis at %s (DB: %d) with password", addr, db)
	} else {
		log.Printf("Connecting to Redis at %s (DB: %d) without password", addr, db)
	}
	
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	
	// Create a context for Redis operations
	ctx := context.Background()
	
	// Ping Redis to verify the connection - this is handled by the caller with retry mechanism
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	
	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from Redis
func (r *RedisClient) Get(key string) ([]byte, error) {
	if r.client == nil {
		// Stub implementation always returns cache miss
		return nil, ErrCacheMiss
	}
	
	val, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("error getting key %s from Redis: %w", key, err)
	}
	return val, nil
}

// Set stores a value in Redis
func (r *RedisClient) Set(key string, value []byte, expiration time.Duration) error {
	if r.client == nil {
		// Stub implementation does nothing
		return nil
	}
	
	err := r.client.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("error setting key %s in Redis: %w", key, err)
	}
	return nil
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
	if r.client == nil {
		// Stub implementation does nothing
		return nil
	}
	
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting key %s from Redis: %w", key, err)
	}
	return nil
}

// Exists checks if a key exists in Redis
func (r *RedisClient) Exists(key string) (bool, error) {
	if r.client == nil {
		// Stub implementation always returns false
		return false, nil
	}
	
	val, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error checking if key %s exists in Redis: %w", key, err)
	}
	return val > 0, nil
}

// Ping checks if the Redis connection is alive
func (r *RedisClient) Ping() error {
	if r.client == nil {
		// Stub implementation always returns success
		return nil
	}
	
	return r.client.Ping(r.ctx).Err()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	if r.client == nil {
		// Stub implementation does nothing
		return nil
	}
	
	log.Println("Closing Redis connection")
	return r.client.Close()
}

// FlushDB removes all keys from the current database
func (r *RedisClient) FlushDB() error {
	if r.client == nil {
		// Stub implementation does nothing
		return nil
	}
	
	err := r.client.FlushDB(r.ctx).Err()
	if err != nil {
		return fmt.Errorf("error flushing Redis database: %w", err)
	}
	return nil
}

// PostsWithTotal represents posts with their total count
type PostsWithTotal struct {
	Posts []*domain.PostWithUser `json:"posts"`
	Total int                    `json:"total"`
}

// PostCache implements caching for posts
type PostCache struct {
	client RedisClientInterface
}

// NewPostCache creates a new post cache
func NewPostCache(client RedisClientInterface) *PostCache {
	return &PostCache{
		client: client,
	}
}

// GetPosts retrieves posts from the cache
func (c *PostCache) GetPosts() ([]*domain.Post, error) {
	// Get posts from Redis
	data, err := c.client.Get("posts")
	if err != nil {
		return nil, err
	}
	
	// Unmarshal posts
	var posts []*domain.Post
	err = json.Unmarshal(data, &posts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling posts: %w", err)
	}
	
	return posts, nil
}

// SetPosts stores posts in the cache
func (c *PostCache) SetPosts(posts []*domain.Post) error {
	// Marshal posts
	data, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("error marshaling posts: %w", err)
	}
	
	// Set posts in Redis
	return c.client.Set("posts", data, 5*time.Minute)
}

// GetPostsWithUser retrieves posts with user information from the cache
func (c *PostCache) GetPostsWithUser() ([]*domain.PostWithUser, int, error) {
	// Get posts from Redis
	data, err := c.client.Get("posts_with_user")
	if err != nil {
		return nil, 0, err
	}
	
	// Unmarshal posts with total
	var postsWithTotal PostsWithTotal
	err = json.Unmarshal(data, &postsWithTotal)
	if err != nil {
		// Try to unmarshal as old format for backward compatibility
		var posts []*domain.PostWithUser
		err2 := json.Unmarshal(data, &posts)
		if err2 != nil {
			return nil, 0, fmt.Errorf("error unmarshaling posts with user: %w", err)
		}
		return posts, len(posts), nil
	}
	
	return postsWithTotal.Posts, postsWithTotal.Total, nil
}

// SetPostsWithUserAndTotal stores posts with user information and total count in the cache
func (c *PostCache) SetPostsWithUserAndTotal(posts []*domain.PostWithUser, total int) error {
	// Create posts with total
	postsWithTotal := PostsWithTotal{
		Posts: posts,
		Total: total,
	}
	
	// Marshal posts with total
	data, err := json.Marshal(postsWithTotal)
	if err != nil {
		return fmt.Errorf("error marshaling posts with user and total: %w", err)
	}
	
	// Set posts in Redis
	return c.client.Set("posts_with_user", data, 5*time.Minute)
}

// SetPostsWithUser stores posts with user information in the cache (for backward compatibility)
func (c *PostCache) SetPostsWithUser(posts []*domain.PostWithUser) error {
	return c.SetPostsWithUserAndTotal(posts, len(posts))
}

// InvalidatePosts invalidates the posts cache
func (c *PostCache) InvalidatePosts() error {
	// Delete posts from Redis
	err1 := c.client.Delete("posts")
	err2 := c.client.Delete("posts_with_user")
	
	if err1 != nil {
		return fmt.Errorf("error deleting posts cache: %w", err1)
	}
	if err2 != nil {
		return fmt.Errorf("error deleting posts with user cache: %w", err2)
	}
	
	return nil
}

// GetPost retrieves a post from the cache
func (c *PostCache) GetPost(id string) (*domain.Post, error) {
	// Get post from Redis
	key := fmt.Sprintf("post:%s", id)
	data, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	
	// Unmarshal post
	var post domain.Post
	err = json.Unmarshal(data, &post)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling post: %w", err)
	}
	
	return &post, nil
}

// SetPost stores a post in the cache
func (c *PostCache) SetPost(post *domain.Post) error {
	// Marshal post
	key := fmt.Sprintf("post:%s", post.ID)
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("error marshaling post: %w", err)
	}
	
	// Set post in Redis
	return c.client.Set(key, data, 5*time.Minute)
}

// InvalidatePost invalidates a post in the cache
func (c *PostCache) InvalidatePost(id string) error {
	// Delete post from Redis
	key := fmt.Sprintf("post:%s", id)
	err := c.client.Delete(key)
	if err != nil {
		return fmt.Errorf("error deleting post cache: %w", err)
	}
	
	return nil
}

// Ping checks if the cache connection is alive
func (c *PostCache) Ping() error {
	return c.client.Ping()
}
