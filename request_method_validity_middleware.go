package aquarius

type RequestMethodValidity struct {
	RegisteredHandler map[string]string
}

func NewRequestMethodValidity() RequestMethodValidity {
	return RequestMethodValidity{}
}

func (methodValidity *RequestMethodValidity) Interceptor(aqua *WebContext) bool {
	return true
}
