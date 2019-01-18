package webhook

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/giantswarm/auto-oncall/server/middleware"
	"github.com/giantswarm/auto-oncall/service"
)

const (
	// Method is the HTTP method this endpoint is registered for.
	Method = "POST"
	// Name identifies the endpoint. It is aligned to the package path.
	Name = "webhook"
	// Path is the HTTP request path this endpoint is registered for.
	Path = "/webhook"
)

// Config represents the configuration used to create a version endpoint.
type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

// New creates a new configured version endpoint.
func New(config Config) (*Endpoint, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}
	if config.Service == nil {
		return nil, microerror.Maskf(invalidConfigError, "service must not be empty")
	}

	newEndpoint := &Endpoint{
		Config: config,
	}

	return newEndpoint, nil
}

type Endpoint struct {
	Config
}

func (e *Endpoint) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return r, nil
	}
}

func (e *Endpoint) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		endpointResponse := response.(*Response)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(endpointResponse.StatusCode)
		return json.NewEncoder(w).Encode(endpointResponse.Body)
	}
}

func (e *Endpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		r := request.(*http.Request)

		response := DefaultResponse()

		h, err := e.Service.Webhook.NewHook(r)
		if err != nil {
			response.Body.Message = err.Error()
			response.StatusCode = http.StatusBadRequest
		} else {
			response.Body.Message = "webhook request received"
			response.StatusCode = http.StatusOK
		}

		go e.Service.Webhook.Process(h)

		return response, nil
	}
}

func (e *Endpoint) Method() string {
	return Method
}

func (e *Endpoint) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

func (e *Endpoint) Name() string {
	return Name
}

func (e *Endpoint) Path() string {
	return Path
}
