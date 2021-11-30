package transport

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func SwaggerHttpHandler(logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	pr.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	pr.Handle("/swagger/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "doc"}
	sh1 := middleware.Redoc(opts1, nil)
	pr.Handle("/swagger/doc", sh1)

	return pr
}
