protoc --proto_path=movie/api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/movie/api/proto/v1 movie-service.proto
protoc --proto_path=movie/api/proto/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/movie/api/proto/v1 movie-service.proto
protoc --proto_path=movie/api/proto/v1 --proto_path=third_party --swagger_out=logtostderr=true:movie/api/swagger/v1 movie-service.proto
