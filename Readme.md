

Package gateway provides a drop-in replacement for net/http's `ListenAndServe` for use in [Tencent SCF](https://cloud.tencent.com/product/scf) & [API Gateway](https://cloud.tencent.com/product/apigw), simply swap it out for `gateway.ListenAndServe`. 


# Installation

```
go get github.com/leonco/gateway
```

# Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/leonco/gateway"
)

func main() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	// example retrieving values from the api gateway proxy request context.
	requestContext, ok := gateway.RequestContext(r.Context())
	if !ok || requestContext.Authorizer["sub"] == nil {
		fmt.Fprint(w, "Hello World from Go")
		return
	}

	userID := requestContext.Authorizer["sub"].(string)
	fmt.Fprintf(w, "Hello %s from Go", userID)
}
```

---

