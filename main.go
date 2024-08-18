package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/MizuchiLabs/mantrae/api"
	"github.com/MizuchiLabs/mantrae/util"
	"github.com/lmittmann/tint"
)

//go:embed all:web/build
var webFS embed.FS

// Set up global logger with specified configuration
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := util.GenerateCreds(); err != nil {
		slog.Error("Failed to generate creds", "error", err)
	}
	go util.FetchTraefikConfig()
}

func main() {
	port := flag.Int("port", 3000, "Port to listen on")
	flag.Parse()

	mux := api.Routes()
	middle := api.Chain(api.Log, api.Cors)

	staticContent, err := fs.Sub(webFS, "web/build")
	if err != nil {
		slog.Error("Sub", "error", err)
		return
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))

	// Start the background sync process
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go util.Sync(wg)

	log.Println("Listening on port", *port)
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), middle(mux)); err != nil {
		slog.Error("ListenAndServe", "error", err)
		return
	}
	wg.Wait()
}
