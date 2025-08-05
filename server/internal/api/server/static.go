package server

import (
	"net/http"

	"github.com/mizuchilabs/mantrae/proto"
	"github.com/mizuchilabs/mantrae/web"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

func (s *Server) WithStatic() {
	uploadsContent := http.FileServer(http.Dir("./data/uploads"))
	s.mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsContent))
	s.mux.Handle("/", statigz.FileServer(
		web.StaticFS,
		brotli.AddEncoding,
		statigz.FSPrefix("build"),
	))
	s.mux.Handle(
		"/openapi.yaml",
		statigz.FileServer(proto.OpenAPIFS, statigz.FSPrefix("gen/openapi")),
	)
}
