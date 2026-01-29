// Copyright 2025 Chronosphere Inc.
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

package links

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type CustomBuilder struct {
	chronosphereURL string
	path            string
	params          url.Values
}

func (b *Builder) Custom(path string) *CustomBuilder {
	return &CustomBuilder{
		path:            path,
		chronosphereURL: b.chronosphereURL,
	}
}

func (b *CustomBuilder) WithParam(key, value string) *CustomBuilder {
	if b.params == nil {
		b.params = url.Values{}
	}
	if value == "" {
		return b
	}
	b.params.Add(key, value)
	return b
}

func (b *CustomBuilder) WithTimeSec(key string, t time.Time) *CustomBuilder {
	if b.params == nil {
		b.params = url.Values{}
	}
	b.params.Add(key, strconv.FormatInt(t.UnixMilli()/1000, 10))
	return b
}

func (b *CustomBuilder) WithParams(params url.Values) *CustomBuilder {
	if b.params == nil {
		b.params = url.Values{}
	}
	for key, values := range params {
		for _, value := range values {
			b.params.Add(key, value)
		}
	}
	return b
}

func (b *CustomBuilder) String() string {
	return fmt.Sprintf("%s%s?%s", b.chronosphereURL, b.path, b.params.Encode())
}
