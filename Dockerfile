FROM golang:1.9

WORKDIR /go/src/github.com/sjeandeaux/github-cmd-go
COPY . .

RUN make tools
RUN make build-all

FROM scratch
COPY --from=0 /go/src/github.com/sjeandeaux/github-cmd-go/target /cmd
ENTRYPOINT ["/cmd/incrementor"] 