package aquarius

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"reflect"
)

type WebContext struct {
	AppInfo        *Application
	ControllerInfo interface{}
	Writer         http.ResponseWriter
	Request        *http.Request
	Controller     string
	MethodFunc     string
	PureMethodFunc string
	Method         string
	Url            string
}

func (aquaWebContext *WebContext) WriteHTML(templatePaths ...string) {
	fileToParse := []string{}

	fullLayoutPath := path.Join(aquaWebContext.AppInfo.ViewsPath, "layout.html")
	if aquaWebContext.AppInfo.Layout != "" {
		fullLayoutPath = path.Join(aquaWebContext.AppInfo.ViewsPath, aquaWebContext.AppInfo.Layout)
	}

	if _, err := os.Stat(fullLayoutPath); os.IsNotExist(err) {
		http.Error(aquaWebContext.Writer, fmt.Sprintf("Could not open layout file %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fileToParse = append(fileToParse, fullLayoutPath)

	for _, templatePath := range templatePaths {
		fullViewPath := path.Join(aquaWebContext.AppInfo.ViewsPath, aquaWebContext.Controller, fmt.Sprintf("%s.html", aquaWebContext.PureMethodFunc))
		if templatePath != "" {
			fullViewPath = path.Join(aquaWebContext.AppInfo.ViewsPath, templatePath)
		}

		if _, err := os.Stat(fullViewPath); os.IsNotExist(err) {
			http.Error(aquaWebContext.Writer, fmt.Sprintf("Could not open template file %s", err.Error()), http.StatusInternalServerError)
			return
		}

		fileToParse = append(fileToParse, fullViewPath)
	}

	tmpl, err := template.ParseFiles(fileToParse...)
	if err != nil {
		http.Error(aquaWebContext.Writer, fmt.Sprintf("Could not parse template file %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(aquaWebContext.Writer, "layout", nil)
	if err != nil {
		http.Error(aquaWebContext.Writer, fmt.Sprintf("Could not execute template file %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (aquaWebContext *WebContext) WriteJSON(data interface{}) {
	checkTheType := reflect.ValueOf(data)
	if checkTheType.Kind() != reflect.Map {
		http.Error(aquaWebContext.Writer, fmt.Sprintf("Could not parse data to json format"), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(data)
	if err != nil {
		http.Error(aquaWebContext.Writer, fmt.Sprintf("Could not parse data to json format %s", err.Error()), http.StatusInternalServerError)
		return
	}

	aquaWebContext.Writer.Write(j)
}
