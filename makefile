build_proto:
	protoc --go_out=. --go-grpc_out=. -I=proto/ proto/vodQuery/*.proto