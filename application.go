package aquarius

type Application struct {
	Name              string
	Path              string
	Layout            string
	Describtion       string
	ViewsPath         string
	Controllers       []interface{}
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

	return Application{Name: appName, GlobalMiddlewares: globalMiddleware, InterceptorSignal: INTERCEPTOR_RUN}
}

func (app *Application) RegisterController(handlerStruct interface{}) {
	app.Controllers = append(app.Controllers, handlerStruct)
}
