FROM golang:1.11.3-alpine3.8 as builder

RUN apk add --update ca-certificates git tar make

WORKDIR /go/src/github.com/lessor/lessor
COPY . .

RUN make deps-dep
RUN make build
RUN cp ./lessor-controller /usr/bin/lessor-controller

FROM alpine:3.8

RUN apk --update add ca-certificates

COPY --from=builder /usr/bin/lessor-controller /usr/bin/lessor-controller

CMD ["lessor-controller"]
