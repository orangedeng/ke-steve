// +build embed

package ui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// content holds our static web server content.
//go:embed ui/*
var staticContent embed.FS

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

func pathExist(path string) bool {
	path = formatPath(path)
	_, err := staticContent.Open(path)
	return err == nil
}

func openFile(path string) (fs.File, error) {
	path = formatPath(path)
	file, err := staticContent.Open(path)
	if err != nil {
		logrus.Errorf("openEmbedFile %s err: %v", path, err)
	}
	return file, err
}

func formatPath(path string) string {
	// To replace _nuxt/_cluster/...
	// For embed, If a pattern names a directory,
	// all files in the subtree rooted at that directory are embedded (recursively),
	// except that files with names beginning with ‘.’ or ‘_’ are excluded.
	return strings.ReplaceAll(path, "_", "")
}

func serveEmbed(basePath string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		logrus.Debugf("serveEmbed name: %s", name)
		assetPath := filepath.Join(basePath, name)
		logrus.Debugf("serveEmbed final path: %s", assetPath)
		return openFile(assetPath)
	})

	return http.FileServer(http.FS(handler))
}

func serveEmbedIndex(basePath string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		path := filepath.Join(basePath, "dashboard", "index.html")
		logrus.Debugf("serveEmbedIndex : %s", path)
		f, _ := staticContent.Open(path)
		io.Copy(rw, f)
		f.Close()
	})
}

func (u *Handler) ServeAsset() http.Handler {
	return u.middleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		serveEmbed(u.pathSetting()).ServeHTTP(rw, req)
	}))
}

func (u *Handler) ServeFaviconDashboard() http.Handler {
	return u.middleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		serveEmbed(filepath.Join(u.pathSetting(), "dashboard")).ServeHTTP(rw, req)
	}))
}

func (u *Handler) IndexFileOnNotFound() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		path := filepath.Join(u.pathSetting(), req.URL.Path)
		if pathExist(path) {
			u.ServeAsset().ServeHTTP(rw, req)
		} else {
			u.IndexFile().ServeHTTP(rw, req)
		}
	})
}

func (u *Handler) IndexFile() http.Handler {
	return u.indexMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		serveEmbedIndex(u.pathSetting()).ServeHTTP(rw, req)
	}))
}
