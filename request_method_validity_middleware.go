package aquarius

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestMethodValidity struct {
	RegisteredHandler map[string]string
}

func NewRequestMethodValidity() RequestMethodValidity {
	return RequestMethodValidity{}
}

func (methodValidity *RequestMethodValidity) Interceptor(aqua *WebContext) bool {

	if strings.HasPrefix(aqua.Method, "post_") {
		if aqua.Request.Method != "POST" {
			http.NotFoundHandler().ServeHTTP(aqua.Writer, aqua.Request)
			return false
		}

	} else if strings.HasPrefix(aqua.Method, "get_") {
		if aqua.Request.Method != "GET" {
			http.NotFoundHandler().ServeHTTP(aqua.Writer, aqua.Request)
			return false
		}
	}

	return true
}

func (methodValidity *RequestMethodValidity) BeforeRegisterHandler(aqua *WebContext) {
	secondUrl := ""
	if strings.HasPrefix(aqua.Method, "post_") {
		secondUrl = strings.Replace(aqua.Method, "post_", "", -1)

	} else if strings.HasPrefix(aqua.Method, "get_") {
		secondUrl = strings.Replace(aqua.Method, "get_", "", -1)

	} else {
		secondUrl = aqua.Method
	}

	aqua.Url = fmt.Sprintf("/%s/%s", aqua.Controller, secondUrl)
}
