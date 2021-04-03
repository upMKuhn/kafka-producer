FROM memominsk/protobuf-alpine
RUN apk add go git \
    && go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go

ENV PATH=$PATH:/root/go/bin