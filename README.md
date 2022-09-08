# Webby

A super simple web application which returns the IP of the host running Webby.

Written in [golang](https://golang.org/), it can be compiled for most major platforms.

# Build Webby
## Current architecture
```go build```

## Linux (if not current)
```GOOS=linux GOARCH=amd64 go build```

# Distribution

The supplied Dockerfile can be used to wrap the webby binary in a docker image -- for distribution/exection into the Docker ecosystem.

The resulting docker image is <7 Mb...so, you should be able to run 100s (or 1000s!) of copies of Webby

WEBBY ALL THE THINGS!

# Roadmap
* TBD
