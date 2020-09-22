// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"reflect"
	"testing"
	"time"
)

func TestMakisu_Plugin_Exec(t *testing.T) {
	// TODO Write test
}

func TestMakisu_Plugin_Unmarshal(t *testing.T) {
	// setup types
	p := &Plugin{
		GlobalRaw: `
  {"cpu": { "profile": true }, "log": { "fmt": "console", "level": "info", "output": "stdout"} }
`,
	}

	want := &Global{
		CPU: &CPU{
			Profile: true,
		},
		Log: &Log{
			Fmt:    "console",
			Level:  "info",
			Output: "stdout",
		},
	}

	err := p.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %w", err)
	}

	if !reflect.DeepEqual(p.Global, want) {
		t.Errorf("Unmarshal is %+v, want %+v", p.Global, want)
	}
}

func TestMakisu_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Registry: &Registry{
			Password: "superSecretPassword",
			Addr:     "index.docker.io",
			Username: "octocat",
		},
		//nolint
		Build: &Build{
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
		},
		Push: &Push{
			Path:           ".",
			Pushes:         []string{"FOO"},
			RegistryConfig: "{}",
			Replicas:       []string{"FOO"},
			Tag:            "latest",
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestMakisu_Plugin_Validate_BadBuild(t *testing.T) {
	// setup types
	p := &Plugin{
		Registry: &Registry{
			Password: "superSecretPassword",
			Addr:     "index.docker.io",
			Username: "octocat",
		},
		Build: &Build{},
		Push: &Push{
			Path:           ".",
			Pushes:         []string{"FOO"},
			RegistryConfig: "{}",
			Replicas:       []string{"FOO"},
			Tag:            "latest",
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Plugin_Validate_BadPush(t *testing.T) {
	// setup types
	p := &Plugin{
		Registry: &Registry{
			Password: "superSecretPassword",
			Addr:     "index.docker.io",
			Username: "octocat",
		},
		//nolint
		Build: &Build{
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
		},
		Push: &Push{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
