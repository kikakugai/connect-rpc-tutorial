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
		connect.WithGRPC(),
	)

	res, err := client.ListFiles(
		context.Background(),
		connect.NewRequest(&filev1.ListFilesRequest{}),
	)
	if err != nil {

		log.Println(err)
		return
	}

	log.Println(res.Msg.Filenames)
}
