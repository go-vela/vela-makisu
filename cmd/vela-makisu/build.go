// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const buildAction = "build"

type (
	// Build represents the plugin configuration for build information.
	//
	// Makisu documents their command usage:
	// https://github.com/uber/makisu/blob/master/docs/COMMAND.md
	// nolint
	Build struct {
		// Build args should be the argument to the dockerfile as per the spec of ARG.
		BuildArgs []string
		// Commit should be set to explicit to only commit at steps with '#!COMMIT' annotations
		Commit string
		// Image compression level, could be 'no', 'speed', 'size', 'default' (default "default")
		Compression string
		// enables setting the context for building an image
		Context string
		// deny list should be the list of locations to ignore within the resulting docker images
		DenyList []string
		// represents object to configure specific flags within Docker
		Docker *Docker
		// Destination of the image tar
		Destination string
		// The absolute path to the dockerfile (default "Dockerfile")
		File string
		// Represents object to configure specific flags within a http cache
		HTTPCache *HTTPCache
		// Load image into docker daemon after build. Requires access to docker socket
		// at location defined by ${DOCKER_HOST}
		Load bool
		// Time-To-Live for local cache (default 168h0m0s)
		LocalCacheTTL time.Duration
		// Allow makisu to modify files outside of its internal storage dir
		ModifyFS bool
		// Copy / in the storage dir and copy it back after build.
		PreserveRoot bool
		// Registries to push an image to
		Pushes []string
		// Represents object to configure specific flags within a redis cache
		RedisCache *RedisCache
		// Set build-time variables
		RegistryConfig string
		// Push targets with alternative full image names "<registry>/<repo>:<tag>"
		Replicas []string
		// Directory that Makisu uses for temp files and cached layers.
		// Mount this path for better caching performance.
		// If modifyfs is set, default to /makisu-storage; Otherwise default to /tmp/makisu-storage
		Storage string
		// Image tag (required)
		Tag string
		// Set the target build stage to build.
		Target string
	}

	// Docker represnets the "docker" prefixed flags within the
	// Makisu build command.
	Docker struct {
		// docker host should be used to load images into the daemon
		Host string
		// docker scheme should be set for api calls to docker daemon (default "http")
		Scheme string
		// docker version should be Docker API version of daemon for loading images
		Version string
	}

	// HTTPCache represnets the "http-cache" prefixed flags within the
	// Makisu build command.
	HTTPCache struct {
		// The address of the http server for cacheID to layer sha mapping
		Addr string
		// Request header for http cache server.
		Headers []string
	}

	// RedisCache represnets the "redis-cache" prefixed flags within the
	// Makisu build command.
	RedisCache struct {
		// The address of a redis server for cacheID to layer sha mapping
		Addr string
		// The password of the Redis server, should match 'requirepass' in redis.conf
		Password string
		// Time-To-Live for redis cache (default 168h0m0s)
		TTL time.Duration
	}
)

// buildFlags represents for config settings on the cli.
var buildFlags = []cli.Flag{}

// Command formats and outputs the Build command from
// the provided configuration to build a Docker image.
// nolint
func (b *Build) Command() *exec.Cmd {
	logrus.Trace("creating makisu build command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if BuildArgs is provided
	if len(b.BuildArgs) > 0 {
		var args string
		for _, arg := range b.BuildArgs {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for BuildArgs from provided build command
		flags = append(flags, fmt.Sprintf("--build-arg \"%s\"", strings.TrimPrefix(args, " ")))
	}

	// check if Commit is provided
	if len(b.Commit) > 0 {
		// add flag for Commit from provided build command
		flags = append(flags, fmt.Sprintf("--commit %s", b.Commit))
	}

	// check if Compression is provided
	if len(b.Compression) > 0 {
		// add flag for Compression from provided build command
		flags = append(flags, fmt.Sprintf("--compression %s", b.Compression))
	}

	// check if DenyList is provided
	if len(b.DenyList) > 0 {
		var args string
		for _, arg := range b.DenyList {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for DenyList from provided build command
		flags = append(flags, fmt.Sprintf("--blacklist \"%s\"", strings.TrimPrefix(args, " ")))
	}

	// add flags for Docker configuration
	flags = append(flags, b.Docker.Flags()...)

	// check if Destination is provided
	if len(b.Destination) > 0 {
		// add flag for Destination from provided build command
		flags = append(flags, fmt.Sprintf("--dest %s", b.Destination))
	}

	// check if File is provided
	if len(b.File) > 0 {
		// add flag for File from provided build command
		flags = append(flags, fmt.Sprintf("--file %s", b.File))
	}

	// add flags for HTTPCache configuration
	flags = append(flags, b.HTTPCache.Flags()...)

	// check if Load is provided
	if b.Load {
		// add flag for Load from provided build command
		flags = append(flags, "--load")
	}

	// check if LocalCacheTTL is provided
	if !isDurationZero(b.LocalCacheTTL) {
		// add flag for LocalCacheTTL from provided build command
		flags = append(flags, fmt.Sprintf("--local-cache-ttl %s", b.LocalCacheTTL))
	}

	// check if ModifyFS is provided
	if b.ModifyFS {
		// add flag for ModifyFS from provided build command
		flags = append(flags, "--modifyfs")
	}

	// check if PreserveRoot is provided
	if b.PreserveRoot {
		// add flag for PreserveRoot from provided build command
		flags = append(flags, "--preserve-root")
	}

	// check if Pushes is provided
	if len(b.Pushes) > 0 {
		var args string
		for _, arg := range b.Pushes {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for Pushes from provided build command
		flags = append(flags, fmt.Sprintf("--push \"%s\"", strings.TrimPrefix(args, " ")))
	}

	// add flags for RedisCache configuration
	flags = append(flags, b.RedisCache.Flags()...)

	// check if RegistryConfig is provided
	if len(b.RegistryConfig) > 0 {
		// add flag for RegistryConfig from provided build command
		flags = append(flags, fmt.Sprintf("--registry-config %s", b.RegistryConfig))
	}

	// check if Replicas is provided
	if len(b.Replicas) > 0 {
		var args string
		for _, arg := range b.Replicas {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for Replicas from provided build command
		flags = append(flags, fmt.Sprintf("--replica \"%s\"", strings.TrimPrefix(args, " ")))
	}

	// check if Tag is provided
	if len(b.Storage) > 0 {
		// add flag for Tag from provided build command
		flags = append(flags, fmt.Sprintf("--storage %s", b.Storage))
	}

	// check if Tag is provided
	if len(b.Tag) > 0 {
		// add flag for Tag from provided build command
		flags = append(flags, fmt.Sprintf("--tag %s", b.Tag))
	}

	// check if Target is provided
	if len(b.Target) > 0 {
		// add flag for Target from provided build command
		flags = append(flags, fmt.Sprintf("--target %s", b.Target))
	}

	// add the required directory param
	flags = append(flags, b.Context)

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_makisu, append([]string{buildAction}, flags...)...)
}

// Flags formats and outputs the flags for
// configuring a Docker.
func (d *Docker) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if Host is provided
	if len(d.Host) > 0 {
		// add flag for Host from provided build command
		flags = append(flags, fmt.Sprintf("--docker-host %s", d.Host))
	}

	// check if Scheme is provided
	if len(d.Scheme) > 0 {
		// add flag for Scheme from provided build command
		flags = append(flags, fmt.Sprintf("--docker-scheme %s", d.Scheme))
	}

	// check if Version is provided
	if len(d.Version) > 0 {
		// add flag for Version from provided build command
		flags = append(flags, fmt.Sprintf("--docker-version %s", d.Version))
	}

	return flags
}

// Flags formats and outputs the flags for
// configuring a http cache.
func (h *HTTPCache) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if Addr is provided
	if len(h.Addr) > 0 {
		// add flag for Addr from provided build command
		flags = append(flags, fmt.Sprintf("--http-cache-addr %s", h.Addr))
	}

	// check if Headers is provided
	if len(h.Headers) > 0 {
		var args string
		for _, arg := range h.Headers {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for BuildArgs from provided build command
		flags = append(flags, fmt.Sprintf("--http-cache-header \"%s\"", strings.TrimPrefix(args, " ")))
	}

	return flags
}

// Flags formats and outputs the flags for
// configuring a redis cache.
func (r *RedisCache) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if Addr is provided
	if len(r.Addr) > 0 {
		// add flag for Addr from provided build command
		flags = append(flags, fmt.Sprintf("--redis-cache-addr %s", r.Addr))
	}

	// check if Password is provided
	if len(r.Password) > 0 {
		// add flag for Password from provided build command
		flags = append(flags, fmt.Sprintf("--redis-cache-password %s", r.Password))
	}

	// check if TTL is provided
	if !isDurationZero(r.TTL) {
		// add flag for TTL from provided build command
		flags = append(flags, fmt.Sprintf("--redis-cache-ttl %s", r.TTL))
	}

	return flags
}

// Exec formats and runs the commands for building a Docker image.
func (b *Build) Exec() error {
	logrus.Trace("running build with provided configuration")

	// create the build command for the file
	cmd := b.Command()

	// run the build command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build plugin configuration")

	// verify tag are provided
	if len(b.Context) == 0 {
		return fmt.Errorf("no build context provided")
	}

	// verify tag are provided
	if len(b.Tag) == 0 {
		return fmt.Errorf("no build tag provided")
	}

	return nil
}

// helper function to check if the time in duration
// is the zero value.
func isDurationZero(t time.Duration) bool {
	zero, err := time.ParseDuration("0")
	if err != nil {
		logrus.Error(err)
		return false
	}

	// if the duration is zero return true
	if strings.EqualFold(t.String(), zero.String()) {
		return true
	}

	return false
}
