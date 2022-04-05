generate:
	bazel build //client/internal/proto:service_go_proto
	cp -f bazel-bin/client/internal/proto/service_go_proto_/github.com/tolikzinovyev/go-algorand-grpc-sdk/client/internal/proto/service.pb.go client/internal/proto/
	chmod -x client/internal/proto/service.pb.go
