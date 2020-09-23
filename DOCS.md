## Description

This plugin enables you to build and publish [Docker](https://www.docker.com/) images in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-makisu

Registry: https://hub.docker.com/r/target/vela-makisu

## Usage

_Notes:_

* The plugin supports reading all parameters via environment variables or files. Values set as a file take precedence over default values set from the environment.
* We do not recommended using latest for pipelines. Users should use pinned images to decrease volatility of external changes to their pipelines. 

Sample of building and publishing an image:

```yaml
steps:
  - name: publish_hello-world
    image: target/vela-makisu:latest
    parameters:
      registry: index.docker.io
      tag: index.docker.io/octocat/hello-world
```

Sample of building an image without publishing:

```diff
steps:
  - name: publish hello world
    image: target/vela-makisu:latest
    parameters:
+     dry_run: true
      registry: index.docker.io
      tag: index.docker.io/octocat/hello-world:latest
```

Sample of building and publishing an image with custom tags:

```diff
steps:
  - name: publish hello world
    image: target/vela-makisu:latest
    parameters:
      registry: index.docker.io
      tag: index.docker.io/octocat/hello-world:latest
+     replicas:
+       - index.docker.io/octocat/hello-world:1
+       - index.docker.io/octocat/hello-world:foobar
```

Sample of building and publishing an image with build arguments:

```diff
steps:
  - name: publish hello world
    image: target/vela-makisu:latest
    pull: true
    parameters:
+     build_args:
+       - FOO=bar
      registry: index.docker.io
      tag: index.docker.io/octocat/hello-world
```

Sample of building and publishing an image with redis caching:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-makisu:latest
    pull: true
    parameters:
+     redis_cache_options: 
+       addr: redis.company.com
+       password: superSecretPassword
+       ttl: 7d
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

## Secrets

**NOTE: Users should refrain from configuring sensitive information in your pipeline in plain text.**

You can use Vela secrets to substitute sensitive values at runtime:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-makisu:latest
    pull: true
+   secrets: [ docker_username, docker_password, redis_cache ]
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
-     redis_cache: 
-       addr: redis.company.com
-       password: superSecretPassword
-       ttl: 7d      
-     username: octocat
-     password: superSecretPassword
```

## Parameters

**NOTE: Vela injects several variables, by default, that this plugin can load in automatically.**

The following parameters are used to configure the build and push process:

| Name              | Description                                                          | Required | Default |
| ----------------- | -------------------------------------------------------------------- | -------- | ------- |
| `build_args`      | build time arguments for the Dockerfile                              | `false`  | `N/A`   |
| `commit`          | commit info for #!COMMIT annotations                                 | `false`  | `N/A`   |
| `compression`     | compression on the tar file built - options: (no|speed|size|default) | `false`  | `N/A`   |
| `context`         | the context for the image to be built                                | `false`  | `.`     |
| `deny_list`       | list of locations to be ignored within docker image                  | `false`  | `N/A`   |
| `docker`          | configuration on the docker daemon                                   | `false`  | `N/A`   |
| `destination`     | the output of the tar file                                           | `false`  | `N/A`   |
| `file`            | a the absolute path to dockerfile                                    | `false`  | `info`  |
| `http_cache`      | custom http options caching                                          | `false`  | `N/A`   |
| `load`            | enables loading a docker image into the docker daemon post build     | `false`  | `N/A`   |
| `local_cache_ttl` | a time to live for the local docker cache (default 168h0m0s)         | `false`  | `N/A`   |
| `modify_fs`       | makisu to modify files outside its internal storage directories      | `false`  | `N/A`   |
| `perserve_root`   | copying storage from root in the storage during and after build      | `false`  | `N/A`   |
| `pushes`          | registries to push the image to                                      | `false`  | `N/A`   |
| `redis_cache`     | custom redis server for caching                                      | `false`  | `N/A`   |
| `replicas`        | pushing image to alternative targets i.e. `<registry>/<repo>:<tag>`  | `false`  | `N/A`   |
| `storage`         | a directory for makisu to use for temp files and cached layers       | `false`  | `N/A`   |
| `tag`             | the tag for an image                                                 | `true`   | `N/A`   |
| `storage`         | the target build stage to build                                      | `false`  | `N/A`   |

The following parameters are used to configure the registry:

| Name            | Description                                                        | Required | Default           |
| --------------- | ------------------------------------------------------------------ | -------- | ----------------- |
| `mirror`        | name of the mirror registry to use                                 | `false`  | `N/A`             |
| `password`      | password for communication with the registry                       | `true`   | `N/A`             |
| `registry`      | name of the registry for the repository                            | `true`   | `index.docker.io` |
| `repo`          | name of the repository for the image                               | `true`   | `N/A`             |
| `username`      | user name for communication with the registry                      | `true`   | `N/A`             |

## Template

COMING SOON!

## Troubleshooting

Below are a list of common problems and how to solve them:
