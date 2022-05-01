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
	"errors"
	"os"
	"strings"

	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
	pipeline "github.com/vicanso/image-pipeline"
	"github.com/vicanso/image-pipeline-server/log"
)

const IP_REDIS = "IP_REDIS"

func init() {
	// 只处理以IP_开头的环境变量
	envPrefix := "IP_"
	envFinderPrefix := "FINDER_"
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, IP_REDIS) ||
			!strings.HasPrefix(env, envPrefix) {
			continue
		}
		env = env[len(envPrefix):]
		if strings.HasPrefix(env, envFinderPrefix) {
			err := addFinder(env[len(envFinderPrefix):])
			if err != nil {
				log.Error(context.Background()).
					Err(err).Msg("")
			}
		} else {
			addTaskAlias(env)
		}
	}
}

func addFinder(env string) error {
	arr := strings.Split(env, "=")
	if len(arr) < 2 {
		return errors.New("env value is invald")
	}
	name := arr[0]
	uri := strings.Join(arr[1:], "=")
	return pipeline.AddFinder(name, uri)
}

func addTaskAlias(env string) {
	arr := strings.Split(env, "=")
	if len(arr) < 2 {
		return
	}
	pipeline.TaskAlias(arr[0], strings.Join(arr[1:], "="))
}

func main() {

	e := elton.New()

	e.Use(middleware.NewLogger(middleware.LoggerConfig{
		Format: middleware.LoggerCombined,
		OnLog: func(s string, c *elton.Context) {
			log.Info(c.Context()).Msg(s)
		},
	}))

	e.Use(middleware.NewDefaultError())
	redisURI := os.Getenv(IP_REDIS)
	if redisURI != "" {
		store, err := newCacheStore(redisURI)
		if err != nil {
			log.Error(context.Background()).
				Err(err).
				Msg("new cache store fail")
		} else {
			e.Use(middleware.NewDefaultCache(store))
			log.Info(context.Background()).
				Msg("cache middleware success")
		}
	}

	e.GET("/", imagePipelineFromQuery)

	addr := ":7001"
	log.Info(context.Background()).
		Str("addr", addr).
		Msg("server is running")

	err := e.ListenAndServe(addr)
	if err != nil {
		log.Error(context.Background()).
			Err(err).
			Msg("")
	}
}
