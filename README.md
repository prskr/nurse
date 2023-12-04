# Nurse

## Usage

Nurse comes currently with 2 different operation modes:

- server
- CLI

The server starts an HTTP server with configurable endpoints which you can use e.g. in Kubernetes environments to
distinguish between:

- startup
- readiness
- liveness

probes.
Every endpoint has a distinguished set of checks that are executed when you hit the endpoint.
Currently, there is no caching in place (and there are also no plans to change that).

The CLI operation mode on the other hand executes all checks that are provided as arguments e.g. in Docker Swarm
environment where the container image has to ship the health check CLI.

### Primer about checks

All checks are executed **in parallel** which means you shouldn't rely on a certain execution

### Global config/options

Nurse comes with the following global options:

| Switch             | Environment variable   | Default value                                                | Description                                                         |
|--------------------|------------------------|--------------------------------------------------------------|---------------------------------------------------------------------|
| `--config`         | `NURSE_CONFIG`         | `$HOME/.nurse.yaml`, `/etc/nurse/config.yaml`,`./nurse.yaml` | path to the config file                                             |
| `--check-timeout`  | `NURSE_CHECK_TIMEOUT`  | `500ms`                                                      | Timeout for executing all checks                                    |
| `--check-attempts` | `NURSE_CHECK_ATTEMPTS` | `20`                                                         | How often checks should be retried before they're considered failed |
| `--log.level`      |                        | `info`                                                       | Default log level                                                   |
| `--servers`        | `NURSE_SERVER_<name>`  |                                                              | Configure server URLs via environment variables                     |

The individual sub-commands come with additional options, like for example configuring endpoints via environment
variables as well.

The [nurse.yaml](./nurse.yaml) describes how to configure Nurse via a configuration file.

The most interesting root nodes are:

- servers
- endpoints

Within `servers` you can configure different servers for further usage in checks.
For example, to configure a Redis server: `redis://localhost:6379/0`.
Depending on the individual protocols there are further configuration options.

Within `endpoints` you can configure different HTTP endpoints the server exposes and which checks should be executed for
which endpoint.

### Server

The `server` sub-command comes with the following additional config options:

| Switch                      | Environment variable             | Default value | Description                                                |
|-----------------------------|----------------------------------|---------------|------------------------------------------------------------|
| `--endpoints`               | `NURSE_ENDPOINT_<name>`          |               | Configure HTTP endpoints via environment variables         |
| `--http.address`            | `NURSE_HTTP_ADDRESS`             | `:8080`       | IP and port the server will be listening on                |
| `--http.read-header-timout` | `NURSE_HTTP_READ_HEADER_TIMEOUT` | `100ms`       | Timeout until when the client has to have sent the headers |

To configure an endpoint via an environment variable, set it like this:

```
NURSE_ENDPOINT_HEALTHZ='http.GET("https://api.chucknorris.io/jokes/random")=>Status(200);redis.PING("local-redis")'
```

The server will print the configured routes when it is starting up.
In the aforementioned case you should see something like:

```
{"time":"xxxxx","level":"INFO","msg":"Configuring route","route":"/healthz"}
```

Multiple checks can be configured by separating them with a `;` into multiple 'expressions'.

### CLI

The CLI has no additional config options compared to the server.
It simply takes all arguments you pass to it, tries to parse them as checks and executes them with the given time limit.
If one of the check fails it will exit with a non-zero exit code.

Multiple checks can either be passed as single argument in `''` separated with a `;` just like in the environment variables, or you can pass multiple arguments.
The result will be the same.