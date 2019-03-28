package aquarius

import (
	"fmt"
	"strings"
)

type RequestMethodValidity struct {
	RegisteredHandler map[string]string
}

func NewRequestMethodValidity() RequestMethodValidity {
	return RequestMethodValidity{}
}

func (methodValidity *RequestMethodValidity) Interceptor(aqua *WebContext) bool {
	return true
}

func (methodValidity *RequestMethodValidity) BeforeRegisterHandler(aqua *WebContext) {
	secondUrl := ""
	if strings.HasPrefix(aqua.MethodFunc, "post_") {
		secondUrl = strings.Replace(aqua.MethodFunc, "post_", "", -1)
		aqua.Method = "POST"

	} else if strings.HasPrefix(aqua.MethodFunc, "get_") {
		secondUrl = strings.Replace(aqua.MethodFunc, "get_", "", -1)
		aqua.Method = "GET"
	} else {
		secondUrl = aqua.MethodFunc
		aqua.Method = "GET"
	}

	route := aqua.ControllerStruct.FieldByName("Route")
	if route.IsValid() {
		routeInterface := route.Interface()
		mapRoute := routeInterface.(map[string]interface{})
		if val, ok := mapRoute[aqua.PureMethodFunc]; ok {
			aqua.Url = val.(string)
			return
		}
	}
	aqua.Url = fmt.Sprintf("/%s/%s", aqua.Controller, secondUrl)

}
