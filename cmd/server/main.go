package main

import (
	"bytes"
	filev1 "connect-rpc-tutorial/gen/file/v1"
	"connect-rpc-tutorial/gen/file/v1/filev1connect"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type FileServer struct{}

// Upload implements filev1connect.FileServiceHandler.
func (fs *FileServer) Upload(ctx context.Context, stream *connect.BidiStream[filev1.UploadRequest, filev1.UploadResponse]) error {
	log.Println("Upload was invoked")

	var buf bytes.Buffer
	for {
		req, err := stream.Receive()
		if err == io.EOF {
			return stream.Send(&filev1.UploadResponse{Size: int32(buf.Len())})
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("Received data(bytes): %v", data)
		log.Printf("Received data(string): %v", string(data))
		buf.Write(data)
	}
}

// Download implements filev1connect.FileServiceHandler.
func (fs *FileServer) Download(ctx context.Context, req *connect.Request[filev1.DownloadRequest], stream *connect.ServerStream[filev1.DownloadResponse]) error {
	log.Println("Download was invoked")

	filename := req.Msg.Filename
	path := "/Users/norikiyo/workspace/connect-rpc-tutorial/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		sendErr := stream.Send(&filev1.DownloadResponse{Data: buf[:n]})
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}

	return nil
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
