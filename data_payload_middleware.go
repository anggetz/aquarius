package aquarius

import (
	"net/http"
	"strings"
)

type DataPayloadMiddleware struct {
	RegisteredHandler map[string]string
}

func NewDataPayloadMiddleware() DataPayloadMiddleware {
	return DataPayloadMiddleware{}
}

func (dataPayload *DataPayloadMiddleware) Interceptor(aqua *WebContext) bool {

	contentType := aqua.Request.Header.Get("Content-Type")

	aqua.Data = map[string]interface{}{}

	if strings.Contains(contentType, "application/json") {
		err := aqua.GetPayloadData(&aqua.Data)
		if err != nil {
			http.Error(aqua.Writer, "Could not parse json format", http.StatusInternalServerError)
			return false
		}
	}
	return true
}

func (dataPayload *DataPayloadMiddleware) BeforeRegisterHandler(aqua *WebContext) {

}
