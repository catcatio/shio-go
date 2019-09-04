
.PHONY: test

format-proto:
	cd proto && prototool format -f -w

clean:
	rm -rf ./dist

release-webhook-dev:
	cd cmd/gcf && go mod vendor && gcloud functions deploy webhook --project shio-go-dev --entry-point HandleRequest --runtime go111 --trigger-http --region asia-east2 --verbosity=debug

release-outgoingevent-dev:
	cd cmd/gcf && go mod vendor && gcloud functions deploy outgoingevent --project shio-go-dev --entry-point HandleOutgoingEventPubsub --runtime go111 --trigger-topic outgoing-event-topic --region asia-east2 --verbosity=debug

release-incoming-dev:
	cd cmd/gcf && go mod vendor && gcloud functions deploy incomingevent --project shio-go-dev --entry-point HandleIncomingEventPubsub --runtime go111 --trigger-topic incoming-event-topic --region asia-east2 --verbosity=debug

release-fulfillment-dev:
	cd cmd/gcf && go mod vendor && gcloud functions deploy fulfillmentevent --project shio-go-dev --entry-point HandleFulfillmentPubsub --runtime go111 --trigger-topic fulfillment-topic --region asia-east2 --verbosity=debug

pubsub-topic-outgoingevent-dev:
	gcloud pubsub topics create outgoing-event-topic --project shio-go-dev

pubsub-topic-incomingevent-dev:
	gcloud pubsub topics create incoming-event-topic --project shio-go-dev

pubsub-topic-fulfillment-dev:
	gcloud pubsub topics create fulfillment-topic --project shio-go-dev

pubsub-topic-dev: pubsub-topic-incomingevent-dev pubsub-topic-outgoingevent-dev pubsub-topic-fulfillment-dev

release-dev: release-webhook-dev release-outgoingevent-dev release-incoming-dev release-fulfillment-dev
