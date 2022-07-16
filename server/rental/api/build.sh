protoc \
    --go_out=gen/v1 --go_opt=paths=source_relative \
    --go-grpc_out=gen/v1 --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=gen/v1 --grpc-gateway_opt=paths=source_relative,grpc_api_configuration=rental.yaml \
    ./rental.proto