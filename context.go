package gateway

import (
	"context"

	"github.com/tencentyun/scf-go-lib/events"
)

// key is the type used for any items added to the request context.
type key int

// requestContextKey is the key for the api gateway proxy `RequestContext`.
const requestContextKey key = iota

// newContext returns a new Context with specific api gateway proxy values.
func newContext(ctx context.Context, e events.APIGatewayRequest) context.Context {
	return context.WithValue(ctx, requestContextKey, e.Context)
}

// RequestContext returns the APIGatewayRequestContext value stored in ctx.
func RequestContext(ctx context.Context) (events.APIGatewayRequestContext, bool) {
	c, ok := ctx.Value(requestContextKey).(events.APIGatewayRequestContext)
	return c, ok
}
