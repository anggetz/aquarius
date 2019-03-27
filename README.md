# Aquarius

A simple web server golang. Super simple..

# Installation
```
go get https://github.com/gorilla/mux
go get https://github.com/anggito12345/aquarius
```

# Quick Start

```
func myHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.Write([]byte(id))
}

func main() {
	aqua := aquarius.NewAquarius()
	aqua.Port = "8089"

	aquaApp := aquarius.NewApplication("NewApp", nil)
	aquaApp.ViewsPath = "views"
	aquaApp.Layout = "layout/layout.html"

	aqua.MuxRouter.HandleFunc("/halo/{id}", myHandler)
	aqua.RegisterApp(aquaApp)
	aqua.Listen()
}

```

# Register Controller

Create file HomeController in folder controllers

Copy this code below
```
type HomeController struct {
	Middleware []interface{}
}

func NewHomeController() *HomeController {
	home := HomeController{}
	return &home
}

func (home *HomeController) Index(Aqua *aquarius.WebContext) {
	Aqua.WriteHTML("")
}

func (home *HomeController) Post_data(Aqua *aquarius.WebContext) {
	Aqua.WriteHTML("")
}

func (home *HomeController) Get_data(Aqua *aquarius.WebContext) {
	Aqua.WriteHTML("")
}
```

And then register the controller:
```
func main() {
	aqua := aquarius.NewAquarius()
	aqua.Port = "8089"

	aquaApp := aquarius.NewApplication("NewApp", nil)
	aquaApp.ViewsPath = "views"
	aquaApp.Layout = "layout/layout.html"
	aquaApp.RegisterController(controllers.NewHomeController())
	aqua.RegisterApp(aquaApp)
	aqua.Listen()
}
```
