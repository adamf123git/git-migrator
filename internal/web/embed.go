package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

// getStaticFS returns the embedded static file system
func getStaticFS() http.FileSystem {
	fsys, _ := fs.Sub(staticFiles, "static")
	return http.FS(fsys)
}
