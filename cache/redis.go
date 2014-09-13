// The MIT License (MIT)
//
// Copyright (c) 2014 Greivin López
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// This work uses "Redigo" package by Gary Burd:
//
//    https://github.com/garyburd/redigo
//
// --------------  Redigo License --------------
//
// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package rcache

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

// The RedisCacher is an implementation of the MemoryCacher interface.
// See more about MemoryCacher here:
//   https://github.com/greivinlopez/skue
// It is a memory caching system based on Redis:
//   http://redis.io/
type RedisCacher struct {
}

var (
	redisPool     *redis.Pool
	redisServer   = "127.0.0.1:6379"
	redisPassword = ""
)
var ErrCantConnect = errors.New("Can't connect to redis")

// New creates a new RedisCacher.
func New() *RedisCacher {
	return &RedisCacher{}
}

// dial establishes a new Redis connection object and creates
// the connection pool if needed. The password for Redis auth
// is fetched from an environment variable following this:
//
//    http://12factor.net/config
//
func (cacher *RedisCacher) dial() (redis.Conn, error) {
	if redisPool == nil {
		// Retrieve the password from OS environment.
		redisPassword = os.Getenv("RCACHE_REDIS_PASS")
		redisPool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", redisServer)
				if err != nil {
					return nil, err
				}
				if _, err := c.Do("AUTH", redisPassword); err != nil {
					c.Close()
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}
	if redisPool != nil {
		conn := redisPool.Get()
		return conn, nil
	}
	return nil, ErrCantConnect
}

// ----------------------------------------------------------------------------
// 			skue.MemoryCacher implementation
// ----------------------------------------------------------------------------

func (cacher *RedisCacher) Set(key interface{}, value interface{}) error {
	c, err := cacher.dial()
	if err != nil {
		return err
	}
	defer c.Close()

	jsonvalue, err := json.MarshalIndent(value, " ", " ")
	if err != nil {
		return err
	}

	c.Do("SET", key, jsonvalue, "EX", 120)
	return nil
}

func (cacher *RedisCacher) Get(key interface{}, entityPointer interface{}) error {
	c, err := cacher.dial()
	if err != nil {
		return err
	}
	defer c.Close()

	jsonvalue, err := redis.String(c.Do("GET", key))
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(jsonvalue), entityPointer)
	return err
}

func (cacher *RedisCacher) Delete(key interface{}) error {
	c, err := cacher.dial()
	if err != nil {
		return err
	}
	defer c.Close()

	c.Do("DEL", key)
	return nil
}

// ----------------------------------------------------------------------------
