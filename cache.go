// Copyright 2022 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vicanso/elton/middleware"
	"github.com/vicanso/go-cache/v2"
)

// 内存缓存默认值
const defaultCacheSizeMB = 20

type store struct {
	cache *cache.Cache
}

func (s *store) Get(ctx context.Context, key string) ([]byte, error) {
	// 获取失败则忽略
	buf, _ := s.cache.GetBytes(ctx, key)
	return buf, nil
}

func (s *store) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return s.cache.SetBytes(ctx, key, data, ttl)
}

func newRedisClient(uri string) (redis.UniversalClient, error) {
	uriInfo, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	query := uriInfo.Query()
	// 获取密码
	password, _ := uriInfo.User.Password()
	username := uriInfo.User.Username()
	// 转换失败则为0
	poolSize, _ := strconv.Atoi(query.Get("poolSize"))
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      strings.Split(uriInfo.Host, ","),
		Username:   username,
		Password:   password,
		PoolSize:   poolSize,
		MasterName: query.Get("master"),
	}), nil
}

func newCacheStore() (middleware.CacheStore, error) {
	redisURI := os.Getenv(IP_REDIS)
	var secondaryStore cache.Store
	if redisURI != "" {
		client, err := newRedisClient(redisURI)
		if err != nil {
			return nil, err
		}
		secondaryStore = cache.NewRedisStore(client)
	}
	cacheSizeValue := os.Getenv(IP_CACHE_SIZE)
	cacheSize := 0
	if cacheSizeValue != "" {
		cacheSize, _ = strconv.Atoi(cacheSizeValue)
	}
	if cacheSize <= 0 {
		cacheSize = defaultCacheSizeMB
	}

	c, err := cache.New(
		// 默认缓存10分钟
		10*time.Minute,
		// 内存缓存大小
		cache.CacheHardMaxCacheSizeOption(cacheSize),
		// image pipeline server cache
		cache.CacheKeyPrefixOption("ipsc:"),
		cache.CacheSecondaryStoreOption(secondaryStore),
	)
	if err != nil {
		return nil, err
	}

	return &store{
		cache: c,
	}, nil
}
