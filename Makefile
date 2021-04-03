
PROTO_GO_OUT=proto/build/go

go-clean-proto:
	docker run --rm --entrypoint /bin/sh -v $$(pwd):/mnt alpine -c "rm -rf /mnt/$(PROTO_GO_OUT)"

go-compile-proto: go-clean-proto
	mkdir $(PROTO_GO_OUT) -p
	#docker run --rm -it --entrypoint /bin/shcd -v $$(pwd):/mnt $$(docker build -q -f ./ops/protoc.Dockerfile .)
	docker run --rm  \
	 -v $$(pwd):/mnt \
	 $$(docker build -q -f ./ops/protoc.Dockerfile .) \
	 --go_out=/mnt/$(PROTO_GO_OUT) -Iproto/ $$(find proto -iname "*.proto")

consume-demo-topic:
	kafkacat -b localhost:9094 -t demo.stocks.trades -C