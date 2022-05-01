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
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vicanso/elton/middleware"
)

type redisStore struct {
	client redis.UniversalClient
}

func (rs *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	// 获取失败则忽略
	buf, _ := rs.client.Get(ctx, key).Bytes()
	return buf, nil
}

func (rs *redisStore) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return rs.client.Set(ctx, key, data, ttl).Err()
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

func newCacheStore(uri string) (middleware.CacheStore, error) {
	client, err := newRedisClient(uri)
	if err != nil {
		return nil, err
	}
	return &redisStore{
		client: client,
	}, nil
}
