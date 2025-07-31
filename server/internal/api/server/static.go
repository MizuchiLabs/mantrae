package server

import (
	"io/fs"
	"log"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mizuchilabs/mantrae/web"
)

func (s *Server) WithStatic() {
	static, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}
	uploadsContent := http.FileServer(http.Dir("./data/uploads"))
	s.mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsContent))
	s.mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalPath := r.URL.Path
		accept := r.Header.Get("Accept-Encoding")

		switch {
		case strings.Contains(accept, "br") && fileExists(static, originalPath+".br"):
			r.URL.Path += ".br"
			w.Header().Set("Content-Encoding", "br")
			setContentType(w, originalPath)

		case strings.Contains(accept, "gzip") && fileExists(static, originalPath+".gz"):
			r.URL.Path += ".gz"
			w.Header().Set("Content-Encoding", "gzip")
			setContentType(w, originalPath)
		}

		http.FileServer(http.FS(static)).ServeHTTP(w, r)
	}))
}

func setContentType(w http.ResponseWriter, origPath string) {
	ext := filepath.Ext(origPath)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
}

func fileExists(fsys fs.FS, name string) bool {
	f, err := fsys.Open(strings.TrimPrefix(name, "/"))
	if err != nil {
		return false
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Error closing file", "error", err)
		}
	}()
	return true
}
