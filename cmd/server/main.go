package main

import (
	filev1 "connect-rpc-tutorial/gen/file/v1"
	"connect-rpc-tutorial/gen/file/v1/filev1connect"
	"context"
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type FileServer struct{}

// Download implements filev1connect.FileServiceHandler.
func (fs *FileServer) Download(context.Context, *connect.Request[filev1.DownloadRequest], *connect.ServerStream[filev1.DownloadResponse]) error {
	panic("unimplemented")
}

// ListFiles implements filev1connect.FileServiceHandler.
func (fs *FileServer) ListFiles(ctx context.Context, req *connect.Request[filev1.ListFilesRequest]) (*connect.Response[filev1.ListFilesResponse], error) {
	log.Println("ListFiles was invoked")

	dir := "/Users/norikiyo/workspace/connect-rpc-tutorial/storage"

	paths, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := connect.NewResponse(&filev1.ListFilesResponse{Filenames: filenames})

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
