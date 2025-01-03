social-stream-gen:
	@protoc \
		--proto_path=protobuf "protobuf/social-stream.proto" \
		--go_out=protobuf/genproto/social-stream --go_opt=paths=source_relative \
  	--go-grpc_out=protobuf/genproto/social-stream --go-grpc_opt=paths=source_relative
	