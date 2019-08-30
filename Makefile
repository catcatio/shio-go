
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
	cd cmd/gcf && go mod vendor && gcloud functions deploy webhook --project shio-go-dev --entry-point HandleRequest --runtime go111 --trigger-http --region asia-east2 --verbosity=debug

release-sendmessage-dev:
	cd cmd/gcf && go mod vendor && gcloud functions deploy sendmessage --project shio-go-dev --entry-point HandleSendMessagePubsub --runtime go111 --trigger-topic send-message-topic --region asia-east2 --verbosity=debug

release-incoming-dev:
	cd cmd/gcf && go mod vendor && gcloud functions deploy incomingevent --project shio-go-dev --entry-point HandleIncomingEventPubsub --runtime go111 --trigger-topic incoming-event-topic --region asia-east2 --verbosity=debug

pubsub-topic-sendmessage-dev:
	gcloud pubsub topics create send-message-topic --project shio-go-dev

pubsub-topic-incomingevent-dev:
	gcloud pubsub topics create incoming-event-topic --project shio-go-dev

pubsub-subscription-sendmessage-dev:
	gcloud pubsub subscriptions create send-message-subscription --topic send-message-topic --project shio-go-dev

pubsub-subscription-incomingevent-dev:
	gcloud pubsub subscriptions create incoming-event-subscription --topic incoming-event-topic --project shio-go-dev

