# http middleware
Common go http server middlewares

# Example

```go
package main

import (
	middleware "github.com/pactus-project/gopkg/middleware/http-mdl"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	middleware.Chain(middleware.Logging(), middleware.Recover())(mux)
	sv := &http.Server{
		Handler: mux,
    }

	sv.ListenAndServe()
}
```