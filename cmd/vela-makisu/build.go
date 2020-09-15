// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
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
		// enables setting build time arguments for the Dockerfile
		BuildArgs []string
		// enables setting compression on the tar file built - options: (no|speed|size|default)
		Commit string
		// Image compression level, could be 'no', 'speed', 'size', 'default' (default "default")
		Compression string
		// enables settting the context for the image to be built
		Context string
		// enables setting list of locations to be ignored within docker image
		DenyList []string
		// Used for translating the raw docker configuration
		Docker *Docker
		// enables setting configuration on the Docker daemon
		DockerRaw string
		// enables setting the output of the tar file
		Destination string
		// enables setting a the absolute path to dockerfile
		File string
		// Used for translating the raw http cache configuration
		HTTPCache *HTTPCache
		// enables setting custom http options caching
		HTTPCacheRaw string
		// enables loading a docker image into the docker daemon post build
		Load bool
		// enables setting a time to live for the local docker cache (default 168h0m0s)
		LocalCacheTTL time.Duration
		// enables setting makisu to modify files outside its internal storage directories
		ModifyFS bool
		// enables setting copying storage from root in the storage during and after build
		PreserveRoot bool
		// enables setting registries to push the image to
		Pushes []string
		// Used for translating the redis cache configuration
		RedisCache *RedisCache
		// enables setting custom redis server for caching
		RedisCacheRaw string
		// enables setting registry configuration for authentication
		RegistryConfig string
		// enables setting pushing image to alternative targets i.e. \"<registry>/<repo>:<tag>\"
		Replicas []string
		// enables setting a directory for makisu to use for temp files and cached layers
		Storage string
		// enables setting the tag for an image
		Tag string
		// enables setting the target build stage to build
		Target string
	}

	// Docker represnets the "docker" prefixed flags within the
	// Makisu build command.
	Docker struct {
		// enables setting the docker host to be used to load images into the daemon
		Host string
		// enables setting the docker scheme to be set for api calls to docker daemon (default "http")
		Scheme string
		// enables setting the docker version to be set for loading images
		Version string
	}

	// HTTPCache represnets the "http-cache" prefixed flags within the
	// Makisu build command.
	HTTPCache struct {
		// enables setting the address of the http server for cacheID to layer sha mapping
		Addr string
		// enables setting the request headers for http cache server.
		Headers []string
	}

	// RedisCache represnets the "redis-cache" prefixed flags within the
	// Makisu build command.
	RedisCache struct {
		// enables setting the address of a redis server for cacheID to layer sha mapping
		Addr string
		// enables seting the password of the Redis server, should match 'requirepass' in redis.conf
		Password string
		// enables setting the Time-To-Live for redis cache (default 168h0m0s)
		TTL time.Duration
	}
)

// buildFlags represents for config settings on the cli.
//nolint
var buildFlags = []cli.Flag{
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_BUILD_ARGS"},
		FilePath: string("/vela/parameters/makisu/build/build_args,/vela/secrets/makisu/build/build_args"),
		Name:     "build.build-args",
		Usage:    "enables setting build time arguments for the dockerfile",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_COMMIT"},
		FilePath: string("/vela/parameters/makisu/build/commit,/vela/secrets/makisu/build/commit"),
		Name:     "build.commit",
		Usage:    "enables setting commit info for #!COMMIT annotations",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_COMPRESSION"},
		FilePath: string("/vela/parameters/makisu/build/compression,/vela/secrets/makisu/build/compression"),
		Name:     "build.compression",
		Usage:    "enables setting compression on the tar file built - options: (no|speed|size|default)",
		Value:    "default",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_CONTEXT"},
		FilePath: string("/vela/parameters/makisu/build/context,/vela/secrets/makisu/build/context"),
		Name:     "build.context",
		Usage:    "enables setting the context for the image to be built",
		Value:    ".",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_DENY_LIST"},
		FilePath: string("/vela/parameters/makisu/build/deny_list,/vela/secrets/makisu/build/deny_list"),
		Name:     "build.deny-list",
		Usage:    "enables setting list of locations to be ignored within docker image",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_DOCKER_OPTIONS"},
		FilePath: string("/vela/parameters/makisu/build/docker_options,/vela/secrets/makisu/build/docker_options"),
		Name:     "build.docker-options",
		Usage:    "enables setting configuration on the docker daemon",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_DESTINATION"},
		FilePath: string("/vela/parameters/makisu/build/destination,/vela/secrets/makisu/build/destination"),
		Name:     "build.destination",
		Usage:    "enables setting the output of the tar file",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_FILE"},
		FilePath: string("/vela/parameters/makisu/build/file,/vela/secrets/makisu/build/file"),
		Name:     "build.file",
		Usage:    "enables setting a the absolute path to dockerfile",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_HTTP_CACHE_OPTIONS"},
		FilePath: string("/vela/parameters/makisu/build/http_cache_options,/vela/secrets/makisu/build/http_cache_options"),
		Name:     "build.http-cache-options",
		Usage:    "enables setting custom http options caching",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_LOAD"},
		FilePath: string("/vela/parameters/makisu/build/load,/vela/secrets/makisu/build/load"),
		Name:     "build.load",
		Usage:    "enables loading a docker image into the docker daemon post build",
	},
	&cli.DurationFlag{
		EnvVars:  []string{"PARAMETER_LOCAL_CACHE_TTL"},
		FilePath: string("/vela/parameters/makisu/build/local_cache_ttl,/vela/secrets/makisu/build/local_cache_ttl"),
		Name:     "build.local-cache-ttl",
		Usage:    "enables setting a time to live for the local docker cache (default 168h0m0s)",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_MODIFY_FS"},
		FilePath: string("/vela/parameters/makisu/build/modify_fs,/vela/secrets/makisu/build/modify_fs"),
		Name:     "build.modify-fs",
		Usage:    "enables setting makisu to modify files outside its internal storage directories",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_PERSERVE_ROOT"},
		FilePath: string("/vela/parameters/makisu/build/perserve_root,/vela/secrets/makisu/build/perserve_root"),
		Name:     "build.perserve-root",
		Usage:    "enables setting copying storage from root in the storage during and after build",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_PUSHES"},
		FilePath: string("/vela/parameters/makisu/build/pushes,/vela/secrets/makisu/build/pushes"),
		Name:     "build.pushes",
		Usage:    "enables setting registries to push the image to",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_REDIS_CACHE_OPTIONS"},
		FilePath: string("/vela/parameters/makisu/build/redis_cache_options,/vela/secrets/makisu/build/redis_cache_options"),
		Name:     "build.redis-cache-options",
		Usage:    "enables setting custom redis server for caching",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_REGISTRY_CONFIG"},
		FilePath: string("/vela/parameters/makisu/build/registry_config,/vela/secrets/makisu/build/registry_config"),
		Name:     "build.registry-config",
		Usage:    "enables setting registry configuration for authentication",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_REPLICAS"},
		FilePath: string("/vela/parameters/makisu/build/replicas,/vela/secrets/makisu/build/replicas"),
		Name:     "build.replicas",
		Usage:    "enables setting pushing image to alternative targets i.e. \"<registry>/<repo>:<tag>\"",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_STORAGE"},
		FilePath: string("/vela/parameters/makisu/build/storage,/vela/secrets/makisu/build/storage"),
		Name:     "build.storage",
		Usage:    "enables setting a directory for makisu to use for temp files and cached layers",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_TAG"},
		FilePath: string("/vela/parameters/makisu/build/tag,/vela/secrets/makisu/build/tag"),
		Name:     "build.tag",
		Usage:    "enables setting the tag for an image",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_TARGET"},
		FilePath: string("/vela/parameters/makisu/build/target,/vela/secrets/makisu/build/target"),
		Name:     "build.target",
		Usage:    "enables setting the target build stage to build",
	},
}

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

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (b *Build) Unmarshal() error {
	logrus.Trace("unmarshaling build options")

	// allocate configuration to structs
	b.Docker = &Docker{}
	b.HTTPCache = &HTTPCache{}
	b.RedisCache = &RedisCache{}

	// cast raw docker options into bytes
	dockerOpts := []byte(b.DockerRaw)

	// serialize raw docker options into expected Props type
	err := json.Unmarshal(dockerOpts, &b.Docker)
	if err != nil {
		return err
	}

	// cast raw http options into bytes
	httpOpts := []byte(b.HTTPCacheRaw)

	// serialize raw http options into expected Props type
	err = json.Unmarshal(httpOpts, &b.HTTPCache)
	if err != nil {
		return err
	}

	// cast raw http options into bytes
	redisOpts := []byte(b.RedisCacheRaw)

	// serialize raw http options into expected Props type
	err = json.Unmarshal(redisOpts, &b.RedisCache)
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
