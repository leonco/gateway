package gateway

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/tencentyun/scf-go-lib/events"
	"github.com/tj/assert"
)

func TestNewRequest_path(t *testing.T) {
	e := events.APIGatewayRequest{
		Path: "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "GET", r.Method)
	assert.Equal(t, `/pets/luna`, r.URL.Path)
	assert.Equal(t, `/pets/luna`, r.URL.String())
	assert.Equal(t, `/pets/luna`, r.RequestURI)
}

func TestNewRequest_method(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "DELETE",
		Path:   "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "DELETE", r.Method)
}

func TestNewRequest_queryString(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "GET",
		Path:   "/pets",
		QueryString: events.APIGatewayQueryString{
			"order":  []string{"desc"},
			"fields": []string{"name,species"},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `/pets?fields=name%2Cspecies&order=desc`, r.URL.String())
	assert.Equal(t, `desc`, r.URL.Query().Get("order"))
}

func TestNewRequest_multiValueQueryString(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "GET",
		Path:   "/pets",
		QueryString: events.APIGatewayQueryString{
			"multi_fields": []string{"name", "species"},
			"multi_arr[]":  []string{"arr1", "arr2"},
			"order":        []string{"desc"},
			"fields":       []string{"name,species"},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.URL.String())
	assert.Equal(t, []string{"name", "species"}, r.URL.Query()["multi_fields"])
	assert.Equal(t, []string{"arr1", "arr2"}, r.URL.Query()["multi_arr[]"])
	assert.Equal(t, `/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.RequestURI)
}

func TestNewRequest_remoteAddr(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "GET",
		Path:   "/pets",
		Context: events.APIGatewayRequestContext{
			SourceIP: "1.2.3.4",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `1.2.3.4`, r.RemoteAddr)
}

func TestNewRequest_header(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "POST",
		Path:   "/pets",
		Body:   `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
		Context: events.APIGatewayRequestContext{
			RequestID: "1234",
			Stage:     "prod",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `example.com`, r.Host)
	assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
}

func TestNewRequest_multiHeader(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "POST",
		Path:   "/pets",
		Body:   `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `example.com`, r.Host)
	assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
}

func TestNewRequest_body(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "POST",
		Path:   "/pets",
		Body:   `{ "name": "Tobi" }`,
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{ "name": "Tobi" }`, string(b))
}

func TestNewRequest_bodyBinary(t *testing.T) {
	e := events.APIGatewayRequest{
		Method: "POST",
		Path:   "/pets",
		Body:   `aGVsbG8gd29ybGQK`,
		Headers: map[string]string{
			"X-Body-Base64": "on",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, "hello world\n", string(b))
}

func TestNewRequest_context(t *testing.T) {
	e := events.APIGatewayRequest{}
	ctx := context.WithValue(context.Background(), "key", "value")
	r, err := NewRequest(ctx, e)
	assert.NoError(t, err)
	v := r.Context().Value("key")
	assert.Equal(t, "value", v)
}
