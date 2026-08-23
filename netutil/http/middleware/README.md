# HTTP middleware

Common go http server middlewares

# Example

```go
package main

import (
	"github.com/pactus-project/gopkg/netutil/http/middleware"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	sv := &http.Server{
		Handler: middleware.Chain(
			middleware.Logging(),
			middleware.Recover())(mux),
	}

	sv.ListenAndServe()
}
```