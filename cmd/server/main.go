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

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"

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
			log.Panicln(err)
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

func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}

func main() {
	mux := http.NewServeMux()
	path, handler := filev1connect.NewFileServiceHandler(&FileServer{})

	mux.Handle(path, handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: withCORS(h2c.NewHandler(mux, &http2.Server{})),
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
