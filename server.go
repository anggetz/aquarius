package aquarius

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var Mux = mux.NewRouter()

type Aquarius struct {
	Server *http.Server
	Port   string
	Header struct {
		Origin []string
	}
	Static struct {
		Path string
		Dir  string
	}
	AppName string
}

func NewAquarius() Aquarius {

	aquarius := Aquarius{}
	aquarius.Header.Origin = []string{"*"}

	return aquarius
}
func (Aqua *Aquarius) Listen() {

	if Aqua.Static.Path != "" {

		Mux.
			PathPrefix(Aqua.Static.Path).
			Handler(http.StripPrefix(Aqua.Static.Path, http.FileServer(http.Dir("."+Aqua.Static.Dir))))

	}

	middlewareInner := func(h http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			originAcc := r.Host

			if strings.Contains(strings.Join(Aqua.Header.Origin, ","), r.Header.Get("Origin")) {
				originAcc = r.Header.Get("Origin")
			}
			w.Header().Set("Access-Control-Allow-Origin", originAcc)

			h.ServeHTTP(w, r)
		})
	}

	Aqua.Server = &http.Server{Addr: fmt.Sprintf(":%s", Aqua.Port), Handler: middlewareInner(Mux)}

	fmt.Printf("[INFO] Starting application: localhost:%s \n", Aqua.Port)

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

	return nil
}
