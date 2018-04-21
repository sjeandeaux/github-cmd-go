FROM golang:1.9

WORKDIR /go/src/github.com/sjeandeaux/toolators
COPY . .

RUN make tools
RUN make build-all

FROM scratch
COPY --from=0 /go/src/github.com/sjeandeaux/toolators/target /cmd
ENTRYPOINT ["/cmd/incrementor"] 