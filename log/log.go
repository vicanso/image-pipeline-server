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

package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var logger = newLogger()

func newLogger() *zerolog.Logger {
	// 全局禁用sampling
	zerolog.DisableSampling(true)
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00"

	l := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	return &l
}

func Info(ctx context.Context) *zerolog.Event {
	return logger.Info()
}

func Error(ctx context.Context) *zerolog.Event {
	return logger.Error()
}

func Debug(ctx context.Context) *zerolog.Event {
	return logger.Debug()
}

func Warn(ctx context.Context) *zerolog.Event {
	return logger.Warn()
}
