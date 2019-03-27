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
	globalMiddleware = append(globalMiddleware, &requestMethodValidity)

	return Application{Name: appName, GlobalMiddlewares: globalMiddleware, InterceptorSignal: INTERCEPTOR_RUN}
}

func (app *Application) RegisterController(handlerStruct interface{}) {
	structValue := reflect.Indirect(reflect.ValueOf(handlerStruct))
	structToRegister := reflect.TypeOf(handlerStruct)

	for i := 0; i < structToRegister.NumMethod(); i++ {
		interceptorFuncs := []reflect.Value{}

		structName := strings.ToLower(strings.Replace(structToRegister.Elem().Name(), "Controller", "", -1))
		method := structToRegister.Method(i)

		fmt.Printf("[INFO] Register route /%s/%s \n", structName, strings.ToLower(method.Name))

		url := fmt.Sprintf("/%s/%s", structName, method)

		webContext := WebContext{
			AppInfo:    app,
			Controller: structName,
			Method:     strings.ToLower(method.Name),
			Url:        url,
		}

		// Check middleware
		middlewareStruct := structValue.FieldByName("Middleware")
		if middlewareStruct.IsValid() {
			for _, m := range webContext.AppInfo.GlobalMiddlewares {
				middlewareStruct = reflect.Append(middlewareStruct, reflect.ValueOf(m))
			}
			for n := 0; n < middlewareStruct.Len(); n++ {
				interceptorFunc := middlewareStruct.Index(n).Elem().MethodByName("Interceptor")
				if interceptorFunc.IsValid() {
					interceptorFuncs = append(interceptorFuncs, interceptorFunc)
				}

				beforeRegisterHandler := middlewareStruct.Index(n).Elem().MethodByName("BeforeRegisterHandler")
				if beforeRegisterHandler.IsValid() {
					beforeRegisterHandler.Call([]reflect.Value{reflect.ValueOf(&webContext)})
				}
			}
		}

		http.HandleFunc(webContext.Url, func(w http.ResponseWriter, req *http.Request) {

			webContext.Writer = w
			webContext.Request = req

			fmt.Printf("[INFO] Incoming request /%s/%s \n", structName, strings.ToLower(method.Name))

			for _, interceptor := range interceptorFuncs {
				returnValues := interceptor.Call([]reflect.Value{reflect.ValueOf(&webContext)})
				if !returnValues[0].Bool() {
					return
				}
			}

			reflect.ValueOf(handlerStruct).MethodByName(method.Name).Call([]reflect.Value{reflect.ValueOf(&webContext)})
		})
	}
}
