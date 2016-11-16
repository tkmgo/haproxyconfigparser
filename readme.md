# haproxyconfigparser

HAProxcy config parser for Golang, but under development.


## Behavior

- It parses HAProxy config and binds to golang struct.
- Analysis ACL, make its reference beteen frontend and backends.


## Usage

main.go

```go
package main

import (
	"fmt"
	"github.com/tkmgo/haproxyconfigparser"
)

func main() {
  config, _ := haproxyconfigparser.ParseFromStdin()
	fmt.Println(config)
}
```

```shell
$ cat haproxy.cfg | go run main.go
```
