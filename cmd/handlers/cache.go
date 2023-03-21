package handlers

import (
	"strings"
	"time"

	"github.com/ShikharY10/gbAUTH/cmd/utils"
	"github.com/go-redis/redis"
)

type Cache struct {
	RedisClient *redis.Client
}

func InitializeCacheHandler(redisClient *redis.Client) *Cache {
	return &Cache{
		RedisClient: redisClient,
	}
}

// ==========================TOKEN==================================

// Return true if token is not expired and saved hash of part of token is match is supplied token hash.
func (c *Cache) IsTokenValid(id string, token string, tokenType string) bool {
	var key string
	if tokenType == "access" {
		key = id + ".accessTokenExpiry"
	} else if tokenType == "refresh" {
		key = id + ".refreshTokenExpiry"
	}

	if key == "" {
		return false
	}

	hash := strings.Split(token, ".")[2]

	result := c.RedisClient.Get(key)
	return result.Val() == hash
}

// Saves token's hash part and set a expiry as specified in Redis Cache.
func (c *Cache) SetAccessTokenExpiry(id string, token string, accessTokenExpiry time.Duration) error {
	hash := strings.Split(token, ".")[2]
	result1 := c.RedisClient.Set(id+".accessTokenExpiry", hash, accessTokenExpiry)
	if result1.Err() != nil {
		return result1.Err()
	}
	return nil
}

// Saves token's hash part and set a expiry as specified in Redis Cache.
func (c *Cache) SetRefreshTokenExpiry(id string, token string, refreshTokenExpiry time.Duration) error {
	hash := strings.Split(token, ".")[2]
	result2 := c.RedisClient.Set(id+".refreshTokenExpiry", hash, refreshTokenExpiry)
	if result2.Err() != nil {
		return result2.Err()
	}
	return nil
}

// Deletes both refresh and access token from redis cache
func (c *Cache) DeleteTokenExpiry(id string) {
	c.RedisClient.Del(id + ".accessTokenExpiry")
	c.RedisClient.Del(id + ".refreshTokenExpiry")
}

// ====================OTP====================================

func (c *Cache) RegisterOTP() (string, string) {
	id64 := utils.GenerateRandomId()
	otp := utils.GenerateOTP(6)
	c.RedisClient.Set(id64, otp, time.Duration(5*time.Minute))
	return id64, otp
}

func (c *Cache) VarifyOTP(id string, otp string) bool {
	res := c.RedisClient.Get(id)
	_otp := res.Val()
	if otp == _otp {
		return true
	} else {
		return false
	}
}
