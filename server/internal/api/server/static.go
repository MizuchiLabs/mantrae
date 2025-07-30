package server

import (
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func CompressedFileHandler(root fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(root))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalPath := r.URL.Path
		accept := r.Header.Get("Accept-Encoding")

		switch {
		case strings.Contains(accept, "br") && fileExists(root, originalPath+".br"):
			r.URL.Path += ".br"
			w.Header().Set("Content-Encoding", "br")
			setContentType(w, originalPath)

		case strings.Contains(accept, "gzip") && fileExists(root, originalPath+".gz"):
			r.URL.Path += ".gz"
			w.Header().Set("Content-Encoding", "gzip")
			setContentType(w, originalPath)
		}

		fileServer.ServeHTTP(w, r)
	})
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
	defer f.Close()
	return true
}
