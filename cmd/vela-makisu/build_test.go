// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

func TestMakisu_Build_Command(t *testing.T) {
	// setup types
	// nolint
	b := &Build{
		BuildArgs:   []string{"FOO"},
		Context:     ".",
		Commit:      "b0bb040e6a6d71ddf98684349c42d36fa6c539ad",
		Compression: "default",
		DenyList:    []string{"FOO"},
		Docker: &Docker{
			Host:    "unix:///var/run/docker.sock",
			Scheme:  "http",
			Version: "1.21",
		},
		Destination: "/path/to/dest",
		File:        "Dockerfile",
		HTTPCache: &HTTPCache{
			Addr:    "http://localhost",
			Headers: []string{"Content-type: Application/json"},
		},
		Load:          true,
		LocalCacheTTL: 1 * time.Minute,
		ModifyFS:      true,
		PreserveRoot:  true,
		Pushes:        []string{"FOO"},
		RedisCache: &RedisCache{
			Addr:     "http://localhost",
			Password: "superSecret123",
			TTL:      "1m0s",
		},
		RegistryConfig: "{}",
		Replicas:       []string{"FOO"},
		Storage:        "foo",
		Tag:            "latest",
		Target:         "dev",
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	want := exec.Command(
		_makisu,
		buildAction,
		fmt.Sprintf("--build-arg=%s", b.BuildArgs[0]),
		fmt.Sprintf("--commit=%s", b.Commit),
		fmt.Sprintf("--compression=%s", b.Compression),
		fmt.Sprintf("--blacklist=%s", b.DenyList[0]),
		fmt.Sprintf("--docker-host=%s", b.Docker.Host),
		fmt.Sprintf("--docker-scheme=%s", b.Docker.Scheme),
		fmt.Sprintf("--docker-version=%s", b.Docker.Version),
		fmt.Sprintf("--dest=%s", b.Destination),
		fmt.Sprintf("--file=%s", b.File),
		fmt.Sprintf("--http-cache-addr=%s", b.HTTPCache.Addr),
		fmt.Sprintf("--http-cache-header=%s", b.HTTPCache.Headers[0]),
		"--load",
		fmt.Sprintf("--local-cache-ttl=%s", b.LocalCacheTTL),
		"--modifyfs",
		"--preserve-root",
		fmt.Sprintf("--push=%s", b.Pushes[0]),
		fmt.Sprintf("--redis-cache-addr=%s", b.RedisCache.Addr),
		fmt.Sprintf("--redis-cache-password=%s", b.RedisCache.Password),
		fmt.Sprintf("--redis-cache-ttl=%s", b.RedisCache.TTL),
		fmt.Sprintf("--registry-config=%s", b.RegistryConfig),
		fmt.Sprintf("--replica=%s", b.Replicas[0]),
		fmt.Sprintf("--storage=%s", b.Storage),
		fmt.Sprintf("--tag=%s", b.Tag),
		fmt.Sprintf("--target=%s", b.Target),
		".",
	)

	got, _ := b.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestMakisu_Build_Exec_Error(t *testing.T) {
	// setup types
	b := &Build{
		Docker:     &Docker{},
		HTTPCache:  &HTTPCache{},
		RedisCache: &RedisCache{},
	}

	err := b.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestMakisu_Build_Unmarshal(t *testing.T) {
	// setup types
	b := &Build{
		DockerRaw: `
  {"host": "unix:///var/run/docker.sock", "scheme": "https", "version": "v1.21.1"}
`,
		HTTPCacheRaw: `
  {"addr": "localhost:8080", "headers": ["CUSTOM_HEADER"]}
`,
		RedisCacheRaw: `
{"addr": "redis.company.com", "password": "superSecretPassword", "ttl": "1m0s"}
`,
	}

	want := &Build{
		Docker: &Docker{
			Host:    "unix:///var/run/docker.sock",
			Scheme:  "https",
			Version: "v1.21.1",
		},
		HTTPCache: &HTTPCache{
			Addr:    "localhost:8080",
			Headers: []string{"CUSTOM_HEADER"},
		},
		RedisCache: &RedisCache{
			Addr:     "redis.company.com",
			Password: "superSecretPassword",
			TTL:      "1m0s",
		},
	}

	err := b.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %v", err)
	}

	if !reflect.DeepEqual(b.Docker, want.Docker) {
		t.Errorf("Unmarshal is %v, want %v", b.Docker, want.Docker)
	}

	if !reflect.DeepEqual(b.HTTPCache, want.HTTPCache) {
		t.Errorf("Unmarshal is %v, want %v", b.HTTPCache, want.HTTPCache)
	}

	if !reflect.DeepEqual(b.RedisCache, want.RedisCache) {
		t.Errorf("Unmarshal is %v, want %v", b.RedisCache, want.RedisCache)
	}
}

func TestMakisu_Build_Unmarshal_FailDockerUnmarshal(t *testing.T) {
	// setup types
	b := &Build{
		DockerRaw: "!@#$%^&*()",
		HTTPCacheRaw: `
  {"addr": "localhost:8080", "headers": ["CUSTOM_HEADER"]}
`,
		RedisCacheRaw: `
{"addr": "redis.company.com", "password": "superSecretPassword", "ttl": "1m0s"}
`,
	}

	err := b.Unmarshal()
	if err == nil {
		t.Errorf("Unmarshal should have returned err")
	}
}

func TestMakisu_Build_Unmarshal_FailHTTPCacheUnmarshal(t *testing.T) {
	// setup types
	b := &Build{
		DockerRaw: `
  {"host": "unix:///var/run/docker.sock", "scheme": "https", "version": "v1.21.1"}
`,
		HTTPCacheRaw: "!@#$%^&*()",
		RedisCacheRaw: `
{"addr": "redis.company.com", "password": "superSecretPassword", "ttl": "1m0s"}
`,
	}

	err := b.Unmarshal()
	if err == nil {
		t.Errorf("Unmarshal should have returned err")
	}
}

func TestMakisu_Build_Unmarshal_FailRedisCacheUnmarshal(t *testing.T) {
	// setup types
	b := &Build{
		DockerRaw: `
  {"host": "unix:///var/run/docker.sock", "scheme": "https", "version": "v1.21.1"}
`,
		HTTPCacheRaw: `
  {"addr": "localhost:8080", "headers": ["CUSTOM_HEADER"]}
`,
		RedisCacheRaw: "!@#$%^&*()",
	}

	err := b.Unmarshal()
	if err == nil {
		t.Errorf("Unmarshal should have returned err")
	}
}

func TestMakisu_Build_Validate(t *testing.T) {
	// setup types
	b := &Build{
		Context: ".",
		Tag:     "latest",
	}

	err := b.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestMakisu_Build_Validate_NoContext(t *testing.T) {
	// setup types
	b := &Build{
		Tag: "latest",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Build_Validate_NoTag(t *testing.T) {
	// setup types
	b := &Build{
		Context: ".",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Docker_Flag(t *testing.T) {
	// setup types
	d := &Docker{
		Host:    "unix:///var/run/docker.sock",
		Scheme:  "https",
		Version: "v1.21.1",
	}

	want := []string{
		fmt.Sprintf("--docker-host=%s", d.Host),
		fmt.Sprintf("--docker-scheme=%s", d.Scheme),
		fmt.Sprintf("--docker-version=%s", d.Version),
	}

	got := d.Flags()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flag is %v, want %v", got, want)
	}
}

func TestMakisu_HTTPCache_Flag(t *testing.T) {
	// setup types
	h := &HTTPCache{
		Addr:    "localhost:8080",
		Headers: []string{"CUSTOM_HEADER"},
	}

	want := []string{
		fmt.Sprintf("--http-cache-addr=%s", h.Addr),
		fmt.Sprintf("--http-cache-header=%s", h.Headers[0]),
	}

	got := h.Flags()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flag is %v, want %v", got, want)
	}
}

func TestMakisu_RedisCache_Flag(t *testing.T) {
	// setup types
	r := &RedisCache{
		Addr:     "redis.company.com",
		Password: "superSecretPassword",
		TTL:      "1m0s",
	}

	want := []string{
		fmt.Sprintf("--redis-cache-addr=%s", r.Addr),
		fmt.Sprintf("--redis-cache-password=%s", r.Password),
		fmt.Sprintf("--redis-cache-ttl=%s", r.TTL),
	}

	got, _ := r.Flags()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flag is %v, want %v", got, want)
	}
}
