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
	"encoding/json"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	pipeline "github.com/vicanso/image-pipeline"
)

type imagePipelineFromBodyParams struct {
	Image []byte `json:"image"`
	Tasks string `json:"tasks"`
}

func imagePipeline(c *elton.Context, tasks string, originalImg *pipeline.Image) error {
	jobs, err := pipeline.Parse(tasks, c.GetRequestHeader("Accept"))
	if err != nil {
		return err
	}
	img, err := pipeline.Do(c.Context(), originalImg, jobs...)
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

func imagePipelineFromQuery(c *elton.Context) error {
	tasks := c.Request.URL.RawQuery
	if tasks == "" {
		return hes.New("Task不能为空")
	}
	return imagePipeline(c, tasks, nil)
}

func imagePipelineFromBody(c *elton.Context) error {
	params := imagePipelineFromBodyParams{}
	err := json.Unmarshal(c.RequestBody, &params)
	if err != nil {
		return err
	}
	if len(params.Image) == 0 || params.Tasks == "" {
		return hes.New("Image and tasks can not be nil")
	}
	img, err := pipeline.NewImageFromBytes(params.Image)
	if err != nil {
		return err
	}
	return imagePipeline(c, params.Tasks, img)
}
