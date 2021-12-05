# Go gRPC Services Course

install protoc (perhaps brew install protobuf)
https://github.com/protocolbuffers/protobuf/releases

note for macs with apple silicon
https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.19.1/protoc-3.19.1-osx-aarch_64.exe
is really the same exact file as the osx-x86_64 version according to https://github.com/protocolbuffers/protobuf/pull/8557

install buf.build components
http
how to install other gRPC related tooling

```bash
go install github/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.6.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.6.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.6.0
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0

ls ~/go/bin/{protoc-gen-go,protoc-gen-go-grpc,protoc-gen-grpc-gateway,protoc-gen-openapiv2}
```

* https://buf.build/
* https://github.com/bufbuild/buf/releases/download/v1.0.0-rc8/buf-Darwin-arm64
* https://github.com/bufbuild/buf/releases/download/v1.0.0-rc8/protoc-gen-buf-breaking-Darwin-arm64
* https://github.com/bufbuild/buf/releases/download/v1.0.0-rc8/protoc-gen-buf-lint-Darwin-arm64

```bash
rename/add to /usr/local/bin/ ( buf | protoc-gen-buf-breaking | protoc-gen-buf-lint )
update permissions ... `rehash` (for zsh) ... run and update security & privacy control panel
buf mod update
buf generate
go generate ./...
```

Working with Protocol Buffers in GoLand
* https://youtrack.jetbrains.com/issue/IDEA-277818
* specify path to your .proto files: Languages & Frameworks | Protocol Buffers
* Uncheck Configure automatically and add a new entry, e.g. if your .proto files are under
  `project-root/api/proto` directory, then specify its path in the settings as well.
* `/Users/ward/.cache/buf/v1/module/data/buf.build/googleapis/googleapis/<hash>`

Make use of gRPC Gateway (while leveraging cmux) to also provide a REST API

```bash
% curl http://localhost:50051/rocket.swagger.json

% curl -d '{"rocket": {"id": "6180ef22-16e8-4f43-8e81-7c8245be69f7", "name": "MISSION-01", "type": "Falcon Heavy"}}' \
  -H 'Content-Type: application/json' -X POST http://localhost:50051/v1/rocket/AddRocket
{"rocket":{"id":"6180ef22-16e8-4f43-8e81-7c8245be69f7", "name":"MISSION-01", "type":"Falcon Heavy"}}%

% curl -d '{"id": "6180ef22-16e8-4f43-8e81-7c8245be69f7"}' -H 'Content-Type: application/json' \
  -X POST http://localhost:50051/v1/rocket/GetRocket
{"rocket":{"id":"6180ef22-16e8-4f43-8e81-7c8245be69f7", "name":"MISSION-01", "type":"Falcon Heavy"}}%

% curl -d '{"rocket": {"id": "6180ef22-16e8-4f43-8e81-7c8245be69f7", "name": "MISSION-01", "type": "Falcon Heavy"}}' \
  -H 'Content-Type: application/json' -X POST http://localhost:50051/v1/rocket/DeleteRocket
{"status":"successfully deleted rocket"}
```