package aquarius

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Application struct {
	Name              string
	Path              string
	Layout            string
	Describtion       string
	ViewsPath         string
	GlobalMiddlewares []interface{}
	InterceptorSignal int
}

const INTERCEPTOR_STOP = -1
const INTERCEPTOR_RUN = 1

func NewApplication(appName string, globalMiddleware []interface{}) Application {
	// Use global middleware?
	if globalMiddleware == nil {
		globalMiddleware = []interface{}{}
	}

	requestMethodValidity := NewRequestMethodValidity()
	dataPayloadMiddleware := NewDataPayloadMiddleware()
	globalMiddleware = append(globalMiddleware, &requestMethodValidity)
	globalMiddleware = append(globalMiddleware, &dataPayloadMiddleware)

	return Application{
		Name:              appName,
		GlobalMiddlewares: globalMiddleware,
		InterceptorSignal: INTERCEPTOR_RUN,
		ViewsPath:         "views",
		Layout:            "layout.html",
	}
}

func (app *Application) RegisterController(handlerStruct interface{}) {
	structValue := reflect.Indirect(reflect.ValueOf(handlerStruct))
	structToRegister := reflect.TypeOf(handlerStruct)

	for i := 0; i < structToRegister.NumMethod(); i++ {
		interceptorFuncs := []reflect.Value{}

		structName := strings.ToLower(strings.Replace(structToRegister.Elem().Name(), "Controller", "", -1))
		method := structToRegister.Method(i)

		url := fmt.Sprintf("/%s/%s", structName, strings.ToLower(method.Name))

		webContext := WebContext{
			AppInfo:          app,
			ControllerStruct: structValue,
			Controller:       structName,
			MethodFunc:       method.Name,
			Url:              url,
		}

		// Check middleware
		middlewareStruct := structValue.FieldByName("Middleware")
		if !middlewareStruct.IsValid() {
			middlewareStruct = reflect.ValueOf([]interface{}{})
		}

		for _, m := range webContext.AppInfo.GlobalMiddlewares {
			middlewareStruct = reflect.Append(middlewareStruct, reflect.ValueOf(m))
		}
		for n := 0; n < middlewareStruct.Len(); n++ {
			interceptorFunc := middlewareStruct.Index(n).Elem().MethodByName("Interceptor")
			if interceptorFunc.IsValid() {
				interceptorFuncs = append(interceptorFuncs, interceptorFunc)
			}

		}

		webContext.MethodValidity()

		fmt.Printf("[INFO] Register route %s method %s \n", webContext.Url, webContext.Method)

		Mux.HandleFunc(webContext.Url, func(w http.ResponseWriter, req *http.Request) {

			webContext.Writer = w
			webContext.Request = req

			fmt.Printf("[INFO] Incoming request %s \n", webContext.Url)

			for _, interceptor := range interceptorFuncs {
				returnValues := interceptor.Call([]reflect.Value{reflect.ValueOf(&webContext)})
				if !returnValues[0].Bool() {
					return
				}
			}
			//

			reflect.ValueOf(handlerStruct).MethodByName(webContext.MethodFunc).Call([]reflect.Value{reflect.ValueOf(&webContext)})

		}).Methods(webContext.Method)

	}
}
