# http content tester
This commandline application makes http requests and prints the address of the request along with the
MD5 hash of the response to STDOUT.

## Installation
Install [Golang 1.17](https://golang.org/doc/install) or higher.

Run `go install github.com/TriAnMan/http-test@latest` and ``cd `go env GOPATH`/bin``

## Run
`./http-test www.yandex.com https://google.com`

Consider `-parallel` parameter to control number of parallel requests (default is 10).
`./http-test -parallel 3 google.com facebook.com yahoo.com yandex.com twitter.com
reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com`

## Design decisions
1. The app exploits the fail-fast methodology.
2. Streams are used to minimize a memory footprint.
3. Pipelines concurrency pattern is utilized to empower the code readability and manageability.
4. App logs to STDERR to facilitate log management (https://12factor.net/logs) and to separate logs from normal output.

## Possible future improvements
1. Parse STDIN instead of command line to overcome various OS limitations for command line environments.
    