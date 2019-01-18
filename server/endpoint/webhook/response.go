package webhook

// Response is a struct that represents what this endpoint returns.
type Response struct {
	Body       Body
	StatusCode int `json:"-"`
}

type Body struct {
	Message string `json:"message"`
}

// DefaultResponse returns empty Response.
func DefaultResponse() *Response {
	return &Response{
		Body: Body{
			Message: "",
		},
	}
}
