// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-vela/vela-makisu/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-makisu"
	app.HelpName = "vela-makisu"
	app.Usage = "Vela img plugin for building and publishing images"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = version.New().Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "VELA_LOG_LEVEL", "MAKISU_LOG_LEVEL"},
			FilePath: string("/vela/parameters/makisu/log_level,/vela/secrets/makisu/log_level"),
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},
	}

	// add build flags
	app.Flags = append(app.Flags, buildFlags...)

	// add config flags
	app.Flags = append(app.Flags, configFlags...)

	// add global flags
	app.Flags = append(app.Flags, globalFlags...)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
	// capture the version information as pretty JSON
	v, err := json.MarshalIndent(version.New(), "", "  ")
	if err != nil {
		return err
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(v))

	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-makisu",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/makisu",
		"registry": "https://hub.docker.com/r/target/vela-makisu",
	}).Info("Vela Makisu Plugin")

	// create the plugin
	p := Plugin{
		Build: &Build{
			BuildArgs:      c.StringSlice("build.build-args"),
			Commit:         c.String("build.commit"),
			Compression:    c.String("build.compression"),
			Context:        c.String("build.context"),
			DenyList:       c.StringSlice("build.deny-list"),
			DockerRaw:      c.String("build.docker-options"),
			Destination:    c.String("build.destination"),
			File:           c.String("build.file"),
			HTTPCacheRaw:   c.String("build.http-cache-options"),
			Load:           c.Bool("build.load"),
			LocalCacheTTL:  c.Duration("build.local-cache-ttl"),
			ModifyFS:       c.Bool("build.modify-fs"),
			PreserveRoot:   c.Bool("build.perserve-root"),
			Pushes:         c.StringSlice("build.pushes"),
			RedisCacheRaw:  c.String("build.redis-cache-options"),
			RegistryConfig: c.String("build.registry-config"),
			Replicas:       c.StringSlice("build.replicas"),
			Storage:        c.String("build.storage"),
			Tag:            c.String("build.tag"),
			Target:         c.String("build.target"),
		},
		GlobalRaw: c.String("global.flags"),
		Registry: &Registry{
			Mirror:   c.String("registry.mirror"),
			Name:     c.String("registry.name"),
			Password: c.String("registry.password"),
			Username: c.String("registry.username"),
		},
	}

	// validate the plugin
	err = p.Validate()
	if err != nil {
		return err
	}

	// unmarshal the plugin
	err = p.Unmarshal()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
