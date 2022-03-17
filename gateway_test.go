package gateway_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/leonco/gateway"
	"github.com/tj/assert"
)

func Example() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World from Go")
}

func TestGateway_Invoke(t *testing.T) {

	e := []byte(`{"version": "1.0", "rawPath": "/pets/luna", "requestContext": {"http": {"method": "POST"}}}`)

	gw := gateway.NewGateway(http.HandlerFunc(hello))

	payload, err := gw.Invoke(context.Background(), e)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"isBase64Encoded":false,"statusCode":200,"headers":{"Content-Type":"text/plain; charset=utf8"},"body":"Hello World from Go\n"}`, string(payload))
}
