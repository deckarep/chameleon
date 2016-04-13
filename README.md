# Chameleon
A small service we can use to emulate healthchecks for other services.

# Usage
By default, chameleon will start a webserver on port 8080. It only responds to GET requests made to `/healthcheck`. It always returns a 200, and will respond with a block of JSON that matches our (sendgrid's) healthcheck v1 format.

If chameleon receives a SIGUSR1, it will flip the "healthiness" of the service.

Here's an example:
```
[~] curl -s localhost:8080/healthcheck | jq .
{
...
  "results": {
    "chameleon": {
      "message": "null",
      "ok": false
    }
  }
}

[~] kill -s SIGUSR1 81040
[~] curl -s localhost:8080/healthcheck | jq .
{
...
  "results": {
    "chameleon": {
      "message": "null",
      "ok": true
    }
  }
}
```

And here's the chameleon logs for this.
```
[~] ./chameleon                                                                                                                                                                            master  ✭ ✱
Running as pid 81040
Masquerade: chameleon
Port: 8080
Healthy: false
MinWait: 0
MaxWait: 0
log1kal-2.local
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /healthcheck              --> main.func·002 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2016/04/13 - 14:48:10 | 200 |     144.916µs | ::1 |   GET     /healthcheck
Received USR1 signal
Masquerade: chameleon
Port: 8080
Healthy: false
MinWait: 0
MaxWait: 0
[GIN] 2016/04/13 - 14:48:32 | 200 |      40.646µs | ::1 |   GET     /healthcheck
```

# Configuration
Chameleon's configuration is environment based, and accepts the following ENV variables.

| Variable | Description | Type |
|----------|-------------|------|
| `CHAMELEON_MASQUERADE` | the service name to masquerade as | string |
| `CHAMELEON_PORT` | the port to listen on | integer |
| `CHAMELEON_HEALTHY` | Is the service healthy? | boolean |
| `CHAMELEON_MINWAIT` | unused | |
| `CHAMELEON_MAXWAIT` | unused | |
| `CHAMELEON_VERSION` | What version of the service is runing? | string |
