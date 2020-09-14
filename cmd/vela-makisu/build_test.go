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

func TestImg_Build_Command(t *testing.T) {
	// setup types
	b := &Build{
		BuildArgs:   []string{"FOO"},
		ContextPath: ".",
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
			TTL:      1 * time.Minute,
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
		fmt.Sprintf("--build-arg \"%s\"", b.BuildArgs[0]),
		fmt.Sprintf("--commit %s", b.Commit),
		fmt.Sprintf("--compression %s", b.Compression),
		fmt.Sprintf("--blacklist \"%s\"", b.DenyList[0]),
		fmt.Sprintf("--docker-host %s", b.Docker.Host),
		fmt.Sprintf("--docker-scheme %s", b.Docker.Scheme),
		fmt.Sprintf("--docker-version %s", b.Docker.Version),
		fmt.Sprintf("--dest %s", b.Destination),
		fmt.Sprintf("--file %s", b.File),
		fmt.Sprintf("--http-cache-addr %s", b.HTTPCache.Addr),
		fmt.Sprintf("--http-cache-header \"%s\"", b.HTTPCache.Headers[0]),
		"--load",
		fmt.Sprintf("--local-cache-ttl %s", b.LocalCacheTTL),
		"--modifyfs",
		"--preserve-root",
		fmt.Sprintf("--push \"%s\"", b.Pushes[0]),
		fmt.Sprintf("--redis-cache-addr %s", b.RedisCache.Addr),
		fmt.Sprintf("--redis-cache-password %s", b.RedisCache.Password),
		fmt.Sprintf("--redis-cache-ttl %s", b.RedisCache.TTL),
		fmt.Sprintf("--registry-config %s", b.RegistryConfig),
		fmt.Sprintf("--replica \"%s\"", b.Replicas[0]),
		fmt.Sprintf("--storage %s", b.Storage),
		fmt.Sprintf("--tag %s", b.Tag),
		fmt.Sprintf("--target %s", b.Target),
		".",
	)

	got := b.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestMakisu_Build_Exec_Error(t *testing.T) {
	// setup types
	b := &Build{}

	err := b.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestMakisu_Build_Validate(t *testing.T) {
	// setup types
	b := &Build{
		ContextPath: ".",
		Tag:         "latest",
	}

	err := b.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestMakisu_Build_Validate_NoContextPath(t *testing.T) {
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
		ContextPath: ".",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
