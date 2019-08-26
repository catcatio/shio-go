
.PHONY: test

format-proto:
	cd proto && prototool format -f -w

clean:
	rm -rf ./dist

prepare:
	go get -u google.golang.org/grpc
	go get -u github.com/gogo/protobuf/gogoproto
	go get -u github.com/favadi/protoc-go-inject-tag
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go

gen-proto:
	prototool generate

release-webhook-dev:
	cd cmd/bot-webhook-fn && go mod vendor && gcloud beta functions deploy bot-webhook --project shio-go-dev --entry-point Webhook --runtime go111 --trigger-http --region asia-east2 --verbosity=debug
