// Package illegal_enum_names provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package illegal_enum_names

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Bar defines model for Bar.
type Bar string

// List of Bar
const (
	Bar_Bar      Bar = "Bar"
	Bar_Foo      Bar = "Foo"
	Bar_Foo_Bar  Bar = "Foo Bar"
	Bar_Foo_Bar1 Bar = "Foo-Bar"
	Bar__Foo     Bar = "1Foo"
	Bar__Foo1    Bar = " Foo"
	Bar__Foo_    Bar = " Foo "
	Bar__Foo_1   Bar = "_Foo_"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetFoo request
	GetFoo(ctx context.Context) (*http.Response, error)
}

func (c *Client) GetFoo(ctx context.Context) (*http.Response, error) {
	req, err := NewGetFooRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewGetFooRequest generates requests for GetFoo
func NewGetFooRequest(server string) (*http.Request, error) {
	var err error

	var queryUrl *url.URL
	queryUrl, err = url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/foo")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	req, err = http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetFoo request
	GetFooWithResponse(ctx context.Context) (*GetFooResponse, error)
}

type GetFooResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Bar
}

// Status returns HTTPResponse.Status
func (r GetFooResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetFooResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetFooWithResponse request returning *GetFooResponse
func (c *ClientWithResponses) GetFooWithResponse(ctx context.Context) (*GetFooResponse, error) {
	rsp, err := c.GetFoo(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetFooResponse(rsp)
}

// ParseGetFooResponse parses an HTTP response from a GetFooWithResponse call
func ParseGetFooResponse(rsp *http.Response) (*GetFooResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetFooResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Bar
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /foo)
	GetFoo(ctx GetFooContext) error
}

type GetFooContext struct {
	echo.Context
}

func (c *GetFooContext) JSON200(resp []Bar) error {
	return c.JSON(200, resp)
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler func(echo.Context) ServerInterface
}

// GetFoo converts echo context to params.
func (w *ServerInterfaceWrapper) GetFoo(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler(ctx).GetFoo(GetFooContext{ctx})
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface, pathPrefix string) {

	wrapper := ServerInterfaceWrapper{
		Handler: func(echo.Context) ServerInterface {
			return si
		},
	}
	wrapper.RegisterHandlers(router, pathPrefix)

}

func (wrapper ServerInterfaceWrapper) RegisterHandlers(router EchoRouter, pathPrefix string) {
	router.GET(path.Join(pathPrefix, "/foo"), wrapper.GetFoo)

}

//go:embed spec.yaml
var spec []byte

// returns a raw spec
func RawSpec() []byte {
	return spec
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathPrefix string) map[string]func() []byte {
	// todo: fix spec validator so that external references are correct;
	// now they can point to api.yaml files whereas the real file name is different
	var res = map[string]func() []byte{
		path.Join(pathPrefix, "spec.yaml"): RawSpec,
		path.Join(pathPrefix, "api.yaml"):  RawSpec,
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.Swagger, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.SwaggerLoader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		if spec, ok := resolvePath[pathToFile]; !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		} else {
			return spec(), nil
		}
	}
	swagger, err = loader.LoadSwaggerFromData(spec)
	if err != nil {
		return
	}
	return
}
