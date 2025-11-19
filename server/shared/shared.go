package shared

type Meta map[string]any //@name Meta

type APIResponse struct {
	Errors  []*APIError `json:"errors,omitempty" swaggerignore:"true"`
	Data    any         `json:"data,omitempty" swaggerignore:"true"`
	Message string      `json:"message,omitempty" swaggerignore:"true"`
	Meta    Meta        `json:"meta,omitempty" swaggerignore:"true"`
} //@name Response

// https://jsonapi.org/examples/#error-objects
type APIError struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Meta   Meta   `json:"meta,omitempty" swaggertype:"object,string" example:"key:value,key2:value2"`
} //@name Error

func (a *APIError) Error() string {
	return a.Detail
}
