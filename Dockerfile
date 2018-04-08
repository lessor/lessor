FROM golang:1.10.1 as builder

WORKDIR /go/src/github.com/lessor/lessor
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only
RUN go build -o /bin/lessor ./cmd/lessor

FROM gcr.io/distroless/base

COPY --from=builder /bin/lessor /bin/lessor

CMD ["lessor"]
