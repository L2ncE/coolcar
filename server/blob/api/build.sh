protoc \
    --go_out=gen/v1 --go_opt=paths=source_relative \
    --go-grpc_out=gen/v1 --go-grpc_opt=paths=source_relative \
    ./blob.proto