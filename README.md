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

# Register Controller and render HTML

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
	Aqua.WriteJSON(map[string]interface{}{
		"data": "Success",
	})
}

func (home *HomeController) Get_data(Aqua *aquarius.WebContext) {
	Aqua.WriteHTML("home/hello_world.html")
}

```

## Template file

Crate layout.html in folder views:
```
{{ define "layout"}}
<!DOCTYPE html>
<html>
    <head>
        <meta  charset='UTF-8'>
        <meta name='viewport' content='width=device-width, initial-scale=1.0'>
        <meta http-equiv='X-UA-Compatible'  content='ie-edge'>
        <title>title</title>
    <body>
    </body>
        {{ template "body"}}
    </head>
</html>
{{ end }}
```

Create 2 template files

views/home/index

```
{{ define "body"}}
    test index
{{ end }}
```
and
views/home/hello_world.html
```
{{ define "body"}}
    test hello world
{{ end }}
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

Above code will automatically register all method in HomeController struct :
The route will be:
home/data 
home/index 


## Routing.

You can modify your own route for each method. 

```
type HomeController struct {
	Middleware []interface{}	
	Route      map[string]interface{}
}

func NewHomeController() *HomeController {
	home := HomeController{}	
	home.Route = map[string]interface{}{
		"Get_data": "/api/data/{name}",
	}
	return &home
}
```

The field route must be map[string]interface. Key same as method name and the value must string type.