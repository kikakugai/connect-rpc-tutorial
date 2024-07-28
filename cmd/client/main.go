package main

import (
	filev1 "connect-rpc-tutorial/gen/file/v1"
	"connect-rpc-tutorial/gen/file/v1/filev1connect"
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
)

func main() {
	client := filev1connect.NewFileServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		// gRPCで通信したい場合は、以下を指定する
		// connect.WithGRPC(),
	)

	callListFiles(client)
	callDownload(client)
}

func callListFiles(client filev1connect.FileServiceClient) {
	res, err := client.ListFiles(
		context.Background(),
		connect.NewRequest(&filev1.ListFilesRequest{}),
	)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(res.Msg.Filenames)
}

func callDownload(client filev1connect.FileServiceClient) {
	req := connect.NewRequest(&filev1.DownloadRequest{Filename: "names.txt"})
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Panicln(err)
	}

	for stream.Receive() {
		res := stream.Msg()
		log.Printf("Response from Download(bytes): %v, %s", res.GetData(), string(res.GetData()))
	}
}
