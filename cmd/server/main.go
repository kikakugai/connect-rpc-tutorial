package main

import (
	filev1 "connect-rpc-tutorial/gen/file/v1"
	"connect-rpc-tutorial/gen/file/v1/filev1connect"
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type FileServer struct{}

// ListFiles implements filev1connect.FileServiceHandler.
func (fs *FileServer) ListFiles(ctx context.Context, req *connect.Request[filev1.ListFilesRequest]) (*connect.Response[filev1.ListFilesResponse], error) {
	log.Println("ListFiles invoked")

	res := connect.NewResponse(&filev1.ListFilesResponse{Filenames: []string{"test"}})
	res.Header().Set("API-Version", "v1")

	return res, nil
}

func main() {
	fs := &FileServer{}
	mux := http.NewServeMux()
	path, handler := filev1connect.NewFileServiceHandler(fs)
	mux.Handle(path, handler)
	log.Println("server started...")
	http.ListenAndServe("localhost:8080", h2c.NewHandler(mux, &http2.Server{}))
}
