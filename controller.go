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
	"bytes"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	pipeline "github.com/vicanso/image-pipeline"
)

func imagePipelineFromQuery(c *elton.Context) error {
	tasks := c.Request.URL.RawQuery
	if tasks == "" {
		return hes.New("Task不能为空")
	}

	jobs, err := pipeline.Parse(tasks, c.GetRequestHeader("Accept"))
	if err != nil {
		return err
	}
	img, err := pipeline.Do(c.Context(), nil, jobs...)
	if err != nil {
		return err
	}
	data, format := img.Bytes()
	if len(data) == 0 {
		buf, err := img.PNG()
		if err != nil {
			return err
		}
		data = buf
		format = "png"
	}
	c.CacheMaxAge(time.Hour, 5*time.Minute)
	c.SetContentTypeByExt("." + format)
	c.BodyBuffer = bytes.NewBuffer(data)

	return nil
}
