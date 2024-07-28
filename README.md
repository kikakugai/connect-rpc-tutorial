# Connect RPC Tutorial

## 環境構築

1. VSCODEへ以下の拡張機能を導入する
   zxh404.vscode-proto3

1. 下記をインストールする

> <https://connectrpc.com/docs/go/getting-started/#install-tools>より抜粋

```sh
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

## 実装

1. protoファイルを新規作成or修正時には以下のコマンドを実行する

```sh
# protoファイルを元にlintを走らせた後、コードの自動生成を実施
make proto-gen
```
