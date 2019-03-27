package aquarius

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Aquarius struct {
	Server       *http.Server
	MuxRouter    http.Handler
	Port         string
	AppName      string
	Applications map[string]Application
}

func NewAquarius() Aquarius {
	return Aquarius{}
}
func (Aqua *Aquarius) Listen() {

	Aqua.Server = &http.Server{Addr: fmt.Sprintf(":%s", Aqua.Port)}

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
	if Aqua.Applications == nil {
		Aqua.Applications = make(map[string]Application)
	}

	if _, ok := Aqua.Applications[app.Name]; !ok {
		Aqua.Applications[app.Name] = app
	} else {
		return fmt.Errorf("Application name already registered")
	}

	return nil
}
