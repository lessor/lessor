FROM golang:1.10.1 as builder

WORKDIR /go/src/github.com/lessor/lessor
COPY . .

RUN make deps
RUN make
RUN mv build/lessor /bin/lessor

FROM gcr.io/distroless/base

COPY --from=builder /bin/lessor /bin/lessor

CMD ["lessor"]
