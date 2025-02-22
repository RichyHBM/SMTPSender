package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetHttpFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("Using live mode")
		return http.FS(os.DirFS("static"))
	}

	log.Print("Using embed mode")
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

type EnsureHeaderMiddleware struct {
	header string
	value  string
}

func MakeEnsureHeaderMiddleware(header string, value string) *EnsureHeaderMiddleware {
	return &EnsureHeaderMiddleware{
		strings.ToLower(header),
		strings.ToLower(value),
	}
}

func (middleware *EnsureHeaderMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if len(middleware.header) == 0 {
		log.Println("Using empty EnsureHeaderMiddleware, can be removed")
		next(w, req)
		return
	}

	if headerValue := req.Header.Get(middleware.header); len(headerValue) > 0 {
		if len(middleware.value) == 0 || strings.ToLower(headerValue) == middleware.value {
			next(w, req)
			return
		}
	}

	w.WriteHeader(http.StatusForbidden)
}
