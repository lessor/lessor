FROM golang:1.10.1-alpine3.7 as builder

RUN apk add --update ca-certificates git tar

WORKDIR /go/src/github.com/lessor/lessor
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only
RUN go build -o /usr/bin/lessor

FROM alpine:3.7

RUN apk --update add ca-certificates

COPY --from=builder /usr/bin/lessor /usr/bin/lessor

CMD ["lessor"]
