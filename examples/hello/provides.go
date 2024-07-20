package main

import (
	"github.com/fluxscope/fleet"
	"github.com/fluxscope/fleet/pkg/encode/json"
	"github.com/fluxscope/fleet/server/http"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	gohttp "net/http"
)

var ProvideSet = wire.NewSet(NewHttpServer, NewApp)

func NewHttpServer() *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/health", func(rw gohttp.ResponseWriter, r *gohttp.Request) {
		resp := map[string]string{
			"health": "ok",
		}
		_, _ = rw.Write(json.SafeMarshal(resp))
	})

	hs := &gohttp.Server{
		Handler: router,
		Addr:    "0.0.0.0:8181",
	}
	return http.NewServer(http.WithHTTPServer(hs))
}

func NewApp(hs *http.Server) *fleet.App {
	return fleet.NewApp(fleet.WithService(hs))
}
