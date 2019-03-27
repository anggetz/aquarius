package aquarius

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

type Aquarius struct {
	Server    *http.Server
	MuxRouter *mux.Router
	Port      string
	AppName   string
}

func NewAquarius() Aquarius {

	return Aquarius{
		MuxRouter: mux.NewRouter(),
	}
}
func (Aqua *Aquarius) Listen() {

	Aqua.Server = &http.Server{Addr: fmt.Sprintf(":%s", Aqua.Port), Handler: Aqua.MuxRouter}

	if err := Aqua.Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func (Aqua *Aquarius) StopServer() {
	if err := Aqua.Server.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func (Aqua *Aquarius) RegisterApp(app Application) error {
	for _, handlerStruct := range app.Controllers {
		structValue := reflect.Indirect(reflect.ValueOf(handlerStruct))
		structToRegister := reflect.TypeOf(handlerStruct)

		for i := 0; i < structToRegister.NumMethod(); i++ {
			interceptorFuncs := []reflect.Value{}

			structName := strings.ToLower(strings.Replace(structToRegister.Elem().Name(), "Controller", "", -1))
			method := structToRegister.Method(i)

			url := fmt.Sprintf("/%s/%s", structName, method)

			webContext := WebContext{
				AppInfo:        &app,
				Controller:     structName,
				MethodFunc:     strings.ToLower(method.Name),
				PureMethodFunc: strings.ToLower(method.Name),
				Url:            url,
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

			fmt.Printf("[INFO] Register route /%s/%s \n", structName, webContext.PureMethodFunc)

			Aqua.MuxRouter.HandleFunc(webContext.Url, func(w http.ResponseWriter, req *http.Request) {

				webContext.Writer = w
				webContext.Request = req

				fmt.Printf("[INFO] Incoming request /%s/%s \n", structName, webContext.PureMethodFunc)

				for _, interceptor := range interceptorFuncs {
					returnValues := interceptor.Call([]reflect.Value{reflect.ValueOf(&webContext)})
					if !returnValues[0].Bool() {
						return
					}
				}

				reflect.ValueOf(handlerStruct).MethodByName(method.Name).Call([]reflect.Value{reflect.ValueOf(&webContext)})
			}).Methods(webContext.Method)
		}
	}

	return nil
}
