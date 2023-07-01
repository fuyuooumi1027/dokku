# Herokuish Buildpacks

> Subcommands new as of 0.15.0

```
buildpacks:add [--index 1] <app> <buildpack>            # Add new app buildpack while inserting into list of buildpacks if necessary
buildpacks:clear <app>                                  # Clear all buildpacks set on the app
buildpacks:list <app>                                   # List all buildpacks for an app
buildpacks:remove <app> <buildpack>                     # Remove a buildpack set on the app
buildpacks:report [<app>] [<flag>]                      # Displays a buildpack report for one or more apps
buildpacks:set [--index 1] <app> <buildpack>            # Set new app buildpack at a given position defaulting to the first buildpack if no index is specified
buildpacks:set-property [--global|<app>] <key> <value>  # Set or clear a buildpacks property for an app
```

```
builder-herokuish:report [<app>] [<flag>]   # Displays a builder-herokuish report for one or more apps
builder-herokuish:set <app> <key> (<value>) # Set or clear a builder-herokuish property for an app
```

> Warning: If using the `buildpacks` plugin, be sure to unset any `BUILDPACK_URL` and remove any such entries from a committed `.env` file. A specified `BUILDPACK_URL` will always override a `.buildpacks` file or the buildpacks plugin.

Dokku normally defaults to using [Heroku buildpacks](https://devcenter.heroku.com/articles/buildpacks) for deployment, though this may be overridden by committing a valid `Dockerfile` to the root of your repository and pushing the repository to your Dokku installation. To avoid this automatic `Dockerfile` deployment detection, you may do one of the following:

- Set a `BUILDPACK_URL` environment variable
  - This can be done via `dokku config:set` or via a committed `.env` file in the root of the repository. See the [environment variable documentation](/docs/configuration/environment-variables.md) for more details.
- Create a `.buildpacks` file in the root of your repository.
  - This can be via a committed `.buildpacks` file or managed via the `buildpacks` plugin commands.

This page will cover usage of the `buildpacks` plugin.

## Usage

### Detection

This builder will be auto-detected in either the following cases:

- The `BUILDPACK_URL` app environment variable is set.
  - This can be done via `dokku config:set` or via a committed `.env` file in the root of the repository. See the [environment variable documentation](/docs/configuration/environment-variables.md) for more details.
- A `.buildpacks` file exists in the root of the app repository.
  - This can be via a committed `.buildpacks` file or managed via the `buildpacks` plugin commands.

The builder can also be specified via the `builder:set` command:

```shell
dokku builder:set node-js-app selected herokuish
```

> Dokku will only select the `dockerfile` builder if both the `herokuish` and `pack` builders are not detected and a Dockerfile exists. See the [dockerfile builder documentation](/docs/deployment/builders/dockerfiles.md) for more information on how that builder functions.

### Listing Buildpacks in Use

The `buildpacks:list` command can be used to show buildpacks that have been set for an app. This will omit any auto-detected buildpacks.

```shell
# running for an app with no buildpacks specified
dokku buildpacks:list node-js-app
```

```
-----> test buildpack urls
```

```shell
# running for an app with two buildpacks specified
dokku buildpacks:list node-js-app
```

```
-----> test buildpack urls
       https://github.com/heroku/heroku-buildpack-python.git
       https://github.com/heroku/heroku-buildpack-nodejs.git
```

### Adding custom buildpacks

> Please check the documentation for your particular buildpack as you may need to include configuration files (such as a Procfile) in your project root.

To add a custom buildpack, use the `buildpacks:add` command:

```shell
dokku buildpacks:add node-js-app https://github.com/heroku/heroku-buildpack-nodejs.git
```

When no buildpacks are currently specified, the specified buildpack will be the only one executed for detection and compilation.

Multiple buildpacks may be specified by using the `buildpacks:add` command multiple times.

```shell
dokku buildpacks:add node-js-app https://github.com/heroku/heroku-buildpack-ruby.git
dokku buildpacks:add node-js-app https://github.com/heroku/heroku-buildpack-nodejs.git
```

Buildpacks are executed in order, may be inserted at a specified index via the `--index` flag. This flag is specified starting at a 1-index value.

```shell
# will add the golang buildpack at the second position, bumping all proceeding ones by 1 position
dokku buildpacks:add --index 2 node-js-app https://github.com/heroku/heroku-buildpack-golang.git
```

### Overwriting a buildpack position

In some cases, it may be necessary to swap out a given buildpack. Rather than needing to re-specify each buildpack, the `buildpacks:set` command can be used to overwrite a buildpack at a given position.

```shell
dokku buildpacks:set node-js-app https://github.com/heroku/heroku-buildpack-ruby.git
```

By default, this will overwrite the _first_ buildpack specified. To specify an index, the `--index` flag may be used. This flag is specified starting at a 1-index value, and defaults to `1`.

```shell
# the following are equivalent commands
dokku buildpacks:set node-js-app https://github.com/heroku/heroku-buildpack-ruby.git
dokku buildpacks:set --index 1 node-js-app https://github.com/heroku/heroku-buildpack-ruby.git
```

If the index specified is larger than the number of buildpacks currently configured, the buildpack will be appended to the end of the list.

```shell
dokku buildpacks:set --index 99 node-js-app https://github.com/heroku/heroku-buildpack-ruby.git
```

### Removing a buildpack

> At least one of a buildpack or index must be specified

A single buildpack can be removed by name via the `buildpacks:remove` command.

```shell
dokku buildpacks:remove node-js-app https://github.com/heroku/heroku-buildpack-ruby.git
```

Buildpacks can also be removed by index via the `--index` flag. This flag is specified starting at a 1-index value.

```shell
dokku buildpacks:remove node-js-app --index 1
```

### Clearing all buildpacks

> This does not affect automatically detected buildpacks, nor does it impact any specified `BUILDPACK_URL` environment variable.

The `buildpacks:clear` command can be used to clear all configured buildpacks for a specified app.

```shell
dokku buildpacks:clear node-js-app
```

### Customizing the Buildpack stack builder

> New as of 0.23.0

The default stack builder in use by Herokuish buildpacks in Dokku is based on `gliderlabs/herokuish:latest`. Typically, this is installed via an OS package which pulls the requisite Docker image. Users may desire to switch the stack builder to a custom version, either to update the operating system or to customize packages included with the stack builder. This can be performed via the `buildpacks:set-property` command.

```shell
dokku buildpacks:set-property node-js-app stack gliderlabs/herokuish:latest
```

The specified stack builder can also be unset by omitting the name of the stack builder when calling `buildpacks:set-property`.

```shell
dokku buildpacks:set-property node-js-app stack
```

A change in the stack builder value will execute the `post-stack-set` trigger.

Finally, stack builders can be set or unset globally as a fallback. This will take precedence over a globally set `DOKKU_IMAGE` environment variable (`gliderlabs/herokuish:latest-20` by default).

```shell
# set globally
dokku buildpacks:set-property --global stack gliderlabs/herokuish:latest

# unset globally
dokku buildpacks:set-property --global stack
```

### Allowing herokuish for non-amd64 platforms

> New as of 0.29.0

By default, the builder-herokuish plugin is not enabled for non-amd64 platforms, and attempting to use it is blocked. This is because the majority of buildpacks are not cross-platform compatible, and thus building apps will either be considerably slower - due to emulating the amd64 platform - or won't work - due to building amd64 packages on arm/arm64 platforms.

To force-enable herokuish on non-amd64 platforms, the `allowed` property can be set via `builder-herokuish:set`. The default value depends on the host platform architecture (`true` on amd64, `false` otherwise).

```shell
dokku builder-herokuish:set node-js-app allowed true
```

The default value may be set by passing an empty value for the option:

```shell
dokku builder-herokuish:set node-js-app allowed
```

The `allowed` property can also be set globally. The global default is platform-dependent, and the global value is used when no app-specific value is set.

```shell
dokku builder-herokuish:set --global allowed true
```

The default value may be set by passing an empty value for the option.

```shell
dokku builder-herokuish:set --global allowed
```

### Displaying buildpack reports for an app

You can get a report about the app's buildpacks status using the `buildpacks:report` command:

```shell
dokku buildpacks:report
```

```
=====> node-js-app buildpacks information
       Buildpacks computed stack:  gliderlabs/herokuish:v0.5.23-20
       Buildpacks global stack:    gliderlabs/herokuish:latest-20
       Buildpacks list:            https://github.com/heroku/heroku-buildpack-nodejs.git
       Buildpacks stack:           gliderlabs/herokuish:v0.5.23-20
=====> python-sample buildpacks information
       Buildpacks computed stack:  gliderlabs/herokuish:latest-20
       Buildpacks global stack:    gliderlabs/herokuish:latest-20
       Buildpacks list:            https://github.com/heroku/heroku-buildpack-nodejs.git,https://github.com/heroku/heroku-buildpack-python.git
       Buildpacks stack:
=====> ruby-sample buildpacks information
       Buildpacks computed stack:  gliderlabs/herokuish:latest-20
       Buildpacks global stack:    gliderlabs/herokuish:latest-20
       Buildpacks list:
       Buildpacks stack:
```

You can run the command for a specific app also.

```shell
dokku buildpacks:report node-js-app
```

```
=====> node-js-app buildpacks information
       Buildpacks list:               https://github.com/heroku/heroku-buildpack-nodejs.git
```

You can pass flags which will output only the value of the specific information you want. For example:

```shell
dokku buildpacks:report node-js-app --buildpacks-list
```

### Displaying builder-herokuish reports for an app

> New as of 0.29.0

You can get a report about the app's storage status using the `builder-herokuish:report` command:

```shell
dokku builder-herokuish:report
```

```
=====> node-js-app builder-herokuish information
       Builder herokuish computed allowed: false
       Builder herokuish global allowed:   true
       Builder herokuish allowed:          false
=====> python-sample builder-herokuish information
       Builder herokuish computed allowed: true
       Builder herokuish global allowed:   true
       Builder herokuish allowed:
=====> ruby-sample builder-herokuish information
       Builder herokuish computed allowed: true
       Builder herokuish global allowed:   true
       Builder herokuish allowed:
```

You can run the command for a specific app also.

```shell
dokku builder-herokuish:report node-js-app
```

```
=====> node-js-app builder-herokuish information
       Builder herokuish computed allowed: false
       Builder herokuish global allowed:   true
       Builder herokuish allowed:          false
```

You can pass flags which will output only the value of the specific information you want. For example:

```shell
dokku builder-herokuish:report node-js-app --builder-herokuish-allowed
```

```
false
```

## Errata

### Switching from Dockerfile deployments

If an application was previously deployed via Dockerfile, the following commands should be run before a buildpack deploy will succeed:

```shell
dokku config:unset --no-restart node-js-app DOKKU_PROXY_PORT_MAP
```

### Using a specific buildpack version

> Always remember to pin your buildpack versions when using the multi-buildpacks method, or you may find deploys changing your deployed environment.

By default, Dokku uses the [gliderlabs/herokuish](https://github.com/gliderlabs/herokuish/) project, which pins all of it's vendored buildpacks. There may be occasions where the pinned version results in a broken deploy, or does not have a particular feature that is required to build your project. To use a more recent version of a given buildpack, the buildpack may be specified _without_ a Git commit SHA like so:

```shell
# using the latest nodejs buildpack
dokku buildpacks:set node-js-app https://github.com/heroku/heroku-buildpack-nodejs
```

This will use the latest commit on the `master` branch of the specified buildpack. To pin to a newer version of a buildpack, a sha may also be specified by using the form `REPOSITORY_URL#COMMIT_SHA`, where `COMMIT_SHA` is any tree-ish git object - usually a git tag.

```shell
# using v87 of the nodejs buildpack
dokku buildpacks:set node-js-app https://github.com/heroku/heroku-buildpack-nodejs#v87
```

### `curl` build timeouts

Certain buildpacks may time out in retrieving dependencies via `curl`. This can happen when your network connection is poor or if there is significant network congestion. You may see a message similar to `gzip: stdin: unexpected end of file` after a `curl` command.

If you see output similar this when deploying , you may need to override the `curl` timeouts to increase the length of time allotted to those tasks. You can do so via the `config` plugin:

```shell
dokku config:set --global CURL_TIMEOUT=1200
dokku config:set --global CURL_CONNECT_TIMEOUT=180
```

### Clearing buildpack cache

See the [repository management documentation](/docs/advanced-usage/repository-management.md#clearing-app-cache) for more information on how to clear buildpack build cache for an application.

### Specifying commands via Procfile

See the [Procfile documentation](/docs/processes/process-management.md#procfile) for more information on how to specify different processes for your app.
